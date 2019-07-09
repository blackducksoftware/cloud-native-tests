package synopsysctl_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	utils "github.com/blackducksoftware/cloud-native-tests/utils"
	k8sutils "github.com/blackducksoftware/cloud-native-tests/utils/k8shelper"
	crdutils "github.com/blackducksoftware/cloud-native-tests/utils/k8shelper/crd"
	namespaceutils "github.com/blackducksoftware/cloud-native-tests/utils/k8shelper/namespace"
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

		Context("deploying Synopsys Operator in cluster scope", func() {

			Specify("all crds can be enabled", func() {
				// BEGIN SETUP
				out, err := mySynopsysCtl.Exec("deploy", "--cluster-scoped", "--enable-alert", "--enable-blackduck", "--enable-opssight", "-i=gcr.io/saas-hub-stg/blackducksoftware/synopsys-operator:release-2019.6.x")
				if err != nil {
					Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
				}
				// END SETUP

				// BEGIN VERIFICATION
				label := labels.NewSelector()
				r, _ := labels.NewRequirement("app", selection.Equals, []string{"synopsys-operator"})
				label.Add(*r)
				_, err = podutils.WaitForPodsWithLabelRunningReady(kc, "synopsys-operator", label, 2, time.Duration(10*time.Second))
				if err != nil {
					Fail(fmt.Sprintf("Operator pods failed to come up: %v", err))
				}
				err = crdutils.BlockUntilCrdIsAdded(apiExtensionClient, "alerts.synopsys.com", 10)
				if err != nil {
					Fail(fmt.Sprintf("alert crd was not added: %v", err))
				}
				err = crdutils.BlockUntilCrdIsAdded(apiExtensionClient, "blackducks.synopsys.com", 10)
				if err != nil {
					Fail(fmt.Sprintf("black duck crd was not added: %v", err))
				}
				err = crdutils.BlockUntilCrdIsAdded(apiExtensionClient, "opssights.synopsys.com", 10)
				if err != nil {
					Fail(fmt.Sprintf("opssight crd was not added: %v", err))
				}
				// END VERIFICATION

				// BEGIN CLEANUP
				kc.CoreV1().Namespaces().Delete("synopsys-operator", &metav1.DeleteOptions{})
				kc.RbacV1().ClusterRoles().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
				kc.RbacV1().ClusterRoleBindings().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
				apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("alerts.synopsys.com", &metav1.DeleteOptions{})
				apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("blackducks.synopsys.com", &metav1.DeleteOptions{})
				apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("opssights.synopsys.com", &metav1.DeleteOptions{})
				//TODO:WAIT FOR NAMESPACE AND CRD TO BE DELETED
				err = namespaceutils.WaitForNamespacesDeleted(kc, []string{"synopsys-operator"}, time.Duration(30*time.Second))
				if err != nil {
					fmt.Printf("WAIT NS ERROR : %+v", err)
				}
				// END CLEANUP
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

		Context("deploying Synopsys Operator in namespace scope", func() {

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
					// BEGIN SETUP
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
					// END SETUP

					// BEGIN VERIFICATION
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
					// END VERIFICATION

					// BEGIN CLEANUP
					cleanErrs := []error{}
					err = kc.CoreV1().Namespaces().Delete("alert", &metav1.DeleteOptions{})
					cleanErrs = append(cleanErrs, err)
					err = apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("alerts.synopsys.com", &metav1.DeleteOptions{})
					cleanErrs = append(cleanErrs, err)
					err = kc.CoreV1().Namespaces().Delete("bd", &metav1.DeleteOptions{})
					cleanErrs = append(cleanErrs, err)
					err = apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("blackducks.synopsys.com", &metav1.DeleteOptions{})
					cleanErrs = append(cleanErrs, err)
					err = kc.CoreV1().Namespaces().Delete("alert-and-bd", &metav1.DeleteOptions{})
					cleanErrs = append(cleanErrs, err)
					err = kc.RbacV1().ClusterRoles().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
					cleanErrs = append(cleanErrs, err)
					err = kc.RbacV1().ClusterRoleBindings().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
					cleanErrs = append(cleanErrs, err)
					err = namespaceutils.WaitForNamespacesDeleted(kc, []string{"alert", "bd", "alert-and-bd"}, time.Duration(30*time.Second))
					cleanErrs = append(cleanErrs, err)
					// if len(cleanErrs) != 0 {
					// 	Fail(fmt.Sprintf("%+v", cleanErrs))
					// }
					// END CLEANUP
				})
				// create three namespaces
				// execute three deploy commands
				// go routine: wait for all operators to be running
				// destroy

			})
		})
	})

	Describe("destroy command", func() {
		Context("destroying Synopsys Operator in cluster scope", func() {
			Specify("all resources are removed", func() {
				// BEGIN SETUP
				out, err := mySynopsysCtl.Exec("deploy", "--cluster-scoped", "--enable-alert", "--enable-blackduck", "--enable-opssight", "-i=gcr.io/saas-hub-stg/blackducksoftware/synopsys-operator:release-2019.6.x")
				if err != nil {
					Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
				}
				// END SETUP

				// BEGIN VERIFICATION
				out, err = mySynopsysCtl.Exec("destroy")
				if err != nil {
					Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
				}
				label := labels.NewSelector()
				r, _ := labels.NewRequirement("app", selection.Equals, []string{"synopsys-operator"})
				label.Add(*r)
				_, err = podutils.WaitForPodsWithLabelDeleted(kc, "synopsys-operator", label)
				if err != nil {
					Fail(fmt.Sprintf("Synopsys Operator pods failed to stop running: %v", err))
				}
				// TODO : check that CRDs are removed
				// END VERIFICATION

				// BEGIN CLEANUP
				kc.CoreV1().Namespaces().Delete("synopsys-operator", &metav1.DeleteOptions{})
				kc.RbacV1().ClusterRoles().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
				kc.RbacV1().ClusterRoleBindings().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
				apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("alerts.synopsys.com", &metav1.DeleteOptions{})
				apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("blackducks.synopsys.com", &metav1.DeleteOptions{})
				apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("opssights.synopsys.com", &metav1.DeleteOptions{})
				//TODO:WAIT FOR NAMESPACE AND CRD TO BE DELETED
				namespaceutils.WaitForNamespacesDeleted(kc, []string{"synopsys-operator"}, time.Duration(30*time.Second))
				// END CLEANUP
			})
		})
		Context("destroying Synopsys Operator in namespace scope", func() {
			Specify("one instance can be destroyed", func() {
				// BEGIN SETUP
				kc.CoreV1().Namespaces().Create(&corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "so-one",
						Name:      "so-one",
					},
				})
				out, err := mySynopsysCtl.Exec("deploy", "--enable-alert", "--enable-blackduck", "--namespace=so-one", "-i=gcr.io/saas-hub-stg/blackducksoftware/synopsys-operator:release-2019.6.x")
				if err != nil {
					Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
				}
				// END SETUP

				// BEGIN VERIFICATION
				out, err = mySynopsysCtl.Exec("destroy", "so-one")
				if err != nil {
					Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
				}
				label := labels.NewSelector()
				r, _ := labels.NewRequirement("app", selection.Equals, []string{"synopsys-operator"})
				label.Add(*r)
				_, err = podutils.WaitForPodsWithLabelDeleted(kc, "so-one", label)
				if err != nil {
					Fail(fmt.Sprintf("Synopsys Operator pods failed to stop running: %v", err))
				}
				// TODO : check that CRDs are removed
				// END VERIFICATION

				// BEGIN CLEANUP
				kc.CoreV1().Namespaces().Delete("so-one", &metav1.DeleteOptions{})
				kc.RbacV1().ClusterRoles().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
				kc.RbacV1().ClusterRoleBindings().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
				apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("alerts.synopsys.com", &metav1.DeleteOptions{})
				apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("blackducks.synopsys.com", &metav1.DeleteOptions{})
				apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("opssights.synopsys.com", &metav1.DeleteOptions{})
				//TODO:WAIT FOR NAMESPACE AND CRD TO BE DELETED
				namespaceutils.WaitForNamespacesDeleted(kc, []string{"so-one"}, time.Duration(30*time.Second))
				// END CLEANUP
			})
			// Specify("multiple instances can be destroyed at once", func() {
			// 	// BEGIN SETUP
			// 	kc.CoreV1().Namespaces().Create(&corev1.Namespace{
			// 		ObjectMeta: metav1.ObjectMeta{
			// 			Namespace: "so-one",
			// 			Name:      "so-one",
			// 		},
			// 	})
			// 	out, err := mySynopsysCtl.Exec("deploy", "--enable-alert", "--enable-blackduck", "--namespace=so-one", "-i=gcr.io/saas-hub-stg/blackducksoftware/synopsys-operator:release-2019.6.x")
			// 	if err != nil {
			// 		Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
			// 	}
			// 	kc.CoreV1().Namespaces().Create(&corev1.Namespace{
			// 		ObjectMeta: metav1.ObjectMeta{
			// 			Namespace: "so-two",
			// 			Name:      "so-two",
			// 		},
			// 	})
			// 	out, err = mySynopsysCtl.Exec("deploy", "--enable-blackduck", "--enable-alert", "--namespace=so-two", "-i=gcr.io/saas-hub-stg/blackducksoftware/synopsys-operator:release-2019.6.x")
			// 	if err != nil {
			// 		Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
			// 	}
			// 	soLabel := labels.NewSelector()
			// 	r, _ := labels.NewRequirement("app", selection.Equals, []string{"synopsys-operator"})
			// 	soLabel.Add(*r)
			// 	_, err = podutils.WaitForPodsWithLabelRunningReady(kc, "so-one", soLabel, 2, time.Duration(10*time.Second))
			// 	if err != nil {
			// 		Fail(fmt.Sprintf("Synopsys Operator pods failed to come up: %v", err))
			// 	}
			// 	_, err = podutils.WaitForPodsWithLabelRunningReady(kc, "so-two", soLabel, 2, time.Duration(10*time.Second))
			// 	if err != nil {
			// 		Fail(fmt.Sprintf("Synopsys Operator pods failed to come up: %v", err))
			// 	}
			// 	err = crdutils.BlockUntilCrdIsAdded(apiExtensionClient, "alerts.synopsys.com", 10)
			// 	if err != nil {
			// 		Fail(fmt.Sprintf("Alert CRD was not added: %v", err))
			// 	}
			// 	err = crdutils.BlockUntilCrdIsAdded(apiExtensionClient, "blackducks.synopsys.com", 10)
			// 	if err != nil {
			// 		Fail(fmt.Sprintf("Black Duck CRD was not added: %v", err))
			// 	}
			// 	// END SETUP

			// 	// BEGIN VERIFICATION
			// 	out, err = mySynopsysCtl.Exec("destroy", "so-one", "so-two")
			// 	if err != nil {
			// 		Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
			// 	}
			// 	_, err = podutils.WaitForPodsWithLabelDeleted(kc, "so-one", soLabel)
			// 	if err != nil {
			// 		Fail(fmt.Sprintf("Synopsys Operator pods failed to stop running: %v", err))
			// 	}
			// 	// TODO : check that CRDs are removed
			// 	_, err = podutils.WaitForPodsWithLabelDeleted(kc, "so-two", soLabel)
			// 	if err != nil {
			// 		Fail(fmt.Sprintf("Synopsys Operator pods failed to stop running: %v", err))
			// 	}
			// 	// TODO : check that CRDs are removed
			// 	// END VERIFICATION

			// 	// BEGIN CLEANUP
			// 	kc.CoreV1().Namespaces().Delete("so-one", &metav1.DeleteOptions{})
			// 	kc.CoreV1().Namespaces().Delete("so-two", &metav1.DeleteOptions{})
			// 	kc.RbacV1().ClusterRoles().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
			// 	kc.RbacV1().ClusterRoleBindings().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
			// 	apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("alerts.synopsys.com", &metav1.DeleteOptions{})
			// 	apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("blackducks.synopsys.com", &metav1.DeleteOptions{})
			// 	apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("opssights.synopsys.com", &metav1.DeleteOptions{})
			// 	//TODO:WAIT FOR NAMESPACE AND CRD TO BE DELETED
			// 	namespaceutils.WaitForNamespacesDeleted(kc, []string{"so-one", "so-two"}, time.Duration(30*time.Second))
			// 	// END CLEANUP
			// })
			Specify("a Synopsys Operator instance can be forcefully destroyed", func() {
				// BEGIN SETUP
				// create the namespace
				kc.CoreV1().Namespaces().Create(&corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "so-one",
						Name:      "so-one",
					},
				})
				// deploy a Synopsys Operator instance
				out, err := mySynopsysCtl.Exec("deploy", "--enable-alert", "--namespace=so-one", "-i=gcr.io/saas-hub-stg/blackducksoftware/synopsys-operator:release-2019.6.x")
				if err != nil {
					Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
				}
				soLabel := labels.NewSelector()
				r, _ := labels.NewRequirement("app", selection.Equals, []string{"synopsys-operator"})
				soLabel.Add(*r)
				_, err = podutils.WaitForPodsWithLabelRunningReady(kc, "so-one", soLabel, 2, time.Duration(10*time.Second))
				if err != nil {
					Fail(fmt.Sprintf("Synopsys Operator pods failed to come up: %v", err))
				}
				err = crdutils.BlockUntilCrdIsAdded(apiExtensionClient, "alerts.synopsys.com", 10)
				if err != nil {
					Fail(fmt.Sprintf("alert crd was not added: %v", err))
				}
				// create an Alert instance
				out, err = mySynopsysCtl.Exec("create", "alert", "alt-one", "--namespace=so-one", "--standalone=false", "--persistent-storage=false")
				if err != nil {
					Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
				}
				altLabel := labels.NewSelector()
				r, _ = labels.NewRequirement("app", selection.Equals, []string{"alert"})
				altLabel.Add(*r)
				_, err = podutils.WaitForPodsWithLabelRunningReady(kc, "so-one", altLabel, 1, time.Duration(10*time.Second))
				if err != nil {
					Fail(fmt.Sprintf("Synopsys Operator pods failed to come up: %v", err))
				}
				// END SETUP

				// BEGIN VERIFICATION
				out, err = mySynopsysCtl.Exec("destroy", "so-one")
				if err != nil {
					Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
				}
				// TODO : Check that instance isn't destroyed
				out, err = mySynopsysCtl.Exec("destroy", "so-one", "--force")
				if err != nil {
					Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
				}
				_, err = podutils.WaitForPodsWithLabelDeleted(kc, "so-one", soLabel)
				if err != nil {
					Fail(fmt.Sprintf("Synopsys Operator pods failed to stop running: %v", err))
				}
				// TODO : check that CRDs are removed
				// END VERIFICATION

				// BEGIN CLEANUP
				kc.CoreV1().Namespaces().Delete("so-one", &metav1.DeleteOptions{})
				kc.RbacV1().ClusterRoles().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
				kc.RbacV1().ClusterRoleBindings().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
				apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("alerts.synopsys.com", &metav1.DeleteOptions{})
				apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("blackducks.synopsys.com", &metav1.DeleteOptions{})
				apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("opssights.synopsys.com", &metav1.DeleteOptions{})
				// TODO : wait for namesapce and crd to be deleted
				// TODO : clean up the Alert instance
				namespaceutils.WaitForNamespacesDeleted(kc, []string{"so-one"}, time.Duration(30*time.Second))
				// END CLEANUP
			})
		})
	})

})
