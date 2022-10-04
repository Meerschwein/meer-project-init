package types

type Replacements struct {
	Name string
}

func get_replacements(config Config) Replacements {
	return Replacements{
		Name: config.Project_name,
	}
}
