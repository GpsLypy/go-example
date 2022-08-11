package main

import (
	"github.com/GpyLypy/go-example/basicUse/context_test/withTimeCTX/flight"
	"github.com/GpyLypy/go-example/basicUse/context_test/withTimeCTX/publish"
)

func main() {
	position := flight.NewPosition(22.1, 34.2, 1000)
	ph := publish.NewPublishHandler()
	ph.PublishPosition(position)
}
