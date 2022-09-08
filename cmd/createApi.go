/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */
package cmd

import (
	"fmt"

	"github.com/BLasan/APKCTL-Demo/impl"
	"github.com/BLasan/APKCTL-Demo/utils"
	"github.com/spf13/cobra"
)

var dpNamespace string
var serviceUrl string
var file string
var isDryRun bool

const CreateAPICmdLiteral = "api"
const CreateAPICmdShortDesc = "Create API and Deploy"
const CreateCAPImdLongDesc = `Create an API available in the namespace specified by flag (--namespace, -n)
Create an API available in the namesapce specified by flag (--namespace, -n)`
const createAPICmdExamples = utils.ProjectName + ` ` + CreateCmdLiteral + ` ` + CreateAPICmdLiteral + ` petstore \
  --service-url http://localhost:9443 -v 1.0.0 -n wso2
  
  ` + utils.ProjectName + ` ` + CreateCmdLiteral + ` ` + CreateCmdLiteral + CreateAPICmdLiteral + ` petstore \
  -f ./swagger.yaml --namespace wso2
`

// CreateApiCmd represents the create API command
var CreateApiCmd = &cobra.Command{
	Use:     CreateAPICmdLiteral,
	Short:   CreateAPICmdShortDesc,
	Long:    CreateCAPImdLongDesc,
	Example: createAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("API Name: ", args[0])
		apiName := args[0]
		fmt.Println("Is dry run: ", isDryRun)
		handleCreateApi(apiName)
	},
}

func handleCreateApi(apiName string) {
	impl.CreateAPI(file, dpNamespace, serviceUrl, apiName, isDryRun)
}

func init() {
	CreateCmd.AddCommand(CreateApiCmd)
	CreateApiCmd.Flags().StringVarP(&dpNamespace, "namespace", "n", "", "Namespace of the API")
	CreateApiCmd.Flags().StringVar(&serviceUrl, "service-url", "", "Backend Service URL")
	CreateApiCmd.Flags().StringVarP(&file, "file", "f", "", "Path to swagger/OAS definition/GraphQL SDL/WSDL")
	CreateApiCmd.Flags().BoolVar(&isDryRun, "dry-run", false, "Generate configuration files")
}
