/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"os"
	"path"

	"github.com/jvzantvoort/tmux-project/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	verbose       bool
	cfgFile       string
	OutputDirOpt  string
	OutputDir     string
	OutputFileOpt string
	OutputFile    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tmux-project",
	Short: messages.GetShort("root"),
	Long:  messages.GetLong("root"),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")

	// Setup logging
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

}

// GetConfigParamName wrapper to get subsection of the configfile
// func GetConfigParamName(instr string) string {
// 	return fmt.Sprintf("%s.%s", "main", instr)
// }

func PrincipalConfigDir() string {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	return path.Join(home, ".config", ApplicationName)
}

func EnvironmentConfigDir(configdir string) string {
	option, ok := os.LookupEnv(EnvConfigDir)
	if ok {
		return option
	}
	return configdir
}

func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in home directory with name ".devb" (without extension).
		viper.AddConfigPath(EnvironmentConfigDir(PrincipalConfigDir()))
		viper.SetConfigType("yaml")
		viper.SetConfigName("config.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("Error: %v", err)
	}

	viper.SetDefault("logdir", DefaultLogDir)
	viper.SetDefault("filename", DefaultLogName)

	// Bind them to environment vars
	viper.BindEnv("logdir", EnvLogDir)
	viper.BindEnv("filename", EnvLogName)

	// print errors if not set
	if !viper.IsSet("logdir") {
		log.Errorf("Option %s (%s) is not set", "logdir", EnvLogDir)
	}

	if !viper.IsSet("filename") {
		log.Errorf("Option %s (%s) is not set", "filename", EnvLogName)
	}

	// get the options
	OutputDir = viper.GetString("logdir")
	OutputFile = viper.GetString("filename")
}
