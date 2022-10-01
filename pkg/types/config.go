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
	CliArguments
	Profiles []Profile
}

func (config Config) project_directory() string {
	return config.Project_name
}

func (config Config) Validate() error {
	all_profiles_names := []string{}
	for _, p := range config.Profiles {
		all_profiles_names = append(all_profiles_names, p.Name)
	}

	errors := Errors{}

	if !slices.Contains(all_profiles_names, config.Selected_profile) {
		errors.appendf("Profile %s is selected but it was not found", config.Selected_profile)
	}

	for _, profile := range config.Profiles {
		// Check if all the profiles this profile is using exist
		for _, needs := range profile.Using {
			if !slices.Contains(all_profiles_names, needs) {
				errors.appendf("Profile %s is using profile %s but it was not found", profile.Name, needs)
			}
		}

		// Check if all sources exist
		for _, to_copy := range profile.Paths {
			_, err := os.Stat(to_copy.Source)
			if err != nil {
				errors.appendf("Profile %s wants to copy the file %s but it was not found", profile.Name, to_copy.Source)
			}
		}
	}

	// FIXME check if a directory with the project name exists
	// if it does error out

	if errors.len() > 0 {
		return errors
	} else {
		return nil
	}
}

func (config Config) selected_profile() Profile {
	return config.Get_profile(config.Selected_profile)
}

func (config Config) Get_profile(name string) Profile {
	for _, profile := range config.Profiles {
		if profile.Name == name {
			return profile
		}
	}

	panic(fmt.Sprintf("unknown profile %s\nconfig %+v\n", config.Selected_profile, config))
}

func (config Config) Execute() (err error) {
	return execute_profile(
		get_executor(config),
		config,
		config.selected_profile(),
	)
}
