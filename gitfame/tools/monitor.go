package tools

import (
	"sync/atomic"
)

type Monitor struct {
	readDirJobs   atomic.Int64
	parseFileJobs atomic.Int64

	finish chan struct{}
}

func NewMonitor() *Monitor {
	return &Monitor{
		finish:        make(chan struct{}, 1),
		readDirJobs:   atomic.Int64{},
		parseFileJobs: atomic.Int64{},
	}
}

func (m *Monitor) checkFinished() {
	if m.readDirJobs.Load() == 0 && m.parseFileJobs.Load() == 0 {
		m.finish <- struct{}{}
	}
}

func (m *Monitor) UpdateReadDirJobs(delta int64) {
	m.readDirJobs.Add(delta)
	m.checkFinished()
}

func (m *Monitor) UpdateParseFileJobs(delta int64) {
	m.parseFileJobs.Add(delta)
	m.checkFinished()
}

func (m *Monitor) Finish() chan struct{} {
	return m.finish
}
