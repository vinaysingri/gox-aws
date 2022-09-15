package gox_aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/devlibx/gox-base"
	errors2 "github.com/devlibx/gox-base/errors"
	"github.com/devlibx/gox-base/util"
)

type awsSessionContext struct {
	session *session.Session
	gox.CrossFunction
}

func (a *awsSessionContext) GetSession() *session.Session {
	return a.session
}

func NewAwsContext(cf gox.CrossFunction, config Config) (ctx AwsContext, err error) {
	_ctx := &awsSessionContext{CrossFunction: cf}
	var region *string = nil
	var endpoint *string = nil
	var creds *credentials.Credentials = nil

	// Set default region
	if util.IsStringEmpty(config.Region) {
		config.Region = "ap-south-1"
	} else {
		region = aws.String(config.Region)
	}

	// use end point if configured
	if len(config.Endpoint) > 0 {
		endpoint = aws.String(config.Endpoint)
	}

	// use credentials if configured
	if !util.IsStringEmpty(config.AwsAccessKey) {
		creds = credentials.NewStaticCredentials(
			config.AwsAccessKey,
			config.AwsSecretKey,
			config.AwsSessionKey)

	}

	_ctx.session, err = session.NewSession(
		&aws.Config{
			Credentials: creds,
			Region:      region,
			Endpoint:    endpoint,
		})

	if err != nil {
		return nil, errors2.Wrap(err, "failed to create aws session: endpoint=%s, region=%s", config.Endpoint, config.Region)
	}

	return _ctx, err
}
