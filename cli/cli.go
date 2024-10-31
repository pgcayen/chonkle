package cli

import (
	"fmt"
	"os"
)

// indication constants
const (
	Above     = "↑"
	Below     = "↓"
	ToRight   = "→"
	ToLeft    = "←"
	Correct   = "✓"
	Incorrect = "x"
)

// print to clear terminal
const (
	Clear = "\033[H\033[2J"
	Exit  = "exit"
)

func PressKeyToContinue() {
	fmt.Println("Press any key to continue!")
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
}
