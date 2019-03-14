package scaffold

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/shurcooL/httpfs/vfsutil"
	log "github.com/sirupsen/logrus"
	tmpl "github.com/snowdrop/generator/pkg/template"
)

const (
	configYamlName      = "generator.yaml"
	dummyDirName        = "dummy"
	allVersionsSelector = "_all_"
)

var (
	config *Config

	// Store template in map by name then version. If there are no version for a given template name, then the template applies
	// to all versions.
	templates = make(templateRegistry)

	simplifiedVersionRegexp = regexp.MustCompile("^(\\d+.\\d+.\\d+)")
)

type versionRegistry map[string][]*template.Template
type templateRegistry map[string]versionRegistry

func (vr versionRegistry) getTemplatesFor(version string) []*template.Template {
	return vr[version]
}

func (vr versionRegistry) addTemplate(version, path string) error {
	templates := vr.getTemplatesFor(version)
	if templates == nil {
		templates = make([]*template.Template, 0, 20)
	}

	// Create a new Template using the File name as key and add it to the array
	t := template.New(path)

	// Read Template's content
	data, err := vfsutil.ReadFile(tmpl.Assets, path)
	if err != nil {
		return err
	}
	t, err = t.Parse(bytes.NewBuffer(data).String())
	if err != nil {
		return err
	}

	templates = append(templates, t)
	vr[version] = templates

	return nil
}

func (tr templateRegistry) getTemplatesFor(name, version string) []*template.Template {
	log.Infof("Retrieving templates for project template '%s' with version '%s'", name, version)

	// extract simplified Spring Boot version from project
	simplifiedVersion := allVersionsSelector
	matches := simplifiedVersionRegexp.FindStringSubmatch(version)
	if matches != nil {
		simplifiedVersion = matches[1]
	}

	// first check if we have templates for this version
	if versions, ok := tr[name]; ok {
		templates := versions.getTemplatesFor(simplifiedVersion)
		if templates == nil {
			log.Infof("No templates were found for '%s' (converted to simplified version: '%s'), attempting default version", version, simplifiedVersion)
			templates = versions.getTemplatesFor(allVersionsSelector)
		}

		return templates
	}

	return nil
}

func (tr templateRegistry) addTemplate(path string) error {
	// first, extract name and version from path
	name, version := extractNameAndVersion(path)

	// check if we already have a versions map for this template or create it otherwise
	versions, ok := templates[name]
	if !ok {
		versions = make(versionRegistry)
		templates[name] = versions
	}

	log.Infof("Adding template %s, version: %s, path: %s", name, version, path)
	return versions.addTemplate(version, path)
}

func extractNameAndVersion(path string) (name, version string) {
	split := strings.Split(path, string(filepath.Separator))
	name = split[1] // split[0] is empty because path starts with a separator
	potentialVersion := split[2]
	// check if the second hierarchy level match a version
	if simplifiedVersionRegexp.MatchString(potentialVersion) {
		version = potentialVersion
	} else {
		// otherwise, use the all version selector as version
		version = allVersionsSelector
	}

	return name, version
}

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

func CollectVfsTemplates() {
	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Printf("can't stat file %s: %v\n", path, err)
			return err
		}

		if fi.IsDir() {
			return nil
		}

		return templates.addTemplate(path)
	}

	errW := vfsutil.Walk(tmpl.Assets, "/", walkFn)
	if errW != nil {
		panic(errW)
	}
}

func ParseSelectedTemplate(project *Project, dir string, outDir string) error {
	templatesFor := templates.getTemplatesFor(project.Template, project.SpringBootVersion)
	if templatesFor == nil {
		return fmt.Errorf("'%s' template is not supported for '%s' Spring Boot version", project.Template, project.SpringBootVersion)
	}

	for _, t := range templatesFor {
		log.Infof("Processed template : %s", t.Name())
		var b bytes.Buffer

		// Enrich project with dependencies if they exist
		if strings.Contains(t.Name(), "pom.xml") {
			if project.Modules != nil {
				addDependenciesToModule(config.Modules, project)
			}
		}

		// Remove duplicate's dependencies from modules
		project.Dependencies = RemoveDuplicates(project.Modules)

		if log.GetLevel() == log.InfoLevel {
			for _, dep := range project.Dependencies {
				log.Infof("Dependency : %s-%s-%s", dep.GroupId, dep.ArtifactId, dep.Version)
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
		log.Debugf("Path converted: %s", pathF)

		// Convert FileName
		fileName := strings.Join([]string{dir, outDir, tFileName}, "/")
		log.Debugf("## File name : %s", fileName)
		fileNameConverted := strings.Replace(fileName, dummyDirName, convertPackageToPath(project.PackageName), -1)
		log.Debugf("File name converted : %s", fileNameConverted)

		// Create missing folders
		log.Debugf("Path to generated file : %s", pathConverted)
		os.MkdirAll(pathConverted, os.ModePerm)

		// Content generated
		log.Debugf("Content generated : %s", b.Bytes())

		err = ioutil.WriteFile(fileNameConverted, b.Bytes(), 0644)
		if err != nil {
			log.Error(err.Error())
		}
	}
	log.Infof("Enriched project %+v", project)
	return nil
}

func RemoveDuplicates(mods []Module) []Dependency {
	keys := make(map[string]bool)
	list := []Dependency{}
	for _, mod := range mods {
		for _, dep := range mod.Dependencies {
			gav := strings.Join([]string{dep.GroupId, dep.ArtifactId, dep.Version}, "-")
			if _, value := keys[gav]; !value {
				keys[gav] = true
				list = append(list, dep)
			}
		}
	}
	return list

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

func addDependenciesToModule(configModules []Module, project *Project) {
	for _, configModule := range configModules {
		for i, pModule := range project.Modules {
			if configModule.Name == pModule.Name {
				// check if the module is available for the project's requested BOM
				sbVersion := project.SpringBootVersion
				if configModule.IsAvailableFor(sbVersion) {
					log.Infof("Match found for project's module %s and modules %+v ", pModule.Name, configModule)
					project.Modules[i].Dependencies = configModule.Dependencies
					project.Modules[i].DependencyManagement = configModule.DependencyManagement
				} else {
					log.Infof("Ignoring module %s matching an existing module not available for SB version %s", pModule.Name, sbVersion)
				}
			}
		}
	}
}

func convertPackageToPath(p string) string {
	c := strings.Replace(p, ".", "/", -1)
	log.Debugf("Converted path : %s", c)
	return c
}
