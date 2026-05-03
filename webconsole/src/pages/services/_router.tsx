/**
 * Service route resolver. Maps /services/:serviceId to the corresponding
 * service page component. Services without a dedicated page show a
 * "Coming soon" placeholder using catalog metadata.
 */
import { useParams } from "react-router";
import { useTranslation } from "react-i18next";
import { SERVICE_MAP } from "@/lib/service-catalog";
import { SSMPage } from "./ssm";
import { S3Page } from "./s3";
import { DynamoDBPage } from "./dynamodb";
import { SQSPage } from "./sqs";
import { SNSPage } from "./sns";
import { LambdaPage } from "./lambda";
import { KinesisPage } from "./kinesis";
import { IAMPage } from "./iam";
import { CloudWatchPage } from "./cloudwatch";
import { CloudWatchLogsPage } from "./cloudwatchlogs";
import { EventBridgePage } from "./eventbridge";
import { AthenaPage } from "./athena";
import { SFNPage } from "./sfn";
import { SESPage } from "./ses";
import { CloudTrailPage } from "./cloudtrail";
import { KMSPage } from "./kms";
import { CognitoPage } from "./cognito";
import { APIGatewayPage } from "./apigateway";
import { CloudFrontPage } from "./cloudfront";
import { STSPage } from "./sts";
import { AppSyncPage } from "./appsync";
import { TimestreamPage } from "./timestream";
import { NeptunePage } from "./neptune";
import { ACMPage } from "./acm";
import { Route53Page } from "./route53";

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
          Unknown service.
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

  return <Page key={serviceId} />;
}
