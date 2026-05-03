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
    proxy: {
      "/admin_auth": "http://127.0.0.1:9090",
      "/admin_config": "http://127.0.0.1:9090",
      "/acm": "http://127.0.0.1:9090",
      "/apigateway": "http://127.0.0.1:9090",
      "/appsync": "http://127.0.0.1:9090",
      "/athena": "http://127.0.0.1:9090",
      "/cloudfront": "http://127.0.0.1:9090",
      "/cloudtrail": "http://127.0.0.1:9090",
      "/cloudwatch": "http://127.0.0.1:9090",
      "/cognitoidentity": "http://127.0.0.1:9090",
      "/dynamodb": "http://127.0.0.1:9090",
      "/iam": "http://127.0.0.1:9090",
      "/kinesis": "http://127.0.0.1:9090",
      "/kms": "http://127.0.0.1:9090",
      "/lambda": "http://127.0.0.1:9090",
      "/neptune": "http://127.0.0.1:9090",
      "/route53": "http://127.0.0.1:9090",
      "/s3": "http://127.0.0.1:9090",
      "/scheduler": "http://127.0.0.1:9090",
      "/secretsmanager": "http://127.0.0.1:9090",
      "/sesv2": "http://127.0.0.1:9090",
      "/sfn": "http://127.0.0.1:9090",
      "/sns": "http://127.0.0.1:9090",
      "/sqs": "http://127.0.0.1:9090",
      "/ssm": "http://127.0.0.1:9090",
      "/sts": "http://127.0.0.1:9090",
      "/timestreamquery": "http://127.0.0.1:9090",
      "/timestreamwrite": "http://127.0.0.1:9090",
      "/waf": "http://127.0.0.1:9090",
      "/wafv2": "http://127.0.0.1:9090",
    },
  },
  build: {
    outDir: "dist",
  },
});
