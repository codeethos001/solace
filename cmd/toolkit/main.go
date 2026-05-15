package main

import (
	"fmt"
	"log"
	"path/filepath" // for forward/backward slash compatibility
	"os"
	
	"solace/internal/api"
	"solace/internal/engine"
	"solace/internal/osops"
)

func main() {
	fmt.Println("Starting Solace (v0.1.0-alpha)")
	fmt.Println("--------------------------------------------------")

	// OS 
	osEngine, err := osops.NewEngine()
	if err != nil {
		log.Fatalf("✖ Critical Error: %v\n", err)
	}
	fmt.Printf("✔ operating system detection: %s\n", osEngine.GetOSName())

	// hardening 
	hEngine := engine.NewHardeningEngine(osEngine)

	rulesPath := filepath.Join("rules", "linux")
	if osEngine.GetOSName() == "Windows" {
		rulesPath = filepath.Join("rules", "windows")
	}

	err = hEngine.LoadRules(rulesPath)
	if err != nil {
		log.Fatalf("✖ Failed to load rules: %v\n", err)
	}
	fmt.Println("✔ YAML Rules loaded successuflly.")
	fmt.Println("--------------------------------------------------")

	if len(os.Args) > 1 && os.Args[1] == "--serve" {
		server := api.NewServer(hEngine, "8080")
		log.Fatal(server.Start())
		return // server blocks indefinitely
	}

	fmt.Println("evaluating...")
	results := hEngine.EvaluateRules()

	passed := 0
	failed := 0

	// print results.
	fmt.Printf("\n%-13s | %-8s | %s\n", "RULE ID", "STATUS", "MESSAGE")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, res := range results {
		if res.Status == "Passed" {
			passed++
			fmt.Printf("✔ %-11s | %-8s | %s\n", res.RuleID, res.Status, res.Message)
		} else if res.Status == "Failed" {
			failed++
			fmt.Printf("✖ %-11s | %-8s | %s\n", res.RuleID, res.Status, res.Message)
		} else {
			fmt.Printf("⚠ %-11s | %-8s | %s\n", res.RuleID, res.Status, res.Message)
		}
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("Audit Complete: %d Passed, %d Failed\n", passed, failed)
}
