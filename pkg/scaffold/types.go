package scaffold

import (
	"github.com/blang/semver"
	"github.com/sirupsen/logrus"
	"strings"
)

const releaseSuffix = ".RELEASE"

type Project struct {
	GroupId     string `yaml:"groupid"           json:"groupid"`
	ArtifactId  string `yaml:"artifactid"        json:"artifactid"`
	Version     string `yaml:"version"           json:"version"`
	PackageName string `yaml:"packagename"       json:"packagename"`
	OutDir      string `yaml:"outdir"            json:"outdir"`
	Template    string `yaml:"template"          json:"template"`

	SnowdropBomVersion string `yaml:"snowdropbom"       json:"snowdropbom"`
	SpringBootVersion  string `yaml:"springbootversion" json:"springbootversion"`
	UseSupported       bool

	Modules         []Module `yaml:"modules"           json:"modules"`
	Dependencies    []Dependency
	ExtraProperties ExtraProperties

	UrlService string `yaml:"urlservice"           json:"urlservice"`
}

func (p *Project) HasDekorate() bool {
	for _, mod := range p.Modules {
		if "dekorate" == mod.Name {
			return true
		}
	}
	return false
}

type Config struct {
	Templates       []Template      `yaml:"templates"    json:"templates"`
	Boms            []Bom           `yaml:"bomversions"  json:"bomversions"`
	Modules         []Module        `yaml:"modules"      json:"modules"`
	ExtraProperties ExtraProperties `yaml:"extraProperties"      json:"extraProperties"`
}

func (c *Config) GetDefaultBom() Bom {
	for _, bom := range c.Boms {
		if bom.Default {
			return bom
		}
	}

	return Bom{}
}

func (c *Config) GetCorrespondingSnowDropBom(version string) Bom {
	for _, bom := range c.Boms {
		if bom.Community == version {
			return bom
		}
	}
	return Bom{}
}

func (c *Config) doesBOMExistFor(version string) bool {
	// Add .RELEASE if it's not present since it's expected in the configuration
	i := strings.Index(version, releaseSuffix)
	if i < 0 {
		version = version + releaseSuffix
	}

	return c.GetCorrespondingSnowDropBom(version).Community == version
}

func (c *Config) GetModulesCompatibleWith(version string) []Module {
	if c.doesBOMExistFor(version) {
		return keepModulesCompatibleWith(c.Modules, version)
	}
	return []Module{}
}

func keepModulesCompatibleWith(modules []Module, version string) []Module {
	compatible := make([]Module, 0, len(modules))
	for _, module := range modules {
		if module.IsAvailableFor(version) {
			compatible = append(compatible, module)
		}
	}
	return compatible
}

type Template struct {
	Name        string `yaml:"name"                     json:"name"`
	Description string `yaml:"description"              json:"description"`
}

type Bom struct {
	Community string `yaml:"community" json:"community"`
	Snowdrop  string `yaml:"snowdrop"  json:"snowdrop"`
	Supported string `yaml:"supported"  json:"supported"`
	Default   bool   `yaml:"default"  json:"default"`
}

type ExtraProperties struct {
	DekorateVersion string `json:"dekorateVersion"`
}

type Module struct {
	Name                 string                 `yaml:"name"                     json:"name"`
	Description          string                 `yaml:"description"              json:"description"`
	Guide                string                 `yaml:"guide_ref"                json:"guide_ref"`
	Dependencies         []Dependency           `yaml:"dependencies"             json:"dependencies"`
	DependencyManagement []DependencyManagement `yaml:"dependencymanagement"     json:"dependencymanagement"`
	Tags                 []string               `yaml:"tags"                     json:"tags"`
	Availability         string                 `yaml:"availability,omitempty"   json:"availability,omitempty"`
}

func (m Module) IsAvailableFor(bomVersion string) bool {
	// remove .RELEASE from BOM version if present since it's not part of semantic versioning
	i := strings.Index(bomVersion, releaseSuffix)
	if i > 0 {
		bomVersion = bomVersion[:i]
	}

	// if provided version is incorrect, module should not be available
	sbVersion, err := semver.Parse(bomVersion)
	if err != nil {
		logrus.Warningf("Invalid input version %s, marking module as unavailable: %v", bomVersion, err)
		return false
	}

	if len(m.Availability) != 0 {
		versionRange, err := semver.ParseRange(m.Availability)
		if err != nil {
			logrus.Warningf("Invalid availability range %s, marking module as unavailable: %v", m.Availability, err)
			return false
		}

		return versionRange(sbVersion)
	}
	return true
}

type DependencyManagement struct {
	Dependencies []Dependency `yaml:"dependencies,omitempty"     json:"dependencies,omitempty"`
}

type Dependency struct {
	GroupId    string `yaml:"groupid,omitempty"           json:"groupid,omitempty"`
	ArtifactId string `yaml:"artifactid,omitempty"        json:"artifactid,omitempty"`
	Scope      string `yaml:"scope,omitempty"   json:"scope,omitempty"`
	Version    string `yaml:"version,omitempty" json:"version,omitempty"`
	Type       string `yaml:"type,omitempty"    json:"type,omitempty"`
}
