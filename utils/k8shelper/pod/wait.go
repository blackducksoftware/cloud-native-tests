// Borrowed from: https://github.com/kubernetes/kubernetes/blob/master/test/e2e/framework/pod/wait.go until upstream/kubernetes testing framework can be imported directly

package pod

import (
	"fmt"
	"time"

	k8sutils "github.com/blackducksoftware/cloud-native-tests/utils/k8shelper"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
)

// NewLogger returns a *logrus.Logger
func NewLogger() *logrus.Logger {
	newLog := logrus.New()
	newLog.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	newLog.SetLevel(logrus.DebugLevel)
	newLog.SetReportCaller(true)
	return newLog
}

const (
	// defaultPodDeletionTimeout is the default timeout for deleting pod.
	defaultPodDeletionTimeout = 3 * time.Minute

	// podListTimeout is how long to wait for the pod to be listable.
	podListTimeout = time.Minute

	podRespondingTimeout = 15 * time.Minute

	// How long pods have to become scheduled onto nodes
	podScheduledBeforeTimeout = podListTimeout + (20 * time.Second)

	// podStartTimeout is how long to wait for the pod to be started.
	// Initial pod start can be delayed O(minutes) by slow docker pulls.
	// TODO: Make this 30 seconds once #4566 is resolved.
	podStartTimeout = 5 * time.Minute

	// poll is how often to poll pods, nodes and claims.
	poll             = 2 * time.Second
	pollShortTimeout = 1 * time.Minute
	pollLongTimeout  = 5 * time.Minute

	// singleCallTimeout is how long to try single API calls (like 'get' or 'list'). Used to prevent
	// transient failures from failing tests.
	// TODO: client should not apply this timeout to Watch calls. Increased from 30s until that is fixed.
	singleCallTimeout = 5 * time.Minute

	// Some pods can take much longer to get ready due to volume attach/detach latency.
	slowPodStartTimeout = 15 * time.Minute
)

type podCondition func(pod *v1.Pod) (bool, error)

// WaitForPodsWithLabelRunningReady waits for exact amount of matching pods to become running and ready.
// Return the list of matching pods.
func WaitForPodsWithLabelRunningReady(c clientset.Interface, ns string, label labels.Selector, num int, timeout time.Duration) (pods *v1.PodList, err error) {
	var current int
	err = wait.Poll(poll, timeout,
		func() (bool, error) {
			pods, err = WaitForPodsWithLabel(c, ns, label)
			if err != nil {
				// e2elog.Logf("Failed to list pods: %v", err)
				fmt.Printf("Failed to list pods: %v", err)
				if k8sutils.IsRetryableAPIError(err) {
					return false, nil
				}
				return false, err
			}
			current = 0
			for _, pod := range pods.Items {
				if flag, err := PodRunningReady(&pod); err == nil && flag == true {
					current++
				}
			}
			if current != num {
				// e2elog.Logf("Got %v pods running and ready, expect: %v", current, num)
				fmt.Printf("Got %v pods running and ready, expect: %v", current, num)
				return false, nil
			}
			return true, nil
		})
	return pods, err
}

// WaitForPodsWithLabelDeleted waits up to podListTimeout for pods with certain label to not exist
// NOTE: Modified WaitForPodsWithLabel below
func WaitForPodsWithLabelDeleted(c clientset.Interface, ns string, label labels.Selector) (pods *v1.PodList, err error) {
	for t := time.Now(); time.Since(t) < podListTimeout; time.Sleep(poll) {
		options := metav1.ListOptions{LabelSelector: label.String()}
		pods, err = c.CoreV1().Pods(ns).List(options)
		if err != nil {
			if k8sutils.IsRetryableAPIError(err) {
				continue
			}
			return
		}
		if len(pods.Items) == 0 {
			break
		}
	}
	if pods == nil || len(pods.Items) > 0 {
		err = fmt.Errorf("Timeout while waiting for pods with label %v", label)
	}
	return
}

// WaitForPodsWithLabel waits up to podListTimeout for getting pods with certain label
func WaitForPodsWithLabel(c clientset.Interface, ns string, label labels.Selector) (pods *v1.PodList, err error) {
	for t := time.Now(); time.Since(t) < podListTimeout; time.Sleep(poll) {
		options := metav1.ListOptions{LabelSelector: label.String()}
		pods, err = c.CoreV1().Pods(ns).List(options)
		if err != nil {
			if k8sutils.IsRetryableAPIError(err) {
				continue
			}
			return
		}
		if len(pods.Items) > 0 {
			break
		}
	}
	if pods == nil || len(pods.Items) == 0 {
		err = fmt.Errorf("Timeout while waiting for pods with label %v", label)
	}
	return
}

// PodRunningReady checks whether pod p's phase is running and it has a ready
// condition of status true.
func PodRunningReady(p *v1.Pod) (bool, error) {
	// Check the phase is running.
	if p.Status.Phase != v1.PodRunning {
		return false, fmt.Errorf("want pod '%s' on '%s' to be '%v' but was '%v'",
			p.ObjectMeta.Name, p.Spec.NodeName, v1.PodRunning, p.Status.Phase)
	}
	// Check the ready condition is true.
	if !IsPodReady(p) {
		return false, fmt.Errorf("pod '%s' on '%s' didn't have condition {%v %v}; conditions: %v",
			p.ObjectMeta.Name, p.Spec.NodeName, v1.PodReady, v1.ConditionTrue, p.Status.Conditions)
	}
	return true, nil
}

// IsPodReady returns true if a pod is ready; false otherwise.
func IsPodReady(pod *v1.Pod) bool {
	return IsPodReadyConditionTrue(pod.Status)
}

// IsPodReadyConditionTrue returns true if a pod is ready; false otherwise.
func IsPodReadyConditionTrue(status v1.PodStatus) bool {
	condition := GetPodReadyCondition(status)
	return condition != nil && condition.Status == v1.ConditionTrue
}

// GetPodReadyCondition extracts the pod ready condition from the given status and returns that.
// Returns nil if the condition is not present.
func GetPodReadyCondition(status v1.PodStatus) *v1.PodCondition {
	_, condition := GetPodCondition(&status, v1.PodReady)
	return condition
}

// GetPodCondition extracts the provided condition from the given status and returns that.
// Returns nil and -1 if the condition is not present, and the index of the located condition.
func GetPodCondition(status *v1.PodStatus, conditionType v1.PodConditionType) (int, *v1.PodCondition) {
	if status == nil {
		return -1, nil
	}
	return GetPodConditionFromList(status.Conditions, conditionType)
}

// GetPodConditionFromList extracts the provided condition from the given list of condition and
// returns the index of the condition and the condition. Returns -1 and nil if the condition is not present.
func GetPodConditionFromList(conditions []v1.PodCondition, conditionType v1.PodConditionType) (int, *v1.PodCondition) {
	if conditions == nil {
		return -1, nil
	}
	for i := range conditions {
		if conditions[i].Type == conditionType {
			return i, &conditions[i]
		}
	}
	return -1, nil
}
