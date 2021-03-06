package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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

type (
	TemplateData struct {
		Contest *models.Contest
		Task    *models.Task
	}
)

func DefaultRootDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return filepath.Join(home, "atcoder")
}

func RootDir() string {
	rootDir := viper.GetString("root")
	if rootDir == "" {
		rootDir = DefaultRootDir()
	}
	return rootDir
}

func ContestDir() (string, error) {
	contestID := viper.GetString("contest.id")
	if contestID == "" {
		return "", fmt.Errorf(MsgContestIDRequired)
	}
	return filepath.Join(RootDir(), contestID), nil
}

func TaskDir(taskName string) (string, error) {
	if taskName == "" {
		return "", fmt.Errorf(MsgTaskNameRequired)
	}
	contestDir, err := ContestDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(contestDir, strings.ToLower(taskName)), nil
}

func TaskSourceFilePath(taskName string) (string, error) {
	if taskName == "" {
		return "", fmt.Errorf(MsgTaskNameRequired)
	}
	taskDir, err := TaskDir(taskName)
	if err != nil {
		return "", err
	}
	return filepath.Join(taskDir, "Main.cpp"), nil
}

func TaskSampleDir(taskName string, isForTestcases bool) (string, error) {
	if taskName == "" {
		return "", fmt.Errorf(MsgTaskNameRequired)
	}
	taskDir, err := TaskDir(taskName)
	if err != nil {
		return "", err
	}
	dirname := "samples"
	if isForTestcases {
		dirname = "testcases"
	}
	return filepath.Join(taskDir, dirname), nil
}

func TaskSampleInDir(taskName string, isForTestcases bool) (string, error) {
	return taskSampleDir(taskName, isForTestcases, true)
}

func TaskSampleOutDir(taskName string, isForTestcases bool) (string, error) {
	return taskSampleDir(taskName, isForTestcases, false)
}

func taskSampleDir(taskName string, isForTestcases bool, isInput bool) (string, error) {
	taskSampleDir, err := TaskSampleDir(taskName, isForTestcases)
	if err != nil {
		return "", err
	}

	var name string
	if isInput {
		name = "in"
	} else {
		name = "out"
	}

	return filepath.Join(taskSampleDir, name), nil
}

func TaskInputFilePath(taskName string, sampleName string, isForTestcases bool) (string, error) {
	return taskFilePath(taskName, sampleName, isForTestcases, true)
}

func TaskOutputFilePath(taskName string, sampleName string, isForTestcases bool) (string, error) {
	return taskFilePath(taskName, sampleName, isForTestcases, false)
}

func taskFilePath(taskName string, sampleName string, isForTestcases bool, isInput bool) (string, error) {
	if taskName == "" {
		return "", fmt.Errorf(MsgTaskNameRequired)
	}
	dir, err := taskSampleDir(taskName, isForTestcases, isInput)
	if err != nil {
		return "", err
	}

	var filePath string
	if isForTestcases {
		filePath = filepath.Join(dir, sampleName)
	} else {
		filePath = filepath.Join(dir, sampleName+".txt")

	}
	return filePath, nil
}

func GetSampleNames(taskName string, isForTestcases bool) ([]string, error) {
	taskSampleInDir, err := TaskSampleInDir(taskName, isForTestcases)
	if err != nil {
		return nil, err
	}
	if _, err = os.Stat(taskSampleInDir); err != nil {
		return nil, err
	}

	var pat string
	if isForTestcases {
		pat = filepath.Join(taskSampleInDir, "*")
	} else {
		pat = filepath.Join(taskSampleInDir, "*.txt")
	}
	g, _ := filepath.Glob(pat)
	sampleNames := make([]string, len(g))
	for i, fpath := range g {
		if isForTestcases {
			sampleNames[i] = filepath.Base(fpath)
		} else {
			sampleNames[i] = strings.TrimSuffix(filepath.Base(fpath), ".txt")
		}
	}
	return sampleNames, nil
}

func CreateDirsForTask(task *models.Task) error {
	if task.Name == "" {
		return fmt.Errorf(MsgTaskNameRequired)
	}
	taskSampleInDir, err := TaskSampleInDir(task.Name, false)
	if err != nil {
		return err
	}
	taskSampleOutDir, err := TaskSampleOutDir(task.Name, false)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(taskSampleInDir, 0755); err != nil {
		return err
	}
	if err = os.MkdirAll(taskSampleOutDir, 0755); err != nil {
		return err
	}
	return nil
}

func CreateSourceFile(contest *models.Contest, task *models.Task) error {
	if err := CreateDirsForTask(task); err != nil {
		return err
	}
	taskSourceFilePath, err := TaskSourceFilePath(task.Name)
	if err != nil {
		return err
	}
	taskSrouceFilePathRel, err := filepath.Rel(RootDir(), taskSourceFilePath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(taskSourceFilePath); err == nil {
		// Already exists
		fmt.Printf("Already exists: %s\n", taskSrouceFilePathRel)
		return nil
	}

	cppTemplatePath := viper.GetString("cppTemplatePath")
	if cppTemplatePath == "" {
		// Create empty source file
		_, err = os.OpenFile(taskSourceFilePath, os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
	} else {
		// Create source file from template
		tmplData := TemplateData{
			Contest: contest,
			Task:    task,
		}
		t, err := template.ParseFiles(cppTemplatePath)
		if err != nil {
			return err
		}
		f, err := os.Create(taskSourceFilePath)
		if err != nil {
			return err
		}
		err = t.Execute(f, tmplData)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Created file  : %s\n", taskSrouceFilePathRel)
	return nil
}

func CreateSampleFiles(task *models.Task) error {
	for _, sample := range task.Samples {
		sampleName := fmt.Sprintf("sample%02d", sample.Num)
		taskInputFilePath, err := TaskInputFilePath(task.Name, sampleName, false)
		if err != nil {
			return err
		}
		taskOutputFilePath, err := TaskOutputFilePath(task.Name, sampleName, false)
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
			filePathRel, err := filepath.Rel(RootDir(), filePath)
			if err != nil {
				return err
			}
			fmt.Printf("Created file  : %s\n", filePathRel)
		}
	}
	return nil
}

func CreateFilesForTask(contest *models.Contest, task *models.Task) error {
	if err := CreateSourceFile(contest, task); err != nil {
		return err
	}
	if err := CreateSampleFiles(task); err != nil {
		return err
	}
	return nil
}

func CreateFilesForTasks(contest *models.Contest, tasks []*models.Task) error {
	for _, task := range tasks {
		if err := CreateFilesForTask(contest, task); err != nil {
			return err
		}
	}
	return nil
}
