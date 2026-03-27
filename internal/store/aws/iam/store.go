// Package iam provides IAM (Identity and Access Management) storage functionality for vorpalstacks.
package iam

import (
	"vorpalstacks/internal/core/storage"
)

// IAMStore provides a unified interface to all IAM-related stores.
type IAMStore struct {
	users                *UserStore
	accessKeys           *AccessKeyStore
	loginProfiles        *LoginProfileStore
	groups               *GroupStore
	userGroups           *UserGroupStore
	roles                *RoleStore
	instanceProfiles     *InstanceProfileStore
	policies             *PolicyStore
	inlinePolicies       *InlinePolicyStore
	attachedPolicies     *AttachedPolicyStore
	mfaDevices           *MFADeviceStore
	passwordPolicy       *PasswordPolicyStore
	accountAlias         *AccountAliasStore
	serverCertificates   *ServerCertificateStore
	signingCertificates  *SigningCertificateStore
	sshPublicKeys        *SSHPublicKeyStore
	serviceSpecificCreds *ServiceSpecificCredentialStore
	samlProviders        *SAMLProviderStore
	oidcProviders        *OpenIDConnectProviderStore
	accountSettings      *AccountSettingsStore
	serviceLastAccessed  *ServiceLastAccessedDetailsJobStore
	arnBuilder           *ARNBuilder
	accountID            string
}

// NewIAMStore creates a new IAMStore instance with the specified storage and account ID.
func NewIAMStore(store storage.BasicStorage, accountID string) *IAMStore {
	s := &IAMStore{
		users:                NewUserStore(store, accountID),
		accessKeys:           NewAccessKeyStore(store),
		loginProfiles:        NewLoginProfileStore(store),
		groups:               NewGroupStore(store, accountID),
		userGroups:           NewUserGroupStore(store),
		roles:                NewRoleStore(store, accountID),
		instanceProfiles:     NewInstanceProfileStore(store, accountID),
		policies:             NewPolicyStore(store, accountID),
		inlinePolicies:       NewInlinePolicyStore(store),
		attachedPolicies:     NewAttachedPolicyStore(store),
		mfaDevices:           NewMFADeviceStore(store, accountID),
		passwordPolicy:       NewPasswordPolicyStore(store),
		accountAlias:         NewAccountAliasStore(store),
		serverCertificates:   NewServerCertificateStore(store, accountID),
		signingCertificates:  NewSigningCertificateStore(store),
		sshPublicKeys:        NewSSHPublicKeyStore(store),
		serviceSpecificCreds: NewServiceSpecificCredentialStore(store, accountID),
		samlProviders:        NewSAMLProviderStore(store, accountID),
		oidcProviders:        NewOpenIDConnectProviderStore(store, accountID),
		accountSettings:      NewAccountSettingsStore(store),
		serviceLastAccessed:  NewServiceLastAccessedDetailsJobStore(store),
		arnBuilder:           NewARNBuilder(accountID),
		accountID:            accountID,
	}
	s.initializeAWSManagedPolicies()
	return s
}

func (s *IAMStore) initializeAWSManagedPolicies() {
	awsPolicies := []AWSManagedPolicy{
		{
			ARN:         "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
			Path:        "/service-role/",
			PolicyName:  "AWSLambdaBasicExecutionRole",
			Description: "Provides write permissions to CloudWatch Logs, and write permissions to X-Ray.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["logs:CreateLogGroup","logs:CreateLogStream","logs:PutLogEvents"],"Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/service-role/AWSLambdaRole",
			Path:        "/service-role/",
			PolicyName:  "AWSLambdaRole",
			Description: "Provides permissions to invoke AWS Lambda functions.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["lambda:InvokeFunction"],"Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/AdministratorAccess",
			Path:        "/",
			PolicyName:  "AdministratorAccess",
			Description: "Provides full access to all AWS services and resources.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/PowerUserAccess",
			Path:        "/",
			PolicyName:  "PowerUserAccess",
			Description: "Provides full access to AWS services and resources, but does not allow management of IAM users and groups.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","NotAction":["iam:*","organizations:*","account:*"],"Resource":"*"},{"Effect":"Allow","Action":["iam:CreateServiceLinkedRole","iam:DeleteServiceLinkedRole","iam:ListRoles","organizations:DescribeOrganization","account:ListRegions"],"Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/ReadOnlyAccess",
			Path:        "/",
			PolicyName:  "ReadOnlyAccess",
			Description: "Provides read-only access to AWS services and resources.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["a4b:Get*","a4b:List*","a4b:Search*","access-analyzer:Get*","access-analyzer:List*","account:Get*","account:List*","acm:Describe*","acm:List*","acm-pca:Describe*","acm-pca:Get*","acm-pca:List*","amplify:Get*","amplify:List*","apigateway:GET","appconfig:Get*","appconfig:List*","appflow:Describe*","appflow:List*","appmesh:Describe*","appmesh:List*","apprunner:Describe*","apprunner:List*","appstream:Describe*","appstream:Get*","appstream:List*","appsync:Get*","appsync:List*","athena:Get*","athena:List*","autoscaling:Describe*","autoscaling-plans:Describe*","autoscaling-plans:Get*","backup:Describe*","backup:Get*","backup:List*","batch:Describe*","batch:List*","cloud9:Describe*","cloud9:List*","clouddirectory:List*","clouddirectory:Get*","cloudformation:Describe*","cloudformation:Get*","cloudformation:List*","cloudfront:Get*","cloudfront:List*","cloudhsm:Describe*","cloudhsm:Get*","cloudhsm:List*","cloudsearch:Describe*","cloudsearch:List*","cloudtrail:Describe*","cloudtrail:Get*","cloudtrail:List*","cloudwatch:Describe*","cloudwatch:Get*","cloudwatch:List*","codeartifact:Describe*","codeartifact:Get*","codeartifact:List*","codebuild:BatchGet*","codebuild:List*","codecommit:BatchGet*","codecommit:Get*","codecommit:List*","codedeploy:BatchGet*","codedeploy:Get*","codedeploy:List*","codeguru-profiler:Describe*","codeguru-profiler:Get*","codeguru-profiler:List*","codeguru-reviewer:Describe*","codeguru-reviewer:Get*","codeguru-reviewer:List*","codepipeline:Get*","codepipeline:List*","codestar:Describe*","codestar:Get*","codestar:List*","cognito-identity:Describe*","cognito-identity:Get*","cognito-identity:List*","cognito-idp:Describe*","cognito-idp:Get*","cognito-idp:List*","cognito-sync:Describe*","cognito-sync:Get*","cognito-sync:List*","comprehend:Describe*","comprehend:List*","config:Describe*","config:Get*","config:List*","connect:Describe*","connect:Get*","connect:List*","cur:Describe*","cur:Get*","cur:Validate*","databrew:Describe*","databrew:List*","dataexchange:Get*","dataexchange:List*","datasync:Describe*","datasync:List*","dax:Describe*","dax:List*","detective:Get*","detective:List*","devicefarm:Get*","devicefarm:List*","directconnect:Describe*","discovery:Describe*","discovery:List*","dlm:Get*","dms:Describe*","dms:List*","ds:Describe*","ds:Get*","ds:List*","dynamodb:Describe*","dynamodb:Get*","dynamodb:List*","ec2:Describe*","ec2:Get*","ecr:BatchGet*","ecr:Describe*","ecr:Get*","ecr:List*","ecs:Describe*","ecs:Get*","ecs:List*","eks:Describe*","eks:Get*","eks:List*","elasticache:Describe*","elasticache:List*","elasticbeanstalk:Check*","elasticbeanstalk:Describe*","elasticbeanstalk:List*","elasticbeanstalk:Request*","elasticfilesystem:Describe*","elasticloadbalancing:Describe*","elasticmapreduce:Describe*","elasticmapreduce:Get*","elasticmapreduce:List*","elastictranscoder:List*","es:Describe*","es:Get*","es:List*","events:Describe*","events:List*","events:Test*","firehose:Describe*","firehose:List*","fms:Describe*","fms:Get*","fms:List*","forecast:Describe*","forecast:Get*","forecast:List*","freertos:Describe*","freertos:Get*","freertos:List*","fsx:Describe*","fsx:List*","gamelift:Describe*","gamelift:Get*","gamelift:List*","glacier:Describe*","glacier:Get*","glacier:List*","globalaccelerator:Describe*","globalaccelerator:List*","glue:BatchGet*","glue:Get*","glue:List*","greengrass:Get*","greengrass:List*","groundstation:Describe*","groundstation:Get*","groundstation:List*","guardduty:Get*","guardduty:List*","health:Describe*","health:Get*","health:List*","iam:Generate*","iam:Get*","iam:List*","identitystore:Describe*","identitystore:Get*","identitystore:List*","imagebuilder:Get*","imagebuilder:List*","inspector:Describe*","inspector:Get*","inspector:List*","iot:Describe*","iot:Get*","iot:List*","iotanalytics:Describe*","iotanalytics:Get*","iotanalytics:List*","iotevents:Describe*","iotevents:Get*","iotevents:List*","kafka:Describe*","kafka:Get*","kafka:List*","kendra:Describe*","kendra:Get*","kendra:List*","kinesis:Describe*","kinesis:Get*","kinesis:List*","kinesisanalytics:Describe*","kinesisanalytics:Get*","kinesisanalytics:List*","kinesisvideo:Describe*","kinesisvideo:Get*","kinesisvideo:List*","kms:Describe*","kms:Get*","kms:List*","lakeformation:Describe*","lakeformation:Get*","lakeformation:List*","lambda:Get*","lambda:List*","lex:Get*","lex:List*","license-manager:Get*","license-manager:List*","lightsail:Get*","lightsail:Is*","logs:Describe*","logs:Get*","logs:List*","machinelearning:Describe*","machinelearning:Get*","macie2:Get*","macie2:List*","managedblockchain:Describe*","managedblockchain:Get*","managedblockchain:List*","mediaconnect:Describe*","mediaconnect:List*","mediaconvert:Describe*","mediaconvert:Get*","mediaconvert:List*","medialive:Describe*","medialive:Get*","medialive:List*","mediapackage:Describe*","mediapackage:List*","mediastore:Describe*","mediastore:Get*","mediastore:List*","mediatailor:Describe*","mediatailor:Get*","mediatailor:List*","mobiletargeting:Get*","mobiletargeting:List*","mobileanalytics:Get*","mq:Describe*","mq:List*","neptune:Describe*","neptune:List*","networkmanager:Describe*","networkmanager:Get*","networkmanager:List*","opsworks:Describe*","opsworks:Get*","opsworks:List*","opsworks-cm:Describe*","opsworks-cm:List*","organizations:Describe*","organizations:List*","outposts:Get*","outposts:List*","personalize:Describe*","personalize:Get*","personalize:List*","pi:Describe*","pi:Get*","polly:Describe*","polly:Get*","polly:List*","pricing:Describe*","pricing:Get*","pricing:List*","qldb:Describe*","qldb:Get*","qldb:List*","quicksight:Describe*","quicksight:Get*","quicksight:List*","ram:Get*","ram:List*","rds:Describe*","rds:List*","redshift:Describe*","redshift:Get*","redshift:List*","rekognition:Describe*","rekognition:Get*","rekognition:List*","resource-groups:Get*","resource-groups:List*","robomaker:Describe*","robomaker:Get*","robomaker:List*","route53:Get*","route53:List*","route53domains:Check*","route53domains:Get*","route53domains:List*","route53resolver:Get*","route53resolver:List*","s3:Get*","s3:List*","sagemaker:Describe*","sagemaker:Get*","sagemaker:List*","schemas:Describe*","schemas:Get*","schemas:List*","sdb:Get*","sdb:List*","secretsmanager:Describe*","secretsmanager:Get*","secretsmanager:List*","securityhub:Describe*","securityhub:Get*","securityhub:List*","serverlessrepo:Get*","serverlessrepo:List*","servicecatalog:Describe*","servicecatalog:Get*","servicecatalog:List*","servicediscovery:Get*","servicediscovery:List*","ses:Get*","ses:List*","shield:Describe*","shield:Get*","shield:List*","signer:Describe*","signer:Get*","signer:List*","sns:Get*","sns:List*","sqs:Get*","sqs:List*","ssm:Describe*","ssm:Get*","ssm:List*","states:Describe*","states:Get*","states:List*","storagegateway:Describe*","storagegateway:Get*","storagegateway:List*","sts:Get*","sumerian:List*","support:Describe*","support:Get*","swf:Count*","swf:Describe*","swf:Get*","swf:List*","textract:Get*","textract:List*","transcribe:Get*","transcribe:List*","transfer:Describe*","transfer:Get*","transfer:List*","translate:Describe*","translate:Get*","translate:List*","waf:Get*","waf:List*","waf-regional:Get*","waf-regional:List*","wafv2:Describe*","wafv2:Get*","wafv2:List*","wellarchitected:Get*","wellarchitected:List*","workdocs:Describe*","workdocs:Get*","workdocs:List*","worklink:Describe*","worklink:List*","workmail:Describe*","workmail:Get*","workmail:List*","workspaces:Describe*","workspaces:List*","xray:BatchGet*","xray:Get*","xray:List*"],"Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/AmazonS3FullAccess",
			Path:        "/",
			PolicyName:  "AmazonS3FullAccess",
			Description: "Provides full access to all buckets via the AWS Management Console.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:*","Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess",
			Path:        "/",
			PolicyName:  "AmazonS3ReadOnlyAccess",
			Description: "Provides read-only access to all buckets via the AWS Management Console.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:Get*","s3:List*"],"Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess",
			Path:        "/",
			PolicyName:  "AmazonDynamoDBFullAccess",
			Description: "Provides full access to Amazon DynamoDB.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"dynamodb:*","Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess",
			Path:        "/",
			PolicyName:  "AmazonDynamoDBReadOnlyAccess",
			Description: "Provides read-only access to Amazon DynamoDB.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["dynamodb:Describe*","dynamodb:List*","dynamodb:Get*","dynamodb:BatchGet*","dynamodb:Query","dynamodb:Scan"],"Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/AmazonSQSFullAccess",
			Path:        "/",
			PolicyName:  "AmazonSQSFullAccess",
			Description: "Provides full access to Amazon SQS.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"sqs:*","Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/AmazonSQSReadOnlyAccess",
			Path:        "/",
			PolicyName:  "AmazonSQSReadOnlyAccess",
			Description: "Provides read-only access to Amazon SQS.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["sqs:Get*","sqs:List*"],"Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/AmazonSNSFullAccess",
			Path:        "/",
			PolicyName:  "AmazonSNSFullAccess",
			Description: "Provides full access to Amazon SNS.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"sns:*","Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/AmazonSNSReadOnlyAccess",
			Path:        "/",
			PolicyName:  "AmazonSNSReadOnlyAccess",
			Description: "Provides read-only access to Amazon SNS.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["sns:Get*","sns:List*"],"Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/AWSLambdaFullAccess",
			Path:        "/",
			PolicyName:  "AWSLambdaFullAccess",
			Description: "Provides full access to AWS Lambda.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"lambda:*","Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/AWSKeyManagementServicePowerUser",
			Path:        "/",
			PolicyName:  "AWSKeyManagementServicePowerUser",
			Description: "Provides full access to AWS KMS.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"kms:*","Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/AmazonKinesisFullAccess",
			Path:        "/",
			PolicyName:  "AmazonKinesisFullAccess",
			Description: "Provides full access to Amazon Kinesis.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"kinesis:*","Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/CloudWatchFullAccess",
			Path:        "/",
			PolicyName:  "CloudWatchFullAccess",
			Description: "Provides full access to Amazon CloudWatch.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"cloudwatch:*","Resource":"*"}]}`,
		},
		{
			ARN:         "arn:aws:iam::aws:policy/CloudWatchLogsFullAccess",
			Path:        "/",
			PolicyName:  "CloudWatchLogsFullAccess",
			Description: "Provides full access to Amazon CloudWatch Logs.",
			Document:    `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"logs:*","Resource":"*"}]}`,
		},
	}

	for _, p := range awsPolicies {
		_ = s.policies.CreateAWSManagedPolicy(p)
	}
}

// Users returns the user store for IAM operations.
func (s *IAMStore) Users() UserStoreInterface {
	return s.users
}

// AccessKeys returns the access key store for IAM operations.
func (s *IAMStore) AccessKeys() AccessKeyStoreInterface {
	return s.accessKeys
}

// LoginProfiles returns the login profile store for IAM operations.
func (s *IAMStore) LoginProfiles() LoginProfileStoreInterface {
	return s.loginProfiles
}

// Groups returns the group store for IAM operations.
func (s *IAMStore) Groups() GroupStoreInterface {
	return s.groups
}

// UserGroups returns the user-group association store for IAM operations.
func (s *IAMStore) UserGroups() UserGroupStoreInterface {
	return s.userGroups
}

// Roles returns the role store for IAM operations.
func (s *IAMStore) Roles() RoleStoreInterface {
	return s.roles
}

// InstanceProfiles returns the instance profile store for IAM operations.
func (s *IAMStore) InstanceProfiles() InstanceProfileStoreInterface {
	return s.instanceProfiles
}

// Policies returns the policy store for IAM operations.
func (s *IAMStore) Policies() PolicyStoreInterface {
	return s.policies
}

// InlinePolicies returns the inline policy store for IAM operations.
func (s *IAMStore) InlinePolicies() InlinePolicyStoreInterface {
	return s.inlinePolicies
}

// AttachedPolicies returns the attached policy store for IAM operations.
func (s *IAMStore) AttachedPolicies() AttachedPolicyStoreInterface {
	return s.attachedPolicies
}

// MFADevices returns the MFA device store for IAM operations.
func (s *IAMStore) MFADevices() MFADeviceStoreInterface {
	return s.mfaDevices
}

// PasswordPolicy returns the password policy store for IAM operations.
func (s *IAMStore) PasswordPolicy() PasswordPolicyStoreInterface {
	return s.passwordPolicy
}

// ARNBuilder returns the ARN builder for IAM resources.
func (s *IAMStore) ARNBuilder() *ARNBuilder {
	return s.arnBuilder
}

// AccountID returns the AWS account ID associated with this store.
func (s *IAMStore) AccountID() string {
	return s.accountID
}

// AccountAlias returns the account alias store.
func (s *IAMStore) AccountAlias() *AccountAliasStore {
	return s.accountAlias
}

// ServerCertificates returns the server certificate store.
func (s *IAMStore) ServerCertificates() *ServerCertificateStore {
	return s.serverCertificates
}

// SigningCertificates returns the signing certificate store.
func (s *IAMStore) SigningCertificates() *SigningCertificateStore {
	return s.signingCertificates
}

// SSHPublicKeys returns the SSH public key store.
func (s *IAMStore) SSHPublicKeys() *SSHPublicKeyStore {
	return s.sshPublicKeys
}

// ServiceSpecificCredentials returns the service-specific credential store.
func (s *IAMStore) ServiceSpecificCredentials() *ServiceSpecificCredentialStore {
	return s.serviceSpecificCreds
}

// SAMLProviders returns the SAML provider store.
func (s *IAMStore) SAMLProviders() *SAMLProviderStore {
	return s.samlProviders
}

// OpenIDConnectProviders returns the OIDC provider store.
func (s *IAMStore) OpenIDConnectProviders() *OpenIDConnectProviderStore {
	return s.oidcProviders
}

// AccountSettings returns the account settings store.
func (s *IAMStore) AccountSettings() *AccountSettingsStore {
	return s.accountSettings
}

// ServiceLastAccessed returns the service last accessed details job store.
func (s *IAMStore) ServiceLastAccessed() *ServiceLastAccessedDetailsJobStore {
	return s.serviceLastAccessed
}
