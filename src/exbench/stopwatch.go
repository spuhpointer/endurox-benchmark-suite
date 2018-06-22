/*
 * @brief Simple stopwatch implementation
 *
 * @file stopwatch.go
 */
package exbench

import (
	"time"
)

//Get UTC milliseconds since epoch
//@return epoch milliseconds
func GetEpochMillis() float64 {
	now := time.Now()
	nanos := now.UnixNano()
	millis := float64(float64(nanos) / float64(1000000))

	return millis
}

//About incoming & outgoing messages:
type StopWatch struct {
	start float64 //Timestamp messag sent
}

//Reset the stopwatch
func (s *StopWatch) Reset() {
	s.start = GetEpochMillis()
}

//Get delta milliseconds
//@return time spent in milliseconds
func (s *StopWatch) GetDeltaMillis() float64 {
	return GetEpochMillis() - s.start
}

//Get delta seconds of the stopwatch
//@return return seconds spent
func (s *StopWatch) GetDetlaSec() float64 {
	return s.GetDeltaMillis() / float64(1000)
}
