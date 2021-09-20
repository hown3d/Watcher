package tui

import (
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	cf "github.com/hown3d/cloudformation-tui/pkg/cloudformation"
)

//newStackModel returns a new viewModel associate with cloudformation stacks
func newStackModel(cfClient *cloudformation.Client, keys *listKeyMap) stackModel {
	stackList, err := getStackItemList(cfClient)
	if err != nil {
		log.Fatalf("Can't create stack list, %v", err)
	}

	stackListModel := list.NewModel(stackList, list.NewDefaultDelegate(), 0, 0)
	stackListModel.Title = "Cloudformation Stacks"
	stackListModel.Styles.Title = titleStyle
	stackListModel.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.togglePagination,
			keys.toggleHelpMenu,
			keys.enterStack,
		}
	}
	return stackModel{
		list: stackListModel,
		keys: keys,
	}
}

// returns a tick, that uses getStackItemList to fetch cloudformation stacks continously
func (m model) fetchStacksFromAWS() tea.Cmd {
	return tea.Tick(time.Second*time.Duration(1), func(t time.Time) tea.Msg {
		stacks, err := getStackItemList(m.cfClient)
		if err != nil {
			return errMsg{err}
		}
		return stackMsg{stacks}
	})
}

//converts types.Stacks to list Items
func getStackItemList(cfClient *cloudformation.Client) ([]list.Item, error) {
	stacks, err := cf.GetStacks(cfClient)
	if err != nil {
		return nil, err
	}

	var stackList []list.Item
	for _, s := range stacks {
		var stackItem stack
		if s.StackStatusReason != nil {
			stackItem.statusReason = *s.StackStatusReason
		}
		stackItem.name = *s.StackName
		stackItem.status = string(s.StackStatus)

		stackList = append(stackList, stackItem)
	}
	return stackList, nil
}

//bubbletea update loop func
func stacksUpdate(msg tea.Msg, m *model) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.stackModel.list.FilterState() == list.Filtering {
			break
		}

		switch {

		case key.Matches(msg, m.stackModel.keys.togglePagination):
			m.stackModel.list.SetShowPagination(!m.stackModel.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.stackModel.keys.toggleHelpMenu):
			m.stackModel.list.SetShowHelp(!m.stackModel.list.ShowHelp())
			return m, nil

		case key.Matches(msg, m.stackModel.keys.enterStack):
			m.isEventView = true
			return m, m.eventInit()
		}
	case stackMsg:
		// This will also call our delegate's update function.
		cmd := m.stackModel.list.SetItems(msg.stacks)
		cmds = append(cmds, m.fetchStacksFromAWS(), cmd)

	case errMsg:
		log.Printf("error accured: %v", msg.err)
		os.Exit(1)
	}
	updatedList, cmd := m.stackModel.list.Update(msg)
	m.stackModel.list = updatedList
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) stacksView() string {
	return appStyle.Render(m.stackModel.list.View())

}
