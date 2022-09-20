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
	"strconv"
	"strings"

	k8sUtils "github.com/BLasan/APKCTL-Demo/k8s"
	"github.com/BLasan/APKCTL-Demo/utils"
	"github.com/go-openapi/spec"
	"gopkg.in/yaml.v2"
)

func CreateAPI(filePath, namespace, serviceUrl, apiName, version string, isDryRun bool) error {
	fmt.Println(filePath)
	_, content, err := resolveYamlOrJSON(filePath)
	// fmt.Println("Content: ", string(content))

	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	var swaggerSpec spec.Swagger

	if content != nil {
		err = json.Unmarshal(content, &swaggerSpec)

		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	// fmt.Println("Swagger Spec: ", swaggerSpec.Paths.Paths)

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

	fmt.Println(version)

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

	var serviceUrlArr []string

	if serviceUrl != "" {
		serviceUrlArr = strings.Split(serviceUrl, ".")
	} else if swaggerSpec.Host != "" {
		serviceUrlArr = strings.Split(swaggerSpec.Host, ".")
	}

	fmt.Println(serviceUrlArr)

	// if swagger path is not defined do not iterate over it
	if filePath == "" {
		apiPath.Type = utils.PathPrefix
		apiPath.Value = "/"
		match.Path = apiPath
		rule.Matches = append(rule.Matches, match)
	} else {
		counter := 1

		for path, pathItem := range swaggerSpec.Paths.Paths {
			// maximum 8 paths are allowed
			if counter > 8 {
				break
			}

			index := strings.IndexAny(path, "{")
			if index >= 0 {
				path = path[:index-1]
			}

			// append "/api/v3" to invoke the petstore apis
			path = "/api/v3" + path

			fmt.Println("Path: ", path)

			// pathArr := strings.Split(path, "/")
			// sort.Strings(pathArr)
			// path = utils.FindPathParam(pathArr)

			apiPath.Type = utils.PathPrefix
			apiPath.Value = path
			match.Path = apiPath

			rule.Matches = append(rule.Matches, match)

			if pathItem.Post != nil {
				fmt.Println("Description Items: ", pathItem.Post.Description)
			}

			counter++

		}
	}

	backendRef.Kind = utils.ServiceKind

	backendRef.Name = serviceUrlArr[0]
	// backendRef.Namespace = serviceUrlArr[1]
	backendRef.Port, err = strconv.Atoi(strings.Split(serviceUrlArr[len(serviceUrlArr)-1], ":")[1])

	rule.BackendRefs = append(rule.BackendRefs, backendRef)
	httpRoute.HttpRouteSpec.Rules = append(httpRoute.HttpRouteSpec.Rules, rule)

	if err != nil {
		return err
	}

	file, err := yaml.Marshal(&httpRoute)

	if err != nil {
		return err
	}

	dirPath := path.Join(utils.GetAPKCTLHomeDir(), utils.APIProjectsDir, apiName)

	os.MkdirAll(dirPath, os.ModePerm)

	desFilePath := path.Join(dirPath, "HTTPRouteConfig.yaml")

	// directory location can be defined in the apkctl config file
	err = ioutil.WriteFile(desFilePath, file, 0644)

	if err != nil {
		return err
	}

	args := []string{"apply", "-f", desFilePath}

	if !isDryRun {
		k8sUtils.ExecuteCommand("kubectl", args...)
	}

	return nil

}

func resolveYamlOrJSON(filename string) (string, []byte, error) {
	// lookup for yaml
	yamlFp := filename
	if info, err := os.Stat(yamlFp); err == nil && !info.IsDir() {
		// utils.Logln(utils.LogPrefixInfo+"Loading", yamlFp)
		// read it
		fn := yamlFp
		yamlContent, err := ioutil.ReadFile(fn)
		if err != nil {
			return "", nil, err
		}
		// load it as yaml
		jsonContent, err := utils.YamlToJson(yamlContent)
		if err != nil {
			return "", nil, err
		}
		return fn, jsonContent, nil
	}

	jsonFp := filename + ".json"
	if info, err := os.Stat(jsonFp); err == nil && !info.IsDir() {
		// utils.Logln(utils.LogPrefixInfo+"Loading", jsonFp)
		// read it
		fn := jsonFp
		jsonContent, err := ioutil.ReadFile(fn)
		if err != nil {
			return "", nil, err
		}
		return fn, jsonContent, nil
	}

	return "", nil, fmt.Errorf("%s was not found as a YAML or JSON", filename)
}
