package smoke_test

import (
	"flag"
	"testing"
	"time"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var Builder string

const procfileBuildpack = "gcr.io/paketo-buildpacks/procfile"

func init() {
	flag.StringVar(&Builder, "name", "", "")
}

func TestSmoke(t *testing.T) {
	Expect := NewWithT(t).Expect

	flag.Parse()

	Expect(Builder).NotTo(Equal(""))

	SetDefaultEventuallyTimeout(60 * time.Second)

	suite := spec.New("Buildpackless Smoke", spec.Parallel(), spec.Report(report.Terminal{}))
	suite("Procfile", testProcfile)
	suite.Run(t)
}
