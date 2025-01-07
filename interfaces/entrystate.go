package interfaces

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	crawler "github.com/serene-brew/ztorrent/crawler"
)

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
					m.currentScreen = CrawlerScreen
					m.focus = inputFocus
				} else if m.choice == "Pick a torrent file" {
					ExecutePickerStub(PickerScreen)
				}
			}

		case CrawlerScreen:
			switch msg.String() {
			case "!":
				m.focus = inputFocus
				m.styles.inputBorder = m.styles.inputBorder.BorderForeground(gloss.Color(m.styles.activeColor))
				m.styles.tableBorder = m.styles.tableBorder.BorderForeground(gloss.Color(m.styles.inactiveColor))
				return m, nil
			case "@":
				noResItem := m.crawlResults[0][6].(string)
				if noResItem == "0" {
					m.focus = inputFocus
					m.styles.inputBorder = m.styles.inputBorder.BorderForeground(gloss.Color(m.styles.activeColor))
					m.styles.tableBorder = m.styles.tableBorder.BorderForeground(gloss.Color(m.styles.inactiveColor))
				} else {

					m.focus = tableFocus
					m.styles.inputBorder = m.styles.inputBorder.BorderForeground(gloss.Color(m.styles.inactiveColor))
					m.styles.tableBorder = m.styles.tableBorder.BorderForeground(gloss.Color(m.styles.activeColor))

				}
				return m, nil
			case "q", "ctrl+c":
				m.currentScreen = ListScreen
				return m, nil
			case "enter":
				if m.focus == inputFocus {
					m.styles.inputBorder = m.styles.inputBorder.BorderForeground(gloss.Color(m.styles.inactiveColor))
					m.styles.tableBorder = m.styles.tableBorder.BorderForeground(gloss.Color(m.styles.activeColor))
					rows := []table.Row{}
					m.searchTerm = m.input.Value()
					m.crawlResults, _ = crawler.GetInfoMediaQuery(m.searchTerm)
					for i, record := range m.crawlResults {
						category := crawler.ClassifyCategory(record[7].(string))
						rows = append(rows, table.Row{
							strconv.Itoa(i + 1),
							record[1].(string),
							crawler.ConvertSize(strconv.Atoi(record[6].(string))),
							record[3].(string),
							record[4].(string),
							category})
					}
					m.table.SetRows(rows)
					m.table.GotoTop()
					m.focus = tableFocus

					noResItem := m.crawlResults[0][6].(string)
					if noResItem == "0" {
						m.focus = inputFocus
						m.styles.inputBorder = m.styles.inputBorder.BorderForeground(gloss.Color(m.styles.activeColor))
						m.styles.tableBorder = m.styles.tableBorder.BorderForeground(gloss.Color(m.styles.inactiveColor))
					}
					return m, nil
				} else if m.focus == tableFocus {
					if len(m.table.Rows()) == 0 {
						m.currentScreen = ListScreen
					} else {
						selectedItem := m.table.SelectedRow()
						pseudoIndex, _ := strconv.Atoi(selectedItem[0])
						index := pseudoIndex - 1
						magnet := crawler.GetMagnet(m.crawlResults[index][2].(string), selectedItem[1])
						m.magnet = magnet
						//tea.Printf("%s", magnet)
						return m, nil
					}
				}

			}

		}

	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.currentScreen == ListScreen {
		m.list, cmd = m.list.Update(msg)
	} else if m.currentScreen == CrawlerScreen {
		if m.focus == inputFocus {
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.focus == tableFocus {
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}
