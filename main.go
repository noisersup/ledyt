package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	yt "github.com/noisersup/ledyt/yt-client"
)

func main() {
	//v := mockVideos(20)
	v := []yt.Video{}

	prog := tea.NewProgram(initialModel(v))
	if err := prog.Start(); err != nil {
		log.Fatal(err)
	}
}

func mockVideos(n int) []yt.Video {
	ch := yt.Channel{"Ledu", "/521matiasda"}
	var out []yt.Video
	for i := 0; i < n; i++ {
		v := yt.Video{
			Title:   fmt.Sprintf("Video %d", i),
			Channel: &ch,
			URL:     fmt.Sprintf("/sfdsdsdf%d", i),
		}
		out = append(out, v)
	}
	return out
}

type model struct {
	videos      []yt.Video
	cursor      int
	searchInput textinput.Model
	client      *yt.Client
}

func initialModel(v []yt.Video) model {
	searchInput := textinput.New()
	searchInput.Placeholder = "Search..."
	searchInput.CharLimit = 50 //TODO: search for max
	searchInput.Width = 20

	return model{
		videos:      v,
		cursor:      0,
		searchInput: searchInput,
		client:      &yt.Client{http.Client{}},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	//key press
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.videos)-1 {
				m.cursor++
			}

		case "/":
			m.searchInput.Focus()
			return m, cmd
		case tea.KeyEsc.String():
			m.searchInput.Blur()

		case tea.KeyEnter.String():
			if m.searchInput.Focused() {
				query := m.searchInput.Value()
				m.searchInput.Blur()

				v, err := m.client.Search(query)
				if err != nil {
					return m, cmd
				}

				m.videos = v
			}
		}
	}

	m.searchInput, cmd = m.searchInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	s := "[LEDYT]\n\n"
	s += m.searchBar()
	s += "\n"
	s += m.searchList()
	return s
}

func (m model) searchBar() string {
	var s string
	var borderStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Width(25)

	s += borderStyle.Render(m.searchInput.View())
	return s
}

func (m model) searchList() string {
	var s string
	var style = lipgloss.NewStyle().
		Bold(true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))
	var selectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("63"))

	for i, video := range m.videos {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			s += selectedStyle.Render(fmt.Sprintf("%s %s [author: %s]", cursor, video.Title, video.Channel.Name))
			s += "\n"
		} else {
			s += fmt.Sprintf("%s %s [author: %s]\n", cursor, video.Title, video.Channel.Name)
		}
	}
	return style.Render(s)
}
