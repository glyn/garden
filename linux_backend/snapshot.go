package linux_backend

import (
	"time"

	"github.com/pivotal-cf-experimental/garden/backend"
	"github.com/pivotal-cf-experimental/garden/linux_backend/network"
)

type ContainerSnapshot struct {
	ID     string
	Handle string

	GraceTime time.Duration

	State  string
	Events []string

	Limits LimitsSnapshot

	Resources ResourcesSnapshot

	Processes []ProcessSnapshot

	NetIns  []NetInSpec
	NetOuts []NetOutSpec
}

type LimitsSnapshot struct {
	Memory    *backend.MemoryLimits
	Disk      *backend.DiskLimits
	Bandwidth *backend.BandwidthLimits
	CPU       *backend.CPULimits
}

type ResourcesSnapshot struct {
	UID     uint32
	Network *network.Network
	Ports   []uint32
}

type ProcessSnapshot struct {
	ID uint32
}
