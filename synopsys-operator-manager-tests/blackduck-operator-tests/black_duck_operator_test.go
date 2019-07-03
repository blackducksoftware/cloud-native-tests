package black_duck_operator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/blackducksoftware/cloud-native-tests/utils"
)

func TestGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Black Duck Operator Test Suite")
}

var _ = Describe("Black Duck Operator", func() {

	mySynopsysCtl := utils.NewSynopsysctl("synopsysctl")

	Context("in Cluster Scope", func() {
		BeforeEach(func() {
			mySynopsysCtl.Exec("deploy", "--cluster-scoped", "--enable-blackduck", "-i=gcr.io/saas-hub-stg/blackducksoftware/synopsys-operator:release-2019.6.x")
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
			Context("delete event for an Alert instance that doesn't exist", func() {
				It("there is no Alert with name 'alt'", func() {})
			})
		})

		/*
			Tests for handling update events to a Custom Resource
		*/

		Describe("Updating the Version", func() {
			Describe("Older to Newer", func() {
				Specify("The correct pods are running", func() {})
			})
			Describe("Newer to Older", func() {
				Specify("The correct pods are running", func() {})
			})
		})

		Describe("Updating the desired state", func() {
			Describe("Stop", func() {
				Specify("All the pods are running", func() {})
			})
			Describe("Start", func() {
				Specify("Only the PVCs exist in the namespace", func() {})
			})
			Describe("db-migrate", func() {
				It("Has Persistent Volume Claims", func() {})
				It("The only pod is for postgres", func() {})
			})
		})

		Describe("Updating the license key", func() {
			Specify("It can update the license key", func() {})
		})

		Describe("Updating the port", func() {
			Specify("It can update the port", func() {})
		})

		Describe("Updating the type", func() {
			Specify("It can update the type", func() {})
		})

		Describe("Updating the environs", func() {
			Specify("It can update the environs", func() {})
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
