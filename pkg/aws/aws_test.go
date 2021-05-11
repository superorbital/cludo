package aws_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/stretchr/testify/assert"
	"github.com/superorbital/cludo/models"
	"github.com/superorbital/cludo/pkg/aws"
)

type mockSTSAPI struct {
	expected *sts.Credentials
}

func (m mockSTSAPI) GetSessionToken(input *sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error) {
	return &sts.GetSessionTokenOutput{
		Credentials: m.expected,
	}, nil
}

func TestGenerateEnvironment(t *testing.T) {
	type test struct {
		name            string
		region          string
		sessionDuration time.Duration
		credentials     *sts.Credentials
		want            *models.ModelsEnvironmentResponse
		wantErr         error
	}

	accessKeyId := "test-access-key-id"
	secretAccessKey := "test-secret-access-key"
	sessionToken := "test-session-token"
	expiration := time.Now().AddDate(0, 0, 1)

	tests := []test{
		{
			name:            "Test basic",
			region:          "test-region",
			sessionDuration: time.Duration(10000),
			credentials: &sts.Credentials{
				AccessKeyId:     &accessKeyId,
				SecretAccessKey: &secretAccessKey,
				SessionToken:    &sessionToken,
				Expiration:      &expiration,
			},
			want: &models.ModelsEnvironmentResponse{
				Bundle: models.ModelsEnvironmentBundle{
					aws.AWSAccessKeyIDEnvVar:     accessKeyId,
					aws.AWSSecretAccessKeyEnvVar: secretAccessKey,
					aws.AWSRegionEnvVar:          "test-region",
					aws.AWSSessionExpiryEnvVar:   expiration.String(),
					aws.AWSSessionTokenEnvVar:    sessionToken,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			plugin := aws.New(tc.region, &mockSTSAPI{expected: tc.credentials}, tc.sessionDuration)
			actual, actualErr := plugin.GenerateEnvironment()

			assert.EqualValues(t, tc.want, actual)
			assert.EqualValues(t, tc.wantErr, actualErr)
		})
	}
}
