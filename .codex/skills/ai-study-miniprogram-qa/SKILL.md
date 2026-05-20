---
name: ai-study-miniprogram-qa
description: Project-specific verification workflow for /data/project/ai-study. Use when changing or validating this WeChat Mini Program, especially when comparing against 原型/screen.png, running simulator functional tests, capturing evidence screenshots, generating preview QR codes, or cleaning wx_login artifacts.
---

# AI Study Mini Program QA

## Project Constants

- Project root: `/data/project/ai-study`
- Mini program root: `/data/project/ai-study/miniprogram`
- Prototype: `/data/project/ai-study/原型/screen.png`
- DevTools CLI: `/data/project/ai-study/.tools/bin/wechat-devtools-cli`
- Automator package location: `/tmp/miniprogram-automator-run/node_modules/miniprogram-automator`
- IDE HTTP port: `3799`
- Automator port: `9420`

## Artifact Rules

Keep `wx_login/` clean. Do not place generated images in its root.

- `wx_login/qr/`: preview QR codes and preview metadata.
- `wx_login/screenshots/prototype/`: prototype comparison screenshots, crops, and reports.
- `wx_login/screenshots/preview/`: project-open or preview-generation screenshots.
- `wx_login/screenshots/selftest/`: automated functional-test screenshots.

Naming:

- Use lowercase kebab-case names.
- Include the screen or scope, then state: `home-bottom-css-arrow-app.png`, `home-full-prototype-diff.png`.
- Use `*-app.png` for mini program screenshots from `miniProgram.screenshot`.
- Use `*-simulator.png` for full DevTools/Xvfb screenshots from `scripts/capture-x11.py`.
- Do not clean screenshots every round by default. Keep recent useful evidence from the last few validation passes.
- Clean only when artifacts become noisy, duplicate, or stale; as a rule of thumb, clean when `wx_login/screenshots/` grows past about 20 images or when old failed attempts make the current evidence hard to identify.
- Always remove temporary debug crops/zooms unless the final response explicitly cites them.

## Required Workflow

1. Start from the existing project state. Read the affected `wxml`, `wxss`, `js`, and the prototype before editing.
2. Make the smallest code/style change that matches the prototype and preserves interaction.
3. Run static checks:

```bash
node --check miniprogram/pages/home/home.js
python3 -m json.tool miniprogram/app.json
python3 -m json.tool miniprogram/pages/home/home.json
```

4. Open and automate the actual simulator:

```bash
.tools/bin/wechat-devtools-cli open --project /data/project/ai-study/miniprogram --port 3799 --disable-gpu --trust-project
.tools/bin/wechat-devtools-cli auto --project /data/project/ai-study/miniprogram --auto-port 9420 --port 3799 --trust-project
```

If `miniprogram-automator` is missing, install it outside the repo:

```bash
npm install --prefix /tmp/miniprogram-automator-run miniprogram-automator
```

5. Run real functional tests with `miniprogram-automator`, not only visual inspection.
6. Capture app screenshots and simulator screenshots for every tested area.
7. Compare against `原型/screen.png`; crop or zoom disputed regions and measure alignment when visual judgment is uncertain.
8. Generate the preview QR after validation:

```bash
./scripts/preview-qr.sh
```

9. Check `wx_login/` according to the artifact rules. Clean only if the screenshot set is noisy or too large, then list remaining files.

## Prototype Comparison

Always compare current screenshots to `原型/screen.png` for the affected UI.

- Use `view_image` on the prototype and new screenshots.
- For full-page layout, compare spacing, card width, button height, typography, selected states, and bottom spacing.
- For disputed alignment, write a small PIL script to compute bounding boxes and center lines. Record the numeric result in the final response.
- For image/icon issues, prefer eliminating fragile image assets when CSS can render the shape exactly, especially for simple arrows.

Minimum prototype evidence for home-page changes:

- `wx_login/screenshots/prototype/home-top-app.png`
- `wx_login/screenshots/prototype/home-bottom-app.png`
- Optional crops such as `home-bottom-arrow-crop.png` only while debugging; delete temporary crops before finishing unless they prove a reported issue.

## Functional Test Checklist

For home page, test and screenshot:

- Phone input.
- Verification-code input.
- Get-code button, including disabled/sent state.
- Login tap.
- Guest continue tap.
- Notification bell tap.
- Stage switch: 小学 and 初中.
- All visible grades for the selected stage.
- Subject switch: 语文, 数学, 英语.
- Bottom `立即开始测评` tap after scrolling to the button.

For learn page, test and screenshot:

- `听`
- `下一个`
- `学会了`

For checkin page, test and screenshot:

- `完成今日打卡`
- `重置示例进度`

When a request only touches one page, still run a quick smoke test on the other pages if shared app config, assets, or global styles changed.

## Automator Pattern

Use this structure for self-tests and adapt selectors as needed:

```js
const path = require("path");
const automator = require("/tmp/miniprogram-automator-run/node_modules/miniprogram-automator");
const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

async function must(page, selector) {
  const el = await page.$(selector);
  if (!el) throw new Error(`missing ${selector}`);
  return el;
}

(async () => {
  const miniProgram = await automator.connect({ wsEndpoint: "ws://127.0.0.1:9420" });
  const exceptions = [];
  miniProgram.on("exception", (err) => exceptions.push(err));

  const page = await miniProgram.reLaunch("/pages/home/home");
  await page.waitFor(".prototype-page");

  const inputs = await page.$$("input");
  await inputs[0].input("13800138000");
  await inputs[1].input("123456");
  await (await must(page, ".code-button")).tap();

  const stages = await page.$$(".stage-option");
  await stages[1].tap();
  await sleep(300);

  const grades = await page.$$(".grade-item");
  await grades[2].tap();
  await sleep(300);

  const subjects = await page.$$(".subject-item");
  await subjects[2].tap();
  await sleep(300);

  await miniProgram.pageScrollTo(808);
  await sleep(600);
  await (await must(page, ".bottom-action")).tap();

  await miniProgram.screenshot({
    path: path.resolve("wx_login/screenshots/selftest/home-interaction-app.png")
  });

  const data = await page.data();
  console.log(JSON.stringify({
    codeButtonText: data.codeButtonText,
    selectedStage: data.selectedStage,
    selectedGrade: data.selectedGrade,
    selectedSubject: data.selectedSubject,
    exceptionCount: exceptions.length
  }));

  if (exceptions.length) throw new Error(JSON.stringify(exceptions));
  await miniProgram.disconnect();
})().catch((err) => {
  console.error(err && err.stack || err);
  process.exit(1);
});
```

Capture the full DevTools view when needed:

```bash
python3 scripts/capture-x11.py wx_login/screenshots/selftest/home-interaction-simulator.png
```

## Artifact Review Checklist

Before final response:

- Keep `wx_login/README.md`.
- Keep current QR outputs in `wx_login/qr/`.
- Allow the latest few useful test rounds to remain; do not delete prior useful evidence just because a new test ran.
- Delete temporary debug zooms/crops unless they are part of the evidence.
- Clean stale retry/blocking/failed-preview screenshots when they make the directory noisy.
- If there are more than about 20 screenshots under `wx_login/screenshots/`, reduce them to the latest 2-3 useful rounds per category (`prototype`, `preview`, `selftest`) plus any screenshots cited in the final response.
- Run:

```bash
find wx_login -maxdepth 3 -type f | sort
```

Report the remaining evidence files and the key test results.
