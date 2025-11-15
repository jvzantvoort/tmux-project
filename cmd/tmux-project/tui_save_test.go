package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/jvzantvoort/tmux-project/project"
)

// TestTUIEditPersistence validates that the TUI edit functionality
// properly writes changes back to the configuration file
func TestTUIEditPersistence(t *testing.T) {
	// Create a unique test project name
	testName := "test-tui-persist"
	
	// Ensure cleanup
	defer func() {
		configFile := project.NewProject(testName).ProjectConfigFile()
		os.Remove(configFile)
	}()
	
	// Step 1: Create initial project
	proj := &project.Project{
		Name:        testName,
		ProjectType: "golang",
		Description: "Original description",
		Directory:   "/tmp/original-dir",
	}
	
	if err := proj.Save(); err != nil {
		t.Fatalf("Failed to create initial project: %v", err)
	}
	
	// Step 2: Load the project (simulating TUI loading)
	loadedProj := project.NewProject(testName)
	if err := loadedProj.Open(); err != nil {
		t.Fatalf("Failed to load project: %v", err)
	}
	
	// Verify initial values
	if loadedProj.Description != "Original description" {
		t.Errorf("Initial description mismatch: got %q, want %q", 
			loadedProj.Description, "Original description")
	}
	
	// Step 3: Modify values (simulating updateProjectFromInputs())
	loadedProj.Description = "Updated via TUI"
	loadedProj.Directory = "/tmp/updated-dir"
	loadedProj.ProjectType = "python"
	
	// Step 4: Save (simulating 's' key press)
	if err := loadedProj.Save(); err != nil {
		t.Fatalf("Failed to save project: %v", err)
	}
	
	// Step 5: Verify file was written correctly
	configFile := loadedProj.ProjectConfigFile()
	rawData, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	
	var fileData map[string]interface{}
	if err := json.Unmarshal(rawData, &fileData); err != nil {
		t.Fatalf("Failed to parse config file: %v", err)
	}
	
	// Verify JSON content
	if desc, ok := fileData["description"].(string); !ok || desc != "Updated via TUI" {
		t.Errorf("File description mismatch: got %v, want %q", fileData["description"], "Updated via TUI")
	}
	if dir, ok := fileData["directory"].(string); !ok || dir != "/tmp/updated-dir" {
		t.Errorf("File directory mismatch: got %v, want %q", fileData["directory"], "/tmp/updated-dir")
	}
	if ptype, ok := fileData["type"].(string); !ok || ptype != "python" {
		t.Errorf("File type mismatch: got %v, want %q", fileData["type"], "python")
	}
	
	// Step 6: Load again to verify persistence
	verifyProj := project.NewProject(testName)
	if err := verifyProj.Open(); err != nil {
		t.Fatalf("Failed to re-load project: %v", err)
	}
	
	// Verify all changes persisted
	if verifyProj.Description != "Updated via TUI" {
		t.Errorf("Description not persisted: got %q, want %q", 
			verifyProj.Description, "Updated via TUI")
	}
	if verifyProj.Directory != "/tmp/updated-dir" {
		t.Errorf("Directory not persisted: got %q, want %q", 
			verifyProj.Directory, "/tmp/updated-dir")
	}
	if verifyProj.ProjectType != "python" {
		t.Errorf("ProjectType not persisted: got %q, want %q", 
			verifyProj.ProjectType, "python")
	}
}

// TestUpdateProjectFromInputs validates the updateProjectFromInputs function
func TestUpdateProjectFromInputs(t *testing.T) {
	// Create a mock model with a project
	proj := &project.Project{
		Name:        "test",
		ProjectType: "golang",
		Description: "original",
		Directory:   "/original",
	}
	
	m := &Model{
		editingProject: proj,
		textInputs:     make([]textinput.Model, 4),
	}
	
	// Set mock input values
	m.textInputs[0].SetValue("new-name")
	m.textInputs[1].SetValue("/new/directory")
	m.textInputs[2].SetValue("new description")
	m.textInputs[3].SetValue("python")
	
	// Call the function
	m.updateProjectFromInputs()
	
	// Verify updates
	if proj.Name != "new-name" {
		t.Errorf("Name not updated: got %q, want %q", proj.Name, "new-name")
	}
	if proj.Directory != "/new/directory" {
		t.Errorf("Directory not updated: got %q, want %q", proj.Directory, "/new/directory")
	}
	if proj.Description != "new description" {
		t.Errorf("Description not updated: got %q, want %q", proj.Description, "new description")
	}
	if proj.ProjectType != "python" {
		t.Errorf("ProjectType not updated: got %q, want %q", proj.ProjectType, "python")
	}
}
