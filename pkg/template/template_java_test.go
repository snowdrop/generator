package template

import (
	"fmt"
	"os"

	"github.com/shurcooL/httpfs/vfsutil"
	"net/http"
	"testing"
)

var (
	templateFiles []string
	project       = "custom"
)

func TestVfsSimpleJavaProject(t *testing.T) {
	tExpectedFiles := []string{
		"custom/pom.xml",
		"custom/src/main/java/dummy/DemoApplication.java",
		"custom/src/main/resources/application-kubernetes.properties",
	}

	tFiles := walkTree()

	for i := range tExpectedFiles {
		if tExpectedFiles[i] != tFiles[i] {
			t.Errorf("Template was incorrect, got: '%s', want: '%s'.", tFiles[i], tExpectedFiles[i])
		}
	}
}

func walkTree() []string {
	var fs http.FileSystem = Assets

	vfsutil.Walk(fs, project, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("can't stat file %s: %v\n", path, err)
			return nil
		}

		if fi.IsDir() {
			return nil
		}

		fmt.Println(path)
		templateFiles = append(templateFiles, path)
		return nil
	})
	return templateFiles
}
