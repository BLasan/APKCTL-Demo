/*
 *  Copyright (c) 2021, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

// Package model contains the implementation of DTOs to convert OpenAPI/Swagger files
// and create a common model which can represent both types.
package oasparser

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/go-openapi/spec"
)

// // NewOperation Creates and returns operation type object
// func NewOperation(method string, security []map[string][]string, extensions map[string]interface{}) *Operation {
// 	tier := ResolveThrottlingTier(extensions)
// 	disableSecurity := ResolveDisableSecurity(extensions)
// 	id := uuid.New().String()
// 	return &Operation{id, method, security, tier, disableSecurity, extensions, OperationPolicies{}, &api.MockedApiConfig{}}
// }

// SetMockedAPIConfigOAS2 generate mock impl endpoint configurations
func (operation *Operation) SetMockedAPIConfigOAS2(openAPIOperation *spec.Operation) {
	if openAPIOperation.Responses != nil && len(openAPIOperation.Responses.StatusCodeResponses) > 0 {
		mockedAPIConfig := &api.MockedApiConfig{
			Responses: make([]*api.MockedResponseConfig, 0),
		}
		// get response codes
		for responseCode, responseRef := range openAPIOperation.Responses.StatusCodeResponses {
			mockedResponse := &api.MockedResponseConfig{
				Code:    strconv.Itoa(responseCode),
				Content: make([]*api.MockedContentConfig, 0),
			}
			for mediaType, content := range responseRef.ResponseProps.Examples {
				//todo(amali) xml payload gen
				example, err := convertToJSON(content)
				if err == nil {
					mockedResponse.Content = append(mockedResponse.Content, &api.MockedContentConfig{
						ContentType: mediaType,
						Examples:    []*api.MockedContentExample{{Ref: "", Body: example}},
					})
				}
			}
			// swagger does not support header example/examples
			if len(mockedResponse.Content) > 0 {
				mockedAPIConfig.Responses = append(mockedAPIConfig.Responses, mockedResponse)
			}
		}
		// get default response examples
		if openAPIOperation.Responses.Default != nil && len(openAPIOperation.Responses.Default.Examples) > 0 {
			mockedResponse := &api.MockedResponseConfig{
				Code:    "default",
				Content: make([]*api.MockedContentConfig, 0),
			}
			for mediaType, content := range openAPIOperation.Responses.Default.Examples {
				example, err := convertToJSON(content)
				if err == nil {
					mockedResponse.Content = append(mockedResponse.Content, &api.MockedContentConfig{
						ContentType: mediaType,
						Examples:    []*api.MockedContentExample{{Ref: "", Body: example}},
					})
				}
			}
			// swagger does not support header example/examples
			if len(mockedResponse.Content) > 0 {
				mockedAPIConfig.Responses = append(mockedAPIConfig.Responses, mockedResponse)
			}
		}
		// if len(mockedAPIConfig.Responses) > 0 {
		// 	operation.mockedAPIConfig = mockedAPIConfig
		// }
	}
}

// convertToJSON parse interface to JSON string. returns error if a null value has passed
func convertToJSON(data interface{}) (string, error) {
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	return "", errors.New("null object passed")
}

// ResolveThrottlingTier extracts the value of x-wso2-throttling-tier and
// x-throttling-tier extension. if x-wso2-throttling-tier is available it
// will be prioritized.
// if both the properties are not available, an empty string is returned.
func ResolveThrottlingTier(vendorExtensions map[string]interface{}) string {
	xTier := ""
	if x, found := vendorExtensions[constants.XWso2ThrottlingTier]; found {
		if val, ok := x.(string); ok {
			xTier = val
		}
	} else if y, found := vendorExtensions[constants.XThrottlingTier]; found {
		if val, ok := y.(string); ok {
			xTier = val
		}
	}
	return xTier
}
