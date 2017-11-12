package main_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	//"path/filepath"

	"github.com/regevbr/keyval-resource/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"path/filepath"
	"github.com/magiconair/properties"
)

var _ = Describe("In", func() {
	var tmpdir string
	var destination string

	var inCmd *exec.Cmd

	BeforeEach(func() {
		var err error

		tmpdir, err = ioutil.TempDir("", "in-destination")
		Expect(err).NotTo(HaveOccurred())

		destination = path.Join(tmpdir, "in-dir")

		inCmd = exec.Command(inPath, destination)
	})

	AfterEach(func() {
		os.RemoveAll(tmpdir)
	})

	Context("when executed", func() {
		var request models.InRequest
		var response models.InResponse

		BeforeEach(func() {

			request = models.InRequest{
				Version: models.Version{
					"a": "1",
					"b": "2",
					"dummy": "dummy",
				},
				Source:  models.Source{},
			}

			response = models.InResponse{}
		})

		JustBeforeEach(func() {
			stdin, err := inCmd.StdinPipe()
			Expect(err).NotTo(HaveOccurred())

			session, err := gexec.Start(inCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			err = json.NewEncoder(stdin).Encode(request)
			Expect(err).NotTo(HaveOccurred())

			<-session.Exited
			Expect(session.ExitCode()).To(Equal(0))

			err = json.Unmarshal(session.Out.Contents(), &response)
			Expect(err).NotTo(HaveOccurred())
		})

		It("reports the version and metadata to be the input version", func() {
			Expect(len(response.Version)).To(Equal(2))
			Expect(response.Version["a"]).To(Equal("1"))
			Expect(response.Version["b"]).To(Equal("2"))

			Expect(len(response.Metadata)).To(Equal(2))
			Expect(response.Metadata[0].Name).To(Equal("a"))
			Expect(response.Metadata[0].Value).To(Equal("1"))
			Expect(response.Metadata[1].Name).To(Equal("b"))
			Expect(response.Metadata[1].Value).To(Equal("2"))
		})

		It("writes the requested data the destination", func() {

			var data = properties.MustLoadFile(filepath.Join(destination, "keyval.properties"),properties.UTF8).Map();

			Expect(len(data)).To(Equal(2))
			Expect(data["a"]).To(Equal("1"))
			Expect(data["b"]).To(Equal("2"))
		})

		Context("when the request has no keys in version", func() {
			BeforeEach(func() {
				request.Version = models.Version{}
			})

			It("reports  empty data", func() {
				Expect(len(response.Version)).To(Equal(0))
				Expect(len(response.Metadata)).To(Equal(0))
				var data = properties.MustLoadFile(filepath.Join(destination, "keyval.properties"),properties.UTF8).Map();
				Expect(len(data)).To(Equal(0))
			})
		})
	})
})
