package main

import (
	"fmt"
	"log"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	yt "github.com/noisersup/ledyt/yt-client"
)

func main() {
	client := yt.Client{http.Client{}}
	v, err := client.Search("minecraft")
	if err != nil {
		log.Fatal(err)
	}

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
	videos []yt.Video
	cursor int
}

func initialModel(v []yt.Video) model {
	return model{
		videos: v,
		cursor: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil // no I/O
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		}

	}
	return m, nil
}

func (m model) View() string {
	s := "[LEDYT]\n\n"
	for i, video := range m.videos {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s [author: %s]\n", cursor, video.Title, video.Channel.Name)
	}
	return s
}
