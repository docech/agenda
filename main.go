package main

import (
	"fmt"
	"os"

    tea "github.com/charmbracelet/bubbletea"
)

func main() {
    if _, err := tea.NewProgram(InitTuiModel()).Run(); err != nil {
        fmt.Printf("Uh oh, there was an error: %v\n", err)
        os.Exit(1)
    }
}
