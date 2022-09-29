package types

import (
	"fmt"

	"github.com/spf13/viper"
)

type CliArguments struct {
	Config_file_path    string
	Base_directory_path string
	Init_directory_path string

	Dry_run bool

	Args []string
}

func (args CliArguments) Parse() (config Config, err error) {
	viper.SetConfigFile(args.Config_file_path)

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	var config_parsing ConfigParsing
	err = viper.Unmarshal(&config_parsing)
	if err != nil {
		return
	}

	config, err = config_parsing.To_config(args)
	if err != nil {
		fmt.Println("error parsing")
		return
	}

	config.Cli_arguments = args

	err = config.Validate()
	if err != nil {
		fmt.Println("error vailadting")
	}

	return
}

type ProfileParsing struct {
	Base_dir        string
	Using           []string
	Commands_before []string
	Commands_after  []string
	Paths           [][]string
}

type ConfigParsing struct {
	Profiles []map[string]ProfileParsing `mapstructure:"profile"`
}

func (parsed ConfigParsing) To_config(args CliArguments) (config Config, err error) {
	if len(args.Args) == 0 {
		return Config{}, fmt.Errorf("no profile selected")
	} else if len(args.Args) > 1 {
		return Config{}, fmt.Errorf("too many profiles selected")
	}

	config.Selected_profile = args.Args[0]

	for _, profiles := range parsed.Profiles {
		for profile_name, profile_properties := range profiles {
			profile := Profile{
				Name:            profile_name,
				Base_dir:        profile_properties.Base_dir,
				Using:           profile_properties.Using,
				Commands_before: profile_properties.Commands_before,
				Commands_after:  profile_properties.Commands_after,
			}

			if profile.Base_dir == "" {
				profile.Base_dir = profile_name
			}

			for _, paths := range profile_properties.Paths {
				to_copy := ToCopy{
					Source: path_to_source(args, profile, paths[0]),
				}

				// if there is only a source the target has the same filename as the source
				if len(paths) == 1 {
					to_copy.Targets = append(to_copy.Targets, path_to_target(args, paths[0]))
				}

				for _, target := range paths[1:] {
					to_copy.Targets = append(to_copy.Targets, path_to_target(args, target))
				}

				profile.Paths = append(profile.Paths, to_copy)
			}

			config.Profiles = append(config.Profiles, profile)
		}
	}

	return
}

func path_to_source(args CliArguments, profile Profile, source string) string {
	return fmt.Sprintf("%s/%s/%s", args.Base_directory_path, profile.Base_dir, source)
}

func path_to_target(args CliArguments, target string) string {
	return fmt.Sprintf("%s/%s", args.Init_directory_path, target)
}
