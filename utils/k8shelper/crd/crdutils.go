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

package crd

import (
	"fmt"

	k8sutils "github.com/blackducksoftware/cloud-native-tests/utils/k8shelper"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// WatchCustomResourceDefinition watches the custom resource defintion
func WatchCustomResourceDefinition(apiExtensionClient *apiextensionsclient.Clientset, name string, timeout int) (watch.Interface, error) {
	return apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Watch(metav1.ListOptions{TimeoutSeconds: k8sutils.IntToInt64Ptr(timeout)})
}

// BlockUntilWatchEventReceived blocks until first event is received, and sees if it matches the wanted eventType
func BlockUntilWatchEventReceived(watchInterface watch.Interface, eventType watch.EventType) error {
	defer watchInterface.Stop()
	watchEvent := <-watchInterface.ResultChan()
	if watchEvent.Type != eventType {
		return fmt.Errorf("The event from the watch did not match the wanted eventType: %v", eventType)
	}
	return nil
}

// BlockUntilCrdIsAdded blocks until first event is received, and sees if it matches the wanted eventType
func BlockUntilCrdIsAdded(apiExtensionClient *apiextensionsclient.Clientset, name string, timeout int) error {
	watchInterface, err := WatchCustomResourceDefinition(apiExtensionClient, name, timeout)
	if err != nil {
		return fmt.Errorf("something went wrong in the watch event: %v", err)
	}
	err = BlockUntilWatchEventReceived(watchInterface, watch.Added)
	if err != nil {
		return fmt.Errorf("%v crd was not added: %v", name, err)
	}
	return nil
}
