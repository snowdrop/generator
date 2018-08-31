package scaffold

type Project struct {
	GroupId            string
	ArtifactId         string
	Version            string
	PackageName        string
	Dependencies	   []string
	OutDir             string
	Template 		   string

	SnowdropBomVersion string
	SpringVersion      string
	Modules            []Module
	Starters		   []Starter

	UrlService  	   string
}

type Config struct {
	Template     []string  `yaml:"templates"`
	Bom          []Bom     `yaml:"boms"`
	Module       []Module  `yaml"module"`
}

type Bom struct {
	Community string `yaml:"community version"`
	Snowdrop  string `yaml:"snowdrop version"`
}

type Module struct {
	Name	     string
	Description  string
	Starters     []Starter
}

type Starter struct {
	GroupId	     string
	ArtifactId	 string
	Scope	     string
}
