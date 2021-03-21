package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dataspects/mwmbubbletea/extensions"
	"github.com/dataspects/mwmbubbletea/snapshots"
	"github.com/dataspects/mwmbubbletea/utils"
	"github.com/fogleman/ease"
	"github.com/muesli/reflow/indent"
)

var (
	subtle = utils.MakeFgStyle("241")
	dot    = utils.ColorFg(" • ", "236")
)

func main() {
	initialModel := model{0, false, 0, false, 0, 0, false, false}
	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

type model struct {
	TopTaskChoice         int
	TopTaskChosen         bool
	ExtensionActionChoice int
	ExtensionActionChosen bool
	Frames                int
	Progress              float64
	Loaded                bool
	Quitting              bool
}

type frameMsg struct{}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return nil
}

// Main update function.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if !m.TopTaskChosen {
		return updateTopTaskChoices(msg, m)
	}
	return updateChosen(msg, m)
}

// Update loop for the second view after a choice has been made
func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg.(type) {

	case frameMsg:
		if !m.Loaded {
			m.Frames += 1
			m.Progress = ease.OutBounce(float64(m.Frames) / float64(100))
			if m.Progress >= 1 {
				m.Progress = 1
				m.Loaded = true
				return m, nil
			}
			return m, frame()
		}

	}

	return m, nil
}

// Update loop for the first view where you're choosing a task.
func updateTopTaskChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	numberOfOptions := 3
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.TopTaskChoice += 1
			if m.TopTaskChoice > numberOfOptions {
				m.TopTaskChoice = numberOfOptions
			}
		case "k", "up":
			m.TopTaskChoice -= 1
			if m.TopTaskChoice < 0 {
				m.TopTaskChoice = 0
			}
		case "enter":
			m.TopTaskChosen = true
			return m, frame()
		}

	}

	return m, nil
}

// The main view, which just calls the appropriate sub-view
func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}
	if !m.TopTaskChosen {
		s = choicesView(m)
	} else {
		s = chosenView(m)
	}
	return indent.String("\n"+s+"\n\n", 2)
}

// Sub-update functions

// Sub-views

// The first view, where you're choosing a task
func choicesView(m model) string {
	c := m.TopTaskChoice

	tpl := "MWStake MediaWiki Manager\n\n"
	tpl += "%s\n\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		utils.Checkbox("Enable extensions", c == 0),
		utils.Checkbox("Disable extensions", c == 1),
		utils.Checkbox("Take snapshot", c == 2),
		utils.Checkbox("View snapshots", c == 3),
	)

	return fmt.Sprintf(tpl, choices)
}

// The second view, after a task has been chosen
func chosenView(m model) string {
	var msg string

	switch m.TopTaskChoice {
	case 0:
		msg = extensions.EnableExtensionsInterface()
	case 1:
		msg = snapshots.ManageSnapshotsInterface()
	}

	return msg
}
