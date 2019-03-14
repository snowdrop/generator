package server

import (
	"archive/zip"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
	if len(s) != 0 {
		port = s
	}

	cm := os.Getenv("CONFIGMAP_PATH")
	if len(cm) != 0 {
		pathConfigMap = cm
	} else {
		pathConfigMap = filepath.Join(currentDir, "conf")
	}

	// Parse Generator Config YAML file to load :
	// - Templates available : crud, rest, simple, ...
	// - Different Snowdrop/Community BOMs
	// - Modules and their dependencies associated / the starters
	scaffold.ParseGeneratorConfigFile(pathConfigMap)

	// Create the Go Templates from the Spring Boot template directory (crud, web, custom, ....)
	scaffold.CollectVfsTemplates()

	rand.Seed(time.Now().UnixNano())
}

func Run(version string, gitcommit string) {
	router := mux.NewRouter()
	router.HandleFunc("/app", CreateZipFile).Methods("GET").Name("Generate zip")
	router.HandleFunc("/modules/{version}", modulesFor).Methods("GET").Name("Get modules compatible with Spring Boot version")
	router.HandleFunc("/config", PopulateJSONConfig).Methods("GET").Name("Config")

	log.Infof("Starting Spring Boot Generator Server on port %s - Version %s (%s)", port, version, gitcommit)
	log.Infof("The following REST endpoints are available : ")
	_ = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		log.Infof("%s: %s", route.GetName(), path)
		return nil
	})

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func asModuleArray(modules []string) []scaffold.Module {
	mod := make([]scaffold.Module, 0)
	hasSeenCore := false
	for _, e := range modules {
		if "core" == e {
			hasSeenCore = true
		}
		mod = append(mod, scaffold.Module{Name: e})
	}

	// if we don't have core in the modules, add it because otherwise the apps won't work
	if !hasSeenCore {
		mod = append(mod, scaffold.Module{Name: "core"})
	}

	return mod
}

func modulesFor(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(r, w)
	vars := mux.Vars(r)
	version := vars["version"]
	modules := []scaffold.Module{}
	if len(version) != 0 {
		config := scaffold.GetConfig()
		modules = config.GetModulesCompatibleWith(version)
	}
	jsonStr, _ := json.Marshal(modules)
	fmt.Fprintf(w, "%s", jsonStr)
}

//Process the HTTP Get request to return as JSON message the Generator config
func PopulateJSONConfig(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(r, w)
	w.Header().Set("Content-Type", "application/json")
	jsonStr, _ := json.Marshal(scaffold.GetConfig())
	fmt.Fprintf(w, "%s", jsonStr)
}

// setCORSHeaders sets CORS Headers on the response if the Origin header exists on the request
func setCORSHeaders(r *http.Request, w http.ResponseWriter) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-Requested-With, remember-me, X-CSRF-Token, Authorization")
	}
}

//Process the HTTP GET Raw Request and populate a zip file as HTTP Response
func CreateZipFile(w http.ResponseWriter, r *http.Request) {

	params, _ := url.ParseQuery(r.URL.RawQuery)
	p := scaffold.GetDefaultProject()

	p.SpringBootVersion = params.Get("springbootversion")
	p.SnowdropBomVersion = params.Get("snowdropbom")

	config := scaffold.GetConfig()
	if len(p.SpringBootVersion) == 0 {
		// if we didn't get
		defaultSBVersion := config.GetDefaultBom().Community
		if len(defaultSBVersion) == 0 {
			respondWith("Must provide at least Spring Boot version", http.StatusBadRequest, w)
			return
		}
		p.SpringBootVersion = defaultSBVersion
	}

	// retrieve bom information associated with the Spring Boot version
	bom := config.GetCorrespondingSnowDropBom(p.SpringBootVersion)

	// If the snowdrop bom version is not defined BUT only the Spring Boot Version, then get the corresponding
	// BOM version using the version of the Spring Boot selected from the Config Bom's Array
	if len(p.SnowdropBomVersion) == 0 {
		p.SnowdropBomVersion = bom.Snowdrop
	}

	b, err := strconv.ParseBool(params.Get("supported"))
	if err == nil {
		p.UseSupported = b
	}
	if p.UseSupported {
		if len(bom.Supported) == 0 {
			respondWith(fmt.Sprintf("%s is not a supported Spring Boot version", p.SpringBootVersion), http.StatusBadRequest, w)
			return
		}
		p.SnowdropBomVersion = bom.Supported
	}

	p.Template = params.Get("template")
	p.GroupId = params.Get("groupid")
	p.ArtifactId = params.Get("artifactid")
	p.Version = params.Get("version")
	p.PackageName = params.Get("packagename")
	p.OutDir = params.Get("outdir")

	if len(params["module"]) > 0 {
		p.Modules = asModuleArray(params["module"])
	}

	// As dependencies and template selection can't be used together, we force the template to be equal to "custom"
	// when a user selects a different template. This is because we would like to avoid to populate a project with starters
	// which are incompatible or not fully tested with the template proposed
	if len(params["module"]) > 0 && p.Template != "custom" {
		p.Template = "custom"
	}

	log.Info("Project : ", p)
	log.Infof("Request received : %s", r.URL)

	// Generate a random temp directory where populated files will be saved
	tmpdir := strings.Join([]string{tmpDirName, randStringRunes(10)}, "/")
	log.Infof("Temp dir %s", tmpdir)

	// Parse the java project's template selected and enrich the scaffold.Project with the dependencies (if they are)
	version, err := scaffold.ParseSelectedTemplate(p, currentDir, tmpdir)
	if err != nil {
		respondWith(err.Error(), http.StatusNotFound, w)
		return
	}

	log.Info("Project generated")

	zipDir := filepath.Join(tmpdir, p.Template)
	handleZip(zipDir, version, w)
	log.Info("Zip populated")

	// Remove temp dir where project has been generated
	removeTempDir(tmpdir)
}

func respondWith(msg string, status int, w http.ResponseWriter) {
	log.Info(msg)
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

func removeTempDir(tmpdir string) {
	err := os.RemoveAll(strings.Join([]string{currentDir, tmpdir}, "/"))
	if err != nil {
		log.Error(err.Error())
	}
}

// Generate Zip file to be returned as HTTP Response
func handleZip(tmpdir string, version string, w http.ResponseWriter) {
	zipFilename := "demo.zip"
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", zipFilename))

	fullPathZipDir := filepath.Join(currentDir, tmpdir)
	log.Info("Zip file path : ", fullPathZipDir)
	errZip := recursiveZip(fullPathZipDir, version, w)
	if errZip != nil {
		respondWith(errZip.Error(), http.StatusInternalServerError, w)
	}
}

// Get Files generated from templates under _temp directory and
// them recursively to the file to be zipped
func recursiveZip(destinationPath string, version string, w http.ResponseWriter) error {
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	err := filepath.Walk(destinationPath, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, destinationPath)
		relPath = strings.TrimPrefix(relPath, "/")

		if len(version) > 0 {
			relPath = strings.Replace(relPath, version+"/", "", -1)
		}
		log.Debugf("calculated rel path: %s", relPath)

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
