package utils

import (
	"flag"
	"testing"
)

var term = flag.Bool("terminal", false, "Is a terminal available?")

func TestDetectTerminal(t *testing.T) {

	if DetectTerminal() == true {
		if *term != true {
			t.Error("[ERROR] terminal flag is false, but a terminal was detected")
		}
	} else {
		if *term == true {
			t.Error("[ERROR] terminal flag is true, but no terminal was detected")
		}
	}
}
