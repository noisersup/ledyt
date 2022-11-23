package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/noisersup/ledyt/backend/common"
	"github.com/noisersup/ledyt/tui"
)

func main() {
	//v := mockVideos(20)
	v := []common.Video{}

	prog := tea.NewProgram(tui.InitialModel(v))
	if err := prog.Start(); err != nil {
		log.Fatal(err)
	}
}
