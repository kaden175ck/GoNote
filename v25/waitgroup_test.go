package v25_test

import (
	v25 "dqq/go/basic/v25"
	"testing"
)

func TestWaitGroup(t *testing.T) {
	v25.Sum()
}

// go test -v ./v25 -run=^TestWaitGroup$ -count=1
