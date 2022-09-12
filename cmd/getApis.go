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
	"github.com/BLasan/APKCTL-Demo/impl"
	"github.com/BLasan/APKCTL-Demo/utils"
	"github.com/spf13/cobra"
)

var outputFormat string
var allNamespaces bool

const GetAPICmdLiteral = "apis"
const GetAPICmdShortDesc = "Get all APIs"
const GetCAPImdLongDesc = `Get all APIs available in a cluster or in a specific namespace defined by (--namespace, -n)
 Get an API available in the namesapce specified by flag (--namespace, -n)`
const GetAPICmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAPICmdLiteral + ` petstore \
   --service-url http://localhost:9443 -v 1.0.0 -n wso2
   
   ` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetCmdLiteral + GetAPICmdLiteral + ` petstore \
   -f ./swagger.yaml --namespace wso2
 `

// GetApiCmd represents the Get API command
var GetApiCmd = &cobra.Command{
	Use:     GetAPICmdLiteral,
	Short:   GetAPICmdShortDesc,
	Long:    GetCAPImdLongDesc,
	Example: GetAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		handleGetApis()
	},
}

func handleGetApis() {
	impl.GetAPIs(dpNamespace, outputFormat, allNamespaces)
}

func init() {
	GetCmd.AddCommand(GetApiCmd)
	GetApiCmd.Flags().StringVarP(&dpNamespace, "namespace", "n", "default", "Namespace of the API")
	GetApiCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output Format of APIs")
	GetApiCmd.Flags().BoolVar(&allNamespaces, "all-namespaces", false, "Get APIs in all namespaces")
}
