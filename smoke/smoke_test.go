package smoke_test

import (
	"fmt"
	"testing"
	"time"

	utils "github.com/blackducksoftware/cloud-native-tests/utils"
	k8sutils "github.com/blackducksoftware/cloud-native-tests/utils/k8shelper"
	crutils "github.com/blackducksoftware/cloud-native-tests/utils/k8shelper/cr"
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
	fmt.Printf("[DEBUG] TestGinkgo\n")
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cloud Native Suite Smoke Tests")
	fmt.Printf("[DEBUG] After RunSpecs\n")
}

var _ = Describe("smoke", func() {
	fmt.Printf("[DEBUG] smoke\n")

	defer GinkgoRecover()

	rc, _ := k8sutils.GetRestConfig()
	kc, _ := k8sutils.GetKubeClient(rc)

	// get Crd
	apiExtensionClient, err := apiextensionsclient.NewForConfig(rc)
	if err != nil {
		Fail(fmt.Sprintf("error creating the api extension client: %+v", err))
	}
	mySynopsysCtl := utils.NewSynopsysctl("synopsysctl")

	Context("tests", func() {
		fmt.Printf("[DEBUG] tests\n")

		Specify("cluster scoped operations", func() {
			fmt.Printf("[DEBUG] cluster scoped operations\n")
			// Deploy Operator in Cluster Scope
			out, err := mySynopsysCtl.Exec("deploy", "--cluster-scoped", "--enable-alert", "--enable-blackduck", "--enable-opssight", "-i=gcr.io/saas-hub-stg/blackducksoftware/synopsys-operator:release-2019.6.x")
			if err != nil {
				Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
			}
			// Wait for operator to be running
			soLabel := labels.NewSelector()
			r, _ := labels.NewRequirement("app", selection.Equals, []string{"synopsys-operator"})
			soLabel.Add(*r)
			_, err = podutils.WaitForPodsWithLabelRunningReady(kc, "synopsys-operator", soLabel, 2, time.Duration(10*time.Second))
			if err != nil {
				Fail(fmt.Sprintf("Synopsys Operator pods failed to come up: %v", err))
			}
			fmt.Printf("[DEBUG] Synopsys Operator is running\n")
			// Wait for CRDs to be running
			err = crdutils.BlockUntilCrdIsAdded(apiExtensionClient, "alerts.synopsys.com", 10)
			if err != nil {
				Fail(fmt.Sprintf("Alert crd was not added: %v", err))
			}
			err = crdutils.BlockUntilCrdIsAdded(apiExtensionClient, "blackducks.synopsys.com", 10)
			if err != nil {
				Fail(fmt.Sprintf("Black Duck crd was not added: %v", err))
			}
			err = crdutils.BlockUntilCrdIsAdded(apiExtensionClient, "opssights.synopsys.com", 10)
			if err != nil {
				Fail(fmt.Sprintf("OpsSight crd was not added: %v", err))
			}
			fmt.Printf("[DEBUG] CRDs exists\n")
			// Create an Alert
			out, err = mySynopsysCtl.Exec("create", "alert", "alt-one", "--standalone=false", "--persistent-storage=false")
			if err != nil {
				Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
			}
			time.Sleep(5 * time.Second)
			alertExists, err := crutils.AlertCRExists(kc.RESTClient(), "alt-one", "alt-one")
			if err != nil {
				Fail(fmt.Sprintf("bad get : %s", err))
			}
			if !alertExists {
				Fail(fmt.Sprintf("Alert CR was not created : %v", alertExists))
			}
			// Create a Black Duck
			out, err = mySynopsysCtl.Exec("create", "blackduck", "bd-one", "--admin-password=blackduck", "--postgres-password=blackduck", "--user-password=blackduck")
			if err != nil {
				Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
			}
			time.Sleep(5 * time.Second)
			blackDuckExists, err := crutils.BlackDuckCRExists(kc.RESTClient(), "bd-one", "bd-one")
			if err != nil {
				Fail(fmt.Sprintf("bad get : %s", err))
			}
			if !blackDuckExists {
				Fail(fmt.Sprintf("Black Duck CR was not created : %v", blackDuckExists))
			}
			// Create an OpsSight
			out, err = mySynopsysCtl.Exec("create", "opssight", "ops-one")
			if err != nil {
				Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
			}
			time.Sleep(5 * time.Second)
			opsSightExists, err := crutils.OpsSightCRExists(kc.RESTClient(), "ops-one", "ops-one")
			if err != nil {
				Fail(fmt.Sprintf("bad get : %s", err))
			}
			if !opsSightExists {
				Fail(fmt.Sprintf("OpsSight CR was not created : %v", opsSightExists))
			}
			// Clean Up
			// TODO : change to cleanup with synopsysctl
			kc.CoreV1().Namespaces().Delete("synopsys-operator", &metav1.DeleteOptions{})
			kc.CoreV1().Namespaces().Delete("alt-one", &metav1.DeleteOptions{})
			kc.CoreV1().Namespaces().Delete("bd-one", &metav1.DeleteOptions{})
			kc.CoreV1().Namespaces().Delete("ops-one", &metav1.DeleteOptions{})
			kc.RbacV1().ClusterRoles().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
			kc.RbacV1().ClusterRoleBindings().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
			apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("alerts.synopsys.com", &metav1.DeleteOptions{})
			apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("blackducks.synopsys.com", &metav1.DeleteOptions{})
			apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("opssights.synopsys.com", &metav1.DeleteOptions{})
			namespaceutils.WaitForNamespacesDeleted(kc, []string{"synopsys-operator", "alt-one", "bd-one", "ops-one"}, time.Duration(30*time.Second))
		})

		Specify("namespace scoped operations", func() {
			// Deploy Operator in Namespace Scope
			kc.CoreV1().Namespaces().Create(&corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "so-test",
					Name:      "so-test",
				},
			})
			out, err := mySynopsysCtl.Exec("deploy", "-n=so-test", "--enable-alert", "--enable-blackduck", "-i=gcr.io/saas-hub-stg/blackducksoftware/synopsys-operator:release-2019.6.x")
			if err != nil {
				Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
			}
			// Wait for Operator to be running
			soLabel := labels.NewSelector()
			r, _ := labels.NewRequirement("app", selection.Equals, []string{"synopsys-operator"})
			soLabel.Add(*r)
			_, err = podutils.WaitForPodsWithLabelRunningReady(kc, "so-test", soLabel, 2, time.Duration(10*time.Second))
			if err != nil {
				Fail(fmt.Sprintf("Synopsys Operator pods failed to come up: %v", err))
			}
			// Wait for CRDs to be running
			err = crdutils.BlockUntilCrdIsAdded(apiExtensionClient, "alerts.synopsys.com", 10)
			if err != nil {
				Fail(fmt.Sprintf("Alert crd was not added: %v", err))
			}
			err = crdutils.BlockUntilCrdIsAdded(apiExtensionClient, "blackducks.synopsys.com", 10)
			if err != nil {
				Fail(fmt.Sprintf("Black Duck crd was not added: %v", err))
			}
			// Create an Alert
			out, err = mySynopsysCtl.Exec("create", "alert", "alt-one", "--namespace=so-test", "--standalone=false", "--persistent-storage=false")
			if err != nil {
				Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
			}
			time.Sleep(5 * time.Second)
			alertExists, err := crutils.AlertCRExists(kc.RESTClient(), "so-test", "alt-one")
			if err != nil {
				Fail(fmt.Sprintf("bad get : %s", err))
			}
			if !alertExists {
				Fail(fmt.Sprintf("Alert CR was not created : %v", alertExists))
			}
			// Create a Black Duck
			out, err = mySynopsysCtl.Exec("create", "blackduck", "bd-one", "--namespace=so-test", "--admin-password=blackduck", "--postgres-password=blackduck", "--user-password=blackduck")
			if err != nil {
				Fail(fmt.Sprintf("Out: %s Error: %v", out, err))
			}
			time.Sleep(5 * time.Second)
			blackDuckExists, err := crutils.BlackDuckCRExists(kc.RESTClient(), "so-test", "bd-one")
			if err != nil {
				Fail(fmt.Sprintf("bad get : %s", err))
			}
			if !blackDuckExists {
				Fail(fmt.Sprintf("Black Duck CR was not created : %v", blackDuckExists))
			}
			// Clean Up
			// TODO : change to cleanup with synopsysctl
			kc.CoreV1().Namespaces().Delete("synopsys-operator", &metav1.DeleteOptions{})
			kc.RbacV1().ClusterRoles().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
			kc.RbacV1().ClusterRoleBindings().Delete("synopsys-operator-admin", &metav1.DeleteOptions{})
			apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("alerts.synopsys.com", &metav1.DeleteOptions{})
			apiExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete("blackducks.synopsys.com", &metav1.DeleteOptions{})
			namespaceutils.WaitForNamespacesDeleted(kc, []string{"so-test"}, time.Duration(30*time.Second))
		})
	})
})
