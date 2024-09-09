package eventlist

import (
	"slices"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	groupTextWidth = 20

    groupLineStyle   = lipgloss.NewStyle().BorderBottom(true).BorderStyle(lipgloss.ThickBorder())
	groupTextStyle   = lipgloss.NewStyle().Width(groupTextWidth)
	whenTextStyle    = lipgloss.NewStyle().Width(10)
	summaryTextStyle = lipgloss.NewStyle().PaddingLeft(5)
	itemStyle        = lipgloss.NewStyle().PaddingLeft(groupTextWidth)
    selectedItemStyle= lipgloss.NewStyle().Bold(true)
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
	groups       []group
	selectedItem int
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
		selectedItem: -1,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "j":
			newSelectedItem := m.selectedItem + 1
			if newSelectedItem < m.itemsCount() {
				m.selectedItem = newSelectedItem
			}
		case "k":
			newSelectedItem := m.selectedItem - 1
			if newSelectedItem >= 0 {
				m.selectedItem = newSelectedItem
			}
		}
	}

	return m, nil
}

func (m model) View() string {
    gIdxAdjustment := 0

    groupsCmp := make([]string, len(m.groups)) 
	for _, g := range m.groups {
        var itemsCmp []string 
		for idx, i := range g.items {
			itemCmp := lipgloss.JoinHorizontal(lipgloss.Left,
				whenTextStyle.Render(i.When.toText()),
				summaryTextStyle.Render(i.Summary),
			)
            if m.selectedItem == gIdxAdjustment + idx {
                itemCmp = selectedItemStyle.Render(itemCmp)
            }
            if idx != 0 {
                itemCmp = itemStyle.Render(itemCmp)
            }

            itemsCmp = append(itemsCmp, itemCmp)
		}

        firstItemCmp, itemsCmp := itemsCmp[0], itemsCmp[1:]

		groupHeadingCmp := lipgloss.JoinHorizontal(lipgloss.Left,
			groupTextStyle.Render(g.Title+" "+g.Description),
			firstItemCmp,
		)
        groupCmp := groupLineStyle.Render(
            lipgloss.JoinVertical(
                lipgloss.Left, 
                slices.Insert(itemsCmp, 0, groupHeadingCmp)...,
            ),
        )
		groupsCmp = append(groupsCmp, groupCmp)
        gIdxAdjustment += len(g.items)
	}

    return lipgloss.JoinVertical(lipgloss.Left, groupsCmp...)
}

func (m model) itemsCount() int {
	count := 0
	for _, g := range m.groups {
		count += len(g.items)
	}
	return count
}
