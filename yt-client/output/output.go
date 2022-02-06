package output

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	//TODO: Style
	logRoot   *Message
	bufLength int
}

type Message struct {
	msg  []byte
	next *Message
}

type UpdateLog struct{}

func New(bufLength ...int) Model {
	length := 20 //default msg buffer length
	if len(bufLength) == 0 {
		length = bufLength[0]
	}

	return Model{
		logRoot:   nil,
		bufLength: length,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.(type) {
	case UpdateLog:
	}
	return m, cmd
}

func (m Model) View() string {
	return m.printMsgs()
}

func (m Model) printMsgs() string {
	out := "console\n\n"

	currentMsg := m.logRoot
	for currentMsg != nil {
		out += fmt.Sprintf("%s\n", currentMsg.msg)
		currentMsg = currentMsg.next
	}

	var style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Bold(true).
		BorderForeground(lipgloss.Color("63"))

	return style.Render(out)
}

func (m *Model) PushString(msg string) tea.Cmd {
	return m.PushBytes([]byte(msg))
}

func (m *Model) PushBytes(msg []byte) tea.Cmd {
	current := m.logRoot
	logLen := 1
	chainMsg := Message{
		msg:  msg,
		next: nil,
	}

	// If It's a first message push it as root
	if current == nil {
		m.logRoot = &chainMsg
		return func() tea.Msg { return UpdateLog{} }
	}

	// If not last message in chain - check for newer message and increment counter
	for current.next != nil {
		current = current.next
		logLen++
	}

	// add new message as the child of newest node
	current.next = &chainMsg

	// if message log length equals length specified in struct remove first message
	if logLen == m.bufLength {
		// Garbage collector should remove first message from memory due to no references in code
		m.logRoot = m.logRoot.next
	}
	return func() tea.Msg { return UpdateLog{} }
}
