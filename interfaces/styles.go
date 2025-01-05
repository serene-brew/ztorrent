package interfaces

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	gloss "github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = gloss.NewStyle().MarginLeft(2)
	itemStyle         = gloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = gloss.NewStyle().PaddingLeft(2).Foreground(gloss.Color("170")).Bold(true)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4)
	quitTextStyle     = gloss.NewStyle().Margin(1, 0, 2, 4)
)

type Crawlerstyles struct {
	list1Border   gloss.Style
	list2Border   gloss.Style
	inputBorder   gloss.Style
	tableBorder   gloss.Style
	activeColor   string
	inactiveColor string
}

func CrawlerStyles() Crawlerstyles {
	return Crawlerstyles{
		inputBorder: gloss.NewStyle().
			Border(gloss.RoundedBorder()).
			BorderForeground(gloss.Color("93")).
			Padding(1),
		tableBorder: gloss.NewStyle().
			Border(gloss.RoundedBorder()).
			BorderForeground(gloss.Color("8")).
			Padding(1),
		activeColor:   "93",
		inactiveColor: "8",
	}
}

func getTableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(gloss.NormalBorder()).
		BorderForeground(gloss.Color("8")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(gloss.Color("229")).
		Background(gloss.Color("57")).
		Bold(true)

	return s
}
