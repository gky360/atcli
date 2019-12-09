package cmd

import (
	"testing"
)

func TestRootCmdExecute(t *testing.T) {
	if err := rootCmd.Execute(); err != nil {
		t.Error(err)
		t.Error("Failed to execute root command")
	}
}
