package gox_aws

import "github.com/aws/aws-sdk-go/aws/session"

// Dynamo DB setup configuration
type Config struct {
	Region        string `json:"region,string" yaml:"region,string"`
	Endpoint      string `json:"endpoint,string" yaml:"endpoint,string"`
	AwsAccessKey  string `json:"aws_access_key,string" yaml:"aws_access_key,string"`
	AwsSecretKey  string `json:"aws_secret_key,string" yaml:"aws_secret_key,string"`
	AwsSessionKey string `json:"session_token,string" yaml:"session_token,string"`
}

type AwsContext interface {
	GetSession() *session.Session
}
