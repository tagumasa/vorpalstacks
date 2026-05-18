import { chromium } from "playwright";
const BASE = "http://localhost:50090";
const res = await fetch(`${BASE}/admin_auth.AdminAuthService/LoginRoot`, {
  method: "POST", headers: { "Content-Type": "application/json", "Connect-Protocol-Version": "1" },
  body: JSON.stringify({ password: "vorpal00" }),
});
const { accessToken: TOKEN } = await res.json();

// Create table with SK
await fetch(`${BASE}/dynamodb.DynamoDBService/CreateTable`, {
  method: "POST", headers: { "Content-Type": "application/json", "Connect-Protocol-Version": "1", Authorization: `Bearer ${TOKEN}` },
  body: JSON.stringify({ tablename: "save-test", keyschema: [{ attributename: "pk", keytype: "HASH" }], attributedefinitions: [{ attributename: "pk", attributetype: 1 }], billingmode: 1 }),
});

const browser = await chromium.launch({ headless: true, args: ["--no-sandbox"] });
const ctx = await browser.newContext({ viewport: { width: 2560, height: 1440 } });
await ctx.addInitScript((t) => {
  localStorage.setItem("vs_access_token", t);
  localStorage.setItem("vs_refresh_token", "x");
  localStorage.setItem("vs_id_token", "x");
}, TOKEN);
const page = await ctx.newPage();
const errors = [];
page.on("console", msg => { if (msg.type() === "error") errors.push(msg.text()); });

await page.goto(`${BASE}/webconsole/services/dynamodb`);
await page.waitForSelector("table", { timeout: 10000 });
await page.waitForTimeout(500);

// Click table
await page.locator("tbody tr").first().click();
await page.waitForTimeout(2000);

// Click "Add Item" button
const addBtn = page.getByRole("button", { name: /Add item/i }).first();
await addBtn.click();
await page.waitForTimeout(1000);

// Check: is modal visible?
const modal = page.locator(".modal-overlay, .modal, [role='dialog']").first();
const modalVisible = await modal.isVisible().catch(() => false);
console.log("1. Modal opened:", modalVisible);
await page.screenshot({ path: "/tmp/save-1-modal.png" });

// Try switching tabs in modal
const jsonTab = page.getByRole("button", { name: /JSON/i }).first();
const structTab = page.getByRole("button", { name: /structured/i }).first();
console.log("  JSON tab visible:", await jsonTab.isVisible().catch(() => false));
console.log("  Structured tab visible:", await structTab.isVisible().catch(() => false));

// Fill PK in structured tab
const nameInputs = page.locator("input[placeholder='Attribute name']");
const count = await nameInputs.count();
console.log("  Name inputs:", count);
if (count > 0) {
  await nameInputs.first().fill("pk");
  const valueInputs = page.locator("input[placeholder='hello']");
  const vc = await valueInputs.count();
  if (vc > 0) await valueInputs.first().fill("test1");
}

// Add a Map attribute
const addAttrBtn = page.getByRole("button", { name: /Add attribute/i }).first();
await addAttrBtn.click();
await page.waitForTimeout(300);
const allNameInputs = page.locator("input[placeholder='Attribute name']");
const ac = await allNameInputs.count();
console.log("  After add attr, name inputs:", ac);
if (ac > 1) {
  await allNameInputs.nth(1).fill("data");
  // Change type to Map
  const selects = page.locator("select");
  const sc = await selects.count();
  console.log("  Selects:", sc);
  if (sc > 1) {
    await selects.nth(1).selectOption("M");
    await page.waitForTimeout(300);
    // Fill map entry
    const keyInputs = page.locator("input[placeholder='key']");
    const kc = await keyInputs.count();
    console.log("  Map key inputs:", kc);
    if (kc > 0) {
      await keyInputs.first().fill("name");
      const valInputs = page.locator("input[placeholder='hello']");
      // Last one should be the map entry value
      const lastVal = valInputs.nth(await valInputs.count() - 1);
      await lastVal.fill("hello");
    }
  }
}

await page.screenshot({ path: "/tmp/save-2-filled.png" });

// Try to save
const saveBtn = page.getByRole("button", { name: /create/i }).first();
const saveDisabled = await saveBtn.isDisabled();
console.log("\n2. Save button disabled:", saveDisabled);
if (!saveDisabled) {
  await saveBtn.click();
  await page.waitForTimeout(2000);
  console.log("  Console errors:", errors);
  
  // Check if item appears in list
  await page.screenshot({ path: "/tmp/save-3-after.png" });
  const rows = page.locator("tbody tr");
  console.log("  Table rows after save:", await rows.count());
}

// Cleanup
await fetch(`${BASE}/dynamodb.DynamoDBService/DeleteTable`, {
  method: "POST", headers: { "Content-Type": "application/json", "Connect-Protocol-Version": "1", Authorization: `Bearer ${TOKEN}` },
  body: JSON.stringify({ tablename: "save-test" }),
});

await browser.close();
