package eventlist

import (
	tea "github.com/charmbracelet/bubbletea"
)


type TimeDefinition interface {
    toText() string
}

type item struct {
    When TimeDefinition
    Summary string
}

type group struct {
	Title       string
	Description string
	items       []item
}

type Model struct {
	groups []group
}

func (m Model) Init() tea.Cmd {
	m.groups = []group {
        group{
            Title: "3",
            Description: "ZÁŘÍ, PÁ",
            items: []item {
                item {When: , Summary: "Koupání"
            },
        },
    }
}
:
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
    str := ""
	for _, g := range m.groups {
		if len(g.items) == 0 {
			continue
		}

        firstItem := g.items[0]
        itemStr := firstItem.When.toText() + " " + firstItem.Summary 

        str += g.Title + " " + g.Description 
        str += " " + itemStr + "\n"

        for _, i := range g.items[1:] {
            str += "    " + i.When.toText() + " " + i.Summary
        }
	}

	return str
}
