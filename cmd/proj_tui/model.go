package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/project"
	"github.com/jvzantvoort/tmux-project/projecttype"
	"github.com/jvzantvoort/tmux-project/tmux"
	"github.com/jvzantvoort/tmux-project/utils"
)

type TabView int

const (
	TabProjects TabView = iota
	TabProjectTypes
)

type DialogMode int

const (
	DialogNone DialogMode = iota
	DialogCreateProject
	DialogCreateProjectType
	DialogHelp
)

type Model struct {
	currentTab           TabView
	projectList          []project.Project
	projectTypeList      []projecttype.ProjectTypeConfig
	selectedProjectIndex int
	selectedTypeIndex    int
	width                int
	height               int
	err                  error
	message              string
	dialogMode           DialogMode
	dialogInputs         []textinput.Model
	dialogFocusIndex     int
	typeDropdownOpen     bool
	typeDropdownSelected int
	availableTypes       []string
}

func NewModel() Model {
	projects := loadProjects()
	projectTypes := loadProjectTypes()
	availableTypes := getProjectTypeNames(projectTypes)

	inputs := make([]textinput.Model, 3)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Project name (required)"
	inputs[0].CharLimit = 50

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Description (optional)"
	inputs[1].CharLimit = 100

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Type"
	inputs[2].CharLimit = 30

	return Model{
		currentTab:           TabProjects,
		projectList:          projects,
		projectTypeList:      projectTypes,
		selectedProjectIndex: 0,
		selectedTypeIndex:    0,
		width:                80,
		height:               24,
		dialogInputs:         inputs,
		availableTypes:       availableTypes,
	}
}

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

func loadProjectTypes() []projecttype.ProjectTypeConfig {
	var types []projecttype.ProjectTypeConfig
	inputdir := config.ConfigDir()

	if err := utils.MkdirAll(inputdir); err != nil {
		return types
	}

	targets, err := os.ReadDir(inputdir)
	if err != nil {
		return types
	}

	for _, target := range targets {
		if target.IsDir() {
			obj, err := projecttype.New(target.Name())
			if err == nil {
				types = append(types, obj)
			}
		}
	}

	return types
}

func getProjectTypeNames(types []projecttype.ProjectTypeConfig) []string {
	names := make([]string, len(types))
	for i, t := range types {
		names[i] = t.ProjectType
	}
	return names
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dialogMode != DialogNone {
		return m.handleDialogUpdate(msg)
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

func (m Model) handleDialogUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.dialogMode == DialogHelp {
			switch msg.String() {
			case "esc", "q", "?":
				m.dialogMode = DialogNone
				return m, nil
			}
			return m, nil
		}

		if m.typeDropdownOpen && m.dialogMode == DialogCreateProject {
			switch msg.String() {
			case "up", "k":
				if m.typeDropdownSelected > 0 {
					m.typeDropdownSelected--
				}
				return m, nil
			case "down", "j":
				if m.typeDropdownSelected < len(m.availableTypes)-1 {
					m.typeDropdownSelected++
				}
				return m, nil
			case "enter":
				m.dialogInputs[2].SetValue(m.availableTypes[m.typeDropdownSelected])
				m.typeDropdownOpen = false
				return m, nil
			case "esc":
				m.typeDropdownOpen = false
				return m, nil
			}
			return m, nil
		}

		switch msg.String() {
		case "esc":
			m.dialogMode = DialogNone
			m.clearDialogInputs()
			m.message = ""
			return m, nil

		case "tab", "down":
			m.dialogFocusIndex++
			if m.dialogFocusIndex >= len(m.dialogInputs) {
				m.dialogFocusIndex = 0
			}
			m.updateDialogFocus()
			return m, nil

		case "shift+tab", "up":
			m.dialogFocusIndex--
			if m.dialogFocusIndex < 0 {
				m.dialogFocusIndex = len(m.dialogInputs) - 1
			}
			m.updateDialogFocus()
			return m, nil

		case "enter":
			if m.dialogMode == DialogCreateProject && m.dialogFocusIndex == 2 {
				m.typeDropdownOpen = true
				for i, t := range m.availableTypes {
					if t == m.dialogInputs[2].Value() {
						m.typeDropdownSelected = i
						break
					}
				}
				return m, nil
			}

		case "ctrl+s":
			return m.handleDialogSubmit()
		}

		var cmd tea.Cmd
		m.dialogInputs[m.dialogFocusIndex], cmd = m.dialogInputs[m.dialogFocusIndex].Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m *Model) handleDialogSubmit() (tea.Model, tea.Cmd) {
	switch m.dialogMode {
	case DialogCreateProject:
		return m.createProject()
	case DialogCreateProjectType:
		return m.createProjectType()
	}
	return m, nil
}

func (m *Model) createProject() (tea.Model, tea.Cmd) {
	name := m.dialogInputs[0].Value()
	description := m.dialogInputs[1].Value()
	projectType := m.dialogInputs[2].Value()

	if name == "" {
		m.message = "Error: Project name is required"
		return m, nil
	}

	if projectType == "" {
		projectType = "default"
	}

	proj := project.NewProject(name)
	proj.SetDescription(description)
	err := proj.InitializeProject(projectType, true)
	if err != nil {
		m.message = fmt.Sprintf("Error creating project: %v", err)
		return m, nil
	}

	m.message = fmt.Sprintf("Project '%s' created successfully", name)
	m.dialogMode = DialogNone
	m.clearDialogInputs()
	m.projectList = loadProjects()

	return m, nil
}

func (m *Model) createProjectType() (tea.Model, tea.Cmd) {
	name := m.dialogInputs[0].Value()

	if name == "" {
		m.message = "Error: Project type name is required"
		return m, nil
	}

	ptc, err := projecttype.New(name)
	if err != nil && err != projecttype.ErrProjectNotExists {
		m.message = fmt.Sprintf("Error: %v", err)
		return m, nil
	}

	err = ptc.Setup()
	if err != nil {
		m.message = fmt.Sprintf("Error creating project type: %v", err)
		return m, nil
	}

	m.message = fmt.Sprintf("Project type '%s' created successfully", name)
	m.dialogMode = DialogNone
	m.clearDialogInputs()
	m.projectTypeList = loadProjectTypes()
	m.availableTypes = getProjectTypeNames(m.projectTypeList)

	return m, nil
}

func (m *Model) clearDialogInputs() {
	for i := range m.dialogInputs {
		m.dialogInputs[i].SetValue("")
	}
	m.dialogFocusIndex = 0
	m.typeDropdownOpen = false
}

func (m *Model) updateDialogFocus() {
	for i := range m.dialogInputs {
		if i == m.dialogFocusIndex {
			m.dialogInputs[i].Focus()
		} else {
			m.dialogInputs[i].Blur()
		}
	}
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "1":
		m.currentTab = TabProjects
		m.message = ""

	case "2":
		m.currentTab = TabProjectTypes
		m.message = ""

	case "tab":
		if m.currentTab == TabProjects {
			m.currentTab = TabProjectTypes
		} else {
			m.currentTab = TabProjects
		}
		m.message = ""

	case "?":
		m.dialogMode = DialogHelp

	case "+":
		if m.currentTab == TabProjects {
			m.dialogMode = DialogCreateProject
			m.dialogFocusIndex = 0
			m.updateDialogFocus()
			if len(m.availableTypes) > 0 {
				m.dialogInputs[2].SetValue(m.availableTypes[0])
			}
		} else {
			m.dialogMode = DialogCreateProjectType
			m.dialogFocusIndex = 0
			m.updateDialogFocus()
		}
		m.message = ""
	}

	if m.currentTab == TabProjects {
		return m.handleProjectsTabKeys(msg)
	}
	return m.handleProjectTypesTabKeys(msg)
}

func (m Model) handleProjectsTabKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.selectedProjectIndex > 0 {
			m.selectedProjectIndex--
		}

	case "down", "j":
		if m.selectedProjectIndex < len(m.projectList)-1 {
			m.selectedProjectIndex++
		}

	case "enter":
		if len(m.projectList) > 0 && m.selectedProjectIndex < len(m.projectList) {
			proj := m.projectList[m.selectedProjectIndex]
			sessionname := proj.Name

			if err := proj.Open(); err == nil {
				proj.UpdateLastActivity()
			}
			configfile := filepath.Join(config.SessionDir(), proj.Name+".rc")

			tmux.Resume(sessionname, configfile)
			return m, tea.Quit
		}

	case "e":
		if len(m.projectList) > 0 && m.selectedProjectIndex < len(m.projectList) {
			proj := m.projectList[m.selectedProjectIndex]
			opts := []string{"-O"}
			opts = append(opts, filepath.Join(config.SessionDir(), proj.Name+".rc"))
			opts = append(opts, filepath.Join(config.SessionDir(), proj.Name+".env"))

			utils.Edit(opts...)
			m.projectList = loadProjects()
		}

	case "r":
		m.projectList = loadProjects()
		m.message = "Projects list refreshed"
	}

	return m, nil
}

func (m Model) handleProjectTypesTabKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.selectedTypeIndex > 0 {
			m.selectedTypeIndex--
		}

	case "down", "j":
		if m.selectedTypeIndex < len(m.projectTypeList)-1 {
			m.selectedTypeIndex++
		}

	case "e":
		if len(m.projectTypeList) > 0 && m.selectedTypeIndex < len(m.projectTypeList) {
			pt := m.projectTypeList[m.selectedTypeIndex]
			opts := []string{"-O"}
			opts = append(opts, filepath.Join(pt.ConfigDir, "config.yml"))
			opts = append(opts, filepath.Join(pt.ConfigDir, "default.rc"))
			opts = append(opts, filepath.Join(pt.ConfigDir, "default.env"))

			utils.Edit(opts...)
			m.projectTypeList = loadProjectTypes()
		}

	case "r":
		m.projectTypeList = loadProjectTypes()
		m.availableTypes = getProjectTypeNames(m.projectTypeList)
		m.message = "Project types list refreshed"
	}

	return m, nil
}

func (m Model) View() string {
	if m.dialogMode != DialogNone {
		return m.renderDialog()
	}

	return m.renderMainView()
}

func (m Model) renderMainView() string {
	var content string

	content += m.renderTabs() + "\n\n"

	if m.currentTab == TabProjects {
		content += m.renderProjectsTab()
	} else {
		content += m.renderProjectTypesTab()
	}

	if m.message != "" {
		content += "\n\n" + messageStyle.Render(m.message)
	}

	return content
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Padding(0, 2)

	selectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7D56F4")).
				Background(lipgloss.Color("#3C3C3C"))

	normalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Width(60)

	dialogTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7D56F4")).
				MarginBottom(1)

	inputLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	focusedInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4"))

	dropdownStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginTop(1)
)

func (m Model) renderTabs() string {
	projectsTab := inactiveTabStyle.Render("Projects")
	typesTab := inactiveTabStyle.Render("Project Types")

	if m.currentTab == TabProjects {
		projectsTab = activeTabStyle.Render("Projects")
	} else {
		typesTab = activeTabStyle.Render("Project Types")
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, projectsTab, typesTab)
}

func (m Model) renderProjectsTab() string {
	active, _ := tmux.ListActive()
	activeMap := make(map[string]bool)
	for _, name := range active {
		activeMap[name] = true
	}

	s := titleStyle.Render(" 📋 Projects ") + "\n\n"

	if len(m.projectList) == 0 {
		s += normalItemStyle.Render("No projects found. Press '+' to create one.") + "\n"
	} else {
		for i, proj := range m.projectList {
			cursor := "  "
			style := normalItemStyle

			if i == m.selectedProjectIndex {
				cursor = "› "
				style = selectedItemStyle
			}

			status := ""
			if activeMap[proj.Name] {
				status = " ●"
			}

			line := fmt.Sprintf("%s%-25s  %-12s  %s%s",
				cursor, proj.Name, proj.ProjectType, proj.Description, status)
			s += style.Render(line) + "\n"
		}
	}

	s += "\n" + helpStyle.Render(
		"↑/k: up • ↓/j: down • Enter: resume • e: edit • +: create • r: refresh • Tab/1/2: switch tab • ?: help • q: quit",
	)

	return s
}

func (m Model) renderProjectTypesTab() string {
	s := titleStyle.Render(" 🔧 Project Types ") + "\n\n"

	if len(m.projectTypeList) == 0 {
		s += normalItemStyle.Render("No project types found. Press '+' to create one.") + "\n"
	} else {
		for i, pt := range m.projectTypeList {
			cursor := "  "
			style := normalItemStyle

			if i == m.selectedTypeIndex {
				cursor = "› "
				style = selectedItemStyle
			}

			line := fmt.Sprintf("%s%-25s  %d actions  %d targets  %d repos",
				cursor, pt.ProjectType, len(pt.SetupActions), len(pt.Targets), len(pt.Repos))
			s += style.Render(line) + "\n"
		}
	}

	s += "\n" + helpStyle.Render(
		"↑/k: up • ↓/j: down • e: edit • +: create • r: refresh • Tab/1/2: switch tab • ?: help • q: quit",
	)

	return s
}

func (m Model) renderDialog() string {
	if m.dialogMode == DialogHelp {
		return m.renderHelpDialog()
	}

	var content string

	if m.dialogMode == DialogCreateProject {
		content = dialogTitleStyle.Render("Create New Project") + "\n\n"

		labels := []string{"Name:", "Description:", "Type:"}
		for i, label := range labels {
			labelStr := inputLabelStyle.Render(label)
			if i == m.dialogFocusIndex {
				labelStr = focusedInputStyle.Render("› " + label)
			} else {
				labelStr = "  " + labelStr
			}

			content += labelStr + "\n"
			content += "  " + m.dialogInputs[i].View() + "\n\n"

			if i == 2 && m.typeDropdownOpen && i == m.dialogFocusIndex {
				dropdownContent := ""
				for j, t := range m.availableTypes {
					if j == m.typeDropdownSelected {
						dropdownContent += selectedItemStyle.Render("› "+t) + "\n"
					} else {
						dropdownContent += normalItemStyle.Render("  "+t) + "\n"
					}
				}
				content += dropdownStyle.Render(dropdownContent) + "\n"
			}
		}

		content += helpStyle.Render("Tab/↑/↓: navigate • Ctrl+S: create • Esc: cancel")

	} else if m.dialogMode == DialogCreateProjectType {
		content = dialogTitleStyle.Render("Create New Project Type") + "\n\n"

		labelStr := inputLabelStyle.Render("Name:")
		if m.dialogFocusIndex == 0 {
			labelStr = focusedInputStyle.Render("› Name:")
		} else {
			labelStr = "  " + labelStr
		}

		content += labelStr + "\n"
		content += "  " + m.dialogInputs[0].View() + "\n\n"

		content += helpStyle.Render("Ctrl+S: create • Esc: cancel")
	}

	if m.message != "" {
		content += "\n" + errorStyle.Render(m.message)
	}

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(content))
}

func (m Model) renderHelpDialog() string {
	content := dialogTitleStyle.Render("Help") + "\n\n"

	content += normalItemStyle.Bold(true).Render("Global Keys:") + "\n"
	content += normalItemStyle.Render("  Tab / 1 / 2    Switch tabs") + "\n"
	content += normalItemStyle.Render("  ?              Show this help") + "\n"
	content += normalItemStyle.Render("  q / Ctrl+C     Quit") + "\n\n"

	content += normalItemStyle.Bold(true).Render("Projects Tab:") + "\n"
	content += normalItemStyle.Render("  ↑/k, ↓/j       Navigate") + "\n"
	content += normalItemStyle.Render("  Enter          Resume project") + "\n"
	content += normalItemStyle.Render("  e              Edit project") + "\n"
	content += normalItemStyle.Render("  +              Create new project") + "\n"
	content += normalItemStyle.Render("  r              Refresh list") + "\n\n"

	content += normalItemStyle.Bold(true).Render("Project Types Tab:") + "\n"
	content += normalItemStyle.Render("  ↑/k, ↓/j       Navigate") + "\n"
	content += normalItemStyle.Render("  e              Edit project type") + "\n"
	content += normalItemStyle.Render("  +              Create new project type") + "\n"
	content += normalItemStyle.Render("  r              Refresh list") + "\n\n"

	content += normalItemStyle.Bold(true).Render("Create Dialogs:") + "\n"
	content += normalItemStyle.Render("  Tab / ↑ / ↓    Navigate fields") + "\n"
	content += normalItemStyle.Render("  Ctrl+S         Submit") + "\n"
	content += normalItemStyle.Render("  Esc            Cancel") + "\n\n"

	content += helpStyle.Render("Press Esc or ? to close")

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Width(70).Render(content))
}
