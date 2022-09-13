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

package impl

import (
	"fmt"

	k8sUtils "github.com/BLasan/APKCTL-Demo/k8s"
	"github.com/BLasan/APKCTL-Demo/utils"
)

func UninstallPlatform() {
	// Uninstall k8s components that were installed through `apkctl install platform` command
	fmt.Println("Platform uninstallation initialized...")

	// Envoy Gateway uninstallation (Data Plane profile)
	// Delete the Gateway
	if err := k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sDelete, "gateway/eg"); err != nil {
		utils.HandleErrorAndExit("Error deleting the Gateway", err)
	}

	// Delete the GatewayClass
	if err := k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sDelete, "gc/eg"); err != nil {
		utils.HandleErrorAndExit("Error deleting the Gateway Class", err)
	}

	// Uninstall Envoy Gateway
	if err := k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sDelete, "-f", envoyGatewayInstallYaml); err != nil {
		utils.HandleErrorAndExit("Error uninstalling Envoy Gateway", err)
	}

	// Uninstall Gateway API CRDs
	if err := k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sDelete, "-f", gatewayAPICRDsYaml); err != nil {
		utils.HandleErrorAndExit("Error uninstalling Gateway API CRDs", err)
	}

	fmt.Println("\nUninstallation completed!")
}
