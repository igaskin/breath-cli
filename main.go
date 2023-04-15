package main

// A simple example that shows how to render an animated progress bar. In this
// example we bump the progress by 25% every two seconds, animating our
// progress bar to its new target state.
//
// It's also possible to render a progress bar in a more static fashion without
// transitions. For details on that approach see the progress-static example.

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 10
	maxWidth = 90
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func main() {
	m := model{
		progress:     progress.New(progress.WithDefaultGradient()),
		breathIn:     true,
		breathCount:  0,
		breathTotal:  3,
		breathString: "",
		pause:        false,
		pausePeriod:  80,
		pauseTick:    0,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

type tickMsg time.Time

type model struct {
	progress     progress.Model
	breathIn     bool
	breathCount  int
	breathTotal  int
	breathString string
	pause        bool
	pauseTick    int
	pausePeriod  int
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		cmd := m.progress.IncrPercent(0.0)

		if m.pauseTick >= m.pausePeriod {
			// unpause and reset the tick
			m.pause = false
			m.pauseTick = 0
		}

		if m.pause {
			m.pauseTick += 1
			return m, tea.Batch(tickCmd(), cmd)
		}

		if m.breathIn {
			m.breathString = "Breath in..."
			cmd = m.progress.IncrPercent(0.01)
		} else {
			m.breathString = "Breath out..."
			cmd = m.progress.DecrPercent(0.01)
		}

		if m.progress.Percent() == 1.0 {
			m.breathIn = false
			m.breathCount += 1
			m.pause = true
			m.breathString = "Hold..."
		}

		if m.progress.Percent() == 0.0 {
			m.breathIn = true
			m.breathCount += 1
			m.pause = true
			m.breathString = "Hold..."
		}

		// if m.breathCount == m.breathTotal {
		// 	return m, tea.Quit
		// }

		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle(m.breathString)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

