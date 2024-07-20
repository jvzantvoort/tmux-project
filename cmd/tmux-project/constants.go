package main

const (
	ApplicationName string = "tmux-project"
	ConfigDirName   string = "tmux-project"
	// ConstDirMode permissions on newly created configuration directories
	ConstDirMode int = 0755
	// ConstFileMode permissions on newly created configuration files
	ConstFileMode int = 0644

	// EnvConfigDir environment pointing to the main config dir
	EnvConfigDir string = "TMUX_PROJECT_CONFIG_DIR"

	EnvLogDir  string = "TMUX_PROJECT_OUTPUT_DIR"
	EnvLogName string = "TMUX_PROJECT_OUTPUT_NAME"

	DefaultLogDir string = "~/Logs"
	DefaultLogName string = "tmux.log"
)
