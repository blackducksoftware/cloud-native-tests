package opssight_operator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/blackducksoftware/cloud-native-tests/utils"
)

func TestGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OpsSight Operator Test Suite")
}

var _ = Describe("OpsSight Duck Operator", func() {

	mySynopsysCtl := utils.NewSynopsysctl("synopsysctl")

	Context("in Cluster Scope", func() {
		BeforeEach(func() {
			mySynopsysCtl.Exec("deploy", "--cluster-scoped", "--enable-opssight", "-i=gcr.io/saas-hub-stg/blackducksoftware/synopsys-operator:release-2019.6.x")
			// TODO: wait until pods are running
		})
		AfterEach(func() {
			mySynopsysCtl.Exec("destroy")
			// TODO: wait until pads are deleted
		})

		/*
			Tests for handling create events of a Custom Resource
		*/
		Describe("Creating an instance", func() {
			Context("namespaced mode", func() {
				It("creates resource in the namespace of the operator", func() {})
			})
			Context("clusterspaced mode", func() {
				It("creates the resource in its own namespace", func() {})
			})
			It("Can create into database migration mode", func() {})
			Context("Persistent Storage", func() {
				It("has Persistent Volume Claims", func() {})
			})
		})

		/*
			Tests for handling delete events to a Custom Resource
		*/
		Describe("Deleting an instance", func() {
			Context("namespaced mode", func() {
				It("leaves the namespace if empty", func() {})
			})
			Context("clusterspaced mode", func() {
				It("deletes the namespace if empty", func() {})
			})
			Context("delete event for an OpsSight instance that doesn't exist", func() {
				It("there is no OpsSight with name 'alt'", func() {})
			})
		})

		/*
			Tests for handling update events to a Custom Resource
		*/
		Describe("Updating the CRD Version", func() {
			Describe("Older to Newer", func() {
				Specify("The correct pods are running", func() {})
			})
			Describe("Newer to Older", func() {
				Specify("The correct pods are running", func() {})
			})
		})

		Describe("Updating the CRD desired state", func() {
			Specify("If stop, all the pods are running", func() {})
			Specify("If start, only the PVCs exist in the namespace", func() {})
		})

		/*
			Tests for maintiaing the state of the Custom Resources
		*/
		Describe("Operator handles deviation of live objects from Custom Resource Spec (Ensure/crdupdater)", func() {
			Specify("It puts labels back if removed", func() {})
			Specify("It puts deployment back", func() {})
		})

	})
})
