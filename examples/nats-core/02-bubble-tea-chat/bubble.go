package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

const gap = "\n\n"

type (
	errMsg error
)

type chatModel struct {
	engine chatEngine

	senderID    string
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func newChatModel(engine chatEngine) chatModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	senderID := uuid.NewString()

	m := chatModel{
		engine:      engine,
		senderID:    senderID,
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff2b")),
		err:         nil,
	}

	m.messages = append(m.messages, m.senderStyle.Render("You joined this conversation.")+m.textarea.Value())

	return m
}

func (m chatModel) Init() tea.Cmd {
	return waitForNatsMessage(m.engine.msgChan)
}

func waitForNatsMessage(ch chan chatMessage) tea.Cmd {
	return func() tea.Msg {
		msg := <-ch
		return msg
	}
}

func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msgType := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msgType.Width
		m.textarea.SetWidth(msgType.Width)
		m.viewport.Height = msgType.Height - m.textarea.Height() - lipgloss.Height(gap)

		if len(m.messages) > 0 {
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msgType.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.engine.sendChatMessage(m.textarea.Value())

			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	case chatMessage:
		chat := msg.(chatMessage)
		senderName := "You"

		if chat.SenderID != m.engine.id {
			senderName = m.engine.id
		}

		m.messages = append(m.messages, m.senderStyle.Render(fmt.Sprintf("%s: ", senderName))+chat.Message)
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))

		// Continue listening for more chat messages
		return m, waitForNatsMessage(m.engine.msgChan)

	case errMsg:
		m.err = msgType
		return m, nil

	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m chatModel) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}
