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

const ProjectName = "apkctl"

// File Names and Paths
const HttpRouteApiVersion = "gateway.networking.k8s.io/v1beta1"
const HttpRouteKind = "HTTPRoute"
const PathPrefix = "PathPrefix"
const ServiceKind = "Service"
const APIProjectsDir = "/target/apis/"

const DefaultNamespace = "default"

// Constants for get APIs command
const APIColumnsOutput = "NAME:.metadata.name,VERSION:.metadata.labels.version,HOSTNAMES:.spec.hostnames"
const K8sOutputWithCustomColumns = "custom-columns=" + APIColumnsOutput
