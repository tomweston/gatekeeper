// Copyright © 2018 The GΛTEKEEPER Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/spf13/cobra"
)

// getCmd represents the get command

var gSecret string
var gRegion string
var gKey string
var gAll bool

// var gKey string
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the value of a secret stored AWS Secrets Manager",
	Long: `Get the value of a secret stored AWS Secrets Manager.

	Access to secrets is granted via fine-grained 
	AWS Identity and Access Management (IAM) policies.`,
	Run: func(cmd *cobra.Command, args []string) {
		getSecret(cmd)
	},
}

func init() {
	getCmd.Flags().StringVarP(&gRegion, "region", "r", "", "The region the secret is stored in")
	getCmd.Flags().StringVarP(&gSecret, "secret", "s", "", "The name of the secret stored in AWS Secrets Manager")
	getCmd.Flags().BoolVarP(&gAll, "all", "a", false, "Display all details for secret")
	getCmd.Flags().StringVarP(&gKey, "key", "k", "", "Key of item in secret string")
	rootCmd.AddCommand(getCmd)
}

func getSecret(cmd *cobra.Command) {
	key := cmd.Flag("key").Value.String()
	region := cmd.Flag("region").Value.String()
	secret := cmd.Flag("secret").Value.String()
	svc := secretsmanager.New(session.New(), aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secret),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
				//fmt.Println("You do not have permission to reveal the requested secret!")
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	displayResult(result, key)
}

func displayResult(result *secretsmanager.GetSecretValueOutput, key string) {
	if gAll {
		fmt.Println(result)
	} else {
		var secretResult map[string]interface{}
		json.Unmarshal([]byte(*result.SecretString), &secretResult)

		// If a string key provided try and return the value of the key
		if key != "" {
			value, exists := secretResult[key]

			// If key exists return value, else respond with invalid key message
			if exists {
				fmt.Println(value)
			} else {
				fmt.Println("Key doesn't exist")
			}
		} else {
			fmt.Println(*result.SecretString)
		}
	}
}
