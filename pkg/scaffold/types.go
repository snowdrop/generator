package scaffold

type Project struct {
	GroupId            string       `yaml:"groupid"        json:"groupid"`
	ArtifactId         string       `yaml:"artifactid"     json:"artifactid"`
	Version            string       `yaml:"version"        json:"version"`
	PackageName        string       `yaml:"packagename"    json:"packagename"`
	Dependencies	   []string     `yaml:"dependencies"   json:"dependencies"`
	OutDir             string       `yaml:"outdir"         json:"outdir"`
	Template 		   string       `yaml:"template"       json:"template"`

	SnowdropBomVersion string       `yaml:"snowdropbom"   json:"snowdropbom"`
	SpringVersion      string       `yaml:"springversion" json:"springversion"`
	Modules            []Module     `yaml:"modules"        json:"modules"`
	Starters	       []Dependency `yaml:"dependency"    json:"dependency"`

	UrlService  	   string    `yaml:"urlservice"    json:"urlservice"`
}

type Config struct {
	Templates     []string  `yaml:"templates"   json:"templates"`
	Boms          []Bom     `yaml:"boms"        json:"boms"`
	Modules       []Module  `yaml:"modules"    json:"modules"`
}

type Bom struct {
	Community string `yaml:"community" json:"community"`
	Snowdrop  string `yaml:"snowdrop"  json:"snowdrop"`
}

type Module struct {
	Name	      string       `yaml:"name"             json:"name"`
	Description   string       `yaml:"description"      json:"description"`
	Guide         string       `yaml:"guide_ref"        json:"guide_ref"`
	Dependencies  []Dependency `yaml:"dependencies"     json:"dependencies"`
	tags		  []string     `yaml:"tags"           json:"tags"`
}

type Dependency struct {
	GroupId	     string `yaml:"groupid"           json:"groupid"`
	ArtifactId	 string `yaml:"artifactid"        json:"artifactid"`
	Scope	     string `yaml:"scope"             json:"scope"`
	Version      string `yaml:"version,omitempty" json:"version,omitempty"`
}
