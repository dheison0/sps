package types

type Root struct {
	Port int `toml:"port"`
}

type FilterConfig struct {
	File        string `toml:"file"`
	EnableRegex bool   `toml:"enable_regex"`
	LessMemory  bool   `toml:"less_memory"`
}

type Config struct {
	Main       Root         `toml:"main"`
	Filter     FilterConfig `toml:"filter"`
	ConfigFile string
}
