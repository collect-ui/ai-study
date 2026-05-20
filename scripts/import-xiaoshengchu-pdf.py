#!/usr/bin/env python3
import argparse
import hashlib
import html
import json
import re
import shutil
import sqlite3
import uuid
from collections import Counter
from datetime import datetime
from pathlib import Path


DEFAULT_TEXT = Path("test-results/full-pdf-import/pdf-full.txt")
DEFAULT_DB = Path("backend/database/ai_study_admin.db")
DEFAULT_OUT = Path("test-results/full-pdf-import")
SOURCE_FILE = "【小升初】英语复习题三十套全国通用（含详细解析）.pdf"
NAMESPACE = uuid.UUID("50d338d3-0df1-4aca-85a4-3719f48c0f94")
FALLBACK_ANSWERS = {
    "XSCT-P31-S02-READING_CHOICE": [
        {"answer": "C", "analysis": "PDF 未提供本题答案；根据第一段 Sam 生病呻吟导致作者无法继续睡觉，选 C。"},
        {"answer": "B", "analysis": "PDF 未提供本题答案；第三段 vehicle 指前文到达的 ambulance，选 B。"},
        {"answer": "B", "analysis": "PDF 未提供本题答案；第四段说明 Sam 因及时送医治疗而不需要手术，选 B。"},
        {"answer": "B", "analysis": "PDF 未提供本题答案；Sam 服药后病情继续恶化，因此 B 项不符合原文。"},
    ]
}


PAPER_RE = re.compile(r"^（[一二三四五六七八九十]+）$")
QUESTION_START_RE = re.compile(r"^(?:[\(（]\s*[\)）]\s*)?(\d+)\s*[\.．、]?\s*(.*)$")


def clean_text(value):
    value = str(value or "")
    value = value.replace("\u3000", " ")
    value = re.sub(r"[ \t]+", " ", value)
    value = re.sub(r"\s+\n", "\n", value)
    value = re.sub(r"\n\s+", "\n", value)
    return value.strip()


def compact_text(lines):
    return clean_text("\n".join(line for line in lines if str(line).strip()))


def inline_text(lines):
    return clean_text(" ".join(line for line in lines if str(line).strip()))


def rich_html(text):
    text = clean_text(text)
    if not text:
        return ""
    parts = [part.strip() for part in text.splitlines()]
    return "".join(f"<p>{html.escape(part)}</p>" for part in parts if part)


def answer_letters(value):
    raw = str(value or "").upper().replace("，", ",")
    letters = re.findall(r"[A-D]", raw)
    return list(dict.fromkeys(letters))


def normalize_answer(value):
    return clean_text(str(value or "").replace("，", ","))


def is_answer_marker(line):
    return line.startswith("答案")


def is_heading(line):
    return bool(
        re.match(r"^(一|二|三|四)、", line)
        or re.match(
            r"^(单项选择|单项选择，|语法精炼|语法精练|用适当的句子完成对话。?|阅读下面文章并回答问题|阅读文章回答问题)$",
            line,
        )
    )


def answer_like(line):
    return bool(
        (
            re.match(r"^(\d+|[\(（]\s*[\)）]\s*\d+)\s*[\.．]?", line)
            and ("解析" in line or len(line) < 140 or re.match(r"^\d+\s*[,A-DTFa-z]", line))
        )
        or line.startswith("Keys:")
        or re.match(r"^\d+[A-D]\s*解析", line)
    )


def split_sections(paper_lines):
    sections = []
    pending = []
    zones = []
    current = None
    answer_lines = []
    state = "exercise"

    def flush_current():
        nonlocal current
        if current and current["body"]:
            current["index"] = len(sections)
            sections.append(current)
            pending.append(current["index"])
        current = None

    def flush_answer():
        nonlocal answer_lines, pending
        if answer_lines:
            zones.append({"section_indexes": pending[:], "lines": answer_lines[:]})
        answer_lines = []
        pending = []

    i = 0
    while i < len(paper_lines):
        line = paper_lines[i].strip()
        if not line:
            i += 1
            continue

        if state == "exercise":
            if is_answer_marker(line):
                flush_current()
                state = "answer"
                answer_lines = []
                i += 1
                continue

            if is_heading(line):
                title = line
                if re.fullmatch(r"[一二三四]、", title) and i + 1 < len(paper_lines):
                    title += paper_lines[i + 1].strip()
                    i += 1
                if current and current["body"]:
                    flush_current()
                current = {"title": title, "body": []}
            else:
                if current is None:
                    current = {"title": line, "body": []}
                else:
                    current["body"].append(line)
        else:
            next_line = ""
            for j in range(i + 1, len(paper_lines)):
                if paper_lines[j].strip():
                    next_line = paper_lines[j].strip()
                    break

            fresh_exercise = len(pending) <= 1 and is_heading(line) and not answer_like(next_line)
            if fresh_exercise:
                flush_answer()
                state = "exercise"
                current = {"title": line, "body": []}
            else:
                answer_lines.append(line)
        i += 1

    if state == "exercise":
        flush_current()
    else:
        flush_answer()

    return sections, zones


def split_answer_groups(lines, section_count):
    lines = [line for line in lines if not is_answer_marker(line)]
    if section_count <= 1:
        return [lines]

    groups = []
    current = []
    saw_heading = False
    for line in lines:
        if is_heading(line):
            saw_heading = True
            if current:
                groups.append(current)
                current = []
            continue
        current.append(line)
    if current:
        groups.append(current)

    if saw_heading:
        return groups[:section_count] + [[] for _ in range(max(0, section_count - len(groups)))]
    return [lines] + [[] for _ in range(section_count - 1)]


def parse_answer_entries(lines):
    entries = []
    by_no = {}
    current = None

    def add_entry(no, answer, analysis=""):
        nonlocal current
        no = str(no)
        answer = normalize_answer(answer)
        if no in by_no:
            entry = by_no[no]
            if answer and not entry.get("answer"):
                entry["answer"] = answer
            if analysis and not entry.get("analysis"):
                entry["analysis"] = clean_text(analysis)
            current = entry
            return entry
        entry = {"no": no, "answer": answer, "analysis": clean_text(analysis)}
        entries.append(entry)
        by_no[no] = entry
        current = entry
        return entry

    for raw in lines:
        line = raw.strip()
        if not line or is_heading(line) or is_answer_marker(line):
            continue

        match = re.search(r"Keys\s*:\s*([A-D\s]+)", line, re.IGNORECASE)
        if match:
            keys = re.sub(r"\s+", "", match.group(1))
            for index, key in enumerate(keys, 1):
                add_entry(index, key)
            current = None
            continue

        match = re.match(
            r"^(\d+)\s*[\.．、]?\s*([A-D](?:\s*/\s*[A-D])?(?:\s*[,，]\s*[A-D])?|true|false|T|F)\s*(?:答案：|解析[:：]?|$)\s*(.*)$",
            line,
            re.IGNORECASE,
        )
        if not match:
            match = re.match(r"^(\d+)\s*([A-D])\s*解析[:：]?\s*(.*)$", line)
        if match:
            add_entry(match.group(1), match.group(2).replace(" ", ""), match.group(3))
            continue

        match = re.match(r"^(\d+)\s*[\.．、]?\s*(.*?)\s*解析[:：]?\s*(.*)$", line)
        if match and len(match.group(2).strip()) < 160:
            add_entry(match.group(1), match.group(2).strip().rstrip("。"), match.group(3))
            continue

        match = re.match(r"^(\d+)\s*[\.．、]?\s*(.+)$", line)
        if match and len(match.group(2)) < 180:
            add_entry(match.group(1), match.group(2).strip())
            continue

        if current:
            current["analysis"] = clean_text(f"{current.get('analysis', '')} {line}")

    return entries


def normalize_option_source(text):
    return (
        text.replace("．", ".")
        .replace(" .", ".")
        .replace("。", "。")
        .replace("B .", "B.")
        .replace("C .", "C.")
        .replace("D .", "D.")
        .replace("A .", "A.")
    )


def option_positions(text):
    text = normalize_option_source(text)
    result = []
    pattern = r"(?:^|\s)([A-D])\s*(?:[\.]\s*|\s+)(?=[A-Za-z0-9_/\(（\u4e00-\u9fff])"
    for match in re.finditer(pattern, text):
        result.append((match.start(1), match.end(), match.group(1)))

    out = []
    seen = set()
    for item in sorted(result):
        if item[2] not in seen:
            out.append(item)
            seen.add(item[2])
    return out


def parse_choice_chunk(no, chunk):
    text = inline_text(chunk)
    text = re.sub(
        r"^(?:[\(（]\s*[\)）]\s*)?" + re.escape(str(no)) + r"\s*[\.．、]?\s*",
        "",
        text,
    ).strip()
    normalized = normalize_option_source(text)
    positions = option_positions(normalized)
    options = {}
    stem = normalized
    if len(positions) >= 2:
        stem = normalized[: positions[0][0]].strip()
        for index, (start, end, label) in enumerate(positions):
            next_start = positions[index + 1][0] if index + 1 < len(positions) else len(normalized)
            options[label] = normalized[end:next_start].strip().strip("。．")
    return {"no": str(no), "stem": clean_text(stem), "options": options, "raw": clean_text(normalized)}


def find_question_starts(body):
    starts = []
    last_number = None
    for index, line in enumerate(body):
        match = QUESTION_START_RE.match(line)
        if not match:
            continue
        number = int(match.group(1))
        if last_number is not None and number < last_number:
            continue
        starts.append((index, match.group(1), match.group(2).strip()))
        last_number = number
    return starts


def parse_numbered_items(body):
    starts = find_question_starts(body)
    items = []
    for start_index, (line_index, no, _rest) in enumerate(starts):
        end = starts[start_index + 1][0] if start_index + 1 < len(starts) else len(body)
        chunk = body[line_index:end]
        text = compact_text(chunk)
        text = re.sub(
            r"^(?:[\(（]\s*[\)）]\s*)?" + re.escape(str(no)) + r"\s*[\.．、]?\s*",
            "",
            text,
        ).strip()
        items.append({"no": str(no), "text": clean_text(text), "chunk": chunk, "choice": parse_choice_chunk(no, chunk)})
    return items


def parse_choice_items(body):
    return [item["choice"] for item in parse_numbered_items(body) if item["choice"]["options"]]


def base_category(title):
    text = title.replace(" ", "")
    if "完形" in text:
        return "cloze_choice", "single_choice"
    if "阅读" in text and ("回答" in text or "短文回答" in text or "文章回答" in text):
        return "reading_short_answer", "short_answer"
    if "阅读" in text or "普通阅读" in text:
        return "reading_choice", "single_choice"
    if "判断" in text:
        return "judge_tf", "judge"
    if "对话" in text:
        return "dialogue_completion", "blank"
    if "智力" in text:
        return "iq", "single_choice"
    if any(key in text for key in ["所给词", "适当形式", "单词拼写", "汉语", "提示完成", "补充单词", "选词"]):
        return "fill_word", "blank"
    if any(key in text for key in ["单项", "单选", "语法", "选出不同"]):
        return "grammar_choice", "single_choice"
    return "normal", "single_choice"


def split_papers(text_path):
    lines = []
    for raw in text_path.read_text(encoding="utf-8").splitlines():
        line = raw.strip()
        if not line or line.startswith("--- PAGE"):
            continue
        lines.append(line)
    if lines and "小升初英语复习题" in lines[0]:
        lines = lines[1:]

    marks = [(index, line) for index, line in enumerate(lines) if PAPER_RE.fullmatch(line)]
    papers = []
    for paper_index, (start, mark) in enumerate(marks, 1):
        end = marks[paper_index][0] if paper_index < len(marks) else len(lines)
        sections, zones = split_sections(lines[start + 1 : end])
        for zone in zones:
            groups = split_answer_groups(zone["lines"], len(zone["section_indexes"]))
            for section_index, answer_lines in zip(zone["section_indexes"], groups):
                sections[section_index]["answers"] = parse_answer_entries(answer_lines)
                sections[section_index]["answer_lines"] = answer_lines
        papers.append({"paper_index": paper_index, "mark": mark, "sections": sections})
    return papers


def answer_for(items, answers, index, no):
    by_no = {entry["no"]: entry for entry in answers}
    if str(no) in by_no:
        return by_no[str(no)]
    if index < len(answers):
        return answers[index]
    return {"no": str(no), "answer": "", "analysis": ""}


def first_choice_line_index(body):
    for item in parse_numbered_items(body):
        if item["choice"]["options"]:
            return body.index(item["chunk"][0])
    return None


def make_question(source_key, sequence_no, title, category, qtype, stem, answer_entry, options=None, choice_items=None, blank_answers=None):
    options = options or {}
    choice_items = choice_items or []
    blank_answers = blank_answers or []
    answer_text = normalize_answer(answer_entry.get("answer", "")) if isinstance(answer_entry, dict) else normalize_answer(answer_entry)
    analysis_text = clean_text(answer_entry.get("analysis", "")) if isinstance(answer_entry, dict) else ""
    letters = answer_letters(answer_text)
    if options and len(letters) > 1:
        qtype = "multiple_choice"
    if category in {"cloze_choice", "reading_choice"}:
        qtype = "single_choice"
    if category == "judge_tf":
        qtype = "judge"

    question_id = str(uuid.uuid5(NAMESPACE, source_key))
    content_hash = hashlib.sha256(f"{stem}\n{answer_text}\n{json.dumps(choice_items, ensure_ascii=False)}".encode("utf-8")).hexdigest()

    return {
        "question_id": question_id,
        "question_code": source_key,
        "title": title,
        "subject": "english",
        "stage": "junior",
        "grade": "grade_7",
        "textbook_version": "pep",
        "unit_id": "",
        "unit_code": "",
        "unit_name": "",
        "question_type": qtype,
        "question_category": category,
        "difficulty": "basic",
        "score": 5,
        "duration_seconds": 0,
        "sequence_no": sequence_no,
        "stem_text": clean_text(stem),
        "stem_html": rich_html(stem),
        "analysis_text": analysis_text,
        "analysis_html": rich_html(analysis_text),
        "analysis_media_url": "",
        "analysis_media_name": "",
        "analysis_media_type": "",
        "option_count": len(options) if options and category not in {"cloze_choice", "reading_choice"} else (2 if category == "judge_tf" else 0),
        "blank_count": len(blank_answers) if blank_answers else (len(choice_items) if category in {"cloze_choice", "reading_choice"} else 0),
        "asset_count": 0,
        "content_hash": content_hash,
        "source": "pdf_import",
        "status": "draft",
        "version": 1,
        "publish_time": "",
        "publish_user": "",
        "last_review_id": "",
        "remark": SOURCE_FILE,
        "answer_text": answer_text,
        "answer_value": json.dumps(letters if letters else ([answer_text] if answer_text else []), ensure_ascii=False),
        "reference_text": json.dumps(choice_items, ensure_ascii=False) if choice_items else "",
        "options": options,
        "blank_answers": blank_answers,
    }


def build_import_rows(papers):
    rows = []
    warnings = []
    sequence = 1

    for paper in papers:
        paper_no = paper["paper_index"]
        for section in paper["sections"]:
            section_no = section["index"] + 1
            title = section["title"]
            category, qtype = base_category(title)
            body = section["body"]
            answers = section.get("answers", [])
            numbered = parse_numbered_items(body)
            choices = [item["choice"] for item in numbered if item["choice"]["options"]]

            if category == "grammar_choice" and not choices:
                category, qtype = "fill_word", "blank"

            if category in {"grammar_choice"}:
                for index, choice in enumerate(choices):
                    answer = answer_for(choices, answers, index, choice["no"])
                    stem = choice["stem"] or "选出不同的一项。"
                    source_key = f"XSCT-P{paper_no:02d}-S{section_no:02d}-Q{choice['no']}"
                    rows.append(
                        make_question(
                            source_key,
                            sequence,
                            f"第{paper_no}套 {title} 第{choice['no']}题",
                            category,
                            qtype,
                            stem,
                            answer,
                            options=choice["options"],
                        )
                    )
                    sequence += 1
                continue

            if category == "iq":
                for index, item in enumerate(numbered):
                    answer = answer_for(numbered, answers, index, item["no"])
                    choice = item["choice"]
                    item_category = "iq"
                    item_type = "single_choice" if choice["options"] else "short_answer"
                    stem = choice["stem"] if choice["options"] else item["text"]
                    source_key = f"XSCT-P{paper_no:02d}-S{section_no:02d}-Q{item['no']}"
                    rows.append(
                        make_question(
                            source_key,
                            sequence,
                            f"第{paper_no}套 {title} 第{item['no']}题",
                            item_category,
                            item_type,
                            stem,
                            answer,
                            options=choice["options"],
                        )
                    )
                    sequence += 1
                continue

            if category in {"fill_word"}:
                for index, item in enumerate(numbered):
                    answer = answer_for(numbered, answers, index, item["no"])
                    blank_answers = []
                    if answer.get("answer"):
                        blank_answers.append(
                            {
                                "blank_index": 1,
                                "standard_answer": answer["answer"],
                                "alternative_answers": "[]",
                                "score": 0,
                                "match_mode": "exact",
                                "case_sensitive": "0",
                            }
                        )
                    source_key = f"XSCT-P{paper_no:02d}-S{section_no:02d}-Q{item['no']}"
                    rows.append(
                        make_question(
                            source_key,
                            sequence,
                            f"第{paper_no}套 {title} 第{item['no']}题",
                            category,
                            "blank",
                            item["text"],
                            answer,
                            blank_answers=blank_answers,
                        )
                    )
                    sequence += 1
                continue

            if category == "dialogue_completion":
                blank_answers = []
                for index, answer in enumerate(answers, 1):
                    if not answer.get("answer"):
                        continue
                    blank_answers.append(
                        {
                            "blank_index": index,
                            "standard_answer": answer["answer"],
                            "alternative_answers": "[]",
                            "score": 0,
                            "match_mode": "exact",
                            "case_sensitive": "0",
                        }
                    )
                source_key = f"XSCT-P{paper_no:02d}-S{section_no:02d}-DIALOGUE"
                rows.append(
                    make_question(
                        source_key,
                        sequence,
                        f"第{paper_no}套 {title}",
                        category,
                        "blank",
                        compact_text(body),
                        {"answer": "；".join(item["standard_answer"] for item in blank_answers), "analysis": ""},
                        blank_answers=blank_answers,
                    )
                )
                sequence += 1
                continue

            if category == "judge_tf":
                first_index = numbered[0]["chunk"][0] if numbered else None
                passage = compact_text(body[: body.index(first_index)]) if first_index in body else ""
                for index, item in enumerate(numbered):
                    answer = answer_for(numbered, answers, index, item["no"])
                    raw = normalize_answer(answer.get("answer", "")).lower()
                    judge_answer = "true" if raw in {"t", "true", "对", "正确"} else "false"
                    source_key = f"XSCT-P{paper_no:02d}-S{section_no:02d}-J{item['no']}"
                    rows.append(
                        make_question(
                            source_key,
                            sequence,
                            f"第{paper_no}套 {title} 第{item['no']}题",
                            category,
                            "judge",
                            clean_text(f"{passage}\n\n{item['text']}"),
                            {"answer": judge_answer, "analysis": answer.get("analysis", "")},
                        )
                    )
                    sequence += 1
                continue

            if category in {"cloze_choice", "reading_choice"}:
                first_idx = first_choice_line_index(body)
                passage = compact_text(body[:first_idx]) if first_idx is not None else compact_text(body)
                choice_items = []
                source_key = f"XSCT-P{paper_no:02d}-S{section_no:02d}-{category.upper()}"
                fallback_answers = FALLBACK_ANSWERS.get(source_key, [])
                for index, choice in enumerate(choices):
                    answer = answer_for(choices, answers, index, choice["no"])
                    if not answer.get("answer") and index < len(fallback_answers):
                        answer = {"no": choice["no"], **fallback_answers[index]}
                    if category == "cloze_choice":
                        row = {
                            "__rowId": f"choice_{index + 1}",
                            "blank_no": choice["no"],
                            "option_a": choice["options"].get("A", ""),
                            "option_b": choice["options"].get("B", ""),
                            "option_c": choice["options"].get("C", ""),
                            "option_d": choice["options"].get("D", ""),
                            "answer_key": normalize_answer(answer.get("answer", "")),
                            "analysis": clean_text(answer.get("analysis", "")),
                        }
                    else:
                        row = {
                            "__rowId": f"reading_{index + 1}",
                            "sub_no": choice["no"],
                            "question_text": choice["stem"],
                            "option_a": choice["options"].get("A", ""),
                            "option_b": choice["options"].get("B", ""),
                            "option_c": choice["options"].get("C", ""),
                            "option_d": choice["options"].get("D", ""),
                            "answer_key": normalize_answer(answer.get("answer", "")),
                            "analysis": clean_text(answer.get("analysis", "")),
                        }
                    choice_items.append(row)
                answer_text = "；".join(
                    f"{item.get('blank_no') or item.get('sub_no')}:{item.get('answer_key', '')}" for item in choice_items if item.get("answer_key")
                )
                rows.append(
                    make_question(
                        source_key,
                        sequence,
                        f"第{paper_no}套 {title}",
                        category,
                        "single_choice",
                        passage,
                        {"answer": answer_text, "analysis": ""},
                        choice_items=choice_items,
                    )
                )
                if fallback_answers:
                    warnings.append(f"{source_key} answers inferred because the PDF answer block is missing")
                if not choice_items or any(not item.get("answer_key") for item in choice_items):
                    warnings.append(f"{source_key} has missing choice item answers")
                sequence += 1
                continue

            if category == "reading_short_answer":
                first_chunk = numbered[0]["chunk"][0] if numbered else None
                passage = compact_text(body[: body.index(first_chunk)]) if first_chunk in body else compact_text(body)
                for index, item in enumerate(numbered):
                    answer = answer_for(numbered, answers, index, item["no"])
                    source_key = f"XSCT-P{paper_no:02d}-S{section_no:02d}-SA{item['no']}"
                    rows.append(
                        make_question(
                            source_key,
                            sequence,
                            f"第{paper_no}套 {title} 第{item['no']}题",
                            category,
                            "short_answer",
                            clean_text(f"{passage}\n\n{item['text']}"),
                            answer,
                        )
                    )
                    sequence += 1
                continue

            warnings.append(f"Skipped section P{paper_no} S{section_no}: {title}")

    return rows, warnings


def insert_rows(db_path, rows, out_dir):
    now = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    stamp = datetime.now().strftime("%Y%m%d-%H%M%S")
    backup = db_path.with_name(f"{db_path.name}.bak.before_xiaoshengchu_import_{stamp}.sqlitebackup")
    shutil.copy2(db_path, backup)

    batch_id = str(uuid.uuid5(NAMESPACE, f"batch:{SOURCE_FILE}:{len(rows)}"))
    conn = sqlite3.connect(db_path)
    try:
        cur = conn.cursor()
        cur.execute("BEGIN")
        for table in [
            "question_knowledge_rel",
            "question_scoring_point",
            "question_asset",
            "question_blank_answer",
            "question_answer",
            "question_option",
            "question_review_record",
            "question_change_log",
            "question_import_row",
            "question_import_batch",
            "question_item",
        ]:
            cur.execute(f"DELETE FROM {table}")

        cur.execute(
            """
            INSERT INTO question_import_batch (
              batch_id, file_name, file_url, subject, stage, grade, textbook_version,
              status, total_count, success_count, fail_count, error_summary,
              create_time, create_user, modify_time, modify_user
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """,
            (
                batch_id,
                SOURCE_FILE,
                "",
                "english",
                "junior",
                "grade_7",
                "pep",
                "success",
                len(rows),
                len(rows),
                0,
                "",
                now,
                "import",
                now,
                "import",
            ),
        )

        question_fields = [
            "question_id",
            "question_code",
            "title",
            "subject",
            "stage",
            "grade",
            "textbook_version",
            "unit_id",
            "unit_code",
            "unit_name",
            "question_type",
            "difficulty",
            "score",
            "duration_seconds",
            "sequence_no",
            "stem_html",
            "stem_text",
            "analysis_html",
            "analysis_text",
            "option_count",
            "blank_count",
            "asset_count",
            "content_hash",
            "source",
            "status",
            "version",
            "publish_time",
            "publish_user",
            "last_review_id",
            "remark",
            "is_delete",
            "create_time",
            "create_user",
            "modify_time",
            "modify_user",
            "question_category",
            "analysis_media_url",
            "analysis_media_name",
            "analysis_media_type",
        ]
        placeholders = ",".join(["?"] * len(question_fields))
        for row_index, row in enumerate(rows, 1):
            cur.execute(
                f"INSERT INTO question_item ({','.join(question_fields)}) VALUES ({placeholders})",
                [row.get(field, "") for field in question_fields[:30]]
                + ["0", now, "import", now, "import", row["question_category"], "", "", ""],
            )

            for option_order, option_key in enumerate(["A", "B", "C", "D"], 1):
                option_text = row.get("options", {}).get(option_key, "")
                if not option_text:
                    continue
                option_id = str(uuid.uuid5(NAMESPACE, f"{row['question_id']}:option:{option_key}"))
                correct_letters = answer_letters(row.get("answer_text", ""))
                cur.execute(
                    """
                    INSERT INTO question_option (
                      option_id, question_id, option_key, option_order, content_mode,
                      option_html, option_text, is_correct, asset_count, is_delete,
                      create_time, create_user, modify_time, modify_user
                    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                    """,
                    (
                        option_id,
                        row["question_id"],
                        option_key,
                        option_order,
                        "rich",
                        rich_html(option_text),
                        option_text,
                        "1" if option_key in correct_letters else "0",
                        0,
                        "0",
                        now,
                        "import",
                        now,
                        "import",
                    ),
                )

            answer_id = str(uuid.uuid5(NAMESPACE, f"{row['question_id']}:answer"))
            cur.execute(
                """
                INSERT INTO question_answer (
                  answer_id, question_id, answer_type, answer_value, answer_text,
                  reference_text, case_sensitive, allow_order_change, auto_grading_rule,
                  is_delete, create_time, create_user, modify_time, modify_user
                ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                """,
                (
                    answer_id,
                    row["question_id"],
                    row["question_type"],
                    row.get("answer_value", "[]"),
                    row.get("answer_text", ""),
                    row.get("reference_text", ""),
                    "0",
                    "0",
                    "",
                    "0",
                    now,
                    "import",
                    now,
                    "import",
                ),
            )

            for blank_index, blank in enumerate(row.get("blank_answers", []), 1):
                blank_id = str(uuid.uuid5(NAMESPACE, f"{row['question_id']}:blank:{blank_index}"))
                cur.execute(
                    """
                    INSERT INTO question_blank_answer (
                      blank_answer_id, question_id, blank_index, standard_answer,
                      alternative_answers, score, match_mode, case_sensitive,
                      is_delete, create_time, create_user, modify_time, modify_user
                    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                    """,
                    (
                        blank_id,
                        row["question_id"],
                        blank_index,
                        blank.get("standard_answer", ""),
                        blank.get("alternative_answers", "[]"),
                        int(blank.get("score", 0) or 0),
                        blank.get("match_mode", "exact"),
                        blank.get("case_sensitive", "0"),
                        "0",
                        now,
                        "import",
                        now,
                        "import",
                    ),
                )

            cur.execute(
                """
                INSERT INTO question_import_row (
                  row_id, batch_id, row_index, raw_json, parsed_json, validate_status,
                  error_msg, question_id, create_time, create_user
                ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                """,
                (
                    str(uuid.uuid5(NAMESPACE, f"{batch_id}:row:{row_index}")),
                    batch_id,
                    row_index,
                    json.dumps({"source": SOURCE_FILE, "question_code": row["question_code"]}, ensure_ascii=False),
                    json.dumps(row, ensure_ascii=False),
                    "success",
                    "",
                    row["question_id"],
                    now,
                    "import",
                ),
            )

        conn.commit()
    except Exception:
        conn.rollback()
        raise
    finally:
        conn.close()

    report = {"backup": str(backup), "batch_id": batch_id, "imported": len(rows)}
    (out_dir / "db-import-report.json").write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding="utf-8")
    return report


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--text", type=Path, default=DEFAULT_TEXT)
    parser.add_argument("--db", type=Path, default=DEFAULT_DB)
    parser.add_argument("--out-dir", type=Path, default=DEFAULT_OUT)
    parser.add_argument("--apply", action="store_true")
    args = parser.parse_args()

    args.out_dir.mkdir(parents=True, exist_ok=True)
    papers = split_papers(args.text)
    rows, warnings = build_import_rows(papers)

    preview_rows = []
    for row in rows:
        preview = dict(row)
        preview["stem_text"] = row["stem_text"][:500]
        preview["analysis_text"] = row["analysis_text"][:300]
        preview_rows.append(preview)

    summary = {
        "papers_detected": len(papers),
        "questions": len(rows),
        "by_category": dict(Counter(row["question_category"] for row in rows)),
        "by_type": dict(Counter(row["question_type"] for row in rows)),
        "warnings": warnings,
    }
    (args.out_dir / "import-preview.json").write_text(json.dumps(preview_rows, ensure_ascii=False, indent=2), encoding="utf-8")
    (args.out_dir / "import-summary.json").write_text(json.dumps(summary, ensure_ascii=False, indent=2), encoding="utf-8")

    result = {"summary": summary}
    if args.apply:
        result["db"] = insert_rows(args.db, rows, args.out_dir)

    print(json.dumps(result, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
