import fs from "node:fs/promises";
import path from "node:path";
import { execFileSync } from "node:child_process";
import playwright from "/data/project/sport-ui/node_modules/playwright/index.js";

const { chromium } = playwright;

const baseURL = "http://127.0.0.1:8026";
const outDir = "/data/project/ai-study/feature/ai-question-import-verify";
const dbPath = "/data/project/ai-study/backend/database/ai_study_admin.db";
const reportPath = path.join(outDir, "report.json");
const screenshotOpen = path.join(outDir, "dialog-open.png");
const screenshotParsed = path.join(outDir, "parsed-table.png");
const screenshotImported = path.join(outDir, "imported-result.png");

const fixedJSON = JSON.stringify({
  questions: [
    {
      title: "AI导入验收题",
      subject: "english",
      stage: "junior",
      grade: "grade_7",
      textbook_version: "pep",
      question_type: "single_choice",
      question_category: "normal",
      difficulty: "basic",
      score: 5,
      stem_text: "AI import acceptance test question: Choose the correct word.",
      option_a_text: "apple",
      option_b_text: "book",
      option_c_text: "chair",
      option_d_text: "desk",
      answer_key: "A",
      analysis_text: "apple is the expected answer for this import verification row."
    },
    {
      title: "AI导入验收阅读理解",
      subject: "english",
      stage: "junior",
      grade: "grade_7",
      textbook_version: "pep",
      question_type: "single_choice",
      question_category: "reading_choice",
      difficulty: "basic",
      score: 5,
      stem_text: "AI import reading passage: Tom is at home and he is reading a book.",
      choice_items: [
        {
          sub_no: "1",
          question_text: "Where is Tom?",
          option_a: "At school",
          option_b: "At home",
          option_c: "In a shop",
          option_d: "In a park",
          answer_key: "B",
          analysis: "The passage says Tom is at home."
        }
      ]
    }
  ]
});

const report = {
  started_at: new Date().toISOString(),
  base_url: baseURL,
  console_errors: [],
  page_errors: [],
  request_failed: [],
  http_errors: [],
  steps: [],
  screenshots: {
    dialog_open: screenshotOpen,
    parsed_table: screenshotParsed,
    imported_result: screenshotImported
  }
};

function pushStep(name, data = {}) {
  report.steps.push({ name, at: new Date().toISOString(), ...data });
}

function sqliteValue(sql) {
  return execFileSync("sqlite3", [dbPath, sql], { encoding: "utf8" }).trim();
}

async function postService(page, service, data) {
  const response = await page.request.post(`${baseURL}/template_data/data?service=${encodeURIComponent(service)}`, {
    data
  });
  const body = await response.json();
  if (!body.success) {
    throw new Error(`${service}: ${body.msg || "failed"}`);
  }
  return body;
}

async function loginIfNeeded(page) {
  await page.goto(`${baseURL}/collect-ui/#/collect-ui/framework/question-bank`, {
    waitUntil: "domcontentloaded"
  });
  await page.waitForLoadState("networkidle", { timeout: 20000 }).catch(() => {});
  const loginButton = page.getByRole("button", { name: /登录/ });
  if (await loginButton.isVisible().catch(() => false)) {
    const inputs = page.locator("input");
    await inputs.nth(0).fill("admin");
    await inputs.nth(1).fill("123456");
    await loginButton.click();
    await page.waitForLoadState("networkidle", { timeout: 20000 }).catch(() => {});
    await page.goto(`${baseURL}/collect-ui/#/collect-ui/framework/question-bank`, {
      waitUntil: "domcontentloaded"
    });
    await page.waitForLoadState("networkidle", { timeout: 20000 }).catch(() => {});
  }
}

async function run() {
  await fs.mkdir(outDir, { recursive: true });
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage({ viewport: { width: 1440, height: 940 } });

  page.on("console", (msg) => {
    if (msg.type() === "error") {
      report.console_errors.push({ text: msg.text(), location: msg.location() });
    }
  });
  page.on("pageerror", (error) => {
    report.page_errors.push({ message: error.message, stack: error.stack });
  });
  page.on("requestfailed", (request) => {
    report.request_failed.push({
      url: request.url(),
      method: request.method(),
      failure: request.failure()?.errorText || ""
    });
  });
  page.on("response", (response) => {
    if (response.status() >= 400) {
      report.http_errors.push({ url: response.url(), status: response.status() });
    }
  });

  try {
    await loginIfNeeded(page);
    if (!(await page.getByText("全学科题库管理").isVisible().catch(() => false))) {
      await page.getByText("题库管理").first().click();
    }
    await page.getByText("全学科题库管理").waitFor({ timeout: 20000 });
    pushStep("opened_question_bank");

    await page.getByRole("button", { name: /AI解析导入/ }).click();
    await page.getByText("AI题目解析导入").waitFor({ timeout: 10000 });
    await page.locator(".qb-ai-import-lowcode").waitFor({ state: "visible", timeout: 10000 });
    await page.getByText("分段解析").waitFor({ timeout: 10000 });
    await page.getByRole("button", { name: /读取示例PDF/ }).waitFor({ timeout: 10000 });
    await page.waitForTimeout(300);
    await page.screenshot({ path: screenshotOpen, fullPage: true });
    pushStep("opened_ai_import_dialog");

    await page.getByRole("button", { name: /读取示例PDF/ }).click();
    const firstTextarea = page.locator(".qb-ai-import-lowcode textarea").nth(0);
    await page.waitForFunction(
      (selector) => {
        const el = document.querySelectorAll(selector)[0];
        return el && el.value && el.value.length > 80;
      },
      ".qb-ai-import-lowcode textarea",
      { timeout: 30000 }
    );
    const pdfTextLength = (await firstTextarea.inputValue()).length;
    if (pdfTextLength <= 24000) {
      throw new Error(`sample PDF text is still truncated: ${pdfTextLength}`);
    }
    pushStep("read_sample_pdf", { pdf_text_length: pdfTextLength });

    const mockParse = await postService(page, "question.ai_parse", {
      provider: "codex",
      mock_response: fixedJSON,
      raw_text: "mock",
      subject: "english",
      stage: "junior",
      grade: "grade_7",
      textbook_version: "pep"
    });
    pushStep("mock_ai_parse_service", {
      row_count: mockParse.data?.row_count,
      fixed_text_length: String(mockParse.data?.fixed_text || "").length
    });
    if (mockParse.data?.row_count !== 2) {
      throw new Error(`mock parse should return 2 rows, got ${mockParse.data?.row_count}`);
    }
    const mockRows = mockParse.data?.rows || [];
    const mockReading = mockRows.find((row) => row.question_category === "reading_choice");
    if (!mockReading || !Array.isArray(mockReading.choice_items) || mockReading.choice_items[0]?.answer_key !== "B") {
      throw new Error("mock reading_choice row did not preserve choice_items");
    }

    await page.locator(".qb-ai-import-lowcode textarea").nth(1).fill(fixedJSON);
    await page.getByRole("button", { name: /解析字符串/ }).click();
    await page.locator(".qb-ai-import-lowcode").getByText("AI import acceptance test question").waitFor({ timeout: 10000 });
    await page.locator(".qb-ai-import-lowcode .ag-cell").getByText("AI import reading passage").waitFor({ timeout: 10000 });
    await page.screenshot({ path: screenshotParsed, fullPage: true });
    pushStep("parsed_fixed_string");

    await page.getByRole("button", { name: /确认导入/ }).click();
    await page.waitForFunction(
      () => Array.from(document.querySelectorAll(".qb-ai-import-lowcode .ag-cell")).filter((el) => el.textContent?.trim() === "已导入").length >= 2,
      null,
      { timeout: 30000 }
    );
    await page.screenshot({ path: screenshotImported, fullPage: true });
    pushStep("confirmed_import");

    const questionID = sqliteValue(
      "select question_id from question_item where ifnull(is_delete, '0') = '0' and stem_text like '%AI import acceptance test question%' order by create_time desc limit 1;"
    );
    if (!questionID) {
      throw new Error("imported question was not found in question_item");
    }
    const readingID = sqliteValue(
      "select question_id from question_item where ifnull(is_delete, '0') = '0' and stem_text like '%AI import reading passage%' order by create_time desc limit 1;"
    );
    if (!readingID) {
      throw new Error("imported reading question was not found in question_item");
    }
    const readingDetail = await postService(page, "question.question_choice_detail", {
      question_id: readingID
    });
    const choiceItems = JSON.parse(readingDetail.data?.choice_items || "[]");
    if (
      readingDetail.data?.question_category !== "reading_choice" ||
      readingDetail.data?.question_type !== "single_choice" ||
      choiceItems.length !== 1 ||
      choiceItems[0]?.question_text !== "Where is Tom?" ||
      choiceItems[0]?.option_b !== "At home" ||
      choiceItems[0]?.answer_key !== "B"
    ) {
      throw new Error(`reading_choice fields mismatch: ${JSON.stringify(readingDetail.data)}`);
    }
    pushStep("verified_imported_question", {
      question_id: questionID,
      reading_id: readingID,
      reading_choice_items: choiceItems.length
    });

    await postService(page, "question.question_choice_delete", {
      question_id: questionID
    });
    await postService(page, "question.question_choice_delete", {
      question_id: readingID
    });
    pushStep("cleaned_imported_question", { question_id: questionID, reading_id: readingID });

    report.ok = true;
  } catch (error) {
    report.ok = false;
    report.error = error?.stack || error?.message || String(error);
    await page.screenshot({ path: path.join(outDir, "failure.png"), fullPage: true }).catch(() => {});
    throw error;
  } finally {
    report.finished_at = new Date().toISOString();
    await fs.writeFile(reportPath, JSON.stringify(report, null, 2), "utf8");
    await browser.close();
  }
}

run().catch((error) => {
  console.error(error?.stack || error);
  process.exit(1);
});
