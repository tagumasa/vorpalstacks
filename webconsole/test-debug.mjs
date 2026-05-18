import { chromium } from "playwright";
const BASE = "http://localhost:50090";
const res = await fetch(`${BASE}/admin_auth.AdminAuthService/LoginRoot`, {
  method: "POST", headers: { "Content-Type": "application/json", "Connect-Protocol-Version": "1" },
  body: JSON.stringify({ password: "vorpal00" }),
});
const { accessToken: TOKEN } = await res.json();
const browser = await chromium.launch({ headless: true, args: ["--no-sandbox"] });
const ctx = await browser.newContext({ viewport: { width: 2560, height: 1440 } });
await ctx.addInitScript((t) => {
  localStorage.setItem("vs_access_token", t);
  localStorage.setItem("vs_refresh_token", "x");
  localStorage.setItem("vs_id_token", "x");
}, TOKEN);
const page = await ctx.newPage();

const consoleErrors = [];
page.on("console", msg => { if (msg.type() === "error") consoleErrors.push(msg.text()); });

await page.goto(`${BASE}/webconsole/services/dynamodb`);
await page.waitForTimeout(5000);
await page.screenshot({ path: "/tmp/debug-page.png", fullPage: false });
const body = await page.locator("body").textContent();
console.log("Body:", body?.substring(0, 300));
console.log("Console errors:", consoleErrors);
await browser.close();
