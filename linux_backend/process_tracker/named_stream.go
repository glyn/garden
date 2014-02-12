package process_tracker

import (
	"github.com/pivotal-cf-experimental/garden/backend"
)

type namedStream struct {
	process *Process
	source  backend.ProcessStreamSource
}

func newNamedStream(process *Process, source backend.ProcessStreamSource) *namedStream {
	return &namedStream{
		process: process,
		source:  source,
	}
}

func (s *namedStream) Write(data []byte) (int, error) {
	s.process.sendToStreams(backend.ProcessStream{
		Source: s.source,
		Data:   data,
	})

	return len(data), nil
}
