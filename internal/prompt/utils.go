package prompt

import (
	"errors"
	"fmt"
	"github.com/cqroot/prompt"
	"os"
)

func Check(err error) {
	if err != nil {
		if errors.Is(err, prompt.ErrUserQuit) {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		} else {
			panic(err)
		}
	}
}
