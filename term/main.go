package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	notes []string
}

type errMsg struct{ err error }

type noteMsg struct{ notes []string }

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
	return func() tea.Msg {
		return fetchNotes()
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
	case errMsg:
	case noteMsg:
		m.notes = msg.notes
	}

	return m, nil
}

func (m model) View() string {
	return "Notes application"
}

func fetchNotes() tea.Msg {
	resp, err := http.Get("http://localhost:8080/notes")
	if err != nil {
		return errMsg{err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errMsg{err}
	}

	var notes []string
	if err := json.Unmarshal(body, &notes); err != nil {
		return errMsg{err}
	}

	return noteMsg{notes}
}
