package scaffold

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/shurcooL/httpfs/vfsutil"
	log "github.com/sirupsen/logrus"

	tmpl "github.com/snowdrop/generator/pkg/template"
)

const (
	configYamlName = "generator.yaml"
	dummyDirName   = "dummy"
)

var (
	templateFiles []string
	config        Config

	assetsJavaTemplates = tmpl.Assets
	templates           = make(map[string]template.Template)
)

func GetConfig() *Config {
	return &config
}

//Create a new Default Project using default values
func NewDefaultProject() *Project {
	springBootBomVersion, snowdropBomVersion := GetDefaultBOM()
	return &Project{
		GroupId: "com.example",
		ArtifactId: "demo",
		Version: "0.0.1-SNAPSHOT",
		SnowdropBomVersion: snowdropBomVersion,
		SpringBootVersion: springBootBomVersion,
		Template: "simple",
	}
}

func GetDefaultBOM() (string, string) {
	cfg := GetConfig()
	for _, bom := range cfg.Boms {
		if bom.Default {
			return bom.Community, bom.Snowdrop
		}
	}
	return "", ""
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

func CollectVfsTemplates() {

	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Printf("can't stat file %s: %v\n", path, err)
			return nil
		}

		if fi.IsDir() {
			return nil
		}

		log.Debug("Path of the file to be added as template : " + path)
		templateFiles = append(templateFiles, path)
		return nil
	}

	errW := vfsutil.Walk(assetsJavaTemplates, "/", walkFn)
	if errW != nil {
		panic(errW)
	}

	for i := range templateFiles {
		log.Info("File template : " + templateFiles[i])

		// Create a new Template using the File name as key and add it to the array
		t := template.New(templateFiles[i])

		// Read Template's content
		data, err := vfsutil.ReadFile(assetsJavaTemplates, templateFiles[i])
		if err != nil {
			log.Error(err)
		}
		t, err = t.Parse(bytes.NewBuffer(data).String())
		if err != nil {
			log.Error(err)
		}
		templates[templateFiles[i]] = *t
	}
}

func ParseTemplateSelected(templateSelected string, dir string, outDir string, project Project) {

	// Pickup from the Map of the Templates, the files corresponding to the type selected by the user
	for key, t := range templates {
		if strings.HasPrefix(key,"/" + templateSelected) {

			log.Infof("Template processed : %s", t.Name())
			var b bytes.Buffer

			// Enrich project with dependencies if they exist
			if strings.Contains(t.Name(), "pom.xml") {
				if project.Modules != nil {
					addDependenciesToModule(config.Modules, &project)
				}
			}

			// Remove duplicate's dependencies from modules
			project.Dependencies = RemoveDuplicates(project.Modules)

			if log.GetLevel() == log.InfoLevel {
				for _, dep := range project.Dependencies {
				 		log.Info("Dependency : %s-%s-$s", dep.GroupId, dep.GroupId, dep.Version)
				}
			}

			// Use template to generate the content
			err := t.Execute(&b, project)
			if err != nil {
				log.Error(err.Error())
			}

			// Convert Path
			tFileName := t.Name()
			pathF := strings.Join([]string{dir, outDir, path.Dir(tFileName)}, "/")
			log.Debugf("## Path : %s", pathF)
			pathConverted := strings.Replace(pathF, dummyDirName, convertPackageToPath(project.PackageName), -1)
			log.Debugf("Path converted: ", pathF)

			// Convert FileName
			fileName := strings.Join([]string{dir, outDir, tFileName}, "/")
			log.Debugf("## File name : %s", fileName)
			fileNameConverted := strings.Replace(fileName, dummyDirName, convertPackageToPath(project.PackageName), -1)
			log.Debugf("File name converted : ", fileNameConverted)

			// Create missing folders
			log.Debugf("Path to generated file : ", pathConverted)
			os.MkdirAll(pathConverted, os.ModePerm)

			// Content generated
			log.Debugf("Content generated : %s", b.Bytes())

			err = ioutil.WriteFile(fileNameConverted, b.Bytes(), 0644)
			if err != nil {
				log.Error(err.Error())
			}

		}
	}
	log.Infof("Project enriched %+v ", project)
}

func addDependenciesToModule(configModules []Module, project *Project) {
	for _, configModule := range configModules {
		for i, pModule := range project.Modules {
			if configModule.Name == pModule.Name {
				log.Infof("Match found for project's module %s and modules %+v ", pModule.Name, configModule)
				project.Modules[i].Dependencies = configModule.Dependencies
			}
		}
	}
}

func RemoveDuplicates(mods []Module) []Dependency{
	keys := make(map[string]bool)
	list := []Dependency{}
	for _, mod := range mods {
		for _, dep := range mod.Dependencies {
			gav := strings.Join([]string{dep.GroupId,dep.ArtifactId,dep.Version},"-")
			if _, value := keys[gav]; !value {
				keys[gav] = true
				list = append(list, dep)
			}
		}
	}
	return list

}

func convertPackageToPath(p string) string {
	c := strings.Replace(p, ".", "/", -1)
	log.Debugf("Converted path : ", c)
	return c
}
