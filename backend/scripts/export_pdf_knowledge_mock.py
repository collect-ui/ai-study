#!/usr/bin/env python3

import argparse
import json
from pathlib import Path

import pymysql

from backfill_pdf_source_images import CONFIG_PATH, load_properties, parse_mysql_dsn


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
    return [row[0] for row in cursor.fetchall()]


def rows_for_sources(cursor, source_doc_ids: list[str]) -> list[dict]:
    if not source_doc_ids:
        return []
    placeholders = ",".join(["%s"] * len(source_doc_ids))
    cursor.execute(
        f"""
        SELECT
          unit.unit_id,
          unit.unit_code,
          unit.unit_name,
          COALESCE(NULLIF(k.subject, ''), unit.subject) AS subject,
          COALESCE(NULLIF(k.stage, ''), unit.stage) AS stage,
          COALESCE(NULLIF(k.grade, ''), unit.grade) AS grade,
          unit.textbook_version,
          unit.order_index AS unit_order,
          k.knowledge_id,
          k.knowledge_code,
          k.knowledge_name,
          k.order_index AS knowledge_order,
          c.content_id,
          c.semantic_type,
          c.section_title,
          c.content_text,
          c.content_html,
          c.order_index AS content_order,
          COALESCE(NULLIF(r.raw_quote, ''), s.raw_text, c.content_text) AS source_quote,
          COALESCE(NULLIF(f.page_no, 0), b.page_no, 1) AS page_no,
          COALESCE(NULLIF(r.confidence, 0), f.confidence, 1) AS confidence
        FROM question_knowledge_content c
        JOIN question_knowledge k ON k.knowledge_id = c.knowledge_id
        LEFT JOIN question_unit unit ON unit.unit_id = k.parent_id
        LEFT JOIN question_source_field_rel r
          ON r.entity_type = 'knowledge_content'
         AND r.entity_id = c.content_id
         AND r.field_name = 'content_text'
        LEFT JOIN question_source_fragment f ON f.source_fragment_id = r.source_fragment_id
        LEFT JOIN question_source_block b ON b.source_block_id = r.source_block_id
        LEFT JOIN question_source_snapshot s ON s.content_id = c.content_id
        WHERE c.source_doc_id IN ({placeholders})
          AND IFNULL(c.status, '') <> 'deleted'
        ORDER BY unit.order_index, unit.unit_code, k.order_index, k.knowledge_code, c.order_index, c.content_id
        """,
        tuple(source_doc_ids),
    )
    columns = [desc[0] for desc in cursor.description]
    return [dict(zip(columns, row)) for row in cursor.fetchall()]


def build_mock(rows: list[dict]) -> dict:
    units = []
    unit_map = {}
    knowledge_map = {}
    for row in rows:
        unit_id = row["unit_id"] or "unit_1"
        if unit_id not in unit_map:
            unit = {
                "unit_id": unit_id,
                "unit_code": row["unit_code"] or unit_id,
                "unit_name": row["unit_name"] or unit_id,
                "subject": row["subject"] or "english",
                "stage": row["stage"] or "primary",
                "grade": row["grade"] or "grade_3",
                "textbook_version": row["textbook_version"] or "pep",
                "order_index": int(row["unit_order"] or len(units) + 1),
                "knowledge": [],
            }
            unit_map[unit_id] = unit
            units.append(unit)

        knowledge_id = row["knowledge_id"]
        knowledge_key = (unit_id, knowledge_id)
        if knowledge_key not in knowledge_map:
            knowledge = {
                "knowledge_id": knowledge_id,
                "knowledge_code": row["knowledge_code"] or knowledge_id,
                "knowledge_name": row["knowledge_name"] or knowledge_id,
                "semantic_type": row["semantic_type"] or "knowledge_summary",
                "order_index": int(row["knowledge_order"] or len(unit_map[unit_id]["knowledge"]) + 1),
                "contents": [],
            }
            knowledge_map[knowledge_key] = knowledge
            unit_map[unit_id]["knowledge"].append(knowledge)

        knowledge_map[knowledge_key]["contents"].append(
            {
                "section_title": row["section_title"] or row["knowledge_name"] or "知识正文",
                "semantic_type": row["semantic_type"] or "knowledge_summary",
                "content_text": row["content_text"] or "",
                "content_html": row["content_html"] or "",
                "source_quote": row["source_quote"] or row["content_text"] or "",
                "source_quotes": [],
                "page_no": int(row["page_no"] or 1),
                "order_index": int(row["content_order"] or 1),
                "confidence": float(row["confidence"] or 1),
            }
        )

    return {
        "units": units,
        "questions": [],
        "issues": [],
        "question_draft_total": 0,
        "acceptance_status": "pass",
        "summary": "mock exported from current PDF knowledge rows",
    }


def main() -> int:
    parser = argparse.ArgumentParser(description="Export current PDF knowledge rows as fixed AI JSON.")
    parser.add_argument("--file-name", required=True)
    parser.add_argument("--out", required=True)
    args = parser.parse_args()

    props = load_properties(CONFIG_PATH)
    connection = pymysql.connect(**parse_mysql_dsn(props.get("dataSourceName") or ""))
    try:
        with connection.cursor() as cursor:
            source_doc_ids = source_doc_ids_for_file(cursor, args.file_name)
            rows = rows_for_sources(cursor, source_doc_ids)
    finally:
        connection.close()

    mock = build_mock(rows)
    out_path = Path(args.out)
    out_path.parent.mkdir(parents=True, exist_ok=True)
    out_path.write_text(json.dumps(mock, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(json.dumps({"source_doc_count": len(source_doc_ids), "row_count": len(rows), "unit_count": len(mock["units"]), "out": str(out_path)}, ensure_ascii=False))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
