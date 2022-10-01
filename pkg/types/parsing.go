package types

import "fmt"

type ConfigParsing struct {
	Profiles []map[string]ProfileParsing `mapstructure:"profile"`
}

type ProfileParsing struct {
	Base_dir        string
	Using           []string
	Commands_before []string
	Commands_after  []string
	Paths           [][]string
}

func (parsed ConfigParsing) To_config(args CliArguments) (config Config, err error) {
	config.CliArguments = args

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
					Source: path_to_source(config, profile, paths[0]),
				}

				// if there is only a source the target has the same filename as the source
				if len(paths) == 1 {
					to_copy.Targets = append(to_copy.Targets, path_to_target(config, paths[0]))
				}

				for _, target := range paths[1:] {
					to_copy.Targets = append(to_copy.Targets, path_to_target(config, target))
				}

				profile.Paths = append(profile.Paths, to_copy)
			}

			config.Profiles = append(config.Profiles, profile)
		}
	}

	return
}

func path_to_source(config Config, profile Profile, source string) string {
	return fmt.Sprintf("%s/%s/%s", config.Path_to_base_directory, profile.Base_dir, source)
}

func path_to_target(config Config, target string) string {
	return fmt.Sprintf("%s/%s", config.project_directory(), target)
}
