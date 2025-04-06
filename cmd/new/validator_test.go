package new

import (
	"strings"
	"testing"

	"github.com/rhysmah/CLI-Note-App/models"
)

func TestValidateNoteTitleLength(t *testing.T) {
	tests := []struct {
		name      string
		noteTitle string
		wantErr   bool
	}{
		{
			name:      "Title Less Than Max Length",
			noteTitle: "",
			wantErr:   true,
		},
		{
			name:      "Title At Min Length",
			noteTitle: "a",
			wantErr:   false,
		},
		{
			name:      "Title At Max Length",
			noteTitle: strings.Repeat("a", noteNameMaxLimit),
			wantErr:   false,
		},
		{
			name:      "Title Great Than Max Length",
			noteTitle: strings.Repeat("a", noteNameMaxLimit*2),
			wantErr:   true,
		},
		{
			name:      "Valid Title",
			noteTitle: "new_note",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			note := models.Note{Title: tt.noteTitle}
			err := validateNoteTitleLength(note)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateNoteTitle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateNoteTitleCharacters(t *testing.T) {
	tests := []struct {
		name      string
		noteTitle string
		wantErr   bool
	}{
		{
			name:      "Title Has Illegal Forward Slash",
			noteTitle: "note/",
			wantErr:   true,
		},
		{
			name:      "Title Has Illegal Back Slash",
			noteTitle: "note\\",
			wantErr:   true,
		},
		{
			name:      "Title Has Illegal Back Colon",
			noteTitle: "note:",
			wantErr:   true,
		},
		{
			name:      "Title Has Illegal Bar",
			noteTitle: "note|",
			wantErr:   true,
		},
		{
			name:      "Title Has Illegal Question Mark",
			noteTitle: "note?",
			wantErr:   true,
		},
		{
			name:      "Title Has Illegal Star",
			noteTitle: "note*",
			wantErr:   true,
		},
		{
			name:      "Title Has Illegal Quote",
			noteTitle: "note\"",
			wantErr:   true,
		},
		{
			name:      "Title Has Illegal Period",
			noteTitle: "note.",
			wantErr:   true,
		},
		{
			name:      "Title Has Illegal Less Than",
			noteTitle: "note<",
			wantErr:   true,
		},
		{
			name:      "Title Has Illegal Greater Than",
			noteTitle: "note>",
			wantErr:   true,
		},
		{
			name:      "Title Has Valid Characters",
			noteTitle: "new_note",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			note := models.Note{Title: tt.noteTitle}
			err := validateNoteTitleCharacters(note)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateNoteTitle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNoteValidator(t *testing.T) {
	validator := newValidator()

	// Test that validator has x number of rules.
	// As of [04-03-2025]: 2 rules.
	if len(validator.Rules) != 2 {
		t.Errorf("Validator has %d rules; expected 2", len(validator.Rules))
	}

	// Test cases
	tests := []struct {
		name    string
		Note    models.Note
		wantErr bool
	}{
		{
			name:    "Valid note",
			Note:    models.Note{Title: "valid_note_name"},
			wantErr: false,
		},
		{
			name: "Valid note with spaces",
			Note: models.Note{Title: "title with spaces"},
			wantErr: false,
		},
		{
			name:    "Title too short (empty)",
			Note:    models.Note{Title: ""},
			wantErr: true,
		},
		{
			name:    "Title too long",
			Note:    models.Note{Title: strings.Repeat("a", noteNameMaxLimit+1)},
			wantErr: true,
		},
		{
			name:    "Title contains invalid characters",
			Note:    models.Note{Title: "note_with!_invalid:_characters/"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Run(tt.Note)

			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.Run error = %v; wanted %v", err, tt.wantErr)
			}
		})
	}
}
