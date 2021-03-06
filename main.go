package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	extensions         []string
	extensionsCursor   int
	selectedExtensions map[int]struct{}
	apps               []string
	appsCursor         int
	selectedApps       map[int]struct{}
}

var initialModel = model{
	// Our to-do list is just a grocery list
	extensions: []string{"Manage extensions", "Manage apps", "Manage upgrade", "Manage snapshots"},
	apps:       []string{"CRM", "Tasks", "Support"},

	// A map which indicates which extensions are selectedExtensions. We're using
	// the  map like a mathematical set. The keys refer to the indexes
	// of the `extensions` slice, above.
	selectedExtensions: make(map[int]struct{}),
	selectedApps:       make(map[int]struct{}),
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the extensionsCursor up
		case "up", "k":
			if m.extensionsCursor > 0 {
				m.extensionsCursor--
			}

		// The "down" and "j" keys move the extensionsCursor down
		case "down", "j":
			if m.extensionsCursor < len(m.extensions)-1 {
				m.extensionsCursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selectedExtensions state for the item that the extensionsCursor is pointing at.
		case "enter", " ":
			_, ok := m.selectedExtensions[m.extensionsCursor]
			if ok {
				delete(m.selectedExtensions, m.extensionsCursor)
			} else {
				m.selectedExtensions[m.extensionsCursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "MWStake MediaWiki Manager\n\n"

	s += "Extensions\n=========\n"
	// Iterate over our extensions
	for i, choice := range m.extensions {

		// Is the extensionsCursor pointing at this choice?
		extensionsCursor := " " // no extensionsCursor
		if m.extensionsCursor == i {
			extensionsCursor = ">" // extensionsCursor!
		}

		// Is this choice selectedExtensions?
		checked := " " // not selectedExtensions
		if _, ok := m.selectedExtensions[i]; ok {
			checked = "x" // selectedExtensions!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", extensionsCursor, checked, choice)
	}

	s += "\nApps\n=========\n"
	// Iterate over our apps
	for i, choice := range m.apps {

		// Is the extensionsCursor pointing at this choice?
		appsCursor := " " // no extensionsCursor
		if m.appsCursor == i {
			appsCursor = ">" // extensionsCursor!
		}

		// Is this choice selectedExtensions?
		checked := " " // not selectedExtensions
		if _, ok := m.selectedApps[i]; ok {
			checked = "x" // selectedExtensions!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", appsCursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(initialModel)
	p.EnterAltScreen()
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	p.ExitAltScreen()
}
