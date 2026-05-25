#!/usr/bin/env python3

import argparse
import datetime as dt
import json
import re
from pathlib import Path

import pymysql

from pdf_pages_to_images import count_pdf_pages, render_pdf_pages


BACKEND_ROOT = Path(__file__).resolve().parents[1]
CONFIG_PATH = BACKEND_ROOT / "conf" / "application.properties"


def load_properties(path: Path) -> dict:
    props = {}
    for raw_line in path.read_text(encoding="utf-8").splitlines():
        line = raw_line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, value = line.split("=", 1)
        props[key.strip()] = value.strip()
    return props


def parse_mysql_dsn(dsn: str) -> dict:
    match = re.match(r"^([^:]+):(.*)@tcp\(([^)]+)\)/([^?]+)", dsn)
    if not match:
        raise ValueError("dataSourceName is not a supported MySQL tcp DSN")
    user, password, address, database = match.groups()
    host, _, port_text = address.partition(":")
    return {
        "host": host,
        "port": int(port_text or "3306"),
        "user": user,
        "password": password,
        "database": database,
        "charset": "utf8mb4",
        "autocommit": False,
    }


def local_file_root(props: dict) -> Path:
    root = Path(props.get("local_file_dir") or "./file_data/files")
    if not root.is_absolute():
        root = BACKEND_ROOT / root
    return root.resolve()


def file_url(props: dict, rel_path: Path) -> str:
    prefix = "/" + (props.get("file_prefix") or "/files").strip("/")
    return prefix + "/" + rel_path.as_posix().strip("/")


def source_document(connection, source_doc_id: str) -> dict:
    with connection.cursor(pymysql.cursors.DictCursor) as cursor:
        cursor.execute(
            """
            SELECT source_doc_id, file_url, page_count
            FROM question_source_document
            WHERE source_doc_id = %s
            """,
            (source_doc_id,),
        )
        row = cursor.fetchone()
    if not row:
        raise ValueError(f"source_doc_id not found: {source_doc_id}")
    return row


def upsert_page(cursor, source_doc_id: str, rendered_page: dict, url: str):
    source_page_id = f"{source_doc_id}-p{rendered_page['page']:03d}"
    cursor.execute(
        """
        INSERT INTO question_source_page (
          source_page_id, source_doc_id, page_no, page_image_url, width, height,
          extract_service, extract_params_json, raw_text, raw_html, extract_meta_json, page_hash
        ) VALUES (%s, %s, %s, %s, %s, %s, 'question.ai_pdf_text', '{}', '', '', '{}', '')
        ON DUPLICATE KEY UPDATE
          source_doc_id = VALUES(source_doc_id),
          page_no = VALUES(page_no),
          page_image_url = VALUES(page_image_url),
          width = VALUES(width),
          height = VALUES(height)
        """,
        (
            source_page_id,
            source_doc_id,
            rendered_page["page"],
            url,
            rendered_page["width"],
            rendered_page["height"],
        ),
    )
    return source_page_id


def backfill_source_doc(props: dict, connection, source_doc_id: str, dpi: int, day: str, dry_run: bool) -> dict:
    document = source_document(connection, source_doc_id)
    pdf_path = Path(document["file_url"])
    if not pdf_path.is_file():
        raise FileNotFoundError(f"source PDF not found: {pdf_path}")
    if dry_run:
        counted = count_pdf_pages(pdf_path)
        page_count = counted.get("page_count") or document.get("page_count") or 0
        return {
            "page_count": page_count,
            "rendered_count": page_count,
            "source_page_rows": 0,
            "source_block_rows": 0,
            "source_fragment_rows": 0,
            "snapshot_rows": 0,
            "field_rel_rows": 0,
        }

    rel_dir = Path("pdf-source") / day / source_doc_id
    out_dir = local_file_root(props) / rel_dir
    rendered = render_pdf_pages(
        pdf_path=pdf_path,
        out_dir=out_dir,
        dpi=dpi,
        max_pages=0,
        image_format="jpg",
    )
    pages = rendered.get("pages") or []
    if not pages:
        raise RuntimeError("PDF rendered zero pages")

    update_counts = {
        "page_count": rendered.get("page_count") or len(pages),
        "rendered_count": len(pages),
        "source_page_rows": 0,
        "source_block_rows": 0,
        "source_fragment_rows": 0,
        "snapshot_rows": 0,
        "field_rel_rows": 0,
    }
    with connection.cursor() as cursor:
        for page in pages:
            rel_path = rel_dir / Path(page["path"]).name
            url = file_url(props, rel_path)
            source_page_id = upsert_page(cursor, source_doc_id, page, url)
            update_counts["source_page_rows"] += cursor.rowcount

            cursor.execute(
                """
                UPDATE question_source_block
                SET block_image_url = %s
                WHERE source_doc_id = %s AND page_no = %s
                """,
                (url, source_doc_id, page["page"]),
            )
            update_counts["source_block_rows"] += cursor.rowcount

            cursor.execute(
                """
                UPDATE question_source_fragment
                SET source_page_id = %s
                WHERE source_doc_id = %s AND page_no = %s
                """,
                (source_page_id, source_doc_id, page["page"]),
            )
            update_counts["source_fragment_rows"] += cursor.rowcount

        cursor.execute(
            """
            UPDATE question_source_snapshot s
            JOIN question_source_block b ON b.source_block_id = s.source_block_id
            SET s.source_page_id = CONCAT(s.source_doc_id, '-p', LPAD(b.page_no, 3, '0'))
            WHERE s.source_doc_id = %s
            """,
            (source_doc_id,),
        )
        update_counts["snapshot_rows"] = cursor.rowcount

        cursor.execute(
            """
            UPDATE question_source_field_rel r
            JOIN question_source_block b ON b.source_block_id = r.source_block_id
            SET r.source_page_id = CONCAT(r.source_doc_id, '-p', LPAD(b.page_no, 3, '0'))
            WHERE r.source_doc_id = %s
            """,
            (source_doc_id,),
        )
        update_counts["field_rel_rows"] = cursor.rowcount

        cursor.execute(
            """
            UPDATE question_source_document
            SET page_count = %s
            WHERE source_doc_id = %s
            """,
            (update_counts["page_count"], source_doc_id),
        )
    connection.commit()
    return update_counts


def main() -> int:
    parser = argparse.ArgumentParser(description="Backfill rendered PDF page images for source traceability.")
    parser.add_argument("--source-doc-id", required=True)
    parser.add_argument("--dpi", type=int, default=144)
    parser.add_argument("--date", default=dt.date.today().isoformat())
    parser.add_argument("--dry-run", action="store_true")
    args = parser.parse_args()

    props = load_properties(CONFIG_PATH)
    if props.get("driverName") != "mysql":
        raise RuntimeError("Only MySQL application.properties is supported")

    connection = pymysql.connect(**parse_mysql_dsn(props.get("dataSourceName") or ""))
    try:
        result = backfill_source_doc(
            props=props,
            connection=connection,
            source_doc_id=args.source_doc_id,
            dpi=args.dpi,
            day=args.date,
            dry_run=args.dry_run,
        )
    except Exception:
        connection.rollback()
        raise
    finally:
        connection.close()

    print(json.dumps(result, ensure_ascii=False))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
