package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	questions []string
	width     int
	height    int

	spinner     spinner.Model
	answerField textinput.Model
}

func New(questions []string) *model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &model{
		questions: questions,
		spinner:   s,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "esc" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return fmt.Sprintf(
			"\n%s Carregando o questionário...\n\n>>> Pressione 'q' para sair\n", m.spinner.View(),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.WellcomeComponent(),
	)
}

func (m model) WellcomeComponent() string {
	welcomeMsg := lipgloss.
		NewStyle().
		AlignHorizontal(lipgloss.Center).
		Padding(0, 8, 0, 8).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("286")).
		Render(
			fmt.Sprintf("\n%s Bem vindo ao questionário\n", m.spinner.View()),
		)

	return welcomeMsg
}

func main() {
	questions := []string{
		"Qual é o seu nome?",
		"Qual é a sua idade?",
		"Qual é o seu e-mail?",
	}

	m := New(questions)

	// Salva logs de execução da CLI para debug
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("não foi possível iniciar arquivo de log: %v", err)
	}
	defer f.Close()

	// Inicializa a CLI
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("não foi possível iniciar aplicação: %v", err)
	}
}
