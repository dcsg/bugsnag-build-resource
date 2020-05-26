package main_test

import (
	"bytes"
	"encoding/json"

	resource "github.com/dcsg/bugsnag-build-resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("In", func() {
	var (
		session    *gexec.Session
		inResponse resource.InResponse
		version    = resource.Version{
			Build: "123",
		}
	)

	It("emits build", func() {
		session = RunIn(resource.Source{}, version)

		err := json.NewDecoder(bytes.NewBuffer(session.Out.Contents())).Decode(&inResponse)
		Expect(err).NotTo(HaveOccurred())

		Expect(inResponse.Version.Build).To(Equal(version.Build))
	})

	It("emits empty metadata", func() {
		session = RunIn(resource.Source{}, version)

		err = json.NewDecoder(bytes.NewBuffer(session.Out.Contents())).Decode(&inResponse)
		Expect(err).NotTo(HaveOccurred())

		Expect(inResponse.Metadata).To(ConsistOf(make([]resource.Metadata, 0)))
	})
})
