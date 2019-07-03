/*
Borrowed from: https://github.com/kubernetes/kubernetes/blob/master/test/e2e/framework/util.go until upstream/kubernetes testing framework can be imported directly

Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	k8sutils "github.com/blackducksoftware/cloud-native-tests/utils/k8shelper"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
)

// WaitForService waits until the service appears (exist == true), or disappears (exist == false)
func WaitForService(c clientset.Interface, namespace, name string, exist bool, interval, timeout time.Duration) error {
	err := wait.PollImmediate(interval, timeout, func() (bool, error) {
		_, err := c.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
		switch {
		case err == nil:
			// e2elog.Logf("Service %s in namespace %s found.", name, namespace)
			return exist, nil
		case apierrs.IsNotFound(err):
			// e2elog.Logf("Service %s in namespace %s disappeared.", name, namespace)
			return !exist, nil
		case !k8sutils.IsRetryableAPIError(err):
			// e2elog.Logf("Non-retryable failure while getting service.")
			return false, err
		default:
			// e2elog.Logf("Get service %s in namespace %s failed: %v", name, namespace, err)
			return false, nil
		}
	})
	if err != nil {
		stateMsg := map[bool]string{true: "to appear", false: "to disappear"}
		return fmt.Errorf("error waiting for service %s/%s %s: %v", namespace, name, stateMsg[exist], err)
	}
	return nil
}

// WaitForServiceWithSelector waits until any service with given selector appears (exist == true), or disappears (exist == false)
func WaitForServiceWithSelector(c clientset.Interface, namespace string, selector labels.Selector, exist bool, interval,
	timeout time.Duration) error {
	err := wait.PollImmediate(interval, timeout, func() (bool, error) {
		services, err := c.CoreV1().Services(namespace).List(metav1.ListOptions{LabelSelector: selector.String()})
		switch {
		case len(services.Items) != 0:
			// e2elog.Logf("Service with %s in namespace %s found.", selector.String(), namespace)
			return exist, nil
		case len(services.Items) == 0:
			// e2elog.Logf("Service with %s in namespace %s disappeared.", selector.String(), namespace)
			return !exist, nil
		case !k8sutils.IsRetryableAPIError(err):
			// e2elog.Logf("Non-retryable failure while listing service.")
			return false, err
		default:
			// e2elog.Logf("List service with %s in namespace %s failed: %v", selector.String(), namespace, err)
			// fmt.Logf("List service with %s in namespace %s failed: %v", selector.String(), namespace, err)
			return false, nil
		}
	})
	if err != nil {
		stateMsg := map[bool]string{true: "to appear", false: "to disappear"}
		return fmt.Errorf("error waiting for service with %s in namespace %s %s: %v", selector.String(), namespace, stateMsg[exist], err)
	}
	return nil
}
