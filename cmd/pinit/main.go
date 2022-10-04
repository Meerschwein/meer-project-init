package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Meerschwein/meer-project-init/pkg/types"
	"github.com/adrg/xdg"
	"github.com/spf13/viper"
	"rsc.io/getopt"
)

var (
	path_to_config_file    = ""
	path_to_base_directory = ""
	selected_profile       = ""
	dry_run                = false

	project_name = ""
)

func init() {
	flag.StringVar(
		&path_to_config_file,
		"config-file",
		xdg.ConfigHome+"/pinit/config.toml",
		"Path to the config file.",
	)

	flag.StringVar(
		&path_to_base_directory,
		"base-directory",
		xdg.ConfigHome+"/pinit",
		"Path to the directory from where the files will be copied.",
	)

	flag.StringVar(
		&selected_profile,
		"profile",
		"",
		"Profile to be executed.",
	)
	getopt.Alias("p", "profile")

	flag.BoolVar(
		&dry_run,
		"dry-run",
		false,
		"Prints the actions that would be taken without actually taking them.",
	)

	getopt.Parse()

	// Check if mandatory flags are set and the the correct number
	// of arguments is specified
	ExitIfError(validate_flags())

	project_name = flag.Arg(0)
}

func ExitIfError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func validate_flags() error {
	if len(flag.Args()) == 0 {
		return fmt.Errorf("you need to specify a project name")
	} else if len(flag.Args()) > 1 {
		return fmt.Errorf("you may only specify one project name")
	}

	if selected_profile == "" {
		return fmt.Errorf("you need to specify a profile")
	}

	return nil
}

func main() {
	args := types.CliArguments{
		Path_to_config_file:    path_to_config_file,
		Path_to_base_directory: path_to_base_directory,
		Selected_profile:       selected_profile,
		Dry_run:                dry_run,
		Project_name:           project_name,
	}

	viper.SetConfigFile(args.Path_to_config_file)

	ExitIfError(viper.ReadInConfig())

	var config_parsing types.ConfigParsing
	err := viper.Unmarshal(&config_parsing)
	ExitIfError(err)

	config, err := config_parsing.To_config(args)
	ExitIfError(err)

	config.CliArguments = args

	ExitIfError(config.Validate())

	ExitIfError(config.Execute())
}
