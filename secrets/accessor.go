package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

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

func Sessions() (*session.Session, error) {
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

func NewSSMClient() *SSM {
	// Create AWS Session
	sess, err := Sessions()
	if err != nil {
		log.Println(err)
		return nil
	}
	ssmsvc := &SSM{ssm.New(sess)}
	// Return SSM client
	return ssmsvc
}

func (ssmsvc *SSM) GetValues() (map[string]string, error) {
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
		s := strings.Split(*parameter.Name, "/")
		key := s[len(s)-1]
		output[key] = *parameter.Value
	}
	return output, nil
}

func main() {
	ssmsvc := NewSSMClient()
	result, err := ssmsvc.GetValues()
	if err != nil {
		log.Fatal(err)
	}
	jsonString, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonString))
}
