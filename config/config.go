package config

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type arrayFlags []string

func (i arrayFlags) String() string {
	var ts []string
	for _, el := range i {
		ts = append(ts, el)
	}
	return strings.Join(ts, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// Config speichert die gesamten Flags und Umgebungsvariablen
type Config struct {
	TfPath           string
	WorkingDir       string
	Name             string
	ZipFileName      string
	IPPort           string
	PlanPath         string
	PlanJSONPath     string
	WorkspaceName    string
	TFCOrgName       string
	TFCWorkspaceName string
	Standalone       bool
	ShowSensitive    bool
	GenImage         bool
	GetVersion       string
	TFCNewRun        bool
	TfVarsFiles      arrayFlags
	TfVars           arrayFlags
	TfBackendConfigs arrayFlags
	Version          string // Konstante für Version
}

// Lade Konfiguration aus Flags
func LoadConfig() *Config {
	config := &Config{}

	log.Println("Loading configuration...")
	// Definiere Flags
	flag.StringVar(&config.TfPath, "tfPath", "/usr/local/bin/terraform", "Path to Terraform binary")
	flag.StringVar(&config.WorkingDir, "workingDir", ".", "Path to Terraform configuration")
	flag.StringVar(&config.Name, "name", "rover", "Configuration name")
	flag.StringVar(&config.ZipFileName, "zipFileName", "rover.zip", "Standalone zip file name")
	flag.StringVar(&config.IPPort, "ipPort", "0.0.0.0:9000", "IP and port for Rover server")
	flag.StringVar(&config.PlanPath, "planPath", "", "Plan file path")
	flag.StringVar(&config.PlanJSONPath, "planJSONPath", "", "Plan JSON file path")
	flag.StringVar(&config.WorkspaceName, "workspaceName", "", "Workspace name")
	flag.StringVar(&config.TFCOrgName, "tfcOrg", "", "Terraform Cloud Organization name")
	flag.StringVar(&config.TFCWorkspaceName, "tfcWorkspace", "", "Terraform Cloud Workspace name")
	flag.StringVar(&config.GetVersion, "version", "0.3.3", "Get current version")
	flag.BoolVar(&config.Standalone, "standalone", false, "Generate standalone HTML files")
	flag.BoolVar(&config.ShowSensitive, "showSensitive", false, "Display sensitive values")
	flag.BoolVar(&config.TFCNewRun, "tfcNewRun", false, "Create new Terraform Cloud run")
	flag.BoolVar(&config.GenImage, "genImage", false, "Generate graph image")

	var tfVarsFiles, tfVars, tfBackendConfigs arrayFlags
	flag.Var(&tfVarsFiles, "tfVarsFile", "Path to *.tfvars files")
	flag.Var(&tfVars, "tfVar", "Terraform variable (key=value)")
	flag.Var(&tfBackendConfigs, "tfBackendConfig", "Path to *.tfbackend files")
	flag.Parse()

	// Lade Flags in Config
	config.TfVarsFiles = tfVarsFiles
	config.TfVars = tfVars
	config.TfBackendConfigs = tfBackendConfigs

	path, err := os.Getwd()
	if err != nil {
		panic(errors.New("unable to get current working directory"))
	}

	// Ensure PlanPath is absolute
	if config.PlanPath != "" {
		if !filepath.IsAbs(config.PlanPath) {
			config.PlanPath = filepath.Join(path, config.PlanPath)
		}
		absPlanPath, err := filepath.Abs(config.PlanPath)
		if err != nil {
			panic(errors.New("unable to get absolute path for PlanPath"))
		}
		config.PlanPath = absPlanPath
	}

	// Ensure PlanJSONPath is absolute
	if config.PlanJSONPath != "" {
		if !filepath.IsAbs(config.PlanJSONPath) {
			config.PlanJSONPath = filepath.Join(path, config.PlanJSONPath)
		}
		absPlanJSONPath, err := filepath.Abs(config.PlanJSONPath)
		if err != nil {
			panic(errors.New("unable to get absolute path for PlanJSONPath"))
		}
		config.PlanJSONPath = absPlanJSONPath
	}

	config.Version = "0.3.3"

	// Version ausgeben
	if config.GetVersion != "" {
		log.Printf("Rover v%s\n", config.Version)
	}

	return config
}
