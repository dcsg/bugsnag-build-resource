package main_test

import (
	"bytes"
	"encoding/json"

	resource "github.com/dcsg/bugsnag-build-resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Out", func() {
	var (
		session *gexec.Session
	)

	Context("when given params", func() {
		var (
			bugsnagBuildParams resource.BugsnagBuildParams
			bugsnagResponse    resource.BugsnagResponse
			outParams          resource.OutParams
			outResponse        resource.OutResponse
			source             = resource.Source{ApiKey: "api-key"}
		)

		JustBeforeEach(func() {
			fakeBugsnagServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/"),
					ghttp.VerifyJSONRepresenting(&bugsnagBuildParams),
					ghttp.RespondWithJSONEncoded(200, bugsnagResponse),
				),
			)
		})

		Context("contains all available options", func() {
			BeforeEach(func() {
				outParams = resource.OutParams{
					AppVersion:   "123",
					ReleaseStage: "prod",
					SourceControl: resource.SourceControl{
						Provider:      "github",
						RepositoryUrl: "https://github.com/dcsg/bugsnag-build-resource",
						Revision:      "123",
					},
				}

				bugsnagBuildParams = resource.BugsnagBuildParams{
					ApiKey:       "api-key",
					AppVersion:   "123",
					ReleaseStage: "prod",
					SourceControl: resource.SourceControl{
						Provider:      "github",
						RepositoryUrl: "https://github.com/dcsg/bugsnag-build-resource",
						Revision:      "123",
					},
				}

				bugsnagResponse = resource.BugsnagResponse{
					Status: "ok",
				}
			})

			It("creates the build via the API", func() {
				session = RunOut(source, outParams)
				Expect(fakeBugsnagServer.ReceivedRequests()).To(HaveLen(1))
			})

			It("emits metadata about the event", func() {
				session = RunOut(source, outParams)

				err = json.NewDecoder(bytes.NewBuffer(session.Out.Contents())).Decode(&outResponse)
				Expect(err).NotTo(HaveOccurred())

				Expect(outResponse.Version).To(Equal(resource.Version{
					Build: "123",
				}))

				Expect(outResponse.Metadata).To(ConsistOf(
					resource.Metadata{Name: "app_version", Value: "123"},
					resource.Metadata{Name: "release_stage", Value: "prod"},
					resource.Metadata{Name: "source_control.provider", Value: "github"},
					resource.Metadata{Name: "source_control.repository", Value: "https://github.com/dcsg/bugsnag-build-resource"},
					resource.Metadata{Name: "source_control.revision", Value: "123"},
				))
			})
		})

		Context("contains only the required options", func() {
			BeforeEach(func() {
				outParams = resource.OutParams{
					AppVersion: "123",
				}

				bugsnagBuildParams = resource.BugsnagBuildParams{
					ApiKey:     "api-key",
					AppVersion: "123",
				}

				bugsnagResponse = resource.BugsnagResponse{
					Status: "ok",
				}
			})

			It("creates the build via the API", func() {
				session = RunOut(source, outParams)
				Expect(fakeBugsnagServer.ReceivedRequests()).To(HaveLen(1))
			})

			It("emits metadata", func() {
				session = RunOut(source, outParams)

				err = json.NewDecoder(bytes.NewBuffer(session.Out.Contents())).Decode(&outResponse)
				Expect(err).NotTo(HaveOccurred())

				Expect(outResponse.Version).To(Equal(resource.Version{
					Build: "123",
				}))

				Expect(outResponse.Metadata).To(ConsistOf(
					resource.Metadata{Name: "app_version", Value: "123"},
				))
			})
		})

		Context("contains source_control params", func() {
			It("panics because repository is missing", func() {
				outParams = resource.OutParams{
					AppVersion: "123",
					SourceControl: resource.SourceControl{
						Provider: "github",
					},
				}

				session = RunOut(source, outParams, Not(gexec.Exit(0)))
				Expect(session.Err).To(gbytes.Say("'repository' is required when source_control is defined"))

				Expect(fakeBugsnagServer.ReceivedRequests()).To(HaveLen(0))
			})

			It("panics because revision is missing", func() {
				outParams = resource.OutParams{
					AppVersion: "123",
					SourceControl: resource.SourceControl{
						Provider:      "github",
						RepositoryUrl: "https://github.com/dcsg/bugsnag-build-resource",
					},
				}

				session = RunOut(source, outParams, Not(gexec.Exit(0)))
				Expect(session.Err).To(gbytes.Say("'revision' param is required when source_control is defined"))

				Expect(fakeBugsnagServer.ReceivedRequests()).To(HaveLen(0))
			})
		})

		Context("when bugsnag returns warnings", func() {

			BeforeEach(func() {
				outParams = resource.OutParams{
					AppVersion: "123",
				}

				bugsnagBuildParams = resource.BugsnagBuildParams{
					ApiKey:     "api-key",
					AppVersion: "123",
				}

				bugsnagResponse = resource.BugsnagResponse{
					Status: "ok",
					Warnings: []string{
						"first-warning",
						"second-warning",
					},
				}
			})

			It("emits the warnings in metadata", func() {
				session = RunOut(source, outParams)

				err = json.NewDecoder(bytes.NewBuffer(session.Out.Contents())).Decode(&outResponse)
				Expect(err).NotTo(HaveOccurred())

				Expect(outResponse.Version).To(Equal(resource.Version{
					Build: "123",
				}))

				Expect(outResponse.Metadata).To(ConsistOf(
					resource.Metadata{Name: "app_version", Value: "123"},
					resource.Metadata{Name: "bugsnag.warnings", Value: "[\"first-warning\",\"second-warning\"]"},
				))
			})
		})
	})
})
