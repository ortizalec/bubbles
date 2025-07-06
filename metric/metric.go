package metric

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UpdateValueMsg struct {
	Index int
	Value string
}

type Metric struct {
	Title string
	Value string
}

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

type Layout struct {
	Title     lipgloss.Style
	Value     lipgloss.Style
	Direction Direction
}

var (
	CardLayout = Layout{
		Title:     lipgloss.NewStyle(),
		Value:     lipgloss.NewStyle().Bold(true),
		Direction: Vertical,
	}
	TagLayout = Layout{
		Title: lipgloss.NewStyle().PaddingRight(1).
			Foreground(lipgloss.Color("#18181b")).Background(lipgloss.Color("#71717a")),
		Value: lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#71717a")),
		Direction: Horizontal,
	}
	ListLayout = Layout{
		Title:     lipgloss.NewStyle().PaddingRight(1),
		Value:     lipgloss.NewStyle().Bold(true),
		Direction: Horizontal,
	}
)

type Model struct {
	metrics   []Metric
	layout    Layout
	direction Direction // control direction of join
	border    bool
	gap       int // spacing between cells
}

type Option func(*Model)

func New(opts ...Option) Model {
	m := Model{
		layout:    CardLayout,
		direction: Horizontal,
		border:    false,
		gap:       1,
	}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

func (m *Model) GetMetric(i int) Metric {
	if i >= 0 && i < len(m.metrics) {
		return m.metrics[i]
	}
	return Metric{}
}

func (m *Model) SetLayout(l Layout) {
	m.layout = l
}

func (m *Model) SetDirection(d Direction) {
	m.direction = d
}

func WithMetrics(metrics []Metric) Option {
	return func(m *Model) {
		m.metrics = metrics
	}
}

func WithLayout(layout Layout) Option {
	return func(m *Model) {
		m.layout = layout
	}
}

func WithDirection(dir Direction) Option {
	return func(m *Model) {
		m.direction = dir
	}
}

func WithBorder(b bool) Option {
	return func(m *Model) {
		m.border = b
	}
}

func WithGap(g int) Option {
	return func(m *Model) {
		m.gap = g
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case UpdateValueMsg:
		if msg.Index >= 0 && msg.Index < len(m.metrics) {
			m.metrics[msg.Index].Value = msg.Value
		}
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	views := make([]string, len(m.metrics))
	_, maxWidth := longestMetric(m.metrics)
	for i, metric := range m.metrics {
		// Render title and value
		var content string
		switch m.layout.Direction {
		case Vertical:
			content = lipgloss.JoinVertical(lipgloss.Top,
				m.layout.Title.Render(metric.Title),
				m.layout.Value.Render(metric.Value),
			)
		case Horizontal:
			content = lipgloss.JoinHorizontal(lipgloss.Top,
				m.layout.Title.Render(metric.Title),
				m.layout.Value.Render(metric.Value),
			)
		}

		// Configure style with dynamic borders
		var style lipgloss.Style

		if m.border {
			style = indexAwareBorder(i, m.direction, len(m.metrics), maxWidth)
		} else {
			style = applyGap(style, i, m.direction, len(m.metrics), m.gap)
		}

		views[i] = style.Render(content)
	}

	// Join metrics based on outer direction
	switch m.direction {
	case Vertical:
		return lipgloss.JoinVertical(lipgloss.Top, views...)
	case Horizontal:
		return lipgloss.JoinHorizontal(lipgloss.Top, views...)
	default:
		return lipgloss.JoinVertical(lipgloss.Top, views...)
	}
}

func applyGap(style lipgloss.Style, i int, d Direction, l int, g int) lipgloss.Style {
	switch d {
	case Vertical:
		switch i {
		case 0:
			return style
		case l - 1:
			style = style.MarginTop(g)
		default:
			style = style.MarginTop(g)
		}
	case Horizontal:
		switch i {
		case 0:
			style = style.MarginRight(g)
		case l - 1:
			return style
		default:
			style = style.MarginRight(g)
		}
	}
	return style
}

func indexAwareBorder(i int, d Direction, l int, w int) lipgloss.Style {
	style := lipgloss.NewStyle()
	switch d {
	case Vertical:
		style = style.Width(w)
		border := lipgloss.NormalBorder()
		switch {
		case l == 1:
			style = style.Border(border)
		case i == 0:
			style = style.Border(border).
				BorderBottom(false)
		case i == l-1:
			border.TopRight = "┤"
			border.TopLeft = "├"
			style = style.Border(border)
		default:
			border.TopRight = "┤"
			border.TopLeft = "├"
			style = style.Border(border).
				BorderBottom(false)
		}
	case Horizontal:
		border := lipgloss.NormalBorder()
		switch {
		case l == 1:
			style = style.Border(border)
		case i == 0:
			border.BottomRight = "┴"
			border.TopRight = "┬"
			style = style.Border(border)
		case i == l-1:
			style = style.Border(border).
				BorderLeft(false)
		default:
			border.BottomRight = "┴"
			border.TopRight = "┬"
			style = style.Border(border).
				BorderLeft(false)
		}
	}
	return style
}

func longestMetric(metrics []Metric) (Metric, int) {
	var longest Metric
	maxLen := 0

	for _, m := range metrics {
		length := len(m.Title) + len(m.Value)
		if length > maxLen {
			maxLen = length
			longest = m
		}
	}

	return longest, maxLen
}
