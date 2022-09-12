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
	"os/exec"
	"strings"

	k8sUtils "github.com/BLasan/APKCTL-Demo/k8s"
	"github.com/BLasan/APKCTL-Demo/utils"
)

const gatewayAPICRDsYaml = "https://github.com/envoyproxy/gateway/releases/download/v0.2.0-rc1/gatewayapi-crds.yaml"
const envoyGatewayInstallYaml = "https://github.com/envoyproxy/gateway/releases/download/v0.2.0-rc1/install.yaml"
const gatewayClassYaml = "https://raw.githubusercontent.com/envoyproxy/gateway/v0.2.0-rc1/examples/kubernetes/gatewayclass.yaml"
const gatewayYaml = "https://raw.githubusercontent.com//envoyproxy/gateway/v0.2.0-rc1/examples/kubernetes/gateway.yaml"

func InstallPlatform() {
	// Install components in K8s default cluster with default namespace

	// Envoy Gateway installation (Data Plane profile)
	// Install the Gateway API CRDs
	if err := k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sApply, "-f", gatewayAPICRDsYaml); err != nil {
		utils.HandleErrorAndExit("Error installing Gateway API CRDs", err)
	}
	// Run Envoy Gateway
	if err := k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sApply, "-f", envoyGatewayInstallYaml); err != nil {
		utils.HandleErrorAndExit("Error installing Envoy Gateway", err)
	}

	// Check pod status of `gateway-api-admission-server` to determine if it is in Running state
	for {
		podStatus := getPodStatus()
		if strings.Trim(podStatus, "\n") == "Running" {
			break
		}
	}

	// Create the GatewayClass
	if err := k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sApply, "-f", gatewayClassYaml); err != nil {
		utils.HandleErrorAndExit("Error creating the Gateway Class", err)
	}
	// Create the Gateway
	if err := k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sApply, "-f", gatewayYaml); err != nil {
		utils.HandleErrorAndExit("Error creating the Gateway", err)
	}

}

func getPodStatus() string {
	podStatus, err := exec.Command(
		"bash","-c",
		"kubectl get pods -n gateway-system --no-headers | awk '{if ($1 ~ \"gateway-api-admission-server-\") print $3}'",
	).Output()
	if err != nil {
		utils.HandleErrorAndExit("Error while checking the pod status of a pod that is required for the Envoy Gateway", err)
	}
	return string(podStatus)
}
