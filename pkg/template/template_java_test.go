package template

import (
	"os"
	"fmt"

	"github.com/shurcooL/httpfs/vfsutil"
	"testing"
	"net/http"
	"github.com/snowdrop/generator/pkg/template"
)

var (
	templateFiles   []string
	project         = "simple"
)

func TestVfsSimpleJavaProject(t *testing.T) {
	tExpectedFiles := []string {
		"simple/pom.xml",
		"simple/src/main/java/dummy/DemoApplication.java",
		"simple/src/main/resources/application.properties",
	}

	tFiles := walkTree()

	for i := range tExpectedFiles {
		if tExpectedFiles[i] != tFiles[i] {
			t.Errorf("Template was incorrect, got: '%s', want: '%s'.", tFiles[i], tExpectedFiles[i])
		}
	}
}

func walkTree() []string {
	var fs http.FileSystem = template.Assets

	vfsutil.Walk(fs, project, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("can't stat file %s: %v\n", path, err)
			return nil
		}

		if fi.IsDir() {
			return nil
		}

		fmt.Println(path)
		templateFiles = append(templateFiles,path)
		return nil
	})
	return templateFiles
}