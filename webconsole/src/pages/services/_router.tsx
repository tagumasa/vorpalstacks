/**
 * Service route resolver. Maps /services/:serviceId to the corresponding
 * service page component. Services without a dedicated page show a
 * "Coming soon" placeholder using catalog metadata.
 *
 * Each service page is loaded via React.lazy() so Vite produces a separate
 * chunk per service, keeping the initial bundle small.
 */
import { Suspense, lazy } from "react";
import { useParams } from "react-router";
import { useTranslation } from "react-i18next";
import { SERVICE_MAP } from "@/lib/service-catalog";

const SSMPage = lazy(() => import("./ssm").then((m) => ({ default: m.SSMPage })));
const S3Page = lazy(() => import("./s3").then((m) => ({ default: m.S3Page })));
const DynamoDBPage = lazy(() => import("./dynamodb").then((m) => ({ default: m.DynamoDBPage })));
const SQSPage = lazy(() => import("./sqs").then((m) => ({ default: m.SQSPage })));
const SNSPage = lazy(() => import("./sns").then((m) => ({ default: m.SNSPage })));
const LambdaPage = lazy(() => import("./lambda").then((m) => ({ default: m.LambdaPage })));
const KinesisPage = lazy(() => import("./kinesis").then((m) => ({ default: m.KinesisPage })));
const IAMPage = lazy(() => import("./iam").then((m) => ({ default: m.IAMPage })));
const CloudWatchPage = lazy(() => import("./cloudwatch").then((m) => ({ default: m.CloudWatchPage })));
const CloudWatchLogsPage = lazy(() => import("./cloudwatchlogs").then((m) => ({ default: m.CloudWatchLogsPage })));
const EventBridgePage = lazy(() => import("./eventbridge").then((m) => ({ default: m.EventBridgePage })));
const AthenaPage = lazy(() => import("./athena").then((m) => ({ default: m.AthenaPage })));
const SFNPage = lazy(() => import("./sfn").then((m) => ({ default: m.SFNPage })));
const SESPage = lazy(() => import("./ses").then((m) => ({ default: m.SESPage })));
const CloudTrailPage = lazy(() => import("./cloudtrail").then((m) => ({ default: m.CloudTrailPage })));
const KMSPage = lazy(() => import("./kms").then((m) => ({ default: m.KMSPage })));
const CognitoPage = lazy(() => import("./cognito").then((m) => ({ default: m.CognitoPage })));
const APIGatewayPage = lazy(() => import("./apigateway").then((m) => ({ default: m.APIGatewayPage })));
const CloudFrontPage = lazy(() => import("./cloudfront").then((m) => ({ default: m.CloudFrontPage })));
const STSPage = lazy(() => import("./sts").then((m) => ({ default: m.STSPage })));
const AppSyncPage = lazy(() => import("./appsync").then((m) => ({ default: m.AppSyncPage })));
const TimestreamPage = lazy(() => import("./timestream").then((m) => ({ default: m.TimestreamPage })));
const NeptunePage = lazy(() => import("./neptune").then((m) => ({ default: m.NeptunePage })));
const ACMPage = lazy(() => import("./acm").then((m) => ({ default: m.ACMPage })));
const Route53Page = lazy(() => import("./route53").then((m) => ({ default: m.Route53Page })));
const SecretsManagerPage = lazy(() => import("./secretsmanager").then((m) => ({ default: m.SecretsManagerPage })));
const WAFv2Page = lazy(() => import("./wafv2").then((m) => ({ default: m.WAFv2Page })));
const CognitoIdentityPage = lazy(() => import("./cognito-identity").then((m) => ({ default: m.CognitoIdentityPage })));

const PAGES: Record<string, React.ComponentType> = {
  ssm: SSMPage,
  s3: S3Page,
  dynamodb: DynamoDBPage,
  sqs: SQSPage,
  sns: SNSPage,
  lambda: LambdaPage,
  kinesis: KinesisPage,
  iam: IAMPage,
  cloudwatch: CloudWatchPage,
  cloudwatchlogs: CloudWatchLogsPage,
  eventbridge: EventBridgePage,
  athena: AthenaPage,
  sfn: SFNPage,
  ses: SESPage,
  cloudtrail: CloudTrailPage,
  kms: KMSPage,
  cognito: CognitoPage,
  apigateway: APIGatewayPage,
  cloudfront: CloudFrontPage,
  sts: STSPage,
  appsync: AppSyncPage,
  timestream: TimestreamPage,
  neptune: NeptunePage,
  acm: ACMPage,
  route53: Route53Page,
  secretsmanager: SecretsManagerPage,
  wafv2: WAFv2Page,
  "cognito-identity": CognitoIdentityPage,
};

export function ServiceRoute() {
  const { t } = useTranslation();
  const { serviceId } = useParams();
  const entry = serviceId ? SERVICE_MAP.get(serviceId) : null;
  const Page = serviceId ? PAGES[serviceId] : null;

  if (!entry) {
    return (
      <div className="content-area" style={{ padding: 24 }}>
        <p style={{ color: "var(--text-muted)" }}>
          {t("common.unknownService")}
        </p>
      </div>
    );
  }

  if (!Page) {
    return (
      <div className="content-area" style={{ padding: 24 }}>
        <div className="page-header">
          <span className="page-icon">{entry.icon}</span>
          <h1>{entry.displayName}</h1>
        </div>
        <div className="loading-state">{t("common.comingSoon")}</div>
      </div>
    );
  }

  return (
    <Suspense fallback={<div className="loading-state">{t("common.loading")}</div>}>
      <Page key={serviceId} />
    </Suspense>
  );
}
