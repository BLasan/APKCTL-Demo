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

package utils

type HTTPRouteConfig struct {
	ApiVersion    string        `yaml:"apiVersion"`
	Kind          string        `yaml:"kind"`
	MetaData      MetaData      `yaml:"metadata"`
	HttpRouteSpec HttpRouteSpec `yaml:"spec"`
}

type MetaData struct {
	Name      string            `yaml:"name,omitempty"`
	Namespace string            `yaml:"namespace,omitempty"`
	Labels    map[string]string `yaml:"labels,omitempty"`
}

type HttpRouteSpec struct {
	ParentRefs []ParentRef `yaml:"parentRefs"`
	HostNames  []string    `yaml:"hostnames"`
	Rules      []Rule      `yaml:"rules"`
}

type ParentRef struct {
	Name string `yaml:"name"`
}

type Rule struct {
	Matches     []Match      `yaml:"matches"`
	BackendRefs []BackendRef `yaml:"backendRefs,omitempty"`
}

type Match struct {
	Path        Path         `yaml:"path,omitempty"`
	Headers     []Header     `yaml:"headers,omitempty"`
	QueryParams []QueryParam `yaml:"queryParams,omitempty"`
	// BackendRefs BackendRef `yaml:"backendRefs,omitempty"`
}

type QueryParam struct {
	Type  string `yaml:"type"`
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Path struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

type Header struct {
	Type  string `yaml:"type"`
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type BackendRef struct {
	Name   string `yaml:"name"`
	Port   int    `yaml:"port"`
	Group  string `yaml:"group"`
	Kind   string `yaml:"kind,omitempty"`
	Weight string `yaml:"weight,omitempty"`
}

//  apiVersion: gateway.networking.k8s.io/v1beta1
// kind: HTTPRoute
// metadata:
//   name: bar-route
// spec:
//   parentRefs:
//   - name: example-gateway
//   hostnames:
//   - "bar.example.com"
//   rules:
//   - matches:
//     - headers:
//       - type: Exact
//         name: env
//         value: canary
//     backendRefs:
//     - name: bar-svc-canary
//       port: 8080
//   - backendRefs:
//     - name: bar-svc
//       port: 8080
