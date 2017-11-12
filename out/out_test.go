package main_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"../models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"bufio"
	"fmt"
)

var _ = Describe("Out", func() {
	var tmpdir string
	var source string

	var outCmd *exec.Cmd

	BeforeEach(func() {
		var err error

		tmpdir, err = ioutil.TempDir("", "out-source")
		Expect(err).NotTo(HaveOccurred())

		source = path.Join(tmpdir, "out-dir")
		os.MkdirAll(source, 0755)
		outCmd = exec.Command(outPath, source)
		fmt.Printf("%s", tmpdir)
	})

	AfterEach(func() {
		os.RemoveAll(tmpdir)
	})

	Context("when executed", func() {
		var request models.OutRequest
		var response models.OutResponse

		BeforeEach(func() {
			request = models.OutRequest{}
			response = models.OutResponse{}
		})

		JustBeforeEach(func() {
			stdin := new(bytes.Buffer)

			err := json.NewEncoder(stdin).Encode(request)
			Expect(err).NotTo(HaveOccurred())

			outCmd.Stdin = stdin

			session, err := gexec.Start(outCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			<-session.Exited
			Expect(session.ExitCode()).To(Equal(0))

			err = json.Unmarshal(session.Out.Contents(), &response)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("output data in properties file", func() {
			BeforeEach(func() {
				file, _ := os.Create(path.Join(source, "keyval.properties"))
				defer file.Close()
				w := bufio.NewWriter(file)
				fmt.Fprintln(w, "a=1")
				fmt.Fprintln(w, "b=2")
				w.Flush()
			})

			It("reports empty data", func() {
				Expect(len(response.Version)).To(Equal(2))
				Expect(response.Version["a"]).To(Equal("1"))
				Expect(response.Version["b"]).To(Equal("2"))

				Expect(len(response.Metadata)).To(Equal(2))
				Expect(response.Metadata[0].Name).To(Equal("a"))
				Expect(response.Metadata[0].Value).To(Equal("1"))
				Expect(response.Metadata[1].Name).To(Equal("b"))
				Expect(response.Metadata[1].Value).To(Equal("2"))
			})
		})

		Context("output no data in properties file", func() {
			BeforeEach(func() {
				file, _ := os.Create(path.Join(source, "keyval.properties"))
				defer file.Close()
			})

			It("reports empty data", func() {
				Expect(len(response.Version)).To(Equal(0))
				Expect(len(response.Metadata)).To(Equal(0))
			})
		})

	})
})
