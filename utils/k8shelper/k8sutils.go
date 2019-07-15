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

package k8shelper

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth" //for auths
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

// GetKubeConfig  will return the kube config
func GetKubeConfig(kubeconfigpath string, insecureSkipTLSVerify bool) (*rest.Config, error) {
	log.Debugf("Getting Kube Rest Config")
	var err error
	var kubeConfig *rest.Config
	// creates the in-cluster config
	kubeConfig, err = rest.InClusterConfig()
	if err != nil {
		log.Debugf("using native config due to %+v", err)
		// Determine Config Paths
		if home := homeDir(); len(kubeconfigpath) == 0 && home != "" {
			kubeconfigpath = filepath.Join(home, ".kube", "config")
		}

		kubeConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{
				ExplicitPath: kubeconfigpath,
			},
			&clientcmd.ConfigOverrides{
				ClusterInfo: clientcmdapi.Cluster{
					Server:                "",
					InsecureSkipTLSVerify: insecureSkipTLSVerify,
				},
			}).ClientConfig()
		if err != nil {
			return nil, err
		}
	}
	return kubeConfig, nil
}

// homeDir determines the user's home directory path
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// GetKubeClientSet will return the kube clientset
func GetKubeClientSet(kubeConfig *rest.Config) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(kubeConfig)
}

func newKubeClientFromOutsideCluster() (*rest.Config, error) {
	var kubeConfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	return config, errors.Annotate(err, "error creating default client config")
}

// GetRestConfig calls protoform.GetKubeConfig
func GetRestConfig() (*rest.Config, error) {
	return GetKubeConfig("", false)
}

// GetKubeClient gets the kubernetes client
func GetKubeClient(kubeConfig *rest.Config) (*kubernetes.Clientset, error) {
	client, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	// logger.Debugf("client: %v \n\n", client)
	return client, nil
}

// GetDynamicClient gets a dynamic client
func GetDynamicClient(config *rest.Config) (dynamic.Interface, error) {
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// IntToInt64Ptr takes an int and returns a *int64
func IntToInt64Ptr(i int) *int64 {
	j := int64(i)
	return &j
}

// IsRetryableAPIError CHANGE
func IsRetryableAPIError(err error) bool {
	// These errors may indicate a transient error that we can retry in tests.
	if apierrs.IsInternalError(err) || apierrs.IsTimeout(err) || apierrs.IsServerTimeout(err) ||
		apierrs.IsTooManyRequests(err) || utilnet.IsProbableEOF(err) || utilnet.IsConnectionReset(err) {
		return true
	}
	// If the error sends the Retry-After header, we respect it as an explicit confirmation we should retry.
	if _, shouldRetry := apierrs.SuggestsClientDelay(err); shouldRetry {
		return true
	}
	return false
}

// APIResponse ...
type APIResponse struct {
	APIVersion string      `json:"apiVersion"`
	Kind       string      `json:"kind"`
	Spec       interface{} `json:"spec,omitempty"`
}

/*
GetResponseFromK8sEndpoint puts the data into the struct that is pointed to by unmarshal
unmarshal - pointer to a struct
*/
func GetResponseFromK8sEndpoint(restcli rest.Interface, requesturi string, unmarshal interface{}) error {
	b, err := restcli.Get().RequestURI(requesturi).DoRaw()
	restcli.Get().Stream()
	if err != nil {
		return err
	}
	fmt.Printf("Response: %s\n", string(b))
	if err := json.Unmarshal(b, unmarshal); err != nil {
		return err
	}
	return nil
}
