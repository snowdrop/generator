package scaffold

import (
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

const (
	configYamlName = "generator.yaml"
)

var (
	config *Config
)

func GetConfig() *Config {
	return config
}

func GetDefaultProject() *Project {
	springBootBomVersion, snowdropBomVersion := getDefaultBOM()
	p := &Project{
		GroupId:            "com.example",
		ArtifactId:         "demo",
		PackageName:        "com.example.demo",
		Version:            "0.0.1-SNAPSHOT",
		SnowdropBomVersion: snowdropBomVersion,
		SpringBootVersion:  springBootBomVersion,
		Template:           "custom",
	}
	p.ExtraProperties = GetConfig().ExtraProperties
	return p
}

func ParseGeneratorConfigFile(pathConfigMap string) {

	configPath := strings.Join([]string{pathConfigMap, configYamlName}, "/")
	log.Infof("Parsing Generator's Config at %s", configPath)

	// Read file and parse it to create a Config's type
	if _, err := os.Stat(configPath); err == nil {
		source, err := ioutil.ReadFile(configPath)
		if err != nil {
			log.Fatal(err.Error())
		}

		err = yaml.Unmarshal(source, &config)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		log.Fatal("No Starters's config file detected !!!")
	}

	if log.GetLevel() == log.DebugLevel {
		log.Debug("-------------------")
		log.Debug("Generator's config")
		log.Debug("-------------------")
		s, _ := yaml.Marshal(&config)
		log.Debug(string(s))
	}
}

func getDefaultBOM() (string, string) {
	cfg := GetConfig()
	for _, bom := range cfg.Boms {
		if bom.Default {
			return bom.Community, bom.Snowdrop
		}
	}
	return "", ""
}
