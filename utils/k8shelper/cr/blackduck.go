/*
Copyright (C) 2019 Synopsys, Inc.

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package cr

import (
	"fmt"

	k8sutils "github.com/blackducksoftware/cloud-native-tests/utils/k8shelper"
	_ "k8s.io/client-go/plugin/pkg/client/auth" //for auths
	"k8s.io/client-go/rest"
)

// BlackDuckCRExists returns true if the CR exists in the cluster
func BlackDuckCRExists(restcli rest.Interface, blackDuckNamespace string, blackDuckName string) (bool, error) {
	a, err := GetBlackDuckCRInNamespace(restcli, blackDuckNamespace, blackDuckName)
	if err != nil {
		return false, err
	}
	if a.Kind == "Blackduck" {
		return true, nil
	}
	return false, nil
}

// GetBlackDuckCRInNamespace ...
func GetBlackDuckCRInNamespace(restcli rest.Interface, blackDuckNamespace string, blackDuckName string) (*k8sutils.APIResponse, error) {
	a := &k8sutils.APIResponse{}
	request := fmt.Sprintf("/apis/%s/v1/namespaces/%s/%s/%s", "synopsys.com", blackDuckNamespace, "blackducks", blackDuckName)
	fmt.Printf("Request: %s\n", request)
	err := k8sutils.GetResponseFromK8sEndpoint(restcli, request, a)
	if err != nil {
		return nil, err
	}
	return a, nil
}
