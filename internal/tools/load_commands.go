package tools

import (
	"os"

	"gopkg.in/yaml.v2"
)

// type CommandsList struct {
type Commands struct {
	CommandList []Command `yaml:"commands"`
}

type Command struct {
	Name             string `yaml:"name"`
	ServiceTypeID    int    `yaml:"ServiceTypeID"`
	MessageSubtypeID int    `yaml:"MessageSubtypeID"`
	TM               string `yaml:"TM"`
}

// `yaml:"commands"`
// }

func load_yaml_to_struct(yaml_path string, struct_data *Commands) *Commands {
	data, err := os.ReadFile(yaml_path)

	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, struct_data); err != nil {
		panic(err)
	}

	// print the fields to the console
	return struct_data
}

func Load_configs(path string) Commands {
	var commands Commands
	load_yaml_to_struct(path, &commands)
	return commands
}
