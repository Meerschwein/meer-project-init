package types

import (
	"fmt"
	"os"

	"golang.org/x/exp/slices"
)

type ToCopy struct {
	Source  string
	Targets []string
}

type Profile struct {
	Name     string
	Base_dir string

	Commands_before []string
	Using           []string
	Paths           []ToCopy
	Commands_after  []string
}

type Config struct {
	Profiles         []Profile
	Selected_profile string
	Cli_arguments    CliArguments
}

func (config Config) Validate() error {
	all_profiles_names := []string{}
	for _, p := range config.Profiles {
		all_profiles_names = append(all_profiles_names, p.Name)
	}

	errors := Errors{}

	if !slices.Contains(all_profiles_names, config.Selected_profile) {
		errors.Appendf("Profile %s is selected but it was not found", config.Selected_profile)
	}

	for _, profile := range config.Profiles {
		// Check if all the profiles this profile is using exist
		for _, needs := range profile.Using {
			if !slices.Contains(all_profiles_names, needs) {
				errors.Appendf("Profile %s is using profile %s but it was not found", profile.Name, needs)
			}
		}

		// Check if all sources exist
		for _, to_copy := range profile.Paths {
			_, err := os.Stat(to_copy.Source)
			if err != nil {
				errors.Appendf("Profile %s wants to copy the file %s but it was not found", profile.Name, to_copy.Source)
			}
		}
	}

	_, err := os.Stat(config.Cli_arguments.Init_directory_path)
	if err == nil {
		errors.Append("Target exists")
	}

	if errors.Len() > 0 {
		return errors
	} else {
		return nil
	}
}

func (config Config) Get_profile(name string) Profile {
	for _, profile := range config.Profiles {
		if profile.Name == name {
			return profile
		}
	}

	panic(fmt.Sprintf("unknown profile %s\nconfig %+v\n", name, config))
}

func (config Config) Execute() (err error) {
	exec := Get_executor(config.Cli_arguments)

	selected_profile := config.Get_profile(config.Selected_profile)

	return Execute_profile(exec, config, selected_profile)
}
