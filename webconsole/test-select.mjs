import { chromium } from "playwright";

const BASE = "http://localhost:50090";

const res = await fetch(`${BASE}/admin_auth.AdminAuthService/LoginRoot`, {
  method: "POST", headers: { "Content-Type": "application/json", "Connect-Protocol-Version": "1" },
  body: JSON.stringify({ password: "vorpal00" }),
});
const { accessToken: TOKEN } = await res.json();

await fetch(`${BASE}/dynamodb.DynamoDBService/CreateTable`, {
  method: "POST", headers: { "Content-Type": "application/json", "Connect-Protocol-Version": "1", Authorization: `Bearer ${TOKEN}` },
  body: JSON.stringify({ tablename: "ux-test", keyschema: [{ attributename: "pk", keytype: "HASH" }, { attributename: "sk", keytype: "RANGE" }], attributedefinitions: [{ attributename: "pk", attributetype: 1 }, { attributename: "sk", attributetype: 1 }], billingmode: 1 }),
});
await fetch(`${BASE}/dynamodb.DynamoDBService/PutItem`, {
  method: "POST", headers: { "Content-Type": "application/json", "Connect-Protocol-Version": "1", Authorization: `Bearer ${TOKEN}` },
  body: JSON.stringify({ tablename: "ux-test", item: { pk: { s: "user1" }, sk: { s: "profile" }, email: { s: "a@b.com" }, count: { n: "5" }, active: { bool: true } } }),
});

const browser = await chromium.launch({ headless: true, args: ["--no-sandbox"] });
const ctx = await browser.newContext({ viewport: { width: 2560, height: 1440 } });
await ctx.addInitScript((t) => {
  localStorage.setItem("vs_access_token", t);
  localStorage.setItem("vs_refresh_token", "x");
  localStorage.setItem("vs_id_token", "x");
}, TOKEN);
const page = await ctx.newPage();

await page.goto(`${BASE}/webconsole/services/dynamodb`);
await page.waitForSelector("table", { timeout: 10000 });
await page.waitForTimeout(500);

// Click table row
const tableRow = page.locator("tbody tr").first();
await tableRow.click();
await page.waitForTimeout(2000);

// Click item
const itemRow = page.locator("tbody tr").first();
await itemRow.click();
await page.waitForTimeout(1500);

await page.screenshot({ path: "/tmp/ux-1-detail.png", fullPage: false });

// 1. Verify PK/SK are at top and disabled
const selects = await page.locator("select").all();
console.log("=== TEST 1: PK/SK at top, disabled ===");
console.log(`Total selects: ${selects.length}`);
const pkSelect = selects[0];
const pkVal = await pkSelect.inputValue();
const pkDisabled = await pkSelect.isDisabled();
console.log(`  [0] (PK) value=${pkVal} disabled=${pkDisabled} — EXPECT: disabled=true`);
if (selects.length > 1) {
  const skVal = await selects[1].inputValue();
  const skDisabled = await selects[1].isDisabled();
  console.log(`  [1] (SK) value=${skVal} disabled=${skDisabled} — EXPECT: disabled=true`);
}

// 2. Add new attribute
const addBtn = page.getByRole("button", { name: /\+.*Add attribute/i }).first();
await addBtn.click();
await page.waitForTimeout(500);

const allSelects = await page.locator("select").all();
console.log(`\n=== TEST 2: Add attribute, type selectable ===`);
const newSelect = allSelects[allSelects.length - 1];
const newDisabled = await newSelect.isDisabled();
const newVis = await newSelect.isVisible();
console.log(`  New select: disabled=${newDisabled} visible=${newVis} — EXPECT: disabled=false`);

// Change type
await newSelect.selectOption("N");
const afterChange = await newSelect.inputValue();
console.log(`  After selectOption('N'): value=${afterChange} — EXPECT: N`);

await page.screenshot({ path: "/tmp/ux-2-added.png", fullPage: false });

// 3. Check delete button on new row
const delBtns = page.locator("button").filter({ hasText: "✕" });
const delCount = await delBtns.count();
console.log(`\n=== TEST 3: Delete button ===`);
console.log(`  Delete buttons count: ${delCount}`);

// 4. Verify sort order: key section has border
const keySection = page.locator("div").filter({ has: page.locator("select").first() }).first();
console.log(`\n=== TEST 4: Key section exists ===`);
const border = await keySection.evaluate(el => getComputedStyle(el).borderBottom);
console.log(`  Key section border: ${border}`);

await page.screenshot({ path: "/tmp/ux-3-final.png", fullPage: false });

// Cleanup
await fetch(`${BASE}/dynamodb.DynamoDBService/DeleteTable`, {
  method: "POST", headers: { "Content-Type": "application/json", "Connect-Protocol-Version": "1", Authorization: `Bearer ${TOKEN}` },
  body: JSON.stringify({ tablename: "ux-test" }),
});

await browser.close();
console.log("\nDone");
