package server

import (
	"archive/zip"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/snowdrop/generator/pkg/scaffold"

	"encoding/json"
	"github.com/snowdrop/generator/pkg/common/logger"
	"net/url"
)

var (
	currentDir, _ = os.Getwd()
	port          = "8000"
	pathConfigMap = ""
	tmpDirName    = "_temp"
	letterRunes   = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func init() {
	// Enable Debug level mode if ENV LOG_LEVEL=debug is defined
	logger.EnableLogLevelDebug()
	log.Print("Log level : ", log.GetLevel())

	// Check env vars
	s := os.Getenv("SERVER_PORT")
	if s != "" {
		port = s
	}

	cm := os.Getenv("CONFIGMAP_PATH")
	if cm != "" {
		pathConfigMap = cm
	}

	// Parse Generator Config YAML file to load :
	// - Templates available : crud, rest, sumple, ...
	// - Different Snowdrop/Community BOMs
	// - Modules and their dependencies associated / the starters
	scaffold.ParseGeneratorConfigFile(pathConfigMap)
	scaffold.CreateDefaultProject()

	// Create the Go Templates from the Spring Boot template directory (crud, web, simple, ....)
	scaffold.CollectVfsTemplates()

	rand.Seed(time.Now().UnixNano())
}

func Run(version string, gitcommit string) {
	log.Infof("Starting Spring Boot Generator Server on port %s - Version %s (%s)", port, version, gitcommit)
	log.Infof("The following REST endpoints are available : ")
	log.Infof("Generate zip : %s", "/app")
	log.Infof("Config : %s", "/config")

	router := mux.NewRouter()
	router.HandleFunc("/app", CreateZipFile).Methods("GET")
	router.HandleFunc("/config", PopulateJSONConfig).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getUrlVal(r *http.Request, k string) string {
	return r.URL.Query().Get(k)
}

func getArrayVal(r *http.Request, k string, params map[string][]string) []string {
	return params[k]
}

func convertArrayToStruct(modules []string) []scaffold.Module {
	mod := make([]scaffold.Module, 0)
	for _, e := range modules {
		mod = append(mod, scaffold.Module{Name: e})
	}
	return mod
}

func getArrayModuleVal(r *http.Request, k string, params map[string][]string) []scaffold.Module {
	return convertArrayToStruct(params[k])
}

//Process the HTTP Get request to return as JSON message the Generator config
func PopulateJSONConfig(w http.ResponseWriter, r *http.Request) {
	// Set CORS Headers
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-Requested-With, remember-me, X-CSRF-Token, Authorization")
	}
	jsonStr, _ := json.Marshal(scaffold.GetConfig())
	fmt.Fprintf(w, "%s", jsonStr)
}

//Process the HTTP GET Raw Request and populate a zip file as HTTP Response
func CreateZipFile(w http.ResponseWriter, r *http.Request) {

	params, _ := url.ParseQuery(r.URL.RawQuery)
	p := scaffold.GetDefaultProject()

	if getUrlVal(r, "template") != "" {
		p.Template = getUrlVal(r, "template")
	}
	if getUrlVal(r, "groupid") != "" {
		p.GroupId = getUrlVal(r, "groupid")
	}
	if getUrlVal(r, "artifactid") != "" {
		p.ArtifactId = getUrlVal(r, "artifactid")
	}
	if getUrlVal(r, "version") != "" {
		p.Version = getUrlVal(r, "version")
	}
	if getUrlVal(r, "packagename") != "" {
		p.PackageName = getUrlVal(r, "packagename")
	}
	if len(getArrayModuleVal(r, "module", params)) > 0 {
		p.Modules = getArrayModuleVal(r, "module", params)
	}
	if getUrlVal(r, "snowdropbom") != "" {
		p.SnowdropBomVersion = getUrlVal(r, "snowdropbom")
	}
	if getUrlVal(r, "springbootversion") != "" {
		p.SpringBootVersion = getUrlVal(r, "springbootversion")
	}
	if getUrlVal(r, "outdir") != "" {
		p.OutDir = getUrlVal(r, "outdir")
	}

	// If the snowdropbom version is not defined BUT only the Spring Boot Version, then get the corresponding
	// BOM version using the version of the Spring Boot selected from the Config Bom's Array
	if getUrlVal(r, "snowdropbom") == "" && getUrlVal(r, "springbootversion") != "" {
		p.SnowdropBomVersion = scaffold.GetCorrespondingSnowDropBom(p.SpringBootVersion)
	}

	// As dependencies and template selection can't be used together, we force the template to be equal to "simple"
	// when a user selects a different template. This is because we would like to avoid to populate a project with starters
	// which are incompatible or not fully tested with the template proposed
	if len(getArrayVal(r, "module", params)) > 0 && p.Template != "simple" {
		p.Template = "simple"
	}

	log.Info("Project : ", p)
	log.Infof("Request received : %s", r.URL)

	// Generate a random temp directory where populated files will be saved
	tmpdir := strings.Join([]string{tmpDirName, randStringRunes(10)}, "/")
	log.Infof("Temp dir %s", tmpdir)

	// Parse the java project's template selected and enrich the scaffold.Project with the dependencies (if they are)
	scaffold.ParseTemplateSelected(p.Template, currentDir, tmpdir, p)
	log.Info("Project generated")

	zipDir := strings.Join([]string{tmpdir, p.Template, "/"}, "/")
	handleZip(w, zipDir)
	log.Info("Zip populated")

	// Remove temp dir where project has been generated
	removeTempDir(tmpdir)
}

func removeTempDir(tmpdir string) {
	err := os.RemoveAll(strings.Join([]string{currentDir, tmpdir}, "/"))
	if err != nil {
		log.Error(err.Error())
	}
}

// Generate Zip file to be returned as HTTP Response
func handleZip(w http.ResponseWriter, tmpdir string) {
	zipFilename := "demo.zip"
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", zipFilename))

	errZip := zipFiles(w, tmpdir)
	if errZip != nil {
		log.Fatal(errZip)
	}
}

// Get Files generated from templates under _temp directory and
// them recursively to the file to be zipped
func zipFiles(w http.ResponseWriter, tmpdir string) error {
	fullPathZipDir := strings.Join([]string{currentDir, tmpdir}, "/")
	log.Info("Zip file path : ", fullPathZipDir)
	err := recursiveZip(w, fullPathZipDir)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func recursiveZip(w http.ResponseWriter, destinationPath string) error {
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	err := filepath.Walk(destinationPath, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(destinationPath))
		relPath = strings.TrimPrefix(relPath, "/")
		log.Debugf("relPath calculated : ", relPath)

		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = zipWriter.Close()
	if err != nil {
		return err
	}
	return nil
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
