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

package pvc

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

// WaitForPersistentVolumeDeleted waits for a PersistentVolume to get deleted or until timeout occurs, whichever comes first.
func WaitForPersistentVolumeDeleted(c clientset.Interface, pvName string, Poll, timeout time.Duration) error {
	// e2elog.Logf("Waiting up to %v for PersistentVolume %s to get deleted", timeout, pvName)
	fmt.Printf("Waiting up to %v for PersistentVolume %s to get deleted", timeout, pvName)
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(Poll) {
		pv, err := c.CoreV1().PersistentVolumes().Get(pvName, metav1.GetOptions{})
		if err == nil {
			// e2elog.Logf("PersistentVolume %s found and phase=%s (%v)", pvName, pv.Status.Phase, time.Since(start))
			fmt.Printf("PersistentVolume %s found and phase=%s (%v)", pvName, pv.Status.Phase, time.Since(start))
			continue
		}
		if apierrs.IsNotFound(err) {
			// e2elog.Logf("PersistentVolume %s was removed", pvName)
			fmt.Printf("PersistentVolume %s was removed", pvName)
			return nil
		}
		// e2elog.Logf("Get persistent volume %s in failed, ignoring for %v: %v", pvName, Poll, err)
		fmt.Printf("Get persistent volume %s in failed, ignoring for %v: %v", pvName, Poll, err)
	}
	return fmt.Errorf("PersistentVolume %s still exists within %v", pvName, timeout)
}

// WaitForPersistentVolumeClaimDeleted waits for a PersistentVolumeClaim to be removed from the system until timeout occurs, whichever comes first.
func WaitForPersistentVolumeClaimDeleted(c clientset.Interface, ns string, pvcName string, Poll, timeout time.Duration) error {
	// e2elog.Logf("Waiting up to %v for PersistentVolumeClaim %s to be removed", timeout, pvcName)
	fmt.Printf("Waiting up to %v for PersistentVolumeClaim %s to be removed", timeout, pvcName)
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(Poll) {
		_, err := c.CoreV1().PersistentVolumeClaims(ns).Get(pvcName, metav1.GetOptions{})
		if err != nil {
			if apierrs.IsNotFound(err) {
				// e2elog.Logf("Claim %q in namespace %q doesn't exist in the system", pvcName, ns)
				fmt.Printf("Claim %q in namespace %q doesn't exist in the system", pvcName, ns)
				return nil
			}
			// e2elog.Logf("Failed to get claim %q in namespace %q, retrying in %v. Error: %v", pvcName, ns, Poll, err)
			fmt.Printf("Failed to get claim %q in namespace %q, retrying in %v. Error: %v", pvcName, ns, Poll, err)
		}
	}
	return fmt.Errorf("PersistentVolumeClaim %s is not removed from the system within %v", pvcName, timeout)
}

// WaitForPersistentVolumeClaimPhase waits for a PersistentVolumeClaim to be in a specific phase or until timeout occurs, whichever comes first.
func WaitForPersistentVolumeClaimPhase(phase v1.PersistentVolumeClaimPhase, c clientset.Interface, ns string, pvcName string, Poll, timeout time.Duration) error {
	return WaitForPersistentVolumeClaimsPhase(phase, c, ns, []string{pvcName}, Poll, timeout, true)
}

// WaitForPersistentVolumeClaimsPhase waits for any (if matchAny is true) or all (if matchAny is false) PersistentVolumeClaims
// to be in a specific phase or until timeout occurs, whichever comes first.
func WaitForPersistentVolumeClaimsPhase(phase v1.PersistentVolumeClaimPhase, c clientset.Interface, ns string, pvcNames []string, Poll, timeout time.Duration, matchAny bool) error {
	if len(pvcNames) == 0 {
		return fmt.Errorf("Incorrect parameter: Need at least one PVC to track. Found 0")
	}
	// e2elog.Logf("Waiting up to %v for PersistentVolumeClaims %v to have phase %s", timeout, pvcNames, phase)
	fmt.Printf("Waiting up to %v for PersistentVolumeClaims %v to have phase %s", timeout, pvcNames, phase)
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(Poll) {
		phaseFoundInAllClaims := true
		for _, pvcName := range pvcNames {
			pvc, err := c.CoreV1().PersistentVolumeClaims(ns).Get(pvcName, metav1.GetOptions{})
			if err != nil {
				// e2elog.Logf("Failed to get claim %q, retrying in %v. Error: %v", pvcName, Poll, err)
				fmt.Printf("Failed to get claim %q, retrying in %v. Error: %v", pvcName, Poll, err)
				continue
			}
			if pvc.Status.Phase == phase {
				// e2elog.Logf("PersistentVolumeClaim %s found and phase=%s (%v)", pvcName, phase, time.Since(start))
				fmt.Printf("PersistentVolumeClaim %s found and phase=%s (%v)", pvcName, phase, time.Since(start))
				if matchAny {
					return nil
				}
			} else {
				// e2elog.Logf("PersistentVolumeClaim %s found but phase is %s instead of %s.", pvcName, pvc.Status.Phase, phase)
				fmt.Printf("PersistentVolumeClaim %s found but phase is %s instead of %s.", pvcName, pvc.Status.Phase, phase)
				phaseFoundInAllClaims = false
			}
		}
		if phaseFoundInAllClaims {
			return nil
		}
	}
	return fmt.Errorf("PersistentVolumeClaims %v not all in phase %s within %v", pvcNames, phase, timeout)
}
