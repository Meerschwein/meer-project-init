package main

import (
	"flag"
	"fmt"

	"github.com/Meerschwein/meer-project-init/pkg/types"
	"github.com/adrg/xdg"
	"rsc.io/getopt"
)

var (
	config_file    = ""
	init_directory = ""
	base_directory = ""
	dry_run        = false
)

func init() {
	flag.StringVar(
		&config_file,
		"config-file",
		xdg.ConfigHome+"/meer-project-init/config.toml",
		"Path to the config file.",
	)

	flag.StringVar(
		&init_directory,
		"init-directory",
		"./",
		"Path to the directory to initialize the project in.",
	)

	flag.StringVar(
		&base_directory,
		"base-directory",
		xdg.ConfigHome+"/meer-project-init",
		"Path to the directory from where the files will be copied.",
	)

	flag.BoolVar(
		&dry_run,
		"dry-run",
		false,
		"",
	)

	getopt.Alias("d", "init-directory")

	getopt.Parse()
}

func main() {
	args := types.CliArguments{
		Config_file_path:    config_file,
		Base_directory_path: base_directory,
		Init_directory_path: init_directory,
		Dry_run:             dry_run,
		Args:                flag.Args(),
	}

	config, err := args.Parse()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = config.Execute()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
