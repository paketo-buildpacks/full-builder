package smoke_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testWebServers(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack().WithVerbose().WithNoColor()
		docker = occam.NewDocker()
	})

	context("detects a HTTPD app", func() {
		var (
			image     occam.Image
			container occam.Container

			name   string
			source string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("builds successfully", func() {
			var err error
			source, err = occam.Source(filepath.Join("testdata", "httpd"))
			Expect(err).NotTo(HaveOccurred())

			var logs fmt.Stringer
			image, logs, err = pack.Build.
				WithPullPolicy("never").
				WithBuilder(Builder).
				Execute(name, source)
			Expect(err).ToNot(HaveOccurred(), logs.String)

			container, err = docker.Container.Run.
				WithEnv(map[string]string{"PORT": "8080"}).
				WithPublish("8080").
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(BeAvailable())

			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Apache HTTP Server")))
		})
	})

	context("detects a NGINX app", func() {
		var (
			image     occam.Image
			container occam.Container

			name   string
			source string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("builds successfully", func() {
			var err error
			source, err = occam.Source(filepath.Join("testdata", "nginx"))
			Expect(err).NotTo(HaveOccurred())

			var logs fmt.Stringer
			image, logs, err = pack.Build.
				WithPullPolicy("never").
				WithBuilder(Builder).
				Execute(name, source)
			Expect(err).ToNot(HaveOccurred(), logs.String)

			container, err = docker.Container.Run.
				WithEnv(map[string]string{"PORT": "8080"}).
				WithPublish("8080").
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(BeAvailable())

			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Nginx Server")))
		})
	})

	context("detects a JavaScript frontend app", func() {
		var (
			image     occam.Image
			container occam.Container

			name   string
			source string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		context("app uses react and httpd", func() {
			it("builds successfully", func() {
				var err error
				source, err = occam.Source(filepath.Join("testdata", "javascript-frontend"))
				Expect(err).NotTo(HaveOccurred())

				var logs fmt.Stringer
				image, logs, err = pack.Build.
					WithPullPolicy("never").
					WithBuilder(Builder).
					WithEnv(map[string]string{
						"BP_NODE_RUN_SCRIPTS":             "build",
						"BP_WEB_SERVER":                   "httpd",
						"BP_WEB_SERVER_ROOT":              "build",
						"BP_WEB_SERVER_ENABLE_PUSH_STATE": "true",
					}).
					Execute(name, source)
				Expect(err).ToNot(HaveOccurred(), logs.String)

				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Node Engine")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for NPM Install")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Node Run Script")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Apache HTTP Server")))

				container, err = docker.Container.Run.
					WithEnv(map[string]string{"PORT": "8080"}).
					WithPublish("8080").
					Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				Eventually(container).Should(Serve(ContainSubstring("<title>Paketo Buildpacks</title>")).OnPort(8080))
			})
		})

		context("app uses react and nginx", func() {
			it("builds successfully", func() {
				var err error
				source, err = occam.Source(filepath.Join("testdata", "javascript-frontend"))
				Expect(err).NotTo(HaveOccurred())

				var logs fmt.Stringer
				image, logs, err = pack.Build.
					WithPullPolicy("never").
					WithBuilder(Builder).
					WithEnv(map[string]string{
						"BP_NODE_RUN_SCRIPTS":             "build",
						"BP_WEB_SERVER":                   "nginx",
						"BP_WEB_SERVER_ROOT":              "build",
						"BP_WEB_SERVER_ENABLE_PUSH_STATE": "true",
					}).
					Execute(name, source)
				Expect(err).ToNot(HaveOccurred(), logs.String)

				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Node Engine")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for NPM Install")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Node Run Script")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Nginx Server")))

				container, err = docker.Container.Run.
					WithEnv(map[string]string{"PORT": "8080"}).
					WithPublish("8080").
					Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				Eventually(container).Should(Serve(ContainSubstring("<title>Paketo Buildpacks</title>")).OnPort(8080))
			})
		})
	})
}
