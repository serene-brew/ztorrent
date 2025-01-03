package interfaces

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	crawler "github.com/serene-brew/ztorrent/crawler"
)

const listHeight = 14

type AppState int

const (
	ListScreen AppState = iota
	TableScreen
)
type focus int
const (
	inputFocus focus = iota
	tableFocus
)


var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	focus 		  focus
	input 		  textinput.Model
	list          list.Model
	choice        string
	quitting      bool
	currentScreen AppState
	table         table.Model
	crawlResults [][]interface{}
	searchTerm string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch m.currentScreen {
		case ListScreen:
			switch msg.String() {
			case "q", "ctrl+c":
				m.quitting = true
				return m, tea.Quit

			case "enter":
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.choice = string(i)
				}
				if m.choice == "Crawl 'n Grab" {
					m.currentScreen = TableScreen
				}
			}

		case TableScreen:
			switch msg.String() {
			case "!":
				m.focus = inputFocus
				return m, nil	
			case "@":
				m.focus = tableFocus
				return m, nil
			case "q", "ctrl+c":
				m.currentScreen = ListScreen
				return m, nil
			case "enter":
				if m.focus == inputFocus{
					rows := []table.Row{}
					m.searchTerm = m.input.Value()
					m.crawlResults, _ = crawler.GetInfoMediaQuery(m.searchTerm)
					for i, record := range m.crawlResults{
						category := crawler.ClassifyCategory(record[7].(string))
						rows = append(rows, table.Row{
								strconv.Itoa(i), 
								record[1].(string), 
								crawler.ConvertSize(strconv.Atoi(record[6].(string))), 
								record[3].(string), 
								record[4].(string), 
								category, 
								record[2].(string)})
						}
					m.table.SetRows(rows)
					m.focus = tableFocus
					return m, nil
				}

				if m.focus == tableFocus{
					if len(m.table.Rows()) == 0{
						m.currentScreen = ListScreen
					} else {
						selectedItem := m.table.SelectedRow()
						magnet := crawler.GetMagnet(selectedItem[6], selectedItem[1])
						return m, tea.Batch(
							tea.Println(magnet),
						)
					}
				}
				
				
			}

		}

	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.list, cmd = m.list.Update(msg)
	if m.focus == inputFocus {
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.focus == tableFocus {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}	

	return m, tea.Batch(cmds...)
	/*m.table, cmd = m.table.Update(msg)
	return m, cmd*/
}

func (m model) View() string {
	switch m.currentScreen {
	case ListScreen:
		if m.choice == "Crawl 'n Grab" {
			m.currentScreen = TableScreen
		}
		if m.quitting {
			return quitTextStyle.Render("aight bruv...That’s cool.")
		}
		return "\n" + m.list.View()
	case TableScreen:
		return lipgloss.JoinVertical(lipgloss.Top, m.input.View(), "\n", m.table.View())
	}
	return ""
}

func ListModel() model {
	items := []list.Item{
		item("Crawl 'n Grab"),
		item("Torrent"),
		item("Magnet"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Pick you'r poison"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	columns := []table.Column{
		{Title: "Sl.No", Width: 10},
		{Title: "Name", Width: 40},
		{Title: "Size", Width: 10},
		{Title: "Seeders", Width: 10},
		{Title: "Leechers", Width: 10},
		{Title: "Category", Width: 15},
		{Title: "Info Hash", Width: 40},
	}
	
	input := textinput.New()
	input.Placeholder = "search your anime"
	input.Focus()

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(20),
	)
	

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	m := model{focus: inputFocus, input: input, list: l, table: t}
	return m
}
func Entrypoint() {
	m := ListModel()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
