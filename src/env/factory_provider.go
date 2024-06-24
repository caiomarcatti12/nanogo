/*
 * Copyright 2023 Caio Matheus Marcatti Calimério
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package env

import (
	"fmt"
	"os"
)

func Factory() IEnv {
	provider := os.Getenv("ENV_PROVIDER")

	if provider == "" {
		provider = "LOCAL"
	}

	switch provider {
	case "AWS":
		region := os.Getenv("AWS_SECRET_MANAGER_REGION")
		accessKey := os.Getenv("AWS_SECRET_MANAGER_ACCESS_KEY")
		secretKey := os.Getenv("AWS_SECRET_MANAGER_SECRET_KEY")
		secretName := os.Getenv("AWS_SECRET_MANAGER_NAME")
		err := AwsLoader(region, accessKey, secretKey, secretName)

		if err != nil {
			panic(err)
		}
	case "LOCAL":
		err := LocalLoader()

		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Sprintf("Invalid env provider: %s", provider))
	}

	return NewEnv()
}
