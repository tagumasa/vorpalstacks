/**
 * Service catalog mapping 29 backend dispatcher services to 27 display entries.
 * Related services are grouped into single entries:
 * - Neptune (Neptune + NeptuneData + NeptuneGraph)
 * - CloudWatch (CloudWatch + CloudWatchLogs)
 * - Cognito (Cognito + CognitoIdentity)
 * - Timestream (TimestreamQuery + TimestreamWrite)
 */

/** Metadata for a single service display entry in the sidebar. */
export interface ServiceEntry {
  id: string;
  displayName: string;
  icon: string;
  category: string;
  backendServices: string[];
  hasPage: boolean;
}

/** Full service catalog with display metadata. */
export const SERVICE_CATALOG: ServiceEntry[] = [
  { id: "acm", displayName: "ACM", icon: "🔐", category: "Security", backendServices: ["ACM"], hasPage: true },
  { id: "apigateway", displayName: "API Gateway", icon: "🌉", category: "Networking", backendServices: ["APIGateway"], hasPage: true },
  { id: "appsync", displayName: "AppSync", icon: "⚡", category: "Compute", backendServices: ["AppSync"], hasPage: true },
  { id: "athena", displayName: "Athena", icon: "🔬", category: "Analytics", backendServices: ["Athena"], hasPage: true },
  { id: "cloudfront", displayName: "CloudFront", icon: "☁️", category: "Networking", backendServices: ["CloudFront"], hasPage: true },
  { id: "cloudtrail", displayName: "CloudTrail", icon: "📋", category: "Management", backendServices: ["CloudTrail"], hasPage: true },
  { id: "cloudwatch", displayName: "CloudWatch", icon: "📊", category: "Management", backendServices: ["CloudWatch", "CloudWatchLogs"], hasPage: true },
  { id: "cognito", displayName: "Cognito", icon: "👤", category: "Security", backendServices: ["Cognito", "CognitoIdentity"], hasPage: true },
  { id: "dynamodb", displayName: "DynamoDB", icon: "🗃️", category: "Database", backendServices: ["DynamoDB"], hasPage: true },
  { id: "eventbridge", displayName: "EventBridge", icon: "📡", category: "Integration", backendServices: ["EventBridge"], hasPage: true },
  { id: "iam", displayName: "IAM", icon: "🔒", category: "Security", backendServices: ["IAM"], hasPage: true },
  { id: "kinesis", displayName: "Kinesis", icon: "🌊", category: "Analytics", backendServices: ["Kinesis"], hasPage: true },
  { id: "kms", displayName: "KMS", icon: "🔑", category: "Security", backendServices: ["KMS"], hasPage: true },
  { id: "lambda", displayName: "Lambda", icon: "⚡", category: "Compute", backendServices: ["Lambda"], hasPage: true },
  { id: "neptune", displayName: "Neptune", icon: "🔮", category: "Database", backendServices: ["Neptune", "NeptuneData", "NeptuneGraph"], hasPage: true },
  { id: "route53", displayName: "Route 53", icon: "🌐", category: "Networking", backendServices: ["Route53"], hasPage: true },
  { id: "s3", displayName: "S3", icon: "📦", category: "Storage", backendServices: ["S3"], hasPage: true },
  { id: "scheduler", displayName: "Scheduler", icon: "⏱️", category: "Management", backendServices: ["Scheduler"], hasPage: false },
  { id: "secretsmanager", displayName: "Secrets Manager", icon: "🗝️", category: "Security", backendServices: ["SecretsManager"], hasPage: false },
  { id: "ses", displayName: "SES", icon: "✉️", category: "Messaging", backendServices: ["SESv2"], hasPage: true },
  { id: "sfn", displayName: "Step Functions", icon: "🔀", category: "Integration", backendServices: ["SFN"], hasPage: true },
  { id: "sns", displayName: "SNS", icon: "📨", category: "Messaging", backendServices: ["SNS"], hasPage: true },
  { id: "sqs", displayName: "SQS", icon: "📋", category: "Messaging", backendServices: ["SQS"], hasPage: true },
  { id: "ssm", displayName: "Systems Manager", icon: "📋", category: "Management", backendServices: ["SSM"], hasPage: true },
  { id: "sts", displayName: "STS", icon: "🎫", category: "Security", backendServices: ["STS"], hasPage: true },
  { id: "timestream", displayName: "Timestream", icon: "⏳", category: "Analytics", backendServices: ["TimestreamQuery", "TimestreamWrite"], hasPage: true },
  { id: "wafv2", displayName: "WAF", icon: "🛡️", category: "Security", backendServices: ["WAFv2"], hasPage: false },
];

/** Display categories used for sidebar section headings, in display order. */
export const SERVICE_CATEGORIES = [
  "Compute",
  "Storage",
  "Database",
  "Networking",
  "Security",
  "Management",
  "Analytics",
  "Integration",
  "Messaging",
] as const;

/** Map from display service ID to its catalog entry. */
export const SERVICE_MAP = new Map(SERVICE_CATALOG.map((s) => [s.id, s]));

/** Map from backend dispatcher service name to display service ID. */
export const BACKEND_TO_DISPLAY = new Map<string, string>();
for (const entry of SERVICE_CATALOG) {
  for (const bs of entry.backendServices) {
    BACKEND_TO_DISPLAY.set(bs, entry.id);
  }
}
