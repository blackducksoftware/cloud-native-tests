package synopsysctl_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	utils "github.com/blackducksoftware/cloud-native-tests/utils"
	k8sutils "github.com/blackducksoftware/cloud-native-tests/utils/k8shelper"
	crdutils "github.com/blackducksoftware/cloud-native-tests/utils/k8shelper/crd"
	podutils "github.com/blackducksoftware/cloud-native-tests/utils/k8shelper/pod"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func TestGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ginkgo Suite")
}

var _ = Describe("synopsysctl", func() {

	defer GinkgoRecover()

	rc, _ := k8sutils.GetRestConfig()
	kc, _ := k8sutils.GetKubeClient(rc)

	// get Crd
	apiExtensionClient, err := apiextensionsclient.NewForConfig(rc)
	if err != nil {
		Fail(fmt.Sprintf("error creating the api extension client: %+v", err))
	}
	mySynopsysCtl := utils.NewSynopsysctl("synopsysctl")

	Describe("--version command", func() {
		Context("--version", func() {
			Specify("the version is 2019.6.0", func() {
				out, err := mySynopsysCtl.Exec("--version")
				if err != nil {
					Fail(fmt.Sprintf("%s", err))
				}
				Expect(strings.TrimSpace(out)).To(Equal("synopsysctl version 2019.6.0"))
			})
		})
	})

	Describe("deploy command", func() {

		Context("deploying the operator in cluster scope", func() {

			Specify("all crds can be enabled", func() {
				// Setup
				out, err := mySynopsysCtl.Exec("deploy", "--cluster-scoped", "--enable-alert", "--enable-blackduck", "--enable-opssight", "-i=gcr.io/saas-hub-stg/blackducksoftware/synopsys-operator:release-2019.6.x")
				if err != nil {
					Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
				}
				fmt.Printf("[DEBUG]: inside before suite, ran deploy command")
				// End Setup

				// Begin Verification
				fmt.Printf("[DEBUG]: inside of checking for pods")
				label := labels.NewSelector()
				r, _ := labels.NewRequirement("app", selection.Equals, []string{"synopsys-operator"})
				label.Add(*r)
				_, err = podutils.WaitForPodsWithLabelRunningReady(kc, "synopsys-operator", label, 2, time.Duration(10*time.Second))
				if err != nil {
					Fail(fmt.Sprintf("Operator pods failed to come up: %v", err))
				}
				fmt.Printf("[DEBUG]: done checking for pods")

				fmt.Printf("[DEBUG]: checking for alert crd")
				err = crdutils.BlockUntilCrdIsAdded(apiExtensionClient, "alerts.synopsys.com", 10)
				if err != nil {
					Fail(fmt.Sprintf("alert crd was not added: %v", err))
				}

				fmt.Printf("[DEBUG]: checking for bd crd")
				err = crdutils.BlockUntilCrdIsAdded(apiExtensionClient, "blackducks.synopsys.com", 10)
				if err != nil {
					Fail(fmt.Sprintf("black duck crd was not added: %v", err))
				}

				fmt.Printf("[DEBUG]: checking for ops crd")
				err = crdutils.BlockUntilCrdIsAdded(apiExtensionClient, "opssights.synopsys.com", 10)
				if err != nil {
					Fail(fmt.Sprintf("opssight crd was not added: %v", err))
				}

				// Begin cleanup
				fmt.Printf("[DEBUG]: cleaning up")
				kc.CoreV1().Namespaces().Delete("synopsys-operator", &metav1.DeleteOptions{})
				kc.RbacV1().ClusterRoles().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
				kc.RbacV1().ClusterRoleBindings().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
				apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("alerts.synopsys.com", &metav1.DeleteOptions{})
				apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("blackducks.synopsys.com", &metav1.DeleteOptions{})
				apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("opssights.synopsys.com", &metav1.DeleteOptions{})
				//TODO:WAIT FOR NAMESPACE AND CRD TO BE DELETED
				time.Sleep(10 * time.Second)
				fmt.Printf("[DEBUG]: finished waiting for cluster cleanup")

			})

			Specify("just alert crd can be enabled", func() {
				//TODO:
			})

			Specify("just black duck crd can be enabled", func() {
				//TODO:
			})

			Specify("just black duck crd can be enabled", func() {
				//TODO:
			})

			Specify("cannot deploy two operators in cluster scope", func() {

			})

		})

		Context("deploying the operator in namespace scope", func() {

			Context("one namespace", func() {

				Specify("it can deploy with only one resource enabled", func() {
					//TODO: create namespace
					//TODO: deploy operator in that namespace with only alert enabled
					//TODO: verify only alert crd is there
					//TODO: cleanup

					//TODO: create namespace
					//TODO: deploy operator in that namespace with only blackduck enabled
					//TODO: verify only blackduck crd is there
					//TODO: cleanup
				})

				Specify("it can deploy with two resources enabled", func() {
					//TODO: create namespace
					//TODO: deploy operator in that namespace with both blackduck and alert enabled
					//TODO: verify both alert and blackduck crd is there
					//TODO: cleanup
				})

			})

			Context("multiple namespaces at the same time", func() {

				Specify("one operator is deployed in 'alert' namespace with alert crd enabled, one operator is deployed in 'bd' namespace with blackduck crd enabled, and one operator is deployed in 'alert-and-bd' namespace with both alert and blackduck crd enabled", func() {
					// Setup
					kc.CoreV1().Namespaces().Create(&corev1.Namespace{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "alert",
							Name:      "alert",
						},
					})
					out, err := mySynopsysCtl.Exec("deploy", "-n=alert", "--enable-alert", "-i=gcr.io/saas-hub-stg/blackducksoftware/synopsys-operator:release-2019.6.x")
					if err != nil {
						Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
					}

					kc.CoreV1().Namespaces().Create(&corev1.Namespace{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "bd",
							Name:      "bd",
						},
					})
					out, err = mySynopsysCtl.Exec("deploy", "-n=bd", "--enable-blackduck", "-i=gcr.io/saas-hub-stg/blackducksoftware/synopsys-operator:release-2019.6.x")
					if err != nil {
						Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
					}

					kc.CoreV1().Namespaces().Create(&corev1.Namespace{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "alert-and-bd",
							Name:      "alert-and-bd",
						},
					})
					out, err = mySynopsysCtl.Exec("deploy", "-n=alert-and-bd", "--enable-alert", "--enable-blackduck", "-i=gcr.io/saas-hub-stg/blackducksoftware/synopsys-operator:release-2019.6.x")
					if err != nil {
						Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
					}
					// End Setup

					label := labels.NewSelector()
					r, _ := labels.NewRequirement("app", selection.Equals, []string{"synopsys-operator"})
					label.Add(*r)
					_, err = podutils.WaitForPodsWithLabelRunningReady(kc, "alert", label, 2, time.Duration(30*time.Second))
					if err != nil {
						Fail(fmt.Sprintf("Operator pods failed to come up: %v", err))
					}
					// TODO: make sure the crd has scope: Namespaced
					err = crdutils.BlockUntilCrdIsAdded(apiExtensionClient, "alerts.synopsys.com", 10)
					if err != nil {
						Fail(fmt.Sprintf("alert crd was not added: %v", err))
					}

					label = labels.NewSelector()
					r, _ = labels.NewRequirement("app", selection.Equals, []string{"synopsys-operator"})
					label.Add(*r)
					_, err = podutils.WaitForPodsWithLabelRunningReady(kc, "bd", label, 2, time.Duration(30*time.Second))
					if err != nil {
						Fail(fmt.Sprintf("Operator pods failed to come up: %v", err))
					}
					// TODO: make sure the crd has scope: Namespaced
					err = crdutils.BlockUntilCrdIsAdded(apiExtensionClient, "blackducks.synopsys.com", 10)
					if err != nil {
						Fail(fmt.Sprintf("black duck crd was not added: %v", err))
					}

					label = labels.NewSelector()
					r, _ = labels.NewRequirement("app", selection.Equals, []string{"synopsys-operator"})
					label.Add(*r)
					_, err = podutils.WaitForPodsWithLabelRunningReady(kc, "alert-and-bd", label, 2, time.Duration(30*time.Second))
					if err != nil {
						Fail(fmt.Sprintf("Operator pods failed to come up: %v", err))
					}

					// Cleanup
					kc.CoreV1().Namespaces().Delete("alert", &metav1.DeleteOptions{})
					apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("alerts.synopsys.com", &metav1.DeleteOptions{})
					kc.CoreV1().Namespaces().Delete("bd", &metav1.DeleteOptions{})
					apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("blackducks.synopsys.com", &metav1.DeleteOptions{})
					kc.CoreV1().Namespaces().Delete("alert-and-bd", &metav1.DeleteOptions{})

				})
				// create three namespaces
				// execute three deploy commands
				// go routine: wait for all operators to be running
				// destroy

			})
		})
	})
})
