package main_test

import (
	"encoding/json"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"

	"../models"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Check", func() {
	var (
		checkCmd *exec.Cmd
	)

	BeforeEach(func() {
		checkCmd = exec.Command(checkPath)
	})

	Context("when executed", func() {
		var source map[string]interface{}
		var version *models.EmptyVersion;
		var response models.CheckResponse

		BeforeEach(func() {
			source = map[string]interface{}{}
			response = models.CheckResponse{}
			version = nil
		})

		JustBeforeEach(func() {
			stdin, err := checkCmd.StdinPipe()
			Expect(err).NotTo(HaveOccurred())

			session, err := gexec.Start(checkCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			err = json.NewEncoder(stdin).Encode(map[string]interface{}{
				"source":  source,
				"version": version,
			})
			Expect(err).NotTo(HaveOccurred())

			<-session.Exited
			Expect(session.ExitCode()).To(Equal(0))

			err = json.Unmarshal(session.Out.Contents(), &response)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when no version is given", func() {
			It("outputs an empty version", func() {
				Expect(response).To(HaveLen(1))
			})
		})

		Context("when version is given", func() {

			BeforeEach(func() {
				version = &models.EmptyVersion{}
			})

			It("outputs an empty version", func() {
				Expect(response).To(HaveLen(1))
			})
		})

	})

	Context("with invalid inputs", func() {
		var session *gexec.Session

		JustBeforeEach(func() {
			stdin, err := checkCmd.StdinPipe()
			Expect(err).NotTo(HaveOccurred())

			session, err = gexec.Start(checkCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			stdin.Close()
		});
		Context("with a missing everything", func() {
			It("returns an error", func() {
				<-session.Exited
				Expect(session.Err).To(gbytes.Say("parse error: EOF"))
				Expect(session.ExitCode()).To(Equal(1))
			})
		})

	})
});
