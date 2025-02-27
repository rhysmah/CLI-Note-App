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
			name:      "Title Has Illegal Space",
			noteTitle: "note ",
			wantErr:   true,
		},
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
			noteTitle: "note\"",
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
			err := validateNoteTitleIllegalCharacters(note)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateNoteTitle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
