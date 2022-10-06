package types

import (
	"fmt"
	"html/template"
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
		return ActualExecutor{config}
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
	config Config
}

func (executor ActualExecutor) Execute_command(command string) error {
	replacement := get_replacements(executor.config)
	actual_command, err := replacement.in_string(command)
	if err != nil {
		return err
	}

	cmd := exec.Command("bash", "-c", actual_command)
	cmd.Dir = executor.config.project_directory()
	out, err := cmd.Output()
	if string(out) != "" {
		fmt.Println(string(out))
	}

	return err
}

func (executor ActualExecutor) Copy_file(source, target string) error {
	replacement := get_replacements(executor.config)

	actual_target, err := replacement.in_string(target)
	if err != nil {
		return err
	}

	os.MkdirAll(path.Dir(actual_target), os.ModePerm)

	target_file, err := os.Create(actual_target)
	if err != nil {
		return err
	}
	defer target_file.Close()

	err = template.Must(template.New("").ParseFiles(source)).ExecuteTemplate(target_file, path.Base(source), replacement)

	return err
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