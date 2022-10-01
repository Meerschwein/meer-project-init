package types

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
)

type Executor interface {
	Execute_command(command string) error
	Copy_file(source, target string) error
	Create_dir(path string)
}

func get_executor(config Config) Executor {
	if config.Dry_run {
		return LogExecutor{}
	} else {
		return ActualExecutor{
			Execution_directory: config.project_directory(),
		}
	}
}

type LogExecutor struct{}

func (LogExecutor) Execute_command(command string) error {
	fmt.Printf("Execute command: %s\n", command)
	return nil
}

func (LogExecutor) Copy_file(source, target string) error {
	fmt.Printf("Copying file from %s to %s\n", source, target)
	return nil
}

func (LogExecutor) Create_dir(path string) {
	fmt.Printf("Creating directory: %s\n", path)
}

type ActualExecutor struct {
	Execution_directory string
}

func (executor ActualExecutor) Execute_command(command string) error {
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = executor.Execution_directory
	out, err := cmd.Output()
	if string(out) != "" {
		fmt.Println(string(out))
	}
	return err
}

func (ActualExecutor) Copy_file(source, target string) error {
	fin, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	defer fin.Close()

	os.MkdirAll(path.Dir(target), os.ModePerm)
	fout, err := os.Create(target)
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (ActualExecutor) Create_dir(path string) {
	os.MkdirAll(path, os.ModePerm)
}

func execute_profile(exec Executor, config Config, profile Profile) error {
	exec.Create_dir(config.project_directory())

	for _, command := range profile.Commands_before {
		err := exec.Execute_command(command)
		if err != nil {
			return fmt.Errorf("error executing command %s", command)
		}
	}

	for _, using := range profile.Using {
		err := execute_profile(exec, config, config.Get_profile(using))
		if err != nil {
			return err
		}
	}

	for _, to_copy := range profile.Paths {
		for _, target := range to_copy.Targets {
			err := exec.Copy_file(to_copy.Source, target)
			if err != nil {
				return err
			}
		}
	}

	for _, command := range profile.Commands_after {
		err := exec.Execute_command(command)
		if err != nil {
			return err
		}
	}

	return nil
}
