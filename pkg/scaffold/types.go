package scaffold

type Project struct {
	GroupId     string `yaml:"groupid"           json:"groupid"`
	ArtifactId  string `yaml:"artifactid"        json:"artifactid"`
	Version     string `yaml:"version"           json:"version"`
	PackageName string `yaml:"packagename"       json:"packagename"`
	OutDir      string `yaml:"outdir"            json:"outdir"`
	Template    string `yaml:"template"          json:"template"`

	SnowdropBomVersion string `yaml:"snowdropbom"       json:"snowdropbom"`
	SpringBootVersion  string `yaml:"springbootversion" json:"springbootversion"`

	Modules      []Module `yaml:"modules"           json:"modules"`
	Dependencies []Dependency

	UrlService string `yaml:"urlservice"           json:"urlservice"`
}

type Config struct {
	Templates []Template `yaml:"templates"    json:"templates"`
	Boms      []Bom      `yaml:"bomversions"  json:"bomversions"`
	Modules   []Module   `yaml:"modules"      json:"modules"`
}

type Template struct {
	Name        string `yaml:"name"                     json:"name"`
	Description string `yaml:"description"              json:"description"`
}

type Bom struct {
	Community string `yaml:"community" json:"community"`
	Snowdrop  string `yaml:"snowdrop"  json:"snowdrop"`
	Default   bool   `yaml:"default"  json:"default"`
}

type Module struct {
	Name                 string                 `yaml:"name"                     json:"name"`
	Description          string                 `yaml:"description"              json:"description"`
	Guide                string                 `yaml:"guide_ref"                json:"guide_ref"`
	Dependencies         []Dependency           `yaml:"dependencies"             json:"dependencies"`
	DependencyManagement []DependencyManagement `yaml:"dependencymanagement"     json:"dependencymanagement"`
	Tags                 []string               `yaml:"tags"                     json:"tags"`
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
