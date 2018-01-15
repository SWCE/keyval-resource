package main_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"github.com/SWCE/keyval-resource/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"bufio"
	"fmt"
	"github.com/onsi/gomega/gbytes"
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
				var outDir = path.Join(source, "out")
				os.MkdirAll(outDir, 0755)
				file, _ := os.Create(path.Join(outDir, "keyval.properties"))
				defer file.Close()
				w := bufio.NewWriter(file)
				fmt.Fprintln(w, "a=1")
				fmt.Fprintln(w, "b=2")
				w.Flush()
				request = models.OutRequest{
					Params: models.OutParams {
						File: path.Join("out", "keyval.properties"),
					},		
				}
			})

			It("reports empty data", func() {
				Expect(len(response.Version)).To(Equal(4))
				Expect(response.Version["a"]).To(Equal("1"))
				Expect(response.Version["b"]).To(Equal("2"))
				Expect(response.Version).To(HaveKey("UPDATED"))
				Expect(response.Version).To(HaveKey("UUID"))
				Expect(response.Version["UPDATED"]).To(Not(BeEmpty()))
				Expect(response.Version["UUID"]).To(Not(BeEmpty()))
			})
		})

		Context("output no data in properties file", func() {
			BeforeEach(func() {
				var outDir = path.Join(source, "out")
				os.MkdirAll(outDir, 0755)
				file, _ := os.Create(path.Join(outDir, "keyval.properties"))
				defer file.Close()
				request = models.OutRequest{
					Params: models.OutParams{
						File: path.Join("out", "keyval.properties"),
					},		
				}
			})

			It("reports empty data", func() {
				Expect(len(response.Version)).To(Equal(2))
				Expect(response.Version).To(HaveKey("UPDATED"))
				Expect(response.Version).To(HaveKey("UUID"))
				Expect(response.Version["UPDATED"]).To(Not(BeEmpty()))
				Expect(response.Version["UUID"]).To(Not(BeEmpty()))
			})
		})

	})

	Context("with invalid inputs", func() {
		var request models.OutRequest
		var response models.OutResponse
		var session *gexec.Session

		BeforeEach(func() {
			request = models.OutRequest{}
			response = models.OutResponse{}
		})

		JustBeforeEach(func() {
			stdin := new(bytes.Buffer)

			err := json.NewEncoder(stdin).Encode(request)
			Expect(err).NotTo(HaveOccurred())

			outCmd.Stdin = stdin

			session, err = gexec.Start(outCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

		})

		Context("no file specified", func() {

			It("reports error", func() {
				<-session.Exited
				Expect(session.Err).To(gbytes.Say("no properties file specified"))
				Expect(session.ExitCode()).To(Equal(1))
			})
		})

	})
})
