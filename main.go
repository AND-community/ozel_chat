package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lucian95511/and/internal/pluginapi"
)

var manifest = pluginapi.Manifest{
	Name:        "ozel_chat",
	Label:       "Özel Chat",
	Version:     "2.0.0",
	Description: "Peer ID ile doğrudan özel mesajlaşma — libp2p stream (/and/dm/1.0.0)",
	Author:      "AND",
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--manifest" {
		data, _ := json.Marshal(manifest)
		fmt.Println(string(data))
		return
	}

	client, err := pluginapi.NewClientFromEnv()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ozel_chat:", err)
		os.Exit(1)
	}

	id, err := client.Identity()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ozel_chat: kimlik sorgulanamadı:", err)
		os.Exit(1)
	}

	m := newModel(client, id)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "ozel_chat:", err)
		os.Exit(1)
	}
}

var (
	stTitle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	stSelf   = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	stOther  = lipgloss.NewStyle()
	stOk     = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	stErr    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	stDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	stBorder = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
)

type dmPollResult struct {
	msgs []pluginapi.DMMsg
	err  error
}
type sendResult struct{ err error }

type ekran int

const (
	ekPeerGirisi ekran = iota
	ekSohbet
)

type model struct {
	client    *pluginapi.Client
	identity  pluginapi.IdentityInfo
	ekran     ekran
	peerInput textinput.Model
	msgInput  textinput.Model
	peerID    string
	vp        viewport.Model
	lines     []string
	notice    string
	isErr     bool
	w, h      int
}

func newModel(c *pluginapi.Client, id pluginapi.IdentityInfo) model {
	pid := textinput.New()
	pid.Placeholder = "hedef peer ID (12D3Koo…)"
	pid.CharLimit = 128
	pid.Focus()

	msg := textinput.New()
	msg.Placeholder = "mesaj…"
	msg.CharLimit = 500

	return model{client: c, identity: id, ekran: ekPeerGirisi, peerInput: pid, msgInput: msg}
}

func (m model) Init() tea.Cmd { return textinput.Blink }

func pollCmd(c *pluginapi.Client) tea.Cmd {
	return func() tea.Msg {
		msgs, err := c.PollDM()
		return dmPollResult{msgs: msgs, err: err}
	}
}

func sendCmd(c *pluginapi.Client, peerID, msg string) tea.Cmd {
	return func() tea.Msg {
		return sendResult{err: c.SendDM(peerID, msg)}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		m.vp.Width = msg.Width - 4
		m.vp.Height = msg.Height - 10
		m.msgInput.Width = msg.Width - 4
		m.syncVP()
		return m, nil

	case dmPollResult:
		if msg.err == nil && len(msg.msgs) > 0 {
			for _, dm := range msg.msgs {
				ts := dm.ReceivedAt.Local().Format("15:04")
				line := stOther.Render(fmt.Sprintf("[%s] %s: %s", ts, dm.From, dm.Text))
				m.lines = append(m.lines, line)
			}
			m.syncVP()
		}
		return m, pollCmd(m.client)

	case sendResult:
		if msg.err != nil {
			m.notice = "Gönderme hatası: " + msg.err.Error()
			m.isErr = true
		} else {
			m.notice = ""
			m.isErr = false
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	if m.ekran == ekPeerGirisi {
		var cmd tea.Cmd
		m.peerInput, cmd = m.peerInput.Update(msg)
		return m, cmd
	}
	if m.ekran == ekSohbet {
		var cmd tea.Cmd
		m.msgInput, cmd = m.msgInput.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.ekran {
	case ekPeerGirisi:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			pid := strings.TrimSpace(m.peerInput.Value())
			if pid == "" {
				m.notice = "Peer ID boş olamaz."
				m.isErr = true
				return m, nil
			}
			m.peerID = pid
			m.ekran = ekSohbet
			m.msgInput.Focus()
			m.peerInput.Blur()
			return m, tea.Batch(textinput.Blink, pollCmd(m.client))
		default:
			var cmd tea.Cmd
			m.peerInput, cmd = m.peerInput.Update(msg)
			return m, cmd
		}

	case ekSohbet:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.ekran = ekPeerGirisi
			m.peerInput.Focus()
			m.msgInput.Blur()
			return m, nil
		case "enter":
			text := strings.TrimSpace(m.msgInput.Value())
			if text == "" {
				return m, nil
			}
			m.msgInput.SetValue("")
			ts := "şimdi"
			line := stSelf.Render(fmt.Sprintf("[%s] %s: %s", ts, m.identity.Name, text))
			m.lines = append(m.lines, line)
			m.syncVP()
			return m, sendCmd(m.client, m.peerID, text)
		default:
			var cmd tea.Cmd
			m.msgInput, cmd = m.msgInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m *model) syncVP() {
	m.vp.SetContent(strings.Join(m.lines, "\n"))
	m.vp.GotoBottom()
}

func (m model) View() string {
	switch m.ekran {
	case ekPeerGirisi:
		return m.viewPeerGirisi()
	case ekSohbet:
		return m.viewSohbet()
	}
	return ""
}

func (m model) viewPeerGirisi() string {
	var sb strings.Builder
	sb.WriteString(stTitle.Render("AND — Özel Chat") + "\n\n")
	sb.WriteString("Bağlanılacak peer ID'yi gir:\n\n")
	sb.WriteString(stBorder.Render(m.peerInput.View()) + "\n\n")
	sb.WriteString(stDim.Render("Kendi peer ID'in: "+m.identity.PeerID) + "\n\n")
	if m.isErr {
		sb.WriteString(stErr.Render(m.notice) + "\n")
	}
	sb.WriteString(stDim.Render("enter bağlan   esc/q çıkış"))
	return sb.String()
}

func (m model) viewSohbet() string {
	var sb strings.Builder
	sb.WriteString(stTitle.Render(fmt.Sprintf("AND — Özel Chat  [→ %s…]", shortPeerID(m.peerID))) + "\n\n")
	sb.WriteString(m.vp.View() + "\n\n")
	sb.WriteString(stBorder.Render(m.msgInput.View()) + "\n")
	if m.isErr {
		sb.WriteString(stErr.Render(m.notice) + "\n")
	} else if m.notice != "" {
		sb.WriteString(stOk.Render(m.notice) + "\n")
	}
	sb.WriteString(stDim.Render("enter gönder   esc geri"))
	return sb.String()
}

func shortPeerID(pid string) string {
	if len(pid) > 16 {
		return pid[:16]
	}
	return pid
}
