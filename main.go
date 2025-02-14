package main

import (
	"chronicl/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.RootCmd.Execute())
}

