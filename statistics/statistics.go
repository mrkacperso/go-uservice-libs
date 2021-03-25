package statistics

import (
	"time"
)

type StatisticLogger struct {
	Service string
}

type LogEntry struct {
	ExecutionDuration time.Duration
	Service           string
	Location          string
	Event             string
}

func (s StatisticLogger) GenerateStatisticLog(start, end time.Time, Location, Event string) LogEntry {
	return LogEntry{
		ExecutionDuration: end.Sub(start),
		Service:           s.Service,
		Location:          Location,
		Event:             Event,
	}
}