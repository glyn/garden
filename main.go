package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	rdebug "runtime/debug"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/pivotal-cf-experimental/garden/backend"
	"github.com/pivotal-cf-experimental/garden/backend/fake_backend"
	"github.com/pivotal-cf-experimental/garden/command_runner"
	"github.com/pivotal-cf-experimental/garden/linux_backend"
	"github.com/pivotal-cf-experimental/garden/linux_backend/container_pool"
	"github.com/pivotal-cf-experimental/garden/linux_backend/network_pool"
	"github.com/pivotal-cf-experimental/garden/linux_backend/port_pool"
	"github.com/pivotal-cf-experimental/garden/linux_backend/quota_manager"
	"github.com/pivotal-cf-experimental/garden/linux_backend/uid_pool"
	"github.com/pivotal-cf-experimental/garden/server"
)

var listenNetwork = flag.String(
	"listenNetwork",
	"unix",
	"how to listen on the address (unix, tcp, etc.)",
)

var listenAddr = flag.String(
	"listenAddr",
	"/tmp/warden.sock",
	"address to listen on",
)

var snapshotsPath = flag.String(
	"snapshots",
	"",
	"directory in which to store container state to persist through restarts",
)

var backendName = flag.String(
	"backend",
	"linux",
	"which backend to use (linux or fake)",
)

var binPath = flag.String(
	"bin",
	"",
	"directory containing backend-specific scripts (i.e. ./create.sh)",
)

var depotPath = flag.String(
	"depot",
	"",
	"directory in which to store containers",
)

var rootFSPath = flag.String(
	"rootfs",
	"",
	"directory of the rootfs for the containers",
)

var disableQuotas = flag.Bool(
	"disableQuotas",
	false,
	"disable disk quotas",
)

var containerGraceTime = flag.Int(
	"containerGraceTime",
	0,
	"time (in seconds) after which to destroy idle containers",
)

var debug = flag.Bool(
	"debug",
	false,
	"show low-level command output",
)

func main() {
     	//dumpFilesystem("Main entry.")
	flag.Parse()

	maxProcs := runtime.NumCPU()
	prevMaxProcs := runtime.GOMAXPROCS(maxProcs)

	log.Println("set GOMAXPROCS to", maxProcs, "was", prevMaxProcs)

	var backend backend.Backend

	switch *backendName {
	case "linux":
		if *binPath == "" {
			log.Fatalln("must specify -bin with linux backend")
		}

		if *depotPath == "" {
			log.Fatalln("must specify -depot with linux backend")
		}

		if *rootFSPath == "" {
			log.Fatalln("must specify -rootfs with linux backend")
		}

		uidPool := uid_pool.New(10000, 256)

		_, ipNet, err := net.ParseCIDR("10.254.0.0/22")
		if err != nil {
			log.Fatalln("error parsing CIDR:", err)
		}

		networkPool := network_pool.New(ipNet)

		// TODO: base on ephemeral port range
		portPool := port_pool.New(61000, 6501)

		var runner command_runner.CommandRunner

		runner = command_runner.New(*debug)

		quotaManager, err := quota_manager.New(*depotPath, *binPath, runner)
		if err != nil {
			log.Fatalln("error creating quota manager:", err)
		}

		if *disableQuotas {
			quotaManager.Disable()
		}

		pool := container_pool.New(
			*binPath,
			*depotPath,
			*rootFSPath,
			uidPool,
			networkPool,
			portPool,
			runner,
			quotaManager,
		)

		backend = linux_backend.New(pool, *snapshotsPath)
	case "fake":
		backend = fake_backend.New()
	}

	log.Println("setting up backend")

	err := backend.Setup()
	if err != nil {
		log.Fatalln("failed to set up backend:", err)
	}

	log.Println("starting server; listening with", *listenNetwork, "on", *listenAddr)

	graceTime := time.Duration(*containerGraceTime) * time.Second

	wardenServer := server.New(*listenNetwork, *listenAddr, graceTime, backend)

	err = wardenServer.Start()
	if err != nil {
		log.Fatalln("failed to start:", err)
	}

	signals := make(chan os.Signal, 1)

	go func() {
		signal := <-signals
		log.Printf("stopping due to %v ...\n", signal)
		wardenServer.Stop()

		//dumpFilesystem("After stop.")
				
		os.Exit(0)
	}()

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	select {} // block the current goroutine
}

var dumpCount int
var mutex sync.Mutex

func dumpFilesystem(tag string) {

        //time.Sleep(2 * time.Second)

        log.Print("dumpFilesystem context:")
	if false {
           rdebug.PrintStack()
	}

	mutex.Lock()
	log.Println(tag, "main.go", os.Getpid(), dumpCount, "Filesystem dump:")
	dumpCount++
	mutex.Unlock()

	dumpGlob("Containers:", "/tmp/garden-server*/containers/*", 4)
	dumpGlob("Snapshots:", "/tmp/garden-server*/snapshots/*", 4)
}

func dumpGlob(title string, pattern string, indent int) {
        ind := strings.Repeat(" ", indent)

	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalln("Glob failed:", err)
	}

	log.Println(ind, title);

	sort.Strings(files)
	for _, file := range files {
	       log.Println(ind, ind, file)
	}

	os.Stderr.Sync()
}
