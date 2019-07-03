package synopsys_operator_test

import (
	"fmt"
	"testing"
	"time"

	utils "github.com/blackducksoftware/cloud-native-tests/utils"
	k8sutils "github.com/blackducksoftware/cloud-native-tests/utils/k8shelper"
	podutils "github.com/blackducksoftware/cloud-native-tests/utils/k8shelper/pod"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

// TestGinkgo TODO
func TestGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Synopsys Operator - Operator Manager Test Suite")
}

var _ = Describe("Synopsys Operator Manager Tests", func() {

	defer GinkgoRecover()

	rc, _ := k8sutils.GetRestConfig()
	kc, _ := k8sutils.GetKubeClient(rc)

	// get Crd
	// apiExtensionClient, err := apiextensionsclient.NewForConfig(rc)
	// if err != nil {
	// Fail(fmt.Sprintf("error creating the api extension client: %+v", err))
	// }

	mySynopsysCtl := utils.NewSynopsysctl("synopsysctl")

	Context("operator in namespace scope", func() {

		Context("creating alert in alert namespace", func() {
			// setup
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

			Specify("default alert", func() {
				// assumes synopsysctl works and we have an operator running in namespace, and crd is made correctly
				// assuming this, create an alert
				_, err := mySynopsysCtl.Exec("create", "alert", "alt", "-n=alert")

				//TODO: wait till pods are running with label app=alert
				label := labels.NewSelector()
				r, _ := labels.NewRequirement("app", selection.Equals, []string{"alert"})
				r2, _ := labels.NewRequirement("name", selection.Equals, []string{"alt"})
				label.Add(*r, *r2)

				_, err = podutils.WaitForPodsWithLabelRunningReady(kc, "alert", label, 2, time.Duration(30*time.Second))
				if err != nil {
					Fail(fmt.Sprintf("alert pods failed to come up: %v", err))
				}
				time.Sleep(5 * time.Second)
				// verify the cr is there
				req := kc.RESTClient().Get().Namespace("alert").Resource("alert")
				if req == nil {
					Fail(fmt.Sprintf("alert cr is not there: %v", err))
				}
			})

			AfterEach(func() {
				// cleanup
				time.Sleep(5 * time.Second)
				mySynopsysCtl.Exec("delete", "alert", "alt", "-n=alert")
				time.Sleep(5 * time.Second)
				mySynopsysCtl.Exec("destroy", "alert")
			})

		})
	})

})
