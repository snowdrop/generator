package scaffold

import (
	"bytes"
	"encoding/json"
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
	configDirName  = "config"
	configYamlName = "starters.yaml"
	dummyDirName   = "dummy"
)

var (
	templateFiles []string
	config        Config

	assetsJavaTemplates = tmpl.Assets
	templates           = make(map[string]template.Template)
)

func NewDefaultScaffoldProject() *Project {
	return &Project{
		GroupId: "com.example",
		ArtifactId: "demo",
		Version: "0.0.1-SNAPSHOT",
		SnowdropBomVersion: "1.5.15.Final",
		SpringVersion: "1.5.15.RELEASE",
	}
}

func ParseStartersConfigFile(pathTemplateDir string) {
	if pathTemplateDir == "" {
		pathTemplateDir = "../scaffold"
	}
	startersPath := strings.Join([]string{pathTemplateDir, configDirName, configYamlName}, "/")
	log.Infof("Parsing Starters's Config at %s", startersPath)

	// Read file and parse it to create a Config's type
	if _, err := os.Stat(startersPath); err == nil {
		source, err := ioutil.ReadFile(startersPath)
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
		log.Debug("Starters's config")
		log.Debug("--------------------")
		s, _ := json.Marshal(&config)
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

			// Enrich project with starters dependencies if they exist
			if strings.Contains(t.Name(), "pom.xml") {
				if project.Dependencies != nil {
					project = convertDependencyToModule(project.Dependencies, config.Modules, project)
				}
			}

			// Remove Starter duplicates
			RemoveDuplicates(&project.Starters)

			// log.Debug("Remove duplicates")
			// for _, starter := range project.Starters {
			//  		log.Info("No duplicate Starter : ", starter.ArtifactId)
			// }

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

func convertDependencyToModule(deps []string, modules []Module, p Project) Project {
	for _, dep := range deps {
		for _, module := range modules {
			if module.Name == dep {
				log.Infof("Match found for dep %s and starters %+v ", dep, module)
				p.Modules = append(p.Modules, module)
				for _, starter := range module.Starters {
					p.Starters = append(p.Starters,starter)
				}
			}
		}
	}
	return p
}

func RemoveDuplicates(starters *[]Starter) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *starters {
		if !found[x.ArtifactId] {
			found[x.ArtifactId] = true
			(*starters)[j] = (*starters)[i]
			j++
		}
	}
	*starters = (*starters)[:j]
}

func convertPackageToPath(p string) string {
	c := strings.Replace(p, ".", "/", -1)
	log.Debugf("Converted path : ", c)
	return c
}
