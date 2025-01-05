package interfaces

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 12

type model struct {
	focus         focus
	input         textinput.Model
	list          list.Model
	choice        string
	quitting      bool
	currentScreen AppState
	table         table.Model
	crawlResults  [][]interface{}
	searchTerm    string
	filepicker    FilePickerModel
	selectedFile  string
	magnet        string
	styles        Crawlerstyles
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	switch m.currentScreen {
	case ListScreen:
		if m.quitting {
			return quitTextStyle.Render(m.magnet)
		}
		return "\n" + m.list.View()
	case CrawlerScreen:
		inputS := m.styles.inputBorder.Render(m.input.View())
		tableS := m.styles.tableBorder.Render(m.table.View())
		return lipgloss.JoinVertical(lipgloss.Top, inputS, tableS)
	}
	return ""
}

func ListModel() model {
	styles := CrawlerStyles()

	items := []list.Item{
		item("Crawl 'n Grab"),
		item("Pick a torrent file"),
		item("Have a magnet?"),
		item("Open download manager"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Pick your poison"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	columns := []table.Column{
		{Title: "Sl.No", Width: 5},
		{Title: "Name", Width: 35},
		{Title: "Size", Width: 10},
		{Title: "Seeders", Width: 10},
		{Title: "Leechers", Width: 10},
		{Title: "Category", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	t.SetStyles(getTableStyles())

	input := textinput.New()
	input.Placeholder = "search for the thing you want..."
	input.Width = 89
	input.Focus()

	m := model{
		focus:  inputFocus,
		input:  input,
		list:   l,
		table:  t,
		styles: styles,
		magnet: "",
	}
	return m
}

func Entrypoint() {
	m := ListModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
