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

package impl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"net/url"

	k8sUtils "github.com/BLasan/APKCTL-Demo/k8s"
	"github.com/BLasan/APKCTL-Demo/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/spec"
	"gopkg.in/yaml.v2"
)

var dirPath string
var desFilePath string

func CreateAPI(filePath, namespace, serviceUrl, apiName, version string, isDryRun bool) {

	// Checking if path to API definition is provided. If not specified, use the default OpenAPI definition
	if filePath == "" {
		dir, err := utils.GetAPKCTLHomeDir()
		if err != nil {
			utils.HandleErrorAndExit("Error getting the working directory", err)
		}
		filePath = filepath.Join(dir, utils.SampleResources, utils.DefaultSwagger)
	}

	apiContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		utils.HandleErrorAndExit("Error encountered while reading API definition file", err)
	}

	definitionJsn, err := utils.ToJSON(apiContent)
	if err != nil {
		utils.HandleErrorAndExit("Error converting API definition file to json", err)
	}

	definitionVersion := utils.FindAPIDefinitionVersion(definitionJsn)

	if definitionVersion == utils.Swagger2 {

		// API definition is a Swagger file
		var swaggerSpec spec.Swagger
		err = json.Unmarshal(definitionJsn, &swaggerSpec)
		if err != nil {
			utils.HandleErrorAndExit("Error unmarshalling swagger", err)
		}
		// updateSwaggerUrl(&swaggerSpec, serviceUrl)
		createAndDeploySwaggerAPI(swaggerSpec, filePath, namespace, serviceUrl, apiName, version, isDryRun)

	} else if definitionVersion == utils.OpenAPI3 {

		// API definition is an OpenAPI Definition file
		var openAPISpec openapi3.T
		err = json.Unmarshal(definitionJsn, &openAPISpec)
		if err != nil {
			utils.HandleErrorAndExit("Error unmarshalling OpenAPI Definition", err)
		}
		createAndDeployOpenAPI(openAPISpec, filePath, namespace, serviceUrl, apiName, version, isDryRun)

	} else {
		utils.HandleErrorAndExit("Error resolving API definition. Provided file kind is not supported or not acceptable.", nil)
	}
}

func createAndDeploySwaggerAPI(swaggerSpec spec.Swagger, filePath, namespace, serviceUrl, apiName, version string, isDryRun bool) {
	httpRoute := utils.HTTPRouteConfig{}
	var parentRef utils.ParentRef

	httpRoute.ApiVersion = utils.HttpRouteApiVersion
	httpRoute.Kind = utils.HttpRouteKind
	httpRoute.HttpRouteSpec.HostNames = append(httpRoute.HttpRouteSpec.HostNames, "www.example.com")
	parentRef.Name = "eg"
	httpRoute.HttpRouteSpec.ParentRefs = append(httpRoute.HttpRouteSpec.ParentRefs, parentRef)
	httpRoute.MetaData.Name = apiName
	// httpRoute.MetaData.Namespace = namespace

	labels := make(map[string]string)

	if version == "" {
		labels["version"] = swaggerSpec.Info.Version
		httpRoute.MetaData.Labels = labels
	} else {
		labels["version"] = version
		httpRoute.MetaData.Labels = labels
	}

	var apiPath utils.Path
	var match utils.Match
	var rule utils.Rule
	var backendRef utils.BackendRef

	// Checking if service URL is provided. If not specified, deduce the service URL using the swagger definition
	if serviceUrl == "" {
		if swaggerSpec.Host != "" {
			urlScheme := ""
			for _, scheme := range swaggerSpec.Schemes {
				if scheme == "https" {
					urlScheme = utils.HttpsURLScheme
					break
				} else if scheme == "http" {
					urlScheme = utils.HttpURLScheme
				} else {
					utils.HandleErrorAndExit("Detected scheme(s) within the swagger definition are not supported", nil)
				}
			}
			serviceUrl = urlScheme + swaggerSpec.Host + swaggerSpec.BasePath
		} else {
			utils.HandleErrorAndExit("Unable to find a valid service URL.", nil)
		}
	}

	// If API definition is not specified, provide the wildcard resource as a PathPrefix
	if filePath == "" {
		apiPath.Type = utils.PathPrefix
		apiPath.Value = "/"
		match.Path = apiPath
		rule.Matches = append(rule.Matches, match)
	} else {
		counter := 1

		// path & path item
		for path := range swaggerSpec.Paths.Paths {
			// maximum 8 paths are allowed
			if counter > 8 {
				break
			}

			index := strings.IndexAny(path, "{")
			if index >= 0 {
				path = path[:index-1]
			}

			if strings.Contains(path, "/*") {
				path = strings.ReplaceAll(path, "/*", "")
			}

			// append "/api/v3" to invoke the petstore apis
			path = "/api/v3" + path

			// pathArr := strings.Split(path, "/")
			// sort.Strings(pathArr)
			// path = utils.FindPathParam(pathArr)

			apiPath.Type = utils.PathPrefix
			apiPath.Value = path
			match.Path = apiPath

			rule.Matches = append(rule.Matches, match)

			// if pathItem.Post != nil {
			// 	fmt.Println("Description Items: ", pathItem.Post.Description)
			// }

			counter++

		}
	}

	backendRef.Kind = utils.ServiceKind

	parsedURL, err := url.ParseRequestURI(serviceUrl)
	if err != nil {
		utils.HandleErrorAndExit("Error while parsing the service URL.", err)
	}

	backendRef.Name = strings.Split(parsedURL.Host, ".")[0]
	// backendRef.Namespace = serviceUrlArr[1]
	if parsedURL.Port() != "" {
		u32, err := strconv.ParseUint(parsedURL.Port(), 10, 32)
		if err != nil {
			fmt.Println("Endpoint port is not in the expected format.", err)
		}
		backendRef.Port = int(uint32(u32))
	} else {
		backendRef.Port = int(uint32(80))
	}

	rule.BackendRefs = append(rule.BackendRefs, backendRef)
	httpRoute.HttpRouteSpec.Rules = append(httpRoute.HttpRouteSpec.Rules, rule)
	if err != nil {
		utils.HandleErrorAndExit("Error extracting port number", err)
	}

	file, err := yaml.Marshal(&httpRoute)
	if err != nil {
		utils.HandleErrorAndExit("Error marshalling httproute file", err)
	}

	var configmap utils.ConfigMap
	// configmap.Name = apiName + "-" + "configmap"
	// configmap.Namespace = namespace
	// configmap.File = filePath
	// configmap.SwaggerContent = readSwaggerDef(filePath)

	if !isDryRun {
		handleDeploy(file, filePath, namespace, apiName, configmap)
	} else {
		handleDryRun(file, filePath, namespace, apiName, configmap)
	}
}

func createAndDeployOpenAPI(openAPISpec openapi3.T, filePath, namespace, serviceUrl, apiName, version string, isDryRun bool) {
	httpRoute := utils.HTTPRouteConfig{}
	var parentRef utils.ParentRef

	httpRoute.ApiVersion = utils.HttpRouteApiVersion
	httpRoute.Kind = utils.HttpRouteKind
	httpRoute.HttpRouteSpec.HostNames = append(httpRoute.HttpRouteSpec.HostNames, "www.example.com")
	parentRef.Name = "eg"
	httpRoute.HttpRouteSpec.ParentRefs = append(httpRoute.HttpRouteSpec.ParentRefs, parentRef)
	httpRoute.MetaData.Name = apiName

	labels := make(map[string]string)

	if version == "" {
		labels["version"] = openAPISpec.Info.Version
		httpRoute.MetaData.Labels = labels
	} else {
		labels["version"] = version
		httpRoute.MetaData.Labels = labels
	}

	var apiPath utils.Path
	var match utils.Match
	var rule utils.Rule
	var backendRef utils.BackendRef

	// Checking if service URL is provided. If not specified, use the service URLs provided under the OpenAPI definition
	if serviceUrl == "" {
		var serviceUrls []string
		for _, serverEntry := range openAPISpec.Servers {
			serviceUrls = append(serviceUrls, serverEntry.URL)
		}
		// We will use the first URL provided under the servers object
		serviceUrl = serviceUrls[0]
	}

	// If API definition is not specified, provide the wildcard resource as a PathPrefix
	if filePath == "" {
		apiPath.Type = utils.PathPrefix
		apiPath.Value = "/"
		match.Path = apiPath
		rule.Matches = append(rule.Matches, match)
	} else {
		counter := 1

		// path & path item
		for path := range openAPISpec.Paths {
			// maximum 8 paths are allowed
			if counter > 8 {
				break
			}

			index := strings.IndexAny(path, "{")
			if index >= 0 {
				path = path[:index-1]
			}

			// remove *
			if strings.Contains(path, "/*") {
				path = strings.ReplaceAll(path, "/*", "")
			}

			// append "/api/v3" to invoke the petstore apis
			path = "/api/v3" + path

			apiPath.Type = utils.PathPrefix
			apiPath.Value = path
			match.Path = apiPath

			rule.Matches = append(rule.Matches, match)

			counter++
		}
	}

	backendRef.Kind = utils.ServiceKind

	parsedURL, err := url.ParseRequestURI(serviceUrl)
	if err != nil {
		utils.HandleErrorAndExit("Error while parsing the service URL.", err)
	}

	backendRef.Name = strings.Split(parsedURL.Host, ".")[0]
	if parsedURL.Port() != "" {
		u32, err := strconv.ParseUint(parsedURL.Port(), 10, 32)
		if err != nil {
			fmt.Println("Endpoint port is not in the expected format.", err)
		}
		backendRef.Port = int(uint32(u32))
	} else {
		backendRef.Port = int(uint32(80))
	}

	rule.BackendRefs = append(rule.BackendRefs, backendRef)
	httpRoute.HttpRouteSpec.Rules = append(httpRoute.HttpRouteSpec.Rules, rule)

	file, err := yaml.Marshal(&httpRoute)
	if err != nil {
		utils.HandleErrorAndExit("Error marshalling httproute file.", err)
	}

	var configmap utils.ConfigMap
	// configmap.Name = apiName + "-" + "configmap"
	// configmap.Namespace = namespace
	// configmap.File = filePath
	// configmap.SwaggerContent = readSwaggerDef(filePath)

	if !isDryRun {
		handleDeploy(file, filePath, namespace, apiName, configmap)
	} else {
		handleDryRun(file, filePath, namespace, apiName, configmap)
	}
}

// Handle API deploy
func handleDeploy(file []byte, filePath, namespace, apiName string, configmap utils.ConfigMap) {
	var err error
	dirPath, err = os.MkdirTemp("", apiName)
	if err != nil {
		utils.HandleErrorAndExit("Error creating the temp directory", err)
	}

	defer os.RemoveAll(dirPath)

	desFilePath = filepath.Join(dirPath, "HTTPRouteConfig.yaml")

	// directory location can be defined in the apkctl config file
	err = ioutil.WriteFile(desFilePath, file, 0644)
	if err != nil {
		utils.HandleErrorAndExit("Error creating HTTPRouteConfig file", err)
	}

	createConfigMap(filePath, dirPath, namespace)
	// utils.CreateConfigMapFromTemplate(configmap, dirPath)

	args := []string{k8sUtils.K8sApply, k8sUtils.FilenameFlag, filepath.Join(dirPath, "")}

	err = k8sUtils.ExecuteCommand(k8sUtils.Kubectl, args...)
	if err != nil {
		utils.HandleErrorAndExit("Error Deploying the API", err)
	}
	os.RemoveAll(dirPath)

	fmt.Println("\nSuccessfully deployed " + apiName + " API into the " + namespace + " namespace")
}

// Handle the `Dry Run` option of create API command
// This will generate an API project based on the provided command and flags
func handleDryRun(file []byte, filePath, namespace, apiName string, configmap utils.ConfigMap) {
	var err error
	dirPath, err = utils.GetAPKCTLHomeDir()
	if err != nil {
		utils.HandleErrorAndExit("Error getting apkctl home directory", err)
	}

	dirPath = path.Join(dirPath, utils.APIProjectsDir, apiName)

	os.MkdirAll(dirPath, os.ModePerm)

	desFilePath = filepath.Join(dirPath, "HTTPRouteConfig.yaml")

	// directory location can be defined in the apkctl config file
	err = ioutil.WriteFile(desFilePath, file, 0644)

	if err != nil {
		utils.HandleErrorAndExit("Error creating HTTPRouteConfig file", err)
	}

	createConfigMap(filePath, dirPath, namespace)
	// utils.CreateConfigMapFromTemplate(configmap, dirPath)

	fmt.Println("Successfully created API project with HttpRouteConfig and ConfigMap files!")
	fmt.Println("API project directory: " + utils.APIProjectsDir + apiName)
}

func createConfigMap(filepath, dirPath, namespace string) {
	configmap := utils.ConfigMap{}
	configmap.ApiVersion = "v1"
	configmap.Kind = "ConfigMap"
	configmap.MetaData.Name = "swagger-configmap"

	if namespace != "" {
		configmap.MetaData.Namespace = namespace
	}

	content := readSwaggerDef(filepath)

	if content == "" {
		fmt.Println("Empty Swagger")
		// handle error and exit
	}

	data := make(map[string]string)

	data["swagger"] = content

	configmap.Data = data

	file, err := yaml.Marshal(&configmap)

	if err != nil {
		utils.HandleErrorAndExit("Error Marshaling", err)
	}

	desFilePath := path.Join(dirPath, "ConfigMap.yaml")

	// directory location can be defined in the apkctl config file
	err = ioutil.WriteFile(desFilePath, file, 0644)

	if err != nil {
		utils.HandleErrorAndExit("Error creating config file", err)
	}
}

func readSwaggerDef(filename string) string {
	if info, err := os.Stat(filename); err == nil && !info.IsDir() {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			utils.HandleErrorAndExit("Error Reading Swagger File", err)
		}
		return string(content)
	}

	return ""
}

// func updateSwaggerUrl(swagger *spec.Swagger, serviceUrl string) {
// 	if serviceUrl != "" {
// 		swagger.Host = serviceUrl
// 	}
// }

// func updateOpenAPIUrl(openapi *openapi3.T, serviceUrl string) {
// 	if serviceUrl != "" {
// 		openapi.Servers[0].URL = serviceUrl
// 	}
// }
