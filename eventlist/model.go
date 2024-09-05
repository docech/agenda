package eventlist

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
) 

var (
    groupTextWidth = 20

    groupTextStyle = lipgloss.NewStyle().Width(groupTextWidth)
    whenTextStyle= lipgloss.NewStyle().Width(10)
    summaryTextStyle = lipgloss.NewStyle().PaddingLeft(5)
    itemStyle = lipgloss.NewStyle().PaddingLeft(groupTextWidth)
)

type timeDefinition interface {
	toText() string
}

type TimeOnlyRange struct {
	From time.Time
	To   time.Time
}

func (tr TimeOnlyRange) toText() string {
	return tr.From.Format(time.TimeOnly) + " - " + tr.To.Format(time.TimeOnly)
}

type WholeDay struct{}

func (wd WholeDay) toText() string {
	return "Celý den"
}

type item struct {
	When    timeDefinition
	Summary string
}

type group struct {
	Title       string
	Description string
	items       []item
}

type model struct {
	groups []group
}

func New() model {
	return model{
		groups: []group{
			{
				Title:       "3",
				Description: "ZÁŘÍ, PÁ",
				items: []item{
					{When: WholeDay{}, Summary: "Koupání"},
					{When: WholeDay{}, Summary: "Návštěva lékaře"},
				},
			},
			{
				Title:       "8",
				Description: "ŘÍJ, PO",
				items: []item{
					{When: WholeDay{}, Summary: "práce"},
				},
			},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	str := ""
	for _, g := range m.groups {
		if len(g.items) == 0 {
			continue
		}

		firstItem := g.items[0]

		str += lipgloss.JoinHorizontal(lipgloss.Left,
            groupTextStyle.Render(g.Title + " " + g.Description),
            whenTextStyle.Render(firstItem.When.toText()),
            summaryTextStyle.Render(firstItem.Summary),
        ) + "\n" 

		for _, i := range g.items[1:] {
			str += itemStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left,
                whenTextStyle.Render(i.When.toText()),
                summaryTextStyle.Render(i.Summary),
            )) + "\n"
		}
	}

	return str
}
