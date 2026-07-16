import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import path from "path";
import { execSync } from "child_process";

function gitVersion(): string {
  try {
    return execSync("git describe --tags --always --dirty", { encoding: "utf-8" }).trim();
  } catch {
    try {
      return execSync("git rev-parse --short HEAD", { encoding: "utf-8" }).trim();
    } catch {
      return "dev";
    }
  }
}

export default defineConfig({
  plugins: [react(), tailwindcss()],
  define: {
    __APP_VERSION__: JSON.stringify(gitVersion()),
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    proxy: Object.fromEntries(
      [
        "admin_auth", "admin_config",
        "acm", "apigateway", "appsync", "athena",
        "cloudfront", "cloudtrail", "cloudwatch",
        "cloudwatchevents", "cloudwatchlogs",
        "cognitoidentity", "cognitoidentityprovider",
        "common", "dynamodb", "iam", "kinesis", "kms",
        "lambda", "neptune", "neptunedata", "neptunegraph",
        "rds", "route53", "s3", "scheduler",
        "secretsmanager", "sesv2", "sfn",
        "sns", "sqs", "ssm", "sts",
        "timestreamquery", "timestreamwrite",
        "waf", "wafv2",
      ].map((prefix) => [`/${prefix}`, "http://127.0.0.1:50090"]),
    ),
  },
  base: "/webconsole/",
  build: {
    outDir: "dist",
  },
});
