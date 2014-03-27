package lifecycle_test

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf-experimental/garden/integration/garden_runner"
	"github.com/vito/cmdtest"
	"github.com/vito/gordon"
)

var runner *garden_runner.GardenRunner
var client gordon.Client

func TestLifecycle(t *testing.T) {
	binPath := "../../linux_backend/bin"
	rootFSPath := os.Getenv("GARDEN_TEST_ROOTFS")

	if rootFSPath == "" {
		log.Println("GARDEN_TEST_ROOTFS undefined; skipping")
		return
	}

	var err error

	tmpdir, err := ioutil.TempDir("", "garden-socket")
	if err != nil {
		log.Fatalln("failed to make dir for socker:", err)
	}

	gardenPath, err := cmdtest.Build("github.com/pivotal-cf-experimental/garden", "-race")
	if err != nil {
		log.Fatalln("failed to compile garden:", err)
	}

	runner, err = garden_runner.New(gardenPath, binPath, rootFSPath, "unix", filepath.Join(tmpdir, "warden.sock"))
	if err != nil {
		log.Fatalln("failed to create runner:", err)
	}

	defer cleanup(runner, tmpdir)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		s := <-c
		cleanup(runner, tmpdir)
		log.Fatalf("Terminated by signal %s\n", s.String())
	}()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Lifecycle Suite")
}

func cleanup(runner *garden_runner.GardenRunner, tmpdir string) {
	err := runner.Stop()
	if err != nil {
		log.Fatalln("garden failed to stop:", err)
	}

	err = runner.TearDown()
	if err != nil {
		log.Fatalln("failed to tear down server:", err)
	}

	err = os.RemoveAll(tmpdir)
	if err != nil {
		log.Fatalln("failed to clean up socket dir:", err)
	}
}

var didRunGarden bool

var _ = BeforeEach(func() {
	if didRunGarden {
		return
	}
	didRunGarden = true
	err := runner.Start()
	if err != nil {
		log.Fatalln("garden failed to start:", err)
	}

	client = runner.NewClient()
})
