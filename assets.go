package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-tfe"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/hashicorp/terraform-exec/tfexec"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

func (r *rover) generateAssets() error {
	// Get Plan
	err := r.getPlan()
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to parse Plan: %s", err))
	}

	// Generate RSO, Map, Graph
	err = r.GenerateResourceOverview()
	if err != nil {
		return err
	}

	err = r.GenerateMap()
	if err != nil {
		return err
	}

	err = r.GenerateGraph()
	if err != nil {
		return err
	}

	return nil
}

func (r *rover) getPlan() error {
	tmpDir, err := ioutil.TempDir("", "rover")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	tf, err := tfexec.NewTerraform(r.WorkingDir, r.TfPath)
	if err != nil {
		return err
	}

	// If user provided path to plan file
	if r.PlanPath != "" {
		log.Println("Using provided plan...")
		r.Plan, err = tf.ShowPlanFile(context.Background(), r.PlanPath)
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to read Plan (%s): %s", r.PlanPath, err))
		}
		return nil
	}

	// If user provided path to plan JSON file
	if r.PlanJSONPath != "" {
		log.Println("Using provided JSON plan...")

		planJsonFile, err := os.Open(r.PlanJSONPath)
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to read Plan (%s): %s", r.PlanJSONPath, err))
		}
		defer planJsonFile.Close()

		planJson, err := ioutil.ReadAll(planJsonFile)
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to read Plan (%s): %s", r.PlanJSONPath, err))
		}

		if err := json.Unmarshal(planJson, &r.Plan); err != nil {
			return errors.New(fmt.Sprintf("Unable to read Plan (%s): %s", r.PlanJSONPath, err))
		}

		return nil
	}

	// If user specified TFC workspace
	if r.TFCWorkspaceName != "" {
		tfcToken := os.Getenv("TFC_TOKEN")

		if tfcToken == "" {
			return errors.New("TFC_TOKEN environment variable not set")
		}

		if r.TFCOrgName == "" {
			return errors.New("Must specify Terraform Cloud organization to retrieve plan from Terraform Cloud")
		}

		config := &tfe.Config{
			Token: tfcToken,
		}

		client, err := tfe.NewClient(config)
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to connect to Terraform Cloud. %s", err))
		}

		// Get TFC Workspace
		ws, err := client.Workspaces.Read(context.Background(), r.TFCOrgName, r.TFCWorkspaceName)
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to list workspace %s in %s organization. %s", r.TFCWorkspaceName, r.TFCOrgName, err))
		}

		// Retrieve all runs from specified TFC workspace
		runs, err := client.Runs.List(context.Background(), ws.ID, tfe.RunListOptions{})
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to retrieve plan from %s in %s organization. %s", r.TFCWorkspaceName, r.TFCOrgName, err))
		}

		run := runs.Items[0]

		// Get most recent plan item
		planID := runs.Items[0].Plan.ID

		// Run hasn't been applied or discarded, therefore is still "actionable" by user
		runIsActionable := run.StatusTimestamps.AppliedAt.IsZero() && run.StatusTimestamps.DiscardedAt.IsZero()

		if runIsActionable && r.TFCNewRun {
			return errors.New(fmt.Sprintf("Did not create new run. %s in %s in %s is still active", run.ID, r.TFCWorkspaceName, r.TFCOrgName))
		}

		// If latest run is not actionable, rover will create new run
		if r.TFCNewRun {
			// Create new run in specified TFC workspace
			newRun, err := client.Runs.Create(context.Background(), tfe.RunCreateOptions{
				Refresh:   &TRUE,
				Workspace: ws,
			})
			if err != nil {
				return errors.New(fmt.Sprintf("Unable to generate new run from %s in %s organization. %s", r.TFCWorkspaceName, r.TFCOrgName, err))
			}

			run = newRun

			log.Printf("Starting new Terraform Cloud run in %s workspace...", r.TFCWorkspaceName)

			// Wait maximum of 5 mins
			for i := 0; i < 30; i++ {
				run, err := client.Runs.Read(context.Background(), newRun.ID)
				if err != nil {
					return errors.New(fmt.Sprintf("Unable to retrieve run from %s in %s organization. %s", r.TFCWorkspaceName, r.TFCOrgName, err))
				}

				if run.Plan != nil {
					planID = run.Plan.ID
					// Add 20 second timeout so plan JSON becomes available
					time.Sleep(20 * time.Second)
					log.Printf("Run %s to completed!", newRun.ID)
					break
				}

				time.Sleep(10 * time.Second)
				log.Printf("Waiting for run %s to complete (%ds)...", newRun.ID, 10*(i+1))
			}

			if planID == "" {
				return errors.New(fmt.Sprintf("Timeout waiting for plan to complete in %s in %s organization. %s", r.TFCWorkspaceName, r.TFCOrgName, err))
			}
		}

		// Get most recent plan file
		planBytes, err := client.Plans.JSONOutput(context.Background(), planID)
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to retrieve plan from %s in %s organization. %s", r.TFCWorkspaceName, r.TFCOrgName, err))
		}
		// If empty plan file
		if string(planBytes) == "" {
			return errors.New(fmt.Sprintf("Empty plan. Check run %s in %s in %s is not pending", run.ID, r.TFCWorkspaceName, r.TFCOrgName))
		}

		if err := json.Unmarshal(planBytes, &r.Plan); err != nil {
			return errors.New(fmt.Sprintf("Unable to parse plan (ID: %s) from %s in %s organization.: %s", planID, r.TFCWorkspaceName, r.TFCOrgName, err))
		}

		return nil
	}

	log.Println("Initializing Terraform...")

	// Create TF Init options
	var tfInitOptions []tfexec.InitOption
	tfInitOptions = append(tfInitOptions, tfexec.Upgrade(true))

	// Add *.tfbackend files
	for _, tfBackendConfig := range r.TfBackendConfigs {
		if tfBackendConfig != "" {
			tfInitOptions = append(tfInitOptions, tfexec.BackendConfig(tfBackendConfig))
		}
	}

	// tfInitOptions = append(tfInitOptions, tfexec.LockTimeout("60s"))

	err = tf.Init(context.Background(), tfInitOptions...)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to initialize Terraform Plan: %s", err))
	}

	if r.WorkspaceName != "" {
		log.Printf("Running in %s workspace...", r.WorkspaceName)
		err = tf.WorkspaceSelect(context.Background(), r.WorkspaceName)
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to select workspace (%s): %s", r.WorkspaceName, err))
		}
	}

	log.Println("Generating plan...")
	planPath := fmt.Sprintf("%s/%s-%v", tmpDir, "roverplan", time.Now().Unix())

	// Create TF Plan options
	var tfPlanOptions []tfexec.PlanOption
	tfPlanOptions = append(tfPlanOptions, tfexec.Out(planPath))

	// Add *.tfvars files
	for _, tfVarsFile := range r.TfVarsFiles {
		if tfVarsFile != "" {
			tfPlanOptions = append(tfPlanOptions, tfexec.VarFile(tfVarsFile))
		}
	}

	// Add Terraform variables
	for _, tfVar := range r.TfVars {
		if tfVar != "" {
			tfPlanOptions = append(tfPlanOptions, tfexec.Var(tfVar))
		}
	}

	_, err = tf.Plan(context.Background(), tfPlanOptions...)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to run Plan: %s", err))
	}

	r.Plan, err = tf.ShowPlanFile(context.Background(), planPath)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to read Plan: %s", err))
	}

	return nil
}

func showJSON(g interface{}) {
	j, err := json.Marshal(g)
	if err != nil {
		log.Printf("Error producing JSON: %s\n", err)
		os.Exit(2)
	}
	log.Printf("%+v", string(j))
}

func showModuleJSON(module *tfconfig.Module) {
	j, err := json.MarshalIndent(module, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error producing JSON: %s\n", err)
		os.Exit(2)
	}
	os.Stdout.Write(j)
	os.Stdout.Write([]byte{'\n'})
}

func saveJSONToFile(prefix string, fileType string, path string, j interface{}) string {
	b, err := json.Marshal(j)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error producing JSON: %s\n", err)
		os.Exit(2)
	}

	newpath := filepath.Join(".", fmt.Sprintf("%s/%s", path, prefix))
	err = os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(fmt.Sprintf("%s/%s-%s.json", newpath, prefix, fileType))
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err = f.WriteString(string(b))
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s/%s-%s.json", newpath, prefix, fileType)
}
