package tui

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	tea "github.com/charmbracelet/bubbletea"
	cf "github.com/hown3d/cloudformation-tui/pkg/cloudformation"
)

type model struct {
	stacks []cf.Stack
	cursor  int      // which to-do list item our cursor is pointing at
	getStacksFunc func(*cloudformation.Client) ([]cf.Stack, error)
	cfClient *cloudformation.Client
}

type errMsg struct{ err error }

type stackMsg struct {
	stacks []cf.Stack
}

func (m model) updateStacks() tea.Cmd {
	return tea.Tick(time.Second*time.Duration(10), func(t time.Time) tea.Msg {
		stacks, err := m.getStacksFunc(m.cfClient)
		if err != nil {
			return errMsg{err}
		}
		return stackMsg{stacks}
	})
}
// InitialModel returns an inital bubbletea model with a cloudformation client pointing to awsEndpointURL
func InitialModel(awsEndpointURL string) model {
	var err error
	m := model{
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		cfClient: cf.NewClient(awsEndpointURL),
		getStacksFunc: cf.GetStacks,
	}
	m.stacks, err = m.getStacksFunc(m.cfClient)
	if err != nil {
		log.Fatalf("Failed to get stacks: %v", err)
	}	
	return m
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.stacks)-1 {
				m.cursor++
			}
		}
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	for i, stack := range m.stacks {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, stack)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
