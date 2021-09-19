package tui

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	cf "github.com/hown3d/cloudformation-tui/pkg/cloudformation"
)

// NewEventModel returns a new viewModel for cloudformation events
func newEventModel(m model, keys *listKeyMap) (*eventModel, tea.Cmd) {
	// get the name of the stack which was selected
	selectedStack := m.stackModel.list.SelectedItem().(stack)
	// get all events associated to that stack
	eventItems, err := getEventItemList(selectedStack.name, m.cfClient)
	if err != nil {
		log.Fatalf("Can't get event Items, %v", err)
	}

	eventListModel := list.NewModel(eventItems, list.NewDefaultDelegate(), 0, 0)
	eventListModel.Styles.Title = titleStyle
	eventListModel.Title = fmt.Sprintf("%v's Resources", selectedStack.name)
	return &eventModel{
		list:  eventListModel,
		keys:  keys,
		stack: selectedStack.name,
	}, m.fetchEventsFromAWS(selectedStack.name)
}

// returns a tick, that uses getStackItemList to fetch cloudformation stacks continously
func (m model) fetchEventsFromAWS(stack string) tea.Cmd {
	return tea.Tick(time.Second*time.Duration(1), func(t time.Time) tea.Msg {
		stacks, err := getEventItemList(stack, m.cfClient)
		if err != nil {
			return errMsg{err}
		}
		return eventMsg{stacks}
	})
}

//converts types.Stacks to list Items
func getEventItemList(stack string, cfClient *cloudformation.Client) ([]list.Item, error) {
	events, err := cf.GetStackEvents(stack, cfClient)
	if err != nil {
		return nil, err
	}

	var eventItems []list.Item
	for _, e := range events {
		item := event{status: string(e.ResourceStatus), resourceType: *e.ResourceType, resource: *e.ResourceType}
		eventItems = append(eventItems, item)
	}
	return eventItems, nil
}

func (m model) eventsView() string {
	return appStyle.Render(m.eventModel.list.View())
}

func eventsUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := appStyle.GetPadding()
		m.eventModel.list.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)

	case tea.KeyMsg:
		switch msg.String() {
		case "backspace":
			// dont return to stackView if user is actively typing a filter
			if m.eventModel.list.FilterState() != list.Filtering {
				m.isEventView = false
				return m, nil
			}
		}
	case eventMsg:
		cmd := m.eventModel.list.SetItems(msg.events)
		cmds = append(cmds, cmd)
	case errMsg:
		log.Printf("error accured: %v", msg.err)
	}
	updatedList, cmd := m.eventModel.list.Update(msg)

	cmds = append(cmds, cmd)
	m.eventModel.list = updatedList
	return m, tea.Batch(cmds...)
}
