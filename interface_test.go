package gox_aws

import (
	"github.com/devlibx/gox-base/serialization"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testConfig = `
aws_config:
  endpoint: http://localhot:8000
  region: ap-south-1
  aws_access_key: key_1234
  aws_secret_key: secret_1234
  session_token: session_1234

`

type testAwsConfig struct {
	Config Config `json:"aws_config" yaml:"aws_config"`
}

func TestConfig(t *testing.T) {
	tc := &testAwsConfig{}
	err := serialization.ReadYamlFromString(testConfig, tc)
	assert.NoError(t, err)
	c := tc.Config
	assert.Equal(t, "http://localhot:8000", c.Endpoint)
	assert.Equal(t, "ap-south-1", c.Region)
	assert.Equal(t, "key_1234", c.AwsAccessKey)
	assert.Equal(t, "secret_1234", c.AwsSecretKey)
	assert.Equal(t, "session_1234", c.AwsSessionKey)
}
