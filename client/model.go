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
	username string
	channel  string
	messages []models.Message
	viewport viewport.Model
	textarea textarea.Model
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

	vp := viewport.New(128, 30)
	vp.SetContent(fmt.Sprintf("Chatroom #%s", channel))

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		username: username,
		channel:  channel,
		messages: messages,
		textarea: ta,
		viewport: vp,
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
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(m.formatMessages()))
		}
		m.viewport.GotoBottom()
	case updateMessage:
		// Fetches new messages since updateMessage.lastUpdate
		newMessages := getMessages(m.username, m.channel, int64(msg.lastUpdate))
		m.messages = append(m.messages, newMessages...)
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(m.formatMessages()))
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			// Handles sending messages
			partialMessage := models.PartialMessage{
				Content:  m.textarea.Value(),
				Username: m.username,
			}
			message := partialMessage.GetMessage(m.channel)

			m.messages = append(m.messages, message)

			postMessage(message)

			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(m.formatMessages()))

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

func (m model) formatMessages() string {
	defaultStyle := lipgloss.NewStyle().Width(m.viewport.Width - 8).BorderStyle(lipgloss.RoundedBorder())
	systemStyle := defaultStyle.Foreground(lipgloss.Color("#ff0000"))

	s := fmt.Sprintf("You're logged in to #%s as %s - Start chatting!\n", m.channel, m.username)
	for _, msg := range m.messages {
		switch msg.Username {
		case "system":
			s += systemStyle.Render(msg.Format())
		case m.username:
			color := calculateColorCode(msg.Username)
			style := defaultStyle.Foreground(lipgloss.Color(color)).Align(lipgloss.Right)
			s += style.Render(msg.Format())
		default:
			color := calculateColorCode(msg.Username)
			style := defaultStyle.Foreground(lipgloss.Color(color)).Align(lipgloss.Left)
			s += style.Render(msg.Format())
		}

		s += "\n"
	}

	return s
}
