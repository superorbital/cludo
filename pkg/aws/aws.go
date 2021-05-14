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

type STSGetSessionTokenAPI interface {
	GetSessionToken(input *sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error)
}

type AWSPlugin struct {
	awsRegion string
	stsClient STSGetSessionTokenAPI

	sessionDuration int64
}

func New(region string, stsClient STSGetSessionTokenAPI, sessionDuration time.Duration) *AWSPlugin {
	return &AWSPlugin{
		awsRegion: region,
		stsClient: stsClient,

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

	// If an AWS Role ARN is specified, assume-role into it first.
	if roleARN != "" {
		output, err := stsClient.AssumeRole(&sts.AssumeRoleInput{
			RoleArn: &roleARN,
		})
		if err != nil {
			return nil, fmt.Errorf("Failed to assume-role %v: %v", err, roleARN)
		}

		sess, err = session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(
				*output.Credentials.AccessKeyId,
				*output.Credentials.SecretAccessKey,
				*output.Credentials.SessionToken),
			Region: aws.String(region),
		})
		if err != nil {
			return nil, fmt.Errorf("Failed to create AWS session: %v", err)
		}
		stsClient = sts.New(sess)
	}

	return New(region, stsClient, sessionDuration), nil
}

func (ap *AWSPlugin) GenerateEnvironment() (*models.ModelsEnvironmentResponse, error) {
	output, err := ap.stsClient.GetSessionToken(&sts.GetSessionTokenInput{
		DurationSeconds: &ap.sessionDuration,
	})
	if err != nil {
		return nil, err
	}
	return &models.ModelsEnvironmentResponse{
		Bundle: models.ModelsEnvironmentBundle{
			AWSAccessKeyIDEnvVar:     *output.Credentials.AccessKeyId,
			AWSSecretAccessKeyEnvVar: *output.Credentials.SecretAccessKey,
			AWSSessionTokenEnvVar:    *output.Credentials.SessionToken,
			AWSSessionExpiryEnvVar:   output.Credentials.Expiration.Format(time.RFC3339),
			AWSRegionEnvVar:          ap.awsRegion,
		},
	}, nil
}
