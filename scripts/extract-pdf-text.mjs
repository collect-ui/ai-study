import fs from "fs";
import path from "path";
import * as pdfjsLib from "/data/project/sport-ui/node_modules/pdfjs-dist/legacy/build/pdf.mjs";

const pdfPath = process.argv[2] || "题目/【小升初】英语复习题三十套全国通用（含详细解析）.pdf";
const outDir = process.argv[3] || "test-results/full-pdf-import";

function normalizeLine(line) {
  return line.replace(/\s+/g, " ").trim();
}

async function main() {
  fs.mkdirSync(outDir, { recursive: true });

  const data = new Uint8Array(fs.readFileSync(pdfPath));
  const pdf = await pdfjsLib.getDocument({ data, disableWorker: true }).promise;
  const pages = [];

  for (let pageNo = 1; pageNo <= pdf.numPages; pageNo += 1) {
    const page = await pdf.getPage(pageNo);
    const content = await page.getTextContent();
    const lines = [];
    let lastY = null;
    let line = [];

    for (const item of content.items) {
      const y = Math.round(item.transform[5]);
      if (lastY !== null && Math.abs(y - lastY) > 4) {
        const text = normalizeLine(line.join(" "));
        if (text) lines.push(text);
        line = [];
      }
      line.push(item.str);
      lastY = y;
    }

    const text = normalizeLine(line.join(" "));
    if (text) lines.push(text);
    pages.push({ page: pageNo, text: lines.join("\n") });
  }

  const pagesPath = path.join(outDir, "pdf-pages.json");
  const textPath = path.join(outDir, "pdf-full.txt");
  fs.writeFileSync(pagesPath, JSON.stringify(pages, null, 2));
  fs.writeFileSync(
    textPath,
    pages.map((page) => `--- PAGE ${page.page} ---\n${page.text}`).join("\n\n"),
  );

  console.log(JSON.stringify({ pages: pdf.numPages, pagesPath, textPath }, null, 2));
}

main().catch((error) => {
  console.error(error && error.stack ? error.stack : error);
  process.exit(1);
});
