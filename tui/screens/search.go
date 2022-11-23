package screens

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/noisersup/ledyt/backend"
	"github.com/noisersup/ledyt/backend/common"
)

type Search struct {
	results []common.Video
	cursor  int
	client  backend.Backend

	input textinput.Model
}

func (m Search) Init() tea.Cmd {
	return textinput.Blink
}

func (m Search) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// TODO: check if tea.Quit is handled properly
	cmd = m.handleMsgCommon(msg)

	// Handle inputs
	if m.input.Focused() {
		cmd = tea.Batch(cmd, m.handleMsgFocused(msg))
	} else {
		cmd = tea.Batch(cmd, m.handleMsg(msg))
	}

	return m, cmd
}

func (m Search) View() string {
	return ""
}

// handleMsg handles msg when textinput is not focused
func (m Search) handleMsg(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	// handle key presses
	case tea.KeyMsg:
		switch msg.String() {

		case "up", "k":
			if m.cursor > 0 {
				m.cursor++
			}

		case "down", "j":
			if m.cursor < len(m.results)-1 {
				m.cursor++
			}

		case "/":
			m.input.Focus()
			//return m, cmd

		}
	}
	return nil
}

// handleMsgFocused handles msg when textinput is focused
func (m Search) handleMsgFocused(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	// handle key presses
	case tea.KeyMsg:
		switch msg.String() {

		case tea.KeyEnter.String():

		case tea.KeyEsc.String():
			m.input.Blur()
		}
	}

	return nil
}

// handleMsgCommon handles msgs common for every other handle func
func (m Search) handleMsgCommon(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	// handle key presses
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return tea.Quit
		}
	}

	return nil
}
