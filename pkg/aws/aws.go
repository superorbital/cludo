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

type AWSPlugin struct {
	awsRegion  string
	awsSession *session.Session
	stsClient  *sts.STS

	sessionDuration int64
}

func NewAWSPlugin(accessKeyID string, secretAccessKey string, region string, sessionDuration time.Duration) (*AWSPlugin, error) {
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

	return &AWSPlugin{
		awsRegion:  region,
		awsSession: sess,
		stsClient:  sts.New(sess),

		sessionDuration: int64(sessionDuration.Seconds()),
	}, nil
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
			AWSSessionExpiryEnvVar:   output.Credentials.Expiration.String(),
			AWSRegionEnvVar:          ap.awsRegion,
		},
	}, nil
}
