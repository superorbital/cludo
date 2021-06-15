package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/superorbital/cludo/models"
)

const AWSRegionEnvVar = "AWS_REGION"
const AWSAccessKeyIDEnvVar = "AWS_ACCESS_KEY_ID"
const AWSSecretAccessKeyEnvVar = "AWS_SECRET_ACCESS_KEY"
const AWSSessionTokenEnvVar = "AWS_SESSION_TOKEN"
const AWSSessionExpiryEnvVar = "AWS_SESSION_EXPIRATION"

// This has to be a var so we can take a pointer for it
var CLUDOD_SESSION = "CLUDOD_SESSION"

type STSGetSessionTokenAPI interface {
	GetSessionToken(input *sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error)
	AssumeRole(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error)
}

type AWSPlugin struct {
	awsRegion       string
	stsClient       STSGetSessionTokenAPI
	roleArn         string
	sessionDuration int64
}

func New(region string, stsClient STSGetSessionTokenAPI, sessionDuration time.Duration, roleArn string) *AWSPlugin {
	return &AWSPlugin{
		awsRegion:       region,
		stsClient:       stsClient,
		roleArn:         roleArn,
		sessionDuration: int64(sessionDuration.Seconds()),
	}
}

func NewAWSPlugin(accessKeyID string, secretAccessKey string, region string, sessionDuration time.Duration, roleARN string) (*AWSPlugin, error) {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentials(
				accessKeyID,
				secretAccessKey,
				"",
			),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create AWS session: %v", err)
	}
	stsClient := sts.New(sess)

	return New(region, stsClient, sessionDuration, roleARN), nil
}

func (ap *AWSPlugin) GenerateEnvironment() (*models.ModelsEnvironmentResponse, error) {
	output, err := ap.stsClient.GetSessionToken(&sts.GetSessionTokenInput{
		DurationSeconds: &ap.sessionDuration,
	})
	if err != nil {
		return nil, err
	}
	credentials := output.Credentials
	if ap.roleArn != "" {
		assumeRoleOutput, err := ap.stsClient.AssumeRole(&sts.AssumeRoleInput{
			RoleArn:         &ap.roleArn,
			RoleSessionName: &CLUDOD_SESSION,
		})
		if err != nil {
			return nil, fmt.Errorf("Failed to assume-role %v: %v", err, ap.roleArn)
		}
		credentials = assumeRoleOutput.Credentials
	}
	return &models.ModelsEnvironmentResponse{
		Bundle: models.ModelsEnvironmentBundle{
			AWSAccessKeyIDEnvVar:     *credentials.AccessKeyId,
			AWSSecretAccessKeyEnvVar: *credentials.SecretAccessKey,
			AWSSessionTokenEnvVar:    *credentials.SessionToken,
			AWSSessionExpiryEnvVar:   output.Credentials.Expiration.Format(time.RFC3339),
			AWSRegionEnvVar:          ap.awsRegion,
		},
	}, nil
}
