package utils

import (
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func TestTaskSourceFilePath(t *testing.T) {
	_, err := TaskSourceFilePath("a")
	if err == nil {
		t.Errorf("TaskSourceFilePath(\"a\") without Contest.id set is expected to fail: %v", err)
	}

	viper.Set("contest.id", "practice")
	home, _ := homedir.Dir()

	actual, err := TaskSourceFilePath("a")
	expected := filepath.Join(home, "atcoder", "practice", "a", "Main.cpp")
	if err != nil || actual != expected {
		t.Errorf("TaskSourceFilePath(\"a\") returned a wrong path: %v", actual)
	}
}
