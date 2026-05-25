#!/usr/bin/env python3

import argparse
import json
import shutil
from pathlib import Path

import pymysql

from backfill_pdf_source_images import CONFIG_PATH, load_properties, local_file_root, parse_mysql_dsn


def source_doc_ids_for_file(cursor, file_name: str) -> list[str]:
    cursor.execute(
        """
        SELECT source_doc_id
        FROM question_source_document
        WHERE file_name = %s OR file_url LIKE %s
        ORDER BY create_time, source_doc_id
        """,
        (file_name, f"%{file_name}"),
    )
    return [str(row[0]) for row in cursor.fetchall() if row and row[0]]


def fetch_column(cursor, sql: str, params: tuple) -> list[str]:
    cursor.execute(sql, params)
    return [str(row[0]) for row in cursor.fetchall() if row and row[0]]


def delete_in(cursor, table: str, key: str, values: list[str]) -> int:
    if not values:
        return 0
    placeholders = ",".join(["%s"] * len(values))
    cursor.execute(f"DELETE FROM {table} WHERE {key} IN ({placeholders})", values)
    return cursor.rowcount


def cleanup_source_doc(connection, source_doc_id: str, apply: bool) -> dict:
    counts = {}
    with connection.cursor() as cursor:
        batch_ids = fetch_column(
            cursor,
            "SELECT import_batch_id FROM question_source_document WHERE source_doc_id = %s",
            (source_doc_id,),
        )
        content_ids = fetch_column(
            cursor,
            "SELECT content_id FROM question_knowledge_content WHERE source_doc_id = %s",
            (source_doc_id,),
        )
        knowledge_ids = fetch_column(
            cursor,
            """
            SELECT DISTINCT knowledge_id
            FROM question_knowledge_content
            WHERE source_doc_id = %s AND knowledge_id <> ''
            """,
            (source_doc_id,),
        )
        unit_ids = fetch_column(
            cursor,
            f"""
            SELECT DISTINCT parent_id
            FROM question_knowledge
            WHERE knowledge_id IN ({",".join(["%s"] * len(knowledge_ids))})
              AND parent_id <> ''
            """,
            tuple(knowledge_ids),
        ) if knowledge_ids else []

        counts.update(
            {
                "batch_ids": len(batch_ids),
                "content_ids": len(content_ids),
                "knowledge_ids": len(knowledge_ids),
                "unit_ids": len(unit_ids),
            }
        )
        if not apply:
            return counts

        for table in [
            "question_source_field_rel",
            "question_source_snapshot",
            "question_source_rel",
            "question_pdf_parse_issue",
        ]:
            cursor.execute(f"DELETE FROM {table} WHERE source_doc_id = %s", (source_doc_id,))
            counts[table] = cursor.rowcount

        cursor.execute("DELETE FROM question_knowledge_content WHERE source_doc_id = %s", (source_doc_id,))
        counts["question_knowledge_content"] = cursor.rowcount

        counts["question_knowledge"] = delete_in(cursor, "question_knowledge", "knowledge_id", knowledge_ids)

        if unit_ids:
            placeholders = ",".join(["%s"] * len(unit_ids))
            cursor.execute(
                f"""
                DELETE u FROM question_unit u
                LEFT JOIN question_knowledge k
                  ON k.parent_id = u.unit_id AND IFNULL(k.is_delete, '0') = '0'
                WHERE u.unit_id IN ({placeholders})
                  AND u.create_user = 'pdf_import'
                  AND k.knowledge_id IS NULL
                """,
                tuple(unit_ids),
            )
            counts["question_unit"] = cursor.rowcount
        else:
            counts["question_unit"] = 0

        for table in [
            "question_source_fragment",
            "question_source_block",
            "question_source_page",
            "question_source_document",
        ]:
            cursor.execute(f"DELETE FROM {table} WHERE source_doc_id = %s", (source_doc_id,))
            counts[table] = cursor.rowcount

        counts["question_import_batch"] = delete_in(cursor, "question_import_batch", "batch_id", batch_ids)
    connection.commit()
    return counts


def remove_source_assets(props: dict, source_doc_id: str) -> list[str]:
    removed = []
    root = local_file_root(props) / "pdf-source"
    if not root.exists():
        return removed
    for candidate in root.glob(f"*/{source_doc_id}"):
        if candidate.is_dir() and candidate.parent.parent == root:
            shutil.rmtree(candidate)
            removed.append(str(candidate))
    return removed


def main() -> int:
    parser = argparse.ArgumentParser(description="Clean PDF knowledge import data for one source document.")
    parser.add_argument("--source-doc-id", action="append", default=[])
    parser.add_argument("--file-name", default="")
    parser.add_argument("--apply", action="store_true")
    parser.add_argument("--remove-assets", action="store_true")
    args = parser.parse_args()

    props = load_properties(CONFIG_PATH)
    connection = pymysql.connect(**parse_mysql_dsn(props.get("dataSourceName") or ""))
    try:
        source_doc_ids = list(args.source_doc_id)
        if args.file_name:
            with connection.cursor() as cursor:
                for source_doc_id in source_doc_ids_for_file(cursor, args.file_name):
                    if source_doc_id not in source_doc_ids:
                        source_doc_ids.append(source_doc_id)
        if not source_doc_ids:
            raise ValueError("provide --source-doc-id or --file-name")

        result = {"source_doc_ids": source_doc_ids, "sources": {}}
        for source_doc_id in source_doc_ids:
            source_result = cleanup_source_doc(connection, source_doc_id, args.apply)
            if args.apply and args.remove_assets:
                source_result["removed_asset_dirs"] = remove_source_assets(props, source_doc_id)
            result["sources"][source_doc_id] = source_result
    except Exception:
        connection.rollback()
        raise
    finally:
        connection.close()

    print(json.dumps(result, ensure_ascii=False))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
