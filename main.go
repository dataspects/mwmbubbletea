package main

// An example demonstrating an application with multiple views.
//
// Note that this example was produced before the Bubbles progress component
// was available (github.com/charmbracelet/bubbles/progress) and thus, we're
// implementing a progress bar from scratch here.

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fogleman/ease"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
)

// General stuff for styling the view
var (
	term          = termenv.ColorProfile()
	keyword       = makeFgStyle("211")
	subtle        = makeFgStyle("241")
	progressEmpty = subtle(progressEmptyChar)
	dot           = colorFg(" • ", "236")

	// Gradient colors we'll use for the progress bar
	ramp = makeRamp("#B14FFF", "#00FFA3", progressBarWidth)
)

func main() {
	initialModel := model{0, false, 0, false, 0, 0, false, false}
	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

type frameMsg struct{}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

type model struct {
	TopFeatureChoice      int
	TopFeatureChosen      bool
	ExtensionActionChoice int
	ExtensionActionChosen bool
	Frames                int
	Progress              float64
	Loaded                bool
	Quitting              bool
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
	if !m.TopFeatureChosen {
		return updateTopFeatureChoices(msg, m)
	}
	if !m.ExtensionActionChosen {
		return updateExtensionActionChoices(msg, m)
	}
	return updateChosen(msg, m)
}

// The main view, which just calls the appropriate sub-view
func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}
	if !m.TopFeatureChosen {
		s = choicesView(m)
	} else {
		s = chosenView(m)
	}
	return indent.String("\n"+s+"\n\n", 2)
}

// Sub-update functions

// Update loop for the first view where you're choosing a task.
func updateTopFeatureChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.TopFeatureChoice += 1
			if m.TopFeatureChoice > 1 {
				m.TopFeatureChoice = 1
			}
		case "k", "up":
			m.TopFeatureChoice -= 1
			if m.TopFeatureChoice < 0 {
				m.TopFeatureChoice = 0
			}
		case "enter":
			m.TopFeatureChosen = true
			return m, frame()
		}

	}

	return m, nil
}

func updateExtensionActionChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.ExtensionActionChoice += 1
			if m.ExtensionActionChoice > 1 {
				m.ExtensionActionChoice = 1
			}
		case "k", "up":
			m.ExtensionActionChoice -= 1
			if m.ExtensionActionChoice < 0 {
				m.ExtensionActionChoice = 0
			}
		case "enter":
			m.ExtensionActionChosen = true
			return m, frame()
		}

	}

	return m, nil
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

// Sub-views

// The first view, where you're choosing a task
func choicesView(m model) string {
	c := m.TopFeatureChoice

	tpl := "MWStake MediaWiki Manager\n\n"
	tpl += "%s\n\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s",
		checkbox("Enable/disable extensions", c == 0),
		checkbox("Take/view snapshots", c == 1),
	)

	return fmt.Sprintf(tpl, choices)
}

func extensionCatalogue(m model) string {
	c := m.ExtensionActionChoice
	ec := fmt.Sprintf(
		"%s\n%s",
		checkbox("Enable", c == 0),
		checkbox("Disable", c == 1),
	)
	return ec
}

func manageExtensionsInterface(m model) string {
	ifce := []string{
		"Manage extensions\n",
		extensionCatalogue(m) + "\n",
		fmt.Sprintf("%s or %s", keyword("enable"), keyword("disable")),
	}
	return strings.Join(ifce, "\n")
}

func manageSnapshotsInterface() string {
	ifce := []string{
		"Manage snapshots",
		fmt.Sprintf("%s or %s", keyword("take"), keyword("view")),
	}
	return strings.Join(ifce, "\n")
}

// The second view, after a task has been chosen
func chosenView(m model) string {
	var msg string

	switch m.TopFeatureChoice {
	case 0:
		msg = manageExtensionsInterface(m)
	case 1:
		msg = manageSnapshotsInterface()
	}

	return msg
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

func progressbar(width int, percent float64) string {
	w := float64(progressBarWidth)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += termenv.String(progressFullChar).Foreground(term.Color(ramp[i])).String()
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}

// Utils

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

// Color a string's foreground and background with the given value.
func makeFgBgStyle(fg, bg string) func(string) string {
	return termenv.Style{}.
		Foreground(term.Color(fg)).
		Background(term.Color(bg)).
		Styled
}

// Generate a blend of colors.
func makeRamp(colorA, colorB string, steps float64) (s []string) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, colorToHex(c))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format compatible with termenv.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
