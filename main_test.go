package loggy

import (
	"flag"
	"testing"
)

var updateGolden = flag.Bool("update", false, "update golden files")

func TestMain(m *testing.M) {
	flag.Parse()
	m.Run()
}
