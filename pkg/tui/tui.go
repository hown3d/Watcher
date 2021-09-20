package tui

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
	listKeys = newListKeyMap()
)

type model struct {
	stackModel  stackModel
	eventModel  eventModel
	isEventView bool
	cfClient    *cloudformation.Client
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
		enterStack: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "show stack details"),
		),
	}
}

// InitialModel returns an inital bubbletea model with a cloudformation client pointing to awsEndpointURL
func InitialModel(cfClient *cloudformation.Client) model {

	return model{
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		cfClient:    cfClient,
		stackModel:  newStackModel(cfClient, listKeys),
		eventModel:  newEventModel(listKeys),
		isEventView: false,
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		stackList, err := getStackItemList(m.cfClient)
		if err != nil {
			return errMsg{err}
		}
		return stackMsg{stackList}
	}
}

func (m model) View() string {
	if m.isEventView {
		return m.eventsView()
	}
	return m.stacksView()

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := appStyle.GetPadding()
		m.stackModel.list.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)
		m.eventModel.list.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)
	}

	if m.isEventView {
		return eventsUpdate(msg, &m)
	}
	return stacksUpdate(msg, &m)
}
