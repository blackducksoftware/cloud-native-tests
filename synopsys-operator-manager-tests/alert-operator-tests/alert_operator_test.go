package alert_operator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/blackducksoftware/cloud-native-tests/utils"
)

func TestGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Alert Operator Test Suite")
}

var _ = Describe("Alert Duck Operator", func() {

	mySynopsysCtl := utils.NewSynopsysctl("synopsysctl")

	Context("in Cluster Scope", func() {
		BeforeEach(func() {
			mySynopsysCtl.Exec("deploy", "--cluster-scoped", "--enable-alert", "-i=gcr.io/saas-hub-stg/blackducksoftware/synopsys-operator:release-2019.6.x")
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

		Describe("Updating the CRD Alert image", func() {
			Specify("The container has the correct image", func() {})
		})

		Describe("Updating the CRD Cfssl image", func() {
			Specify("The container has the correct image ", func() {})
		})

		Describe("Updating the CRD Expose service", func() {
			Specify("The correct service appears", func() {})
		})

		Describe("Updating the CRD stand alone", func() {
			Specify("If true, there is a Cfssl pod", func() {})
			Specify("If false, there is not a Cfssl pod", func() {})
		})

		Describe("Updating the CRD port", func() {
			Specify("The port is set correctly", func() {})
		})

		Describe("Updating the CRD EncryptionPassword", func() {
			Specify("It's correctly added to the Config Map", func() {})
		})

		Describe("Updating the CRD EncryptionGlobalSalt", func() {
			Specify("It's correctly added to the Config Map", func() {})
		})

		Describe("Updating the CRD Environs", func() {
			Specify("It's correctly added to the Config Map", func() {})
		})

		Describe("Updating the CRD PersistentStorage", func() {
			Specify("It's correctly updates persistent storage", func() {})
		})

		Describe("Updating the CRD PVCName", func() {
			Specify("It's correctly updates the PVC name", func() {})
		})

		Describe("Updating the CRD PVCStorageClass", func() {
			Specify("It's correctly updates the PVC storage class", func() {})
		})

		Describe("Updating the CRD PVCSize", func() {
			Specify("It's correctly updates the PVC size", func() {})
		})

		Describe("Updating the CRD AlertMemory", func() {
			Specify("It's correctly updates the Alert memory", func() {})
		})

		Describe("Updating the CRD CfsslMemory", func() {
			Specify("It's correctly updates the Cfssl memory", func() {})
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
