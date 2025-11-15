package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jvzantvoort/tmux-project/project"
)

// TestModel_Init verifies model initialization
func TestModel_Init(t *testing.T) {
	m := NewModel()

	if m.viewMode != ViewProjectList {
		t.Errorf("Initial view mode should be ViewProjectList, got %v", m.viewMode)
	}

	if m.selectedIndex != 0 {
		t.Errorf("Initial selectedIndex should be 0, got %d", m.selectedIndex)
	}

	if len(m.textInputs) != 5 {
		t.Errorf("Should have 5 text inputs, got %d", len(m.textInputs))
	}

	if m.isEditing {
		t.Error("Should not be editing initially")
	}
}

// TestModel_TextInputsInitialized verifies text inputs are set up correctly
func TestModel_TextInputsInitialized(t *testing.T) {
	m := NewModel()

	expectedPlaceholders := []string{
		"Project name",
		"Project directory",
		"Description",
		"Project type",
		"Status",
	}

	for i, input := range m.textInputs {
		if input.Placeholder != expectedPlaceholders[i] {
			t.Errorf("Input %d placeholder mismatch: got %s, want %s",
				i, input.Placeholder, expectedPlaceholders[i])
		}
	}
}

// TestModel_LoadProjectToInputs verifies project data loads into inputs
func TestModel_LoadProjectToInputs(t *testing.T) {
	m := NewModel()

	proj := &project.Project{
		Name:        "test-project",
		Directory:   "/tmp/test",
		Description: "Test description",
		ProjectType: "golang",
		Status:      "active",
	}

	m.editingProject = proj
	m.loadProjectToInputs()

	if m.textInputs[0].Value() != "test-project" {
		t.Errorf("Name input mismatch: got %s, want test-project", m.textInputs[0].Value())
	}
	if m.textInputs[1].Value() != "/tmp/test" {
		t.Errorf("Directory input mismatch: got %s, want /tmp/test", m.textInputs[1].Value())
	}
	if m.textInputs[2].Value() != "Test description" {
		t.Errorf("Description input mismatch: got %s, want 'Test description'", m.textInputs[2].Value())
	}
	if m.textInputs[3].Value() != "golang" {
		t.Errorf("Type input mismatch: got %s, want golang", m.textInputs[3].Value())
	}
	if m.textInputs[4].Value() != "active" {
		t.Errorf("Status input mismatch: got %s, want active", m.textInputs[4].Value())
	}
}

// TestModel_UpdateProjectFromInputs verifies inputs update project
func TestModel_UpdateProjectFromInputs(t *testing.T) {
	m := NewModel()

	proj := &project.Project{}
	m.editingProject = proj

	m.textInputs[0].SetValue("new-project")
	m.textInputs[1].SetValue("/tmp/new")
	m.textInputs[2].SetValue("New description")
	m.textInputs[3].SetValue("python")
	m.textInputs[4].SetValue("archived")

	m.updateProjectFromInputs()

	if proj.Name != "new-project" {
		t.Errorf("Name not updated: got %s, want new-project", proj.Name)
	}
	if proj.Directory != "/tmp/new" {
		t.Errorf("Directory not updated: got %s, want /tmp/new", proj.Directory)
	}
	if proj.Description != "New description" {
		t.Errorf("Description not updated: got %s, want 'New description'", proj.Description)
	}
	if proj.ProjectType != "python" {
		t.Errorf("Type not updated: got %s, want python", proj.ProjectType)
	}
	if proj.Status != "archived" {
		t.Errorf("Status not updated: got %s, want archived", proj.Status)
	}
}

// TestModel_RoundTripProjectData verifies data integrity through load/update cycle
func TestModel_RoundTripProjectData(t *testing.T) {
	m := NewModel()

	original := &project.Project{
		Name:        "roundtrip",
		Directory:   "/home/user/roundtrip",
		Description: "Round trip test",
		ProjectType: "javascript",
		Status:      "active",
	}

	// Load to inputs
	m.editingProject = original
	m.loadProjectToInputs()

	// Create new project and update from inputs
	updated := &project.Project{}
	m.editingProject = updated
	m.updateProjectFromInputs()

	// Compare
	if updated.Name != original.Name {
		t.Errorf("Name mismatch: got %s, want %s", updated.Name, original.Name)
	}
	if updated.Directory != original.Directory {
		t.Errorf("Directory mismatch: got %s, want %s", updated.Directory, original.Directory)
	}
	if updated.Description != original.Description {
		t.Errorf("Description mismatch: got %s, want %s", updated.Description, original.Description)
	}
	if updated.ProjectType != original.ProjectType {
		t.Errorf("Type mismatch: got %s, want %s", updated.ProjectType, original.ProjectType)
	}
	if updated.Status != original.Status {
		t.Errorf("Status mismatch: got %s, want %s", updated.Status, original.Status)
	}
}

// TestModel_NavigationKeys verifies navigation in project list
func TestModel_NavigationKeys(t *testing.T) {
	m := NewModel()

	// Add some mock projects for navigation
	m.projects = []project.Project{
		{Name: "project1"},
		{Name: "project2"},
		{Name: "project3"},
	}
	m.selectedIndex = 0

	// Test down navigation
	msg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := m.handleProjectListKeys(msg)
	m = newModel.(Model)

	if m.selectedIndex != 1 {
		t.Errorf("Down key should move to index 1, got %d", m.selectedIndex)
	}

	// Test up navigation
	msg = tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = m.handleProjectListKeys(msg)
	m = newModel.(Model)

	if m.selectedIndex != 0 {
		t.Errorf("Up key should move to index 0, got %d", m.selectedIndex)
	}

	// Test boundary - up at top
	msg = tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = m.handleProjectListKeys(msg)
	m = newModel.(Model)

	if m.selectedIndex != 0 {
		t.Errorf("Up at top should stay at 0, got %d", m.selectedIndex)
	}

	// Test boundary - down at bottom
	m.selectedIndex = 2
	msg = tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ = m.handleProjectListKeys(msg)
	m = newModel.(Model)

	if m.selectedIndex != 2 {
		t.Errorf("Down at bottom should stay at 2, got %d", m.selectedIndex)
	}
}

// TestModel_VimNavigationKeys verifies vim-style j/k navigation
func TestModel_VimNavigationKeys(t *testing.T) {
	m := NewModel()
	m.projects = []project.Project{
		{Name: "project1"},
		{Name: "project2"},
	}
	m.selectedIndex = 0

	// Test 'j' for down
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	newModel, _ := m.handleProjectListKeys(msg)
	m = newModel.(Model)

	if m.selectedIndex != 1 {
		t.Errorf("'j' key should move down, got index %d", m.selectedIndex)
	}

	// Test 'k' for up
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	newModel, _ = m.handleProjectListKeys(msg)
	m = newModel.(Model)

	if m.selectedIndex != 0 {
		t.Errorf("'k' key should move up, got index %d", m.selectedIndex)
	}
}

// TestModel_EditFieldNavigation verifies navigation between edit fields
func TestModel_EditFieldNavigation(t *testing.T) {
	m := NewModel()
	m.viewMode = ViewProjectEdit
	m.editingProject = &project.Project{Name: "test"}
	m.editField = 0
	m.isEditing = false

	// Test moving down through fields
	msg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := m.handleProjectEditKeys(msg)
	m = newModel.(Model)

	if m.editField != 1 {
		t.Errorf("Should move to field 1, got %d", m.editField)
	}

	// Test moving up
	msg = tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = m.handleProjectEditKeys(msg)
	m = newModel.(Model)

	if m.editField != 0 {
		t.Errorf("Should move to field 0, got %d", m.editField)
	}
}

// TestModel_EditFieldNavigationWhileEditing verifies navigation is disabled while editing
func TestModel_EditFieldNavigationWhileEditing(t *testing.T) {
	m := NewModel()
	m.viewMode = ViewProjectEdit
	m.editingProject = &project.Project{Name: "test"}
	m.editField = 1
	m.isEditing = true // Currently editing

	// Try to move down - should not work
	msg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := m.handleProjectEditKeys(msg)
	m = newModel.(Model)

	if m.editField != 1 {
		t.Errorf("Should stay at field 1 while editing, got %d", m.editField)
	}
}

// TestModel_EnterEditMode verifies entering edit mode
func TestModel_EnterEditMode(t *testing.T) {
	m := NewModel()
	m.viewMode = ViewProjectEdit
	m.editingProject = &project.Project{Name: "test"}
	m.editField = 0
	m.isEditing = false

	// Press enter to start editing
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.handleProjectEditKeys(msg)
	m = newModel.(Model)

	if !m.isEditing {
		t.Error("Should be in editing mode after pressing enter")
	}
}

// TestModel_ViewModeTransitions verifies view mode changes
func TestModel_ViewModeTransitions(t *testing.T) {
	m := NewModel()

	// Start in project list
	if m.viewMode != ViewProjectList {
		t.Errorf("Should start in ViewProjectList, got %v", m.viewMode)
	}

	// Go to config
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
	newModel, _ := m.handleProjectListKeys(msg)
	m = newModel.(Model)

	if m.viewMode != ViewConfig {
		t.Errorf("'c' should go to ViewConfig, got %v", m.viewMode)
	}

	// Back from config
	msg = tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ = m.handleConfigKeys(msg)
	m = newModel.(Model)

	if m.viewMode != ViewProjectList {
		t.Errorf("Esc from config should return to ViewProjectList, got %v", m.viewMode)
	}

	// Go to help
	m.viewMode = ViewProjectList
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	newModel, _ = m.handleProjectListKeys(msg)
	m = newModel.(Model)

	if m.viewMode != ViewHelp {
		t.Errorf("'?' should go to ViewHelp, got %v", m.viewMode)
	}
}

// TestModel_EmptyProjectList verifies handling of empty project list
func TestModel_EmptyProjectList(t *testing.T) {
	m := NewModel()
	m.projects = []project.Project{} // Empty

	// Try to navigate - should not crash
	msg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := m.handleProjectListKeys(msg)
	m = newModel.(Model)

	if m.selectedIndex != 0 {
		t.Errorf("Empty list should keep index at 0, got %d", m.selectedIndex)
	}

	// Try to select - should not crash
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ = m.handleProjectListKeys(msg)
	m = newModel.(Model)

	// Should still be in list view, not edit
	if m.viewMode != ViewProjectList {
		t.Errorf("Enter on empty list should stay in ViewProjectList, got %v", m.viewMode)
	}
}

// TestModel_RefreshProjects verifies refresh functionality
func TestModel_RefreshProjects(t *testing.T) {
	m := NewModel()
	initialCount := len(m.projects)

	// Press 'r' to refresh
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	newModel, _ := m.handleProjectListKeys(msg)
	m = newModel.(Model)

	if m.message != "Projects refreshed" {
		t.Errorf("Refresh should set message, got: %s", m.message)
	}

	if m.selectedIndex != 0 {
		t.Errorf("Refresh should reset selectedIndex to 0, got %d", m.selectedIndex)
	}

	// Project count may vary, just verify it doesn't crash
	_ = initialCount
}

// TestModel_WindowResize verifies window resize handling
func TestModel_WindowResize(t *testing.T) {
	m := NewModel()

	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	newModel, _ := m.Update(msg)
	m = newModel.(Model)

	if m.width != 120 {
		t.Errorf("Width should be 120, got %d", m.width)
	}
	if m.height != 40 {
		t.Errorf("Height should be 40, got %d", m.height)
	}
}

// TestModel_LoadProjectToInputsNilProject verifies nil handling
func TestModel_LoadProjectToInputsNilProject(t *testing.T) {
	m := NewModel()
	m.editingProject = nil

	// Should not crash
	m.loadProjectToInputs()

	// Inputs should be unchanged
	for _, input := range m.textInputs {
		if input.Value() != "" {
			t.Errorf("Input should be empty with nil project, got: %s", input.Value())
		}
	}
}

// TestModel_UpdateProjectFromInputsNilProject verifies nil handling
func TestModel_UpdateProjectFromInputsNilProject(t *testing.T) {
	m := NewModel()
	m.editingProject = nil

	// Set some values
	m.textInputs[0].SetValue("test")

	// Should not crash
	m.updateProjectFromInputs()
	// No assertions needed, just verify it doesn't panic
}
