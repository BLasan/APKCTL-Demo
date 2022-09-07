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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BLasan/APKCTL-Demo/utils"
	"github.com/go-openapi/spec"
)

func CreateAPI(filePath, namespace, serviceUrl string) error {
	fmt.Println(filePath)
	_, content, err := resolveYamlOrJSON(filePath)
	// fmt.Println("Content: ", string(content))

	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	var swaggerSpec spec.Swagger
	err = json.Unmarshal(content, &swaggerSpec)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Swagger Spec: ", swaggerSpec)

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
