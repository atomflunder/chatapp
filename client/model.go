package main

import (
	"fmt"

	"github.com/atomflunder/chatapp/models"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	username    string
	channel     string
	messages    []models.Message
	viewport    viewport.Model
	textarea    textarea.Model
	senderStyle lipgloss.Style
}

type updateMessage struct {
	lastUpdate int64
}

func initialModel(username string, channel string) model {
	messages := []models.Message{}

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = fmt.Sprintf("%s: ", username)
	ta.CharLimit = 128

	ta.SetWidth(128)
	ta.SetHeight(1)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(128, 25)
	vp.SetContent(fmt.Sprintf("Chatroom %s", channel))

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		username:    username,
		channel:     channel,
		messages:    messages,
		textarea:    ta,
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height("\n\n")

		if len(m.messages) > 0 {
			s := ""
			for _, msg := range m.messages {
				s += fmt.Sprintf("%s\n", msg.Format())
			}
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(s))
		}
		m.viewport.GotoBottom()
	case updateMessage:
		newMessages := getMessages(m.username, m.channel, int64(msg.lastUpdate))
		m.messages = append(m.messages, newMessages...)

		s := ""
		for _, msg := range m.messages {
			s += fmt.Sprintf("%s\n", msg.Format())
		}

		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(s))
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			partialMessage := models.PartialMessage{
				Content:  m.textarea.Value(),
				Username: m.username,
			}
			message := partialMessage.GetMessage(m.channel)

			m.messages = append(m.messages, message)

			postMessage(message)

			s := ""
			for _, msg := range m.messages {
				s += fmt.Sprintf("%s\n", msg.Format())
			}
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(s))

			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	}

	return m, tea.Batch(tiCmd, vpCmd)

}

func (m model) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		"\n\n",
		m.textarea.View(),
	)
}
