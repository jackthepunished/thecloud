package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/poyraz/cloud/pkg/sdk"
)

var apiURL = "http://localhost:8080"
var apiKey string
var sdkClient *sdk.Client

func main() {
	// 1. Auth Setup
	apiKey = os.Getenv("MINIAWS_API_KEY")
	if apiKey == "" {
		fmt.Println("⚠️  MINIAWS_API_KEY not set.")
		createDemo := false
		prompt := &survey.Confirm{
			Message: "Would you like to generate a temporary key for this session?",
			Default: true,
		}
		survey.AskOne(prompt, &createDemo)

		if createDemo {
			var name string
			namePrompt := &survey.Input{
				Message: "Enter a name for your key (e.g. demo-user):",
				Default: "demo-user",
			}
			survey.AskOne(namePrompt, &name)

			tempClient := sdk.NewClient(apiURL, "")
			key, err := tempClient.CreateKey(name)
			if err == nil {
				apiKey = key
				fmt.Printf("🔑 Generated Key: %s\n\n", apiKey)
			} else {
				fmt.Println("❌ Failed to generate key. Falling back to manual input.")
			}
		}

		if apiKey == "" {
			manualPrompt := &survey.Input{
				Message: "Enter your API Key:",
			}
			survey.AskOne(manualPrompt, &apiKey)
		}
	}

	sdkClient = sdk.NewClient(apiURL, apiKey)

	for {
		mode := ""
		prompt := &survey.Select{
			Message: "☁️  Cloud CLI Control Panel - What would you like to do?",
			Options: []string{"List Instances", "Launch Instance", "Stop Instance", "Remove Instance", "View Logs", "View Details", "Exit"},
		}
		if err := survey.AskOne(prompt, &mode); err != nil {
			fmt.Println("Bye!")
			return
		}

		switch mode {
		case "List Instances":
			listInstances()
		case "Launch Instance":
			launchInstance()
		case "Stop Instance":
			stopInstance()
		case "Remove Instance":
			removeInstance()
		case "View Logs":
			viewLogs()
		case "View Details":
			showInstance()
		case "Exit":
			fmt.Println("👋 See you in the cloud!")
			return
		}
		fmt.Println("")
	}
}

func listInstances() {
	instances, err := sdkClient.ListInstances()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Print("\033[H\033[2J") // Clear screen
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "NAME", "IMAGE", "STATUS", "ACCESS"})

	for _, inst := range instances {
		access := "-"
		if inst.Ports != "" && inst.Status == "RUNNING" {
			pList := strings.Split(inst.Ports, ",")
			var mappings []string
			for _, mapping := range pList {
				parts := strings.Split(mapping, ":")
				if len(parts) == 2 {
					mappings = append(mappings, fmt.Sprintf("localhost:%s->%s", parts[0], parts[1]))
				}
			}
			access = strings.Join(mappings, ", ")
		}

		table.Append([]string{
			inst.ID[:8],
			inst.Name,
			inst.Image,
			inst.Status,
			access,
		})
	}
	table.Render()
}

func launchInstance() {
	qs := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Instance Name:"},
			Validate: survey.Required,
		},
		{
			Name: "image",
			Prompt: &survey.Select{
				Message: "Choose Image:",
				Options: []string{"alpine", "nginx:alpine", "ubuntu", "redis:alpine"},
				Default: "alpine",
			},
		},
		{
			Name: "ports",
			Prompt: &survey.Input{
				Message: "Port Mappings (host:container, optional):",
				Help:    "e.g. 8080:80",
			},
		},
	}

	answers := struct {
		Name  string
		Image string
		Ports string
	}{}

	if err := survey.Ask(qs, &answers); err != nil {
		return
	}

	inst, err := sdkClient.LaunchInstance(answers.Name, answers.Image, answers.Ports)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Launched %s (%s) successfully!\n", inst.Name, inst.Image)
}

func selectInstance(message string, statusFilter string) *sdk.Instance {
	instances, err := sdkClient.ListInstances()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil
	}

	var options []string
	instMap := make(map[string]sdk.Instance)

	for _, inst := range instances {
		if statusFilter != "" && inst.Status != statusFilter {
			continue
		}
		display := fmt.Sprintf("%s (%s) [%s]", inst.Name, inst.ID[:8], inst.Status)
		options = append(options, display)
		instMap[display] = inst
	}

	if len(options) == 0 {
		fmt.Println("⚠️  No matching instances found.")
		return nil
	}

	var selected string
	prompt := &survey.Select{
		Message: message,
		Options: options,
	}
	if err := survey.AskOne(prompt, &selected); err != nil {
		return nil
	}

	inst := instMap[selected]
	return &inst
}

func stopInstance() {
	inst := selectInstance("Select instance to stop:", "RUNNING")
	if inst == nil {
		return
	}

	if err := sdkClient.StopInstance(inst.ID); err != nil {
		fmt.Printf("❌ Failed to stop: %v\n", err)
		return
	}

	fmt.Printf("🛑 Stopping %s...\n", inst.Name)
}

func removeInstance() {
	inst := selectInstance("Select instance to REMOVE (permanent):", "")
	if inst == nil {
		return
	}

	if err := sdkClient.TerminateInstance(inst.ID); err != nil {
		fmt.Printf("❌ Failed to remove: %v\n", err)
		return
	}

	fmt.Printf("🗑️  %s removed successfully.\n", inst.Name)
}

func viewLogs() {
	inst := selectInstance("Select instance to view logs:", "")
	if inst == nil {
		return
	}

	logs, err := sdkClient.GetInstanceLogs(inst.ID)
	if err != nil {
		fmt.Printf("❌ Failed to fetch logs: %v\n", err)
		return
	}

	fmt.Println("📜 --- Logs Start ---")
	fmt.Print(logs)
	fmt.Println("📜 --- Logs End ---")
}

func showInstance() {
	inst := selectInstance("Select instance to view details:", "")
	if inst == nil {
		return
	}

	// Fetch fresh details
	details, err := sdkClient.GetInstance(inst.ID)
	if err != nil {
		fmt.Printf("❌ Failed to fetch details: %v\n", err)
		return
	}

	fmt.Printf("\n☁️  Instance Details\n")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("%-15s %v\n", "ID:", details.ID)
	fmt.Printf("%-15s %v\n", "Name:", details.Name)
	fmt.Printf("%-15s %v\n", "Status:", details.Status)
	fmt.Printf("%-15s %v\n", "Image:", details.Image)
	fmt.Printf("%-15s %v\n", "Ports:", details.Ports)
	fmt.Printf("%-15s %v\n", "Created At:", details.CreatedAt)
	fmt.Printf("%-15s %v\n", "Version:", details.Version)
	fmt.Printf("%-15s %v\n", "Container ID:", details.ContainerID)
	fmt.Println(strings.Repeat("-", 40))
}
