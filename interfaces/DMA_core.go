package interfaces

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type FilePickerModel struct {
	currentScreen DMAState
	filepicker    filepicker.Model
	selectedFile  string
	quitting      bool
	err           error
}

func (m FilePickerModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m FilePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	mainModel := model{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.currentScreen {
		case PickerScreen:
			switch msg.String() {
			case "ctrl+c", "q", "esc":
				m.quitting = true
				mainModel.currentScreen = ListScreen
				return m, tea.Quit
			}
		}

	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.selectedFile = path
		return m, tea.Quit
	}

	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m FilePickerModel) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedFile == "" {
		s.WriteString("Pick a file:")
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}

func ExecutePickerStub(screen DMAState) {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".torrent"}
	fp.CurrentDirectory, _ = os.UserHomeDir()

	m := FilePickerModel{
		currentScreen: screen,
		filepicker:    fp,
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	//tm, _ := tea.NewProgram(&m, tea.WithAltScreen()).Run()
	//mm := tm.(FilePickerModel)
	//fmt.Println("\n  You selected: " + m.filepicker.Styles.Selected.Render(mm.selectedFile) + "\n")
}
