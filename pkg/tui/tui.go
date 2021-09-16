package tui

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	cf "github.com/hown3d/cloudformation-tui/pkg/cloudformation"

	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

// Stack has a name, which resembles the cloudformation stack name and its status
type Stack struct {
	status       string
	statusReason string
	name         string
}

// Title to Match bubbles list item interface
func (s Stack) Title() string { return s.name }

// Description to Match bubbles list item interface
func (s Stack) Description() string { return fmt.Sprintf("%v : %v", s.status, s.statusReason) }

// FilterValue to Match bubbles list item interface
func (s Stack) FilterValue() string { return s.name }

type model struct {
	stackList list.Model
	items     []list.Item
	keys      *listKeyMap
	cfClient  *cloudformation.Client
}

type errMsg struct{ err error }

type stackMsg struct {
	stacks []list.Item
}

type listKeyMap struct {
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
}

func (m model) updateStacks() tea.Cmd {
	return tea.Tick(time.Second*time.Duration(1), func(t time.Time) tea.Msg {
		stacks, err := getStackItemList(m.cfClient)
		if err != nil {
			return errMsg{err}
		}
		return stackMsg{stacks}
	})
}

func getStackItemList(cfClient *cloudformation.Client) ([]list.Item, error) {
	stacks, err := cf.GetStacks(cfClient)
	if err != nil {
		return nil, err
	}

	var stackList []list.Item
	for _, stack := range stacks {
		item := Stack{status: string(stack.StackStatus), statusReason: *stack.StackStatusReason, name: *stack.StackName}
		stackList = append(stackList, item)
	}
	return stackList, nil
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
	}
}

// InitialModel returns an inital bubbletea model with a cloudformation client pointing to awsEndpointURL
func InitialModel(cfClient *cloudformation.Client) model {
	var (
		listKeys = newListKeyMap()
	)

	stackList, err := getStackItemList(cfClient)
	if err != nil {
		log.Fatalf("Can't create stack list, %v", err)
	}


	stackListModel := list.NewModel(stackList, list.NewDefaultDelegate(), 0, 0)
	stackListModel.Title = "Cloudformation Stacks"
	stackListModel.Styles.Title = titleStyle
	stackListModel.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}
	return model{
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		cfClient:  cfClient,
		stackList: stackListModel,
		keys:      listKeys,
		items:     stackList,
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := appStyle.GetPadding()
		m.stackList.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.stackList.FilterState() == list.Filtering {
			break
		}

		switch {

		case key.Matches(msg, m.keys.togglePagination):
			m.stackList.SetShowPagination(!m.stackList.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.stackList.SetShowHelp(!m.stackList.ShowHelp())
			return m, nil

		}
	case stackMsg:
		// This will also call our delegate's update function.
		cmd := m.stackList.SetItems(msg.stacks)
		cmds = append(cmds, m.updateStacks(), cmd)
		return m, tea.Batch(cmds...)

	case errMsg:
		log.Printf("error accured: %v", msg.err)
		os.Exit(1)
	}
	newListModel, cmd := m.stackList.Update(msg)
	m.stackList = newListModel
	cmds = append(cmds, m.updateStacks(), cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.stackList.View())

}
