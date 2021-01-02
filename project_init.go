package tmuxproject

import (
	"github.com/jvzantvoort/tmux-project/config"
)

// func yamlStringSettings() string {
// 	c := viper.AllSettings()
// 	bs, err := yaml.Marshal(c)
// 	if err != nil {
// 		log.Fatalf("unable to marshal config to YAML: %v", err)
// 	}
// 	return string(bs)
// }

func CreateProjectType(projecttype string) error {
	var config config.ProjectTypeConfig
	config.Init(mainconfig.ProjTypeConfigDir, projecttype)

	// viper.SetConfigType("yaml")
	// viper.SetConfigName("config")
	// viper.AddConfigPath(projtypeconfigdir)

	// configData := viper.AllSettings()
	// configDataYAML, err := yaml.Marshal(configData)
	// viper.WriteConfig()
	return nil
}
