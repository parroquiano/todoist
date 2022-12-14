package todoist

import (
	"context"

	"github.com/google/uuid"
)

// SectionsService handles communication with the sections related
// methods of the Todoist API.
//
// Todoist API docs: https://developer.todoist.com/sync/v8/?shell#sections
type SectionsService service

// Section represents a Todoist section.
type Section struct {
	// The ID of the section.
	ID int `json:"id"`

	// The name of the section.
	Name string `json:"name"`

	// Project that the section resides in.
	ProjectID int `json:"project_id"`

	// Legacy project ID for the project that the section resides in.
	// (only shown for objects created before 1 April 2017)
	LegacyProjectID *int `json:"legacy_project_id"`

	// The order of the section. Defines the position of the section among all the sections in the project.
	SectionOrder int `json:"section_order"`

	// Whether the section's tasks are collapsed (a true or false value).
	Collapsed bool `json:"collapsed"`

	// A special ID for shared sections (a number or null if not set). Used internally and can be ignored.
	SyncID *int `json:"sync_id"`

	// Whether the section is marked as deleted (a true or false value).
	IsDeleted bool `json:"is_deleted"`

	// Whether the section is marked as archived (a true or false value).
	IsArchived bool `json:"is_archived"`

	// The date when the section was archived (or null if not archived).
	DateArchived *string `json:"date_archived"`

	// The date when the section was created.
	DateAdded string `json:"date_added"`
}

func (s *SectionsService) List(ctx context.Context, syncToken string) ([]Section, ReadResponse, error) {
	s.client.Logln("---------- Sections.List")

	req, err := s.client.NewRequest(syncToken, []string{"sections"}, nil)
	if err != nil {
		return nil, ReadResponse{}, err
	}

	var readResponse ReadResponse
	_, err = s.client.Do(ctx, req, &readResponse)
	if err != nil {
		return nil, readResponse, err
	}

	return readResponse.Sections, readResponse, nil
}

type AddSection struct {
	// The name of the section.
	Name string `json:"name"`

	// The ID of the parent project.
	ProjectID int `json:"project_id"`

	// The order of the section. Defines the position of the section among all the sections in the project.
	SectionOrder int `json:"section_order,omitempty"`

	TempID string `json:"-"`
}

// Add a new section to a project.
func (s *SectionsService) Add(ctx context.Context, syncToken string, addSection AddSection) ([]Section, CommandResponse, error) {
	s.client.Logln("---------- Sections.Add")

	id := uuid.New().String()
	tempID := addSection.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	addCommand := Command{
		Type:   "section_add",
		Args:   addSection,
		UUID:   id,
		TempID: tempID,
	}

	commands := []Command{addCommand}

	req, err := s.client.NewRequest(syncToken, []string{"sections"}, commands)
	if err != nil {
		return nil, CommandResponse{}, err
	}

	var commandResponse CommandResponse
	_, err = s.client.Do(ctx, req, &commandResponse)
	if err != nil {
		return nil, commandResponse, err
	}

	return commandResponse.Sections, commandResponse, nil
}

type UpdateSection struct {
	// The ID of the section.
	ID string `json:"id"`

	// The name of the section.
	Name string `json:"name,omitempty"`

	// Whether the section's tasks are collapsed (a true or false value).
	Collapsed bool `json:"collapsed"`

	TempID string `json:"-"`
}

// Updates section attributes.
func (s *SectionsService) Update(ctx context.Context, syncToken string, updateSection UpdateSection) ([]Section, CommandResponse, error) {
	s.client.Logln("---------- Sections.Update")

	id := uuid.New().String()
	tempID := updateSection.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	updateCommand := Command{
		Type:   "section_update",
		Args:   updateSection,
		UUID:   id,
		TempID: tempID,
	}

	commands := []Command{updateCommand}

	req, err := s.client.NewRequest(syncToken, []string{"sections"}, commands)
	if err != nil {
		return nil, CommandResponse{}, err
	}

	var commandResponse CommandResponse
	_, err = s.client.Do(ctx, req, &commandResponse)
	if err != nil {
		return nil, commandResponse, err
	}

	return commandResponse.Sections, commandResponse, nil
}

type MoveSection struct {
	// The ID of the section.
	ID string `json:"id"`

	// ID of the destination project.
	ProjectID string `json:"project_id,omitempty"`

	TempID string `json:"-"`
}

// Move section to a different location.
func (s *SectionsService) Move(ctx context.Context, syncToken string, moveSection MoveSection) ([]Section, CommandResponse, error) {
	s.client.Logln("---------- Sections.Move")

	id := uuid.New().String()
	tempID := moveSection.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	moveCommand := Command{
		Type:   "section_move",
		Args:   moveSection,
		UUID:   id,
		TempID: tempID,
	}

	commands := []Command{moveCommand}

	req, err := s.client.NewRequest(syncToken, []string{"sections"}, commands)
	if err != nil {
		return nil, CommandResponse{}, err
	}

	var commandResponse CommandResponse
	_, err = s.client.Do(ctx, req, &commandResponse)
	if err != nil {
		return nil, commandResponse, err
	}

	return commandResponse.Sections, commandResponse, nil
}

type ReorderedSection struct {
	// ID of the section to update.
	ID string `json:"id"`

	// The new order.
	SectionOrder int `json:"section_order"`
}

type ReorderSections struct {
	// An array of objects to update. Each object contains two attributes: id of the section to update and section_order, the new order.
	Sections []ReorderedSection `json:"sections"`

	TempID string `json:"-"`
}

// The command updates section_order properties of sections in bulk.
func (s *SectionsService) Reorder(ctx context.Context, syncToken string, reorderSections ReorderSections) ([]Section, CommandResponse, error) {
	s.client.Logln("---------- Sections.Reorder")

	id := uuid.New().String()
	tempID := reorderSections.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	reorderCommand := Command{
		Type:   "section_reorder",
		Args:   reorderSections,
		UUID:   id,
		TempID: tempID,
	}

	commands := []Command{reorderCommand}

	req, err := s.client.NewRequest(syncToken, []string{"sections"}, commands)
	if err != nil {
		return nil, CommandResponse{}, err
	}

	var commandResponse CommandResponse
	_, err = s.client.Do(ctx, req, &commandResponse)
	if err != nil {
		return nil, commandResponse, err
	}

	return commandResponse.Sections, commandResponse, nil
}

type DeleteSection struct {
	// ID of the section to delete.
	ID string `json:"id"`

	TempID string `json:"-"`
}

// Delete a section and all its child tasks.
func (s *SectionsService) Delete(ctx context.Context, syncToken string, deleteSection DeleteSection) ([]Section, CommandResponse, error) {
	s.client.Logln("---------- Sections.Delete")

	id := uuid.New().String()
	tempID := deleteSection.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	deleteCommand := Command{
		Type:   "section_delete",
		Args:   deleteSection,
		UUID:   id,
		TempID: tempID,
	}

	commands := []Command{deleteCommand}

	req, err := s.client.NewRequest(syncToken, []string{"sections"}, commands)
	if err != nil {
		return nil, CommandResponse{}, err
	}

	var commandResponse CommandResponse
	_, err = s.client.Do(ctx, req, &commandResponse)
	if err != nil {
		return nil, commandResponse, err
	}

	return commandResponse.Sections, commandResponse, nil
}

type ArchiveSection struct {
	// Section ID to archive.
	ID string `json:"id"`

	TempID string `json:"-"`
}

// Archive a section and all its child tasks.
func (s *SectionsService) Archive(ctx context.Context, syncToken string, archiveSection ArchiveSection) ([]Section, CommandResponse, error) {
	s.client.Logln("---------- Sections.Archive")

	id := uuid.New().String()
	tempID := archiveSection.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	archiveCommand := Command{
		Type:   "section_archive",
		Args:   archiveSection,
		UUID:   id,
		TempID: tempID,
	}

	commands := []Command{archiveCommand}

	req, err := s.client.NewRequest(syncToken, []string{"sections"}, commands)
	if err != nil {
		return nil, CommandResponse{}, err
	}

	var commandResponse CommandResponse
	_, err = s.client.Do(ctx, req, &commandResponse)
	if err != nil {
		return nil, commandResponse, err
	}

	return commandResponse.Sections, commandResponse, nil
}

type UnarchiveSection struct {
	// Section ID to unarchive
	ID string `json:"id"`

	TempID string `json:"-"`
}

// This command is used to unarchive a section.
func (s *SectionsService) Unarchive(ctx context.Context, syncToken string, unarchiveSection UnarchiveSection) ([]Section, CommandResponse, error) {
	s.client.Logln("---------- Sections.Unarchive")

	id := uuid.New().String()
	tempID := unarchiveSection.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	unarchiveCommand := Command{
		Type:   "section_unarchive",
		Args:   unarchiveSection,
		UUID:   id,
		TempID: tempID,
	}

	commands := []Command{unarchiveCommand}

	req, err := s.client.NewRequest(syncToken, []string{"sections"}, commands)
	if err != nil {
		return nil, CommandResponse{}, err
	}

	var commandResponse CommandResponse
	_, err = s.client.Do(ctx, req, &commandResponse)
	if err != nil {
		return nil, commandResponse, err
	}

	return commandResponse.Sections, commandResponse, nil
}
