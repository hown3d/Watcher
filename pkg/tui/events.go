package tui

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	cf "github.com/hown3d/cloudformation-tui/pkg/cloudformation"
)

// NewEventModel returns a new viewModel for cloudformation events
func newEventModel(keys *listKeyMap) *eventModel {
	// get the name of the stack which was selected
	eventListModel := list.NewModel([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	eventListModel.Styles.Title = titleStyle
	return &eventModel{
		list: eventListModel,
		keys: keys,
	}
}

// returns a tick, that uses getStackItemList to fetch cloudformation stacks continously
func (m model) fetchEventsFromAWS() tea.Cmd {
	return tea.Tick(time.Second*time.Duration(1), func(t time.Time) tea.Msg {
		// return nil if stack isn't set
		if m.eventModel.stack == nil {
			return nil
		}
		stacks, err := m.getEventItemList()
		if err != nil {
			return errMsg{err}
		}
		return eventMsg{stacks}
	})
}

//converts types.Stacks to list Items
func (m model) getEventItemList() ([]list.Item, error) {
	events, err := cf.GetStackEvents(m.eventModel.stack.name, m.cfClient)
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

func (m model) eventInit() tea.Cmd {
	return func() tea.Msg {
		selectedStack := m.stackModel.list.SelectedItem().(stack)
		m.eventModel.stack = &selectedStack
		m.eventModel.list.Title = fmt.Sprintf("%v's Events", selectedStack.name)
		eventItems, err := m.getEventItemList()
		if err != nil {
			return errMsg{err}
		}
		return eventMsg{eventItems}
	}
}

func (m model) eventsView() string {
	return appStyle.Render(m.eventModel.list.View())
}

func eventsUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "backspace":
			// dont return to stackView if user is actively typing a filter
			if m.eventModel.list.FilterState() != list.Filtering {
				m.eventModel.stack = nil
				m.isEventView = false
				return m, nil
			}
		}
	case eventMsg:
		cmd := m.eventModel.list.SetItems(msg.events)
		cmds = append(cmds, cmd, m.fetchEventsFromAWS())
	case errMsg:
		log.Printf("error accured: %v", msg.err)
	}
	updatedList, cmd := m.eventModel.list.Update(msg)

	cmds = append(cmds, cmd)
	m.eventModel.list = updatedList
	return m, tea.Batch(cmds...)
}
