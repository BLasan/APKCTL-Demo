/*
 * Copyright (c) 2022, WSO2 LLC. (https://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package cmd

import (
	"github.com/BLasan/APKCTL-Demo/utils"
	"github.com/spf13/cobra"
)

const GetCmdLiteral = "get"
const GetCmdShortDesc = "Get API and Deploy"
const GetCmdLongDesc = `Get new API and Deploy onto the Kubernetes Cluster`
const GetCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAPICmdLiteral + ` petstore \
	--namespace  https://localhost:9443 -v 1.0.0
	
	NOTE: The flag --environment (-e) is mandatory.
	You can either provide only the flag --apim , or all the other 4 flags (--registration --publisher --devportal --admin) without providing --apim flag.
	If you are omitting any of --registration --publisher --devportal --admin flags, you need to specify --apim flag with the API Manager endpoint. In both of the
	cases --token flag is optional and use it to specify the gateway token endpoint. This will be used for "apictl get-keys" operation.
	To add a micro integrator instance to an environment you can use the --mi flag.`

// GetCmd represents the Get command
var GetCmd = &cobra.Command{
	Use:     GetCmdLiteral,
	Short:   GetCmdShortDesc,
	Long:    GetCmdLongDesc,
	Example: GetCmdExamples,
}

// func init() {
// 	RootCmd.AddCommand(GetCmd)
// }
