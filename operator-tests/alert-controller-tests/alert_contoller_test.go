package alert_controller_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/blackducksoftware/cloud-native-tests/utils"
)

func TestGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Alert Controller Test Suite")
}

var _ = Describe("Alert Controller", func() {

	mySynopsysCtl := utils.NewSynopsysctl("synopsysctl")
	mySynopsysCtl.Exec("deploy", "--cluster-scoped", "--enable-alert", "--enable-blackduck", "--enable-opssight", "-i=gcr.io/saas-hub-stg/blackducksoftware/synopsys-operator:release-2019.6.x")

	Describe("Creating Alert", func() {
		Context("Initially", func() {
			It("there are no Alert instances", func() {})
		})
		Context("create event", func() {
			It("there is one Alert instance", func() {})
			It("the Alert instance has name 'alt'", func() {})
		})
		Context("another create event", func() {
			It("there are two Alert instances", func() {})
			It("the Alert instance has name 'alt'", func() {})
		})
	})

	Describe("Updating Alert", func() {
		Context("Initially", func() {
			It("has an Alert instance", func() {})
			It("has an Alert instance with name alt", func() {})
			It("has an Alert instance with port <old_port>", func() {})
		})
		Context("update event", func() {
			It("has an Alert instance", func() {})
			It("has an Alert instance with port <new_port>", func() {})
		})
	})

	Describe("Deleting Alert", func() {
		Context("delete event for an Alert instance that doesn't exist", func() {
			It("there is no alert with name 'alt'", func() {})
		})
		Context("delete event for an existing Alert instance", func() {
			It("there is no alert with name 'alt'", func() {})
		})
		Context("delete event with no Alert instance", func() {
			It("nothing happens", func() {})
		})
	})

})
