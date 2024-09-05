package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docech/agenda/eventlist"
)

func main() {
    if _, err := tea.NewProgram(eventlist.New()).Run(); err != nil {
        fmt.Printf("Uh oh, there was an error: %v\n", err)
        os.Exit(1)
    }
}
