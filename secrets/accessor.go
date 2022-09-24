package awssecrets

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// SSM is a SSM API client.
type SSM struct {
	client ssmiface.SSMAPI
}

func sessions() (*session.Session, error) {
	key := "AWS_REGION"
	awsRegion, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("%s not set, exiting", key)
		os.Exit(1)
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewEnvCredentials(),
	},
	)
	return sess, err
}

func newSSMClient() *SSM {
	// Create AWS Session
	sess, err := sessions()
	if err != nil {
		log.Println(err)
		return nil
	}
	ssmsvc := &SSM{ssm.New(sess)}
	// Return SSM client
	return ssmsvc
}

func GetValues() (map[string]string, error) {
	ssmsvc := newSSMClient()
	ssmOpts := ssm.GetParametersByPathInput{
		Path:           aws.String("/"),
		WithDecryption: aws.Bool(true),
	}
	params, err := ssmsvc.client.GetParametersByPath(&ssmOpts)
	if err != nil {
		return nil, err
	}
	output := make(map[string]string)
	for _, parameter := range params.Parameters {
		output[*parameter.Name] = *parameter.Value
	}
	return output, nil
}
