package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

type API struct {
	IAM iamiface.IAMAPI
	SSM ssmiface.SSMAPI
	STS stsiface.STSAPI
}

type Config struct {
	RoleARN    string  `yaml:"roleARN"`
	Endpoint   *string `yaml:"endpoint"`
	Region     *string `yaml:"region"`
	DisableSSL *bool   `yaml:"disableSSL"`
	MaxRetries *int    `yaml:"maxRetries"`
}

func MergeConfigs(c1, c2 *Config) (c *Config) {
	c = &Config{
		RoleARN:    c1.RoleARN,
		Endpoint:   c1.Endpoint,
		Region:     c1.Region,
		DisableSSL: c1.DisableSSL,
		MaxRetries: c1.MaxRetries,
	}

	if strings.TrimSpace(c2.RoleARN) != "" {
		c.RoleARN = c2.RoleARN
	}
	if c2.Endpoint != nil && strings.TrimSpace(*c2.Endpoint) != "" {
		c.Endpoint = c2.Endpoint
	}
	if c2.Region != nil && strings.TrimSpace(*c2.Region) != "" {
		c.Region = c2.Region
	}
	if c2.DisableSSL != nil {
		c.DisableSSL = c2.DisableSSL
	}
	if c2.MaxRetries != nil {
		c.MaxRetries = c2.MaxRetries
	}

	return
}

func (c *Config) ToAWSClientConfig() *aws.Config {
	return &aws.Config{
		Endpoint:   c.Endpoint,
		Region:     c.Region,
		DisableSSL: c.DisableSSL,
		MaxRetries: c.MaxRetries,
	}
}

func NewSession(config *Config) (*session.Session, error) {
	return session.NewSession(config.ToAWSClientConfig())
}

func NewAPI(sess *session.Session, config *Config) *API {
	awsClientConfig := config.ToAWSClientConfig()
	if strings.TrimSpace(config.RoleARN) != "" {
		awsClientConfig.Credentials = stscreds.NewCredentials(sess, config.RoleARN)
	}

	return &API{
		IAM: iam.New(sess, awsClientConfig),
		SSM: ssm.New(sess, awsClientConfig),
		STS: sts.New(sess, awsClientConfig),
	}
}
