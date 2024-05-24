package main

import (
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
)

var pw progress.Writer

func newTracker(message string) (work func(v int64), done func()) {
	tracker := &progress.Tracker{
		ExpectedDuration: 1 * time.Second,
		Message:          message,
		Total:            1000,
		Units:            progress.UnitsDefault,
	}

	pw.AppendTracker(tracker)

	return func(v int64) { tracker.Increment(v) },
		func() { tracker.MarkAsDone() }
}

func setupProgress() {
	if pw != nil {
		pw.Stop()
		time.Sleep(time.Millisecond)
		go pw.Render()
		return
	}

	pw = progress.NewWriter()
	pw.SetAutoStop(false)
	pw.SetMessageLength(40)
	// pw.SetStyle(progress.StyleBlocks)
	pw.SetTrackerLength(40)
	pw.SetTrackerPosition(progress.PositionLeft)
	pw.Style().Chars = progress.StyleCharsBlocks
	pw.Style().Colors = progress.StyleColorsExample
	pw.Style().Name = "flower"
	pw.Style().Options.DoneString = "✓"
	pw.Style().Options.ErrorString = "✗"
	pw.Style().Options.Separator = "   "
	pw.Style().Options.TimeDonePrecision = time.Millisecond
	pw.Style().Options.TimeOverallPrecision = time.Millisecond
	pw.Style().Visibility.ETA = false
	pw.Style().Visibility.ETAOverall = false
	pw.Style().Visibility.Speed = false
	pw.Style().Visibility.SpeedOverall = false
	pw.Style().Visibility.Time = true
	pw.Style().Visibility.Value = false
	pw.SetUpdateFrequency(time.Millisecond * 5)
	go pw.Render()
}
