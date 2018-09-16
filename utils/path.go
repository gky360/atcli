package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gky360/atsrv/models"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	MsgContestIDRequired = "Contest id is required. Try '--help' option for help."
	MsgTaskNameRequired  = "Task name is required. Try '--help' option for help."
	MsgSbmIDRequired     = "Submission id is required. Try '--help' option for help."
	MsgSbmSourceRequired = "Submission source is required. Try '--help' option for help."
)

func DefaultRootPath() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return filepath.Join(home, "atcoder")
}

func RootPath() string {
	rootPath := viper.GetString("root")
	if rootPath == "" {
		rootPath = DefaultRootPath()
	}
	return rootPath
}

func ContestPath() (string, error) {
	contestID := viper.GetString("contest.id")
	if contestID == "" {
		return "", fmt.Errorf(MsgContestIDRequired)
	}
	return filepath.Join(RootPath(), contestID), nil
}

func TaskPath(taskName string) (string, error) {
	if taskName == "" {
		return "", fmt.Errorf(MsgTaskNameRequired)
	}
	contestPath, err := ContestPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(contestPath, strings.ToLower(taskName)), nil
}

func TaskSourceFilePath(taskName string) (string, error) {
	if taskName == "" {
		return "", fmt.Errorf(MsgTaskNameRequired)
	}
	taskPath, err := TaskPath(taskName)
	if err != nil {
		return "", err
	}
	return filepath.Join(taskPath, "Main.cpp"), nil
}

func TaskSamplePath(taskName string) (string, error) {
	if taskName == "" {
		return "", fmt.Errorf(MsgTaskNameRequired)
	}
	taskPath, err := TaskPath(taskName)
	if err != nil {
		return "", err
	}
	return filepath.Join(taskPath, "samples"), nil
}

func TaskInputFilePath(taskName string, sampleNum int) (string, error) {
	if taskName == "" {
		return "", fmt.Errorf(MsgTaskNameRequired)
	}
	taskSamplePath, err := TaskSamplePath(taskName)
	if err != nil {
		return "", err
	}
	return filepath.Join(taskSamplePath, fmt.Sprintf("%02d.in.txt", sampleNum)), nil
}

func TaskOutputFilePath(taskName string, sampleNum int) (string, error) {
	if taskName == "" {
		return "", fmt.Errorf(MsgTaskNameRequired)
	}
	taskSamplePath, err := TaskSamplePath(taskName)
	if err != nil {
		return "", err
	}
	return filepath.Join(taskSamplePath, fmt.Sprintf("%02d.out.txt", sampleNum)), nil
}

func CreateDirsForTask(task *models.Task) error {
	if task.Name == "" {
		return fmt.Errorf(MsgTaskNameRequired)
	}
	taskSamplePath, err := TaskSamplePath(task.Name)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(taskSamplePath, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func CreateSourceFile(task *models.Task) error {
	if err := CreateDirsForTask(task); err != nil {
		return err
	}
	taskSourceFilePath, err := TaskSourceFilePath(task.Name)
	if err != nil {
		return err
	}
	_, err = os.OpenFile(taskSourceFilePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	taskSrouceFilePathRel, err := filepath.Rel(RootPath(), taskSourceFilePath)
	if err != nil {
		return err
	}
	fmt.Printf("Created file: %s\n", taskSrouceFilePathRel)
	return nil
}

func CreateSampleFiles(task *models.Task) error {
	for _, sample := range task.Samples {
		taskInputFilePath, err := TaskInputFilePath(task.Name, sample.Num)
		if err != nil {
			return err
		}
		taskOutputFilePath, err := TaskOutputFilePath(task.Name, sample.Num)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(taskInputFilePath, []byte(sample.Input), 0644); err != nil {
			return err
		}
		if err := ioutil.WriteFile(taskOutputFilePath, []byte(sample.Output), 0644); err != nil {
			return err
		}

		for _, filePath := range []string{taskInputFilePath, taskOutputFilePath} {
			filePathRel, err := filepath.Rel(RootPath(), filePath)
			if err != nil {
				return err
			}
			fmt.Printf("Created file: %s\n", filePathRel)
		}
	}
	return nil
}

func CreateFilesForTask(task *models.Task) error {
	if err := CreateSourceFile(task); err != nil {
		return err
	}
	if err := CreateSampleFiles(task); err != nil {
		return err
	}
	return nil
}

func CreateFilesForTasks(tasks []*models.Task) error {
	for _, task := range tasks {
		if err := CreateFilesForTask(task); err != nil {
			return err
		}
	}
	return nil
}
