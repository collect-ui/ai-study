#!/usr/bin/env python3
import argparse
import datetime as dt
import http.cookiejar
import json
import os
from pathlib import Path
import re
import sys
import time
import urllib.error
import urllib.request


DEFAULT_URL = "http://127.0.0.1:8026"
DEFAULT_MANIFEST = Path(__file__).resolve().parents[1] / "testdata" / "question_ai_fragments" / "manifest.json"


def load_json(path):
    with open(path, "r", encoding="utf-8") as file:
        return json.load(file)


def post_json(opener, base_url, service, payload, timeout):
    url = base_url.rstrip("/") + "/template_data/data?service=" + service
    request = urllib.request.Request(
        url,
        data=json.dumps(payload, ensure_ascii=False).encode("utf-8"),
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    with opener.open(request, timeout=timeout) as response:
        return json.loads(response.read().decode("utf-8"))


def login(opener, args):
    response = post_json(
        opener,
        args.base_url,
        "system.login",
        {"username": args.username, "password": args.password},
        args.timeout,
    )
    if not response.get("success"):
        raise RuntimeError("login failed: " + json.dumps(response, ensure_ascii=False))


def row_summary(row):
    choice_items = row.get("choice_items")
    if isinstance(choice_items, str):
        try:
            choice_items = json.loads(choice_items)
        except json.JSONDecodeError:
            choice_items = []
    if not isinstance(choice_items, list):
        choice_items = []
    return {
        "title": row.get("title", ""),
        "question_type": row.get("question_type", ""),
        "question_category": row.get("question_category", ""),
        "stem_text": row.get("stem_text", ""),
        "answer_key": row.get("answer_key", ""),
        "analysis_text": row.get("analysis_text", ""),
        "choice_item_count": len(choice_items),
        "choice_item_answers": [
            item.get("answer_key", "") for item in choice_items if isinstance(item, dict)
        ],
        "choice_item_analyses": [
            first_text(item, "analysis", "analysis_text", "explanation")
            for item in choice_items
            if isinstance(item, dict)
        ],
    }


def normalize_for_diff(rows):
    return [
        {
            "question_type": row.get("question_type", ""),
            "question_category": row.get("question_category", ""),
            "answer_key": normalize_answer_for_diff(row.get("answer_key", "")),
            "analysis_text": normalize_analysis_for_diff(row.get("analysis_text", "")),
            "choice_item_count": row.get("choice_item_count", 0),
            "choice_item_answers": [
                normalize_answer_for_diff(item) for item in row.get("choice_item_answers", [])
            ],
            "choice_item_analyses": [
                normalize_analysis_for_diff(item) for item in row.get("choice_item_analyses", [])
            ],
        }
        for row in rows
    ]


def first_text(mapping, *keys):
    for key in keys:
        value = mapping.get(key)
        if value is not None and str(value).strip():
            return str(value)
    return ""


def normalize_answer_for_diff(value):
    text = str(value or "").strip()
    text = text.rstrip(".。")
    text = text.replace("’", "'").replace("‘", "'")
    text = text.replace("；", ";").replace("，", ",")
    text = re.sub(r"\s+", " ", text)
    if ";" in text or "," in text:
        parts = [part.strip() for part in re.split(r"[;,]", text) if part.strip()]
        return ";".join(parts)
    return text


def normalize_analysis_for_diff(value):
    text = str(value or "").strip()
    text = text.replace("’", "'").replace("‘", "'")
    text = text.replace("“", '"').replace("”", '"')
    text = text.replace("（", "(").replace("）", ")")
    text = re.sub(r"^[\\[【]?(解析|分析)[\\]】]?[:：]?", "", text)
    text = text.lower()
    text = re.sub(r"[^0-9a-z\u4e00-\u9fff]+", "", text)
    text = text.replace("solas", "soas")
    text = text.replace("便宜得多", "")
    text = re.sub(r"d应该改成[a-z]+", "", text)
    return text


def case_payload(case, raw_text, provider, args):
    payload = {
        "provider": provider,
        "parse_mode": "chunked",
        "enable_cache": False,
        "chunk_chars": args.chunk_chars,
        "max_chars": 24000,
        "max_chunks": 40,
        "source_max_chars": 120000,
        "subject": "english",
        "stage": "primary",
        "grade": "grade_6",
        "textbook_version": "pep",
        "question_type": "single_choice",
        "question_category": "normal",
        "difficulty": "basic",
        "score": 5,
        "raw_text": raw_text,
    }
    payload.update(case.get("params") or {})
    if provider == "deepseek":
        api_key = os.environ.get(args.deepseek_key_env, "").strip()
        if not api_key:
            raise RuntimeError(f"{args.deepseek_key_env} is not set")
        payload["api_key"] = api_key
    return payload


def run_provider_case(opener, args, case, raw_text, provider):
    started = time.time()
    try:
        response = post_json(
            opener,
            args.base_url,
            "question.ai_parse",
            case_payload(case, raw_text, provider, args),
            args.timeout,
        )
    except (urllib.error.URLError, TimeoutError, RuntimeError) as error:
        return {
            "success": False,
            "error": str(error),
            "seconds": round(time.time() - started, 2),
        }

    body = response.get("data") if isinstance(response, dict) else None
    rows = (body or {}).get("rows") or []
    summaries = [row_summary(row) for row in rows if isinstance(row, dict)]
    expected_count = case.get("expected_count")
    expected_choice_items = case.get("expected_choice_items")
    count_ok = expected_count is None or len(rows) == expected_count
    choice_item_ok = True
    if expected_choice_items is not None:
        choice_item_ok = bool(summaries) and summaries[0]["choice_item_count"] == expected_choice_items
    return {
        "success": bool(response.get("success")),
        "msg": response.get("msg", ""),
        "seconds": round(time.time() - started, 2),
        "model": (body or {}).get("model", ""),
        "source": (body or {}).get("source", ""),
        "row_count": (body or {}).get("row_count", len(rows)),
        "expected_count": expected_count,
        "count_ok": count_ok,
        "expected_choice_items": expected_choice_items,
        "choice_item_ok": choice_item_ok,
        "chunk_summaries": (body or {}).get("chunk_summaries", []),
        "rows": summaries,
    }


def compare_results(provider_results):
    successful = {
        provider: result
        for provider, result in provider_results.items()
        if result.get("success") and "rows" in result
    }
    if len(successful) < 2:
        return {"comparable": False}
    providers = sorted(successful)
    base_provider = providers[0]
    base_rows = normalize_for_diff(successful[base_provider]["rows"])
    diffs = {"comparable": True, "base_provider": base_provider, "providers": providers}
    for provider in providers:
        rows = normalize_for_diff(successful[provider]["rows"])
        diffs[provider] = {
            "same_as_base": rows == base_rows,
            "row_count_delta": len(rows) - len(base_rows),
        }
    return diffs


def run(args):
    manifest_path = Path(args.manifest).resolve()
    manifest = load_json(manifest_path)
    fixture_dir = manifest_path.parent
    providers = [item.strip() for item in args.providers.split(",") if item.strip()]
    opener = urllib.request.build_opener(
        urllib.request.ProxyHandler({}),
        urllib.request.HTTPCookieProcessor(http.cookiejar.CookieJar()),
    )
    login(opener, args)

    report = {
        "generated_at": dt.datetime.now(dt.timezone.utc).isoformat(),
        "base_url": args.base_url,
        "manifest": str(manifest_path),
        "providers": providers,
        "repeat": args.repeat,
        "cases": [],
    }

    case_ids = {item.strip() for item in args.case_id.split(",") if item.strip()}
    for case in manifest.get("cases", []):
        if case_ids and case.get("id") not in case_ids:
            continue
        raw_text = (fixture_dir / case["file"]).read_text(encoding="utf-8")
        case_runs = []
        for iteration in range(1, args.repeat + 1):
            provider_results = {}
            for provider in providers:
                provider_results[provider] = run_provider_case(opener, args, case, raw_text, provider)
            case_runs.append(
                {
                    "iteration": iteration,
                    "providers": provider_results,
                    "diff": compare_results(provider_results),
                }
            )
        report["cases"].append(
            {
                "id": case["id"],
                "name": case["name"],
                "file": case["file"],
                "expected_count": case.get("expected_count"),
                "expected_choice_items": case.get("expected_choice_items"),
                "runs": case_runs,
            }
        )
    return report


def print_summary(report):
    for case in report["cases"]:
        print(f"{case['id']} {case['name']}")
        for run in case["runs"]:
            print(f"  iteration {run['iteration']}")
            for provider, result in run["providers"].items():
                if not result.get("success"):
                    print(f"    {provider}: ERROR {result.get('error') or result.get('msg')}")
                    continue
                status = "OK" if (
                    result.get("count_ok")
                    and result.get("choice_item_ok")
                ) else "CHECK"
                print(
                    f"    {provider}: {status} rows={result.get('row_count')} "
                    f"model={result.get('model')} seconds={result.get('seconds')}"
                )
            diff = run.get("diff") or {}
            if diff.get("comparable"):
                same = all(
                    item.get("same_as_base")
                    for key, item in diff.items()
                    if isinstance(item, dict) and key != "providers"
                )
                print(f"    diff: {'same' if same else 'different'}")
            else:
                print("    diff: not comparable")


def main():
    parser = argparse.ArgumentParser(description="Compare question.ai_parse output across providers.")
    parser.add_argument("--base-url", default=DEFAULT_URL)
    parser.add_argument("--manifest", default=str(DEFAULT_MANIFEST))
    parser.add_argument("--providers", default="codex,deepseek")
    parser.add_argument("--deepseek-key-env", default="DEEPSEEK_API_KEY")
    parser.add_argument("--username", default=os.environ.get("AI_STUDY_ADMIN_USER", "admin"))
    parser.add_argument("--password", default=os.environ.get("AI_STUDY_ADMIN_PASSWORD", "123456"))
    parser.add_argument("--timeout", type=int, default=360)
    parser.add_argument("--chunk-chars", type=int, default=6000)
    parser.add_argument("--repeat", type=int, default=1)
    parser.add_argument("--case-id", default="", help="Comma-separated case ids to run from the manifest.")
    parser.add_argument("--out", default="")
    parser.add_argument("--summary", action="store_true")
    args = parser.parse_args()

    report = run(args)
    if args.out:
        out_path = Path(args.out)
        out_path.parent.mkdir(parents=True, exist_ok=True)
        out_path.write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding="utf-8")
    if args.summary:
        print_summary(report)
    else:
        print(json.dumps(report, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    try:
        main()
    except Exception as error:
        print(str(error), file=sys.stderr)
        raise SystemExit(1)
