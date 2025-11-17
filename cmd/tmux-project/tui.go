package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/project"
	"github.com/jvzantvoort/tmux-project/utils"
)

// ViewMode represents the current view in the TUI
type ViewMode int

const (
	ViewProjectList ViewMode = iota
	ViewProjectEdit
	ViewConfig
	ViewHelp
)

// Status options for dropdown
var statusOptions = []string{
	"in progress",
	"to do",
	"in review",
	"completed",
	"inactive",
	"archived",
	"parked",
}

// Model represents the main TUI state
type Model struct {
	projects         []project.Project
	selectedIndex    int
	viewMode         ViewMode
	width            int
	height           int
	err              error
	editingProject   *project.Project
	editField        int
	message          string
	textInputs       []textinput.Model
	isEditing        bool
	dropdownOpen     bool
	dropdownSelected int
}

// NewModel creates a new TUI model
func NewModel() Model {
	projects := loadProjects()

	// Initialize text inputs
	inputs := make([]textinput.Model, 5)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Project name"
	inputs[0].CharLimit = 50

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Project directory"
	inputs[1].CharLimit = 200

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Description"
	inputs[2].CharLimit = 100

	inputs[3] = textinput.New()
	inputs[3].Placeholder = "Project type"
	inputs[3].CharLimit = 30

	inputs[4] = textinput.New()
	inputs[4].Placeholder = "Status"
	inputs[4].CharLimit = 30

	return Model{
		projects:      projects,
		selectedIndex: 0,
		viewMode:      ViewProjectList,
		width:         80,
		height:        24,
		textInputs:    inputs,
		isEditing:     false,
	}
}

// loadProjects loads all available projects
func loadProjects() []project.Project {
	var projects []project.Project
	names := project.ListConfigs()

	for _, name := range names {
		proj := project.NewProject(name)
		if err := proj.Open(); err == nil {
			projects = append(projects, *proj)
		}
	}

	return projects
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If dropdown is open, handle dropdown navigation
	if m.dropdownOpen && m.viewMode == ViewProjectEdit {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "up", "k":
				if m.dropdownSelected > 0 {
					m.dropdownSelected--
				}
				return m, nil
			case "down", "j":
				if m.dropdownSelected < len(statusOptions)-1 {
					m.dropdownSelected++
				}
				return m, nil
			case "enter":
				m.textInputs[4].SetValue(statusOptions[m.dropdownSelected])
				m.dropdownOpen = false
				m.updateProjectFromInputs()
				return m, nil
			case "esc":
				m.dropdownOpen = false
				return m, nil
			}
		}
		return m, nil
	}

	// If we're editing a field, update the text input
	if m.isEditing && m.viewMode == ViewProjectEdit {
		var cmd tea.Cmd
		m.textInputs[m.editField], cmd = m.textInputs[m.editField].Update(msg)

		// Check for enter or esc to finish editing
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "enter":
				m.isEditing = false
				m.updateProjectFromInputs()
				return m, cmd
			case "esc":
				m.isEditing = false
				m.loadProjectToInputs()
				return m, cmd
			}
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case error:
		m.err = msg
		return m, nil
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.viewMode {
	case ViewProjectList:
		return m.handleProjectListKeys(msg)
	case ViewProjectEdit:
		return m.handleProjectEditKeys(msg)
	case ViewConfig:
		return m.handleConfigKeys(msg)
	case ViewHelp:
		return m.handleHelpKeys(msg)
	}
	return m, nil
}

// handleProjectListKeys handles keys in project list view
func (m Model) handleProjectListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		if m.selectedIndex > 0 {
			m.selectedIndex--
		}

	case "down", "j":
		if m.selectedIndex < len(m.projects)-1 {
			m.selectedIndex++
		}

	case "enter", " ":
		if len(m.projects) > 0 {
			m.editingProject = &m.projects[m.selectedIndex]
			m.viewMode = ViewProjectEdit
			m.editField = 0
			m.loadProjectToInputs()
		}

	case "n":
		// Create new project
		m.message = "Create new project (not yet implemented)"

	case "d":
		// Delete/archive project
		if len(m.projects) > 0 {
			m.message = "Archive project (not yet implemented)"
		}

	case "r":
		// Refresh project list
		m.projects = loadProjects()
		m.selectedIndex = 0
		m.message = "Projects refreshed"

	case "?", "h":
		m.viewMode = ViewHelp

	case "c":
		m.viewMode = ViewConfig
	}

	return m, nil
}

// handleProjectEditKeys handles keys in project edit view
func (m Model) handleProjectEditKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.viewMode = ViewProjectList
		m.editingProject = nil
		m.message = ""
		m.isEditing = false

	case "up", "k":
		if !m.isEditing && m.editField > 0 {
			m.editField--
		}

	case "down", "j":
		if !m.isEditing && m.editField < len(m.textInputs)-1 {
			m.editField++
		}

	case "enter", " ":
		if !m.isEditing {
			// Special handling for Status field (index 4) - open dropdown
			if m.editField == 4 {
				m.dropdownOpen = true
				// Find current status in options for initial selection
				currentStatus := m.textInputs[4].Value()
				m.dropdownSelected = 0
				for i, opt := range statusOptions {
					if opt == currentStatus {
						m.dropdownSelected = i
						break
					}
				}
			} else {
				m.isEditing = true
				m.textInputs[m.editField].Focus()
			}
		}

	case "s", "ctrl+s":
		// Save project
		if m.editingProject != nil && !m.isEditing {
			m.updateProjectFromInputs()
			if err := m.editingProject.Save(); err != nil {
				m.message = fmt.Sprintf("Error saving: %v", err)
			} else {
				m.message = "Project saved successfully"
				m.projects = loadProjects()
			}
		}
	}

	return m, nil
}

// loadProjectToInputs loads project data into text inputs
func (m *Model) loadProjectToInputs() {
	if m.editingProject == nil {
		return
	}
	m.textInputs[0].SetValue(m.editingProject.Name)
	m.textInputs[1].SetValue(m.editingProject.Directory)
	m.textInputs[2].SetValue(m.editingProject.Description)
	m.textInputs[3].SetValue(m.editingProject.ProjectType)
	m.textInputs[4].SetValue(m.editingProject.Status)
}

// updateProjectFromInputs updates project from text inputs
func (m *Model) updateProjectFromInputs() {
	if m.editingProject == nil {
		return
	}
	m.editingProject.Name = m.textInputs[0].Value()
	m.editingProject.Directory = m.textInputs[1].Value()
	m.editingProject.Description = m.textInputs[2].Value()
	m.editingProject.ProjectType = m.textInputs[3].Value()
	m.editingProject.Status = m.textInputs[4].Value()
}

// handleConfigKeys handles keys in config view
func (m Model) handleConfigKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.viewMode = ViewProjectList
	}
	return m, nil
}

// handleHelpKeys handles keys in help view
func (m Model) handleHelpKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q", "?":
		m.viewMode = ViewProjectList
	}
	return m, nil
}

// View renders the current view
func (m Model) View() string {
	switch m.viewMode {
	case ViewProjectList:
		return m.viewProjectList()
	case ViewProjectEdit:
		return m.viewProjectEdit()
	case ViewConfig:
		return m.viewConfig()
	case ViewHelp:
		return m.viewHelp()
	}
	return ""
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			Background(lipgloss.Color("235"))

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("82")).
			MarginTop(1)
)

// viewProjectList renders the project list view
func (m Model) viewProjectList() string {
	s := titleStyle.Render("📋 tmux-project - Project Manager") + "\n\n"

	if len(m.projects) == 0 {
		s += normalStyle.Render("No projects found. Press 'n' to create one.") + "\n"
	} else {
		for i, proj := range m.projects {
			cursor := " "
			style := normalStyle

			if i == m.selectedIndex {
				cursor = "›"
				style = selectedStyle
			}

			line := fmt.Sprintf("%s %-20s  %s", cursor, proj.Name, proj.Description)
			s += style.Render(line) + "\n"
		}
	}

	if m.message != "" {
		s += "\n" + messageStyle.Render(m.message)
	}

	s += "\n" + helpStyle.Render(
		"↑/k: up • ↓/j: down • Enter: edit • n: new • d: archive • r: refresh • c: config • ?: help • q: quit",
	)

	return s
}

// viewProjectEdit renders the project edit view
func (m Model) viewProjectEdit() string {
	if m.editingProject == nil {
		return "No project selected"
	}

	proj := m.editingProject
	s := titleStyle.Render(fmt.Sprintf("✏️  Editing: %s", proj.Name)) + "\n\n"

	labels := []string{"Name", "Directory", "Description", "Type", "Status"}

	for i := 0; i < len(m.textInputs); i++ {
		cursor := " "
		labelStyle := normalStyle

		if i == m.editField {
			cursor = "›"
			labelStyle = selectedStyle
		}

		label := labelStyle.Render(fmt.Sprintf("%s %-15s:", cursor, labels[i]))

		if m.isEditing && i == m.editField {
			s += label + " " + m.textInputs[i].View() + "\n"
		} else {
			value := m.textInputs[i].Value()
			if value == "" {
				value = normalStyle.Faint(true).Render("(empty)")
			}
			s += label + " " + normalStyle.Render(value) + "\n"
		}

		// Show dropdown for Status field when dropdown is open
		if i == 4 && m.dropdownOpen && i == m.editField {
			s += "\n"
			dropdownStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("170")).
				Padding(0, 1)

			dropdownContent := ""
			for j, opt := range statusOptions {
				if j == m.dropdownSelected {
					dropdownContent += selectedStyle.Render("› "+opt) + "\n"
				} else {
					dropdownContent += normalStyle.Render("  "+opt) + "\n"
				}
			}
			s += dropdownStyle.Render(dropdownContent) + "\n"
		}
	}

	if m.message != "" {
		s += "\n" + messageStyle.Render(m.message)
	}

	editHelp := ""
	if m.dropdownOpen {
		editHelp = "↑/k: up • ↓/j: down • Enter: select • Esc: cancel"
	} else if m.isEditing {
		editHelp = "Enter: save field • Esc: cancel edit"
	} else {
		editHelp = "↑/k: up • ↓/j: down • Enter: edit field • s: save project • Esc: back"
	}

	s += "\n" + helpStyle.Render(editHelp)

	return s
}

// viewConfig renders the configuration view
func (m Model) viewConfig() string {
	s := titleStyle.Render("⚙️  Configuration") + "\n\n"

	s += normalStyle.Render(fmt.Sprintf("Session Directory: %s", config.SessionDir())) + "\n"
	s += normalStyle.Render(fmt.Sprintf("Config Directory: %s", config.ConfigDir())) + "\n"

	s += "\n" + helpStyle.Render("Esc/q: back")

	return s
}

// viewHelp renders the help view
func (m Model) viewHelp() string {
	s := titleStyle.Render("❓ Help") + "\n\n"

	s += normalStyle.Render("Project List View:") + "\n"
	s += normalStyle.Render("  ↑/k, ↓/j   - Navigate") + "\n"
	s += normalStyle.Render("  Enter      - Edit project") + "\n"
	s += normalStyle.Render("  n          - New project") + "\n"
	s += normalStyle.Render("  d          - Archive project") + "\n"
	s += normalStyle.Render("  r          - Refresh list") + "\n"
	s += normalStyle.Render("  c          - Configuration") + "\n"
	s += normalStyle.Render("  ?/h        - This help") + "\n"
	s += normalStyle.Render("  q          - Quit") + "\n\n"

	s += normalStyle.Render("Project Edit View:") + "\n"
	s += normalStyle.Render("  ↑/k, ↓/j   - Navigate fields") + "\n"
	s += normalStyle.Render("  s          - Save changes") + "\n"
	s += normalStyle.Render("  Esc        - Cancel") + "\n"

	s += "\n" + helpStyle.Render("Press Esc or ? to return")

	return s
}

// RunTUIApp starts the TUI application
func RunTUIApp() {
	utils.LogStart()
	defer utils.LogEnd()

	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
