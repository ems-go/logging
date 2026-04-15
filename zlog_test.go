package logging

import (
	"fmt"
	"testing"
)

func TestZlog(t *testing.T) {
	Debugf("test:%v", fmt.Errorf("hello world"))
}
