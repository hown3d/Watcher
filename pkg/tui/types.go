package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

type stackModel struct {
	list list.Model
	keys *listKeyMap
}
type eventModel struct {
	list  list.Model
	keys  *listKeyMap
	stack *stack
	ready bool
}
type stack struct {
	status       string
	statusReason string
	name         string
}

// Title to Match bubbles list item interface
func (s stack) Title() string { return s.name }

// Description to Match bubbles list item interface
func (s stack) Description() string { return fmt.Sprintf("%v : %v", s.status, s.statusReason) }

// FilterValue to Match bubbles list item interface
func (s stack) FilterValue() string { return s.name }

type event struct {
	status       string
	resourceType string
	resource     string
	statusReason string
}

// Title to Match bubbles list item interface
func (e event) Title() string { return e.resource }

// Description to Match bubbles list item interface
func (e event) Description() string { return fmt.Sprintf("%v : %v", e.resourceType, e.status) }

// FilterValue to Match bubbles list item interface
func (e event) FilterValue() string { return e.resourceType }

type errMsg struct{ err error }

type stackMsg struct {
	stacks []list.Item
}

type eventMsg struct {
	events []list.Item
}

type notReadyMsg string 

func (msg notReadyMsg) String() string {
	return string(msg)
}

type listKeyMap struct {
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	enterStack       key.Binding
}
