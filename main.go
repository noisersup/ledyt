package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	yt "github.com/noisersup/ledyt/yt-client"
	"github.com/noisersup/ledyt/yt-client/output"
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
	client      *yt.Client
	searchInput textinput.Model
	spinner     spinner.Model
	logger      output.Model
	log         chan []byte
	loading     bool
	load        chan bool
	showVideo   chan []yt.Video
}

func initialModel(v []yt.Video) model {
	searchInput := textinput.New()
	searchInput.Placeholder = "Search..."
	searchInput.CharLimit = 50 //TODO: search for max
	searchInput.Width = 20

	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	out := output.New(5)

	return model{
		videos:      v,
		cursor:      0,
		searchInput: searchInput,
		client:      &yt.Client{http.Client{}},
		spinner:     s,
		logger:      out,
		log:         make(chan []byte),
		loading:     false,
		load:        make(chan bool),
		showVideo:   make(chan []yt.Video),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, textinput.Blink, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	select {
	case l := <-m.load:
		m.loading = l
	case v := <-m.showVideo:
		m.videos = v
	case msg := <-m.log:
		m.logger.PushBytes(msg)
	default:
	}

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
			if m.loading {
				break
			}
			if m.searchInput.Focused() {
				query := m.searchInput.Value()
				m.searchInput.Blur()

				m.loading = true
				go func() {
					v, err := m.client.Search(query)
					m.load <- false
					if err != nil {
						return
					}
					m.showVideo <- v
				}()
			} else {
				go func() {
					if(len(m.videos) > m.cursor) {
						m.load <- true
						var buf bytes.Buffer
						mw := io.MultiWriter(&buf)
						go func() {
							for {
								if buf.Bytes() != nil {
									m.load <- false
									return
								}
							}
						}()
						cmd := exec.Command("mpv", m.videos[m.cursor].URL)
						cmd.Stdout = mw
						if err := cmd.Run(); err != nil {
							m.log <- []byte(err.Error())
						}
					}
				}()
			}
		}
	}

	var searchCmd tea.Cmd
	var spinnerCmd tea.Cmd
	var logCmd tea.Cmd

	m.logger, logCmd = m.logger.Update(msg)
	m.searchInput, searchCmd = m.searchInput.Update(msg)
	m.spinner, spinnerCmd = m.spinner.Update(msg)
	return m, tea.Batch(cmd, spinnerCmd, searchCmd, logCmd)
}

func (m model) View() string {
	s := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Render(printLogo())
	s += "\n\n"
	s += m.searchBar()
	s += "\n"
	if m.loading {
		s += m.spinner.View()
	}
	s += "\n"
	s += m.searchList()
	s += "\n\n"
	s += m.logger.View()
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
	s := "results\n\n"

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

func printLogo() string {
	return `
    __     ______  ____  __  __  ______
   / /    / ____/ / __ \ \ \/ / /_  __/
  / /    / __/   / / / /  \  /   / /   
 / /___ / /___  / /_/ /   / /   / /    
/_____//_____/ /_____/   /_/   /_/     
                                     `
}
