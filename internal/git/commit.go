package git

import (
	"fmt"
	"os/exec"
)

// Commit runs `git commit`
func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = exec.Command("git", "status").Stdout
	cmd.Stderr = exec.Command("git", "status").Stderr
	return cmd.Run()
}
