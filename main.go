package main

import (
	"embed"
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
	"io/fs"
	"log"
	"net/http"
	"rover/config"
	"strings"
)

var TRUE = true

//go:embed ui/dist
var frontend embed.FS

type rover struct {
	Name             string
	WorkingDir       string
	TfPath           string
	TfVarsFiles      []string
	TfVars           []string
	TfBackendConfigs []string
	PlanPath         string
	PlanJSONPath     string
	WorkspaceName    string
	TFCOrgName       string
	TFCWorkspaceName string
	ShowSensitive    bool
	GenImage         bool
	TFCNewRun        bool
	Plan             *tfjson.Plan
	RSO              *ResourcesOverview
	Map              *Map
	Graph            Graph
}

func main() {
	cfg := config.LoadConfig()
	r := createRoverFromConfig(*cfg)
	runApp(r, *cfg)
}

func createRoverFromConfig(cfg config.Config) rover {
	parsedTfVarsFiles := strings.Split(cfg.TfVarsFiles.String(), ",")
	parsedTfVars := strings.Split(cfg.TfVars.String(), ",")
	parsedTfBackendConfigs := strings.Split(cfg.TfBackendConfigs.String(), ",")
	return rover{
		Name:             cfg.Name,
		WorkingDir:       cfg.WorkingDir,
		TfPath:           cfg.TfPath,
		PlanPath:         cfg.PlanPath,
		PlanJSONPath:     cfg.PlanJSONPath,
		ShowSensitive:    cfg.ShowSensitive,
		GenImage:         cfg.GenImage,
		TfVarsFiles:      parsedTfVarsFiles,
		TfVars:           parsedTfVars,
		TfBackendConfigs: parsedTfBackendConfigs,
		WorkspaceName:    cfg.WorkspaceName,
		TFCOrgName:       cfg.TFCOrgName,
		TFCWorkspaceName: cfg.TFCWorkspaceName,
		TFCNewRun:        cfg.TFCNewRun,
	}
}

func runApp(r rover, cfg config.Config) {
	log.Println("Starting Rover...")
	// Generate assets
	var err = r.generateAssets()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Done generating assets.")

	// Save to file (debug)
	// saveJSONToFile(name, "plan", "output", r.Plan)
	// saveJSONToFile(name, "rso", "output", r.Plan)
	// saveJSONToFile(name, "map", "output", r.Map)
	// saveJSONToFile(name, "graph", "output", r.Graph)

	// Embed frontend
	fe, err := fs.Sub(frontend, "ui/dist")
	if err != nil {
		log.Fatalln(err)
	}
	frontendFS := http.FileServer(http.FS(fe))

	if cfg.Standalone {
		err = r.generateZip(fe, fmt.Sprintf("%s.zip", cfg.ZipFileName))
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Generated zip file: %s.zip\n", cfg.ZipFileName)
		return
	}

	err = r.startServer(cfg.IPPort, frontendFS)
	if err != nil {
		// http.Serve() returns error on shutdown
		if cfg.GenImage {
			log.Println("Server shut down.")
		} else {
			log.Fatalf("Could not start server: %s\n", err.Error())
		}
	}
}
