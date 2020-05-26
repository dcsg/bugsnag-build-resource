package main_test

import (
	"encoding/json"
	"os/exec"

	resource "github.com/dcsg/bugsnag-build-resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"bytes"
	"io/ioutil"
	"testing"

	"fmt"
	"os"

	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/types"
)

func TestOut(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Out Suite")
}

var (
	binPath           string
	err               error
	fakeBugsnagServer *ghttp.Server

	tmpDir string
)

var _ = BeforeEach(func() {
	fakeBugsnagServer = ghttp.NewServer()

	if _, err = os.Stat("/opt/resource/out"); err == nil {
		binPath = "/opt/resource/out"
	} else {
		binPath, err = gexec.Build("github.com/dcsg/bugsnag-build-resource/cmd/out")
		Expect(err).NotTo(HaveOccurred())
	}


	tmpDir, err = ioutil.TempDir("", "bugsnag_build_resource_out")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterEach(func() {
	fakeBugsnagServer.Close()
	err := os.RemoveAll(tmpDir)
	Expect(err).NotTo(HaveOccurred())
})

func RunOut(s resource.Source, p resource.OutParams, matchers ...types.GomegaMatcher) *gexec.Session {
	payload := resource.OutRequest{
		Source:    s,
		OutParams: p,
	}

	b, err := json.Marshal(&payload)
	Expect(err).NotTo(HaveOccurred())

	c := exec.Command(binPath, tmpDir)
	c.Stdin = bytes.NewBuffer(b)
	envVars := map[string]string{
		"BUGSNAG_HOST": "http://" + fakeBugsnagServer.Addr(),
	}
	for k, v := range envVars {
		c.Env = append(c.Env, fmt.Sprintf("%s=%s", k, v))
	}
	c.Stdin = bytes.NewBuffer(b)
	sess, err := gexec.Start(c, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	<-sess.Exited
	if len(matchers) == 0 {
		Expect(sess).To(gexec.Exit(0), fmt.Sprintf("Expected session to exit 0, exited with %d.\n\nStdout: %s\n\nStderr: %s", sess.ExitCode(), sess.Out.Contents(), sess.Err.Contents()))
	} else {
		for _, matcher := range matchers {
			Expect(sess).To(matcher)
		}
	}

	return sess
}
