package gox_aws

import "github.com/aws/aws-sdk-go/aws/session"

// Dynamo DB setup configuration
type Config struct {
	Region        string `json:"region" yaml:"region"`
	Endpoint      string `json:"endpoint" yaml:"endpoint"`
	AwsAccessKey  string `json:"aws_access_key" yaml:"aws_access_key"`
	AwsSecretKey  string `json:"aws_secret_key" yaml:"aws_secret_key"`
	AwsSessionKey string `json:"session_token" yaml:"session_token"`
}

type AwsContext interface {
	GetSession() *session.Session
}
