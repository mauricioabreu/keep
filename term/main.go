package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type note struct {
	title   string
	content string
}

type notesMsg struct {
	notes []note
}

type model struct {
	notes []note
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Failed to start app: %v", err)
		os.Exit(1)
	}
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return fetchNotes
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case notesMsg:
		m.notes = msg.notes
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString("Notes\n\n")

	for _, note := range m.notes {
		s.WriteString(fmt.Sprintf("Title: %s\n", note.title))
	}

	return s.String()
}

func fetchNotes() tea.Msg {
	return notesMsg{[]note{
		{
			title:   "First Note",
			content: "This is the first note",
		},
		{
			title:   "Second Note",
			content: "This is the second note",
		},
	}}
}
