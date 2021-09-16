package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	cf "github.com/hown3d/cloudformation-tui/pkg/cloudformation"
	"github.com/hown3d/cloudformation-tui/pkg/tui"
)

var endpointURL string

func main() {
	flag.StringVar(&endpointURL, "endpoint", "", "specify a custom endpoint for aws, for example localstack")
	flag.Parse()

	if !(strings.HasPrefix(endpointURL, "http://") || strings.HasPrefix(endpointURL, "https://")) && endpointURL != "" {
		log.Fatal("Endpoint URL is malformed!")
	}

	client := cf.NewClient(endpointURL)
	p := tea.NewProgram(tui.InitialModel(client))
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
