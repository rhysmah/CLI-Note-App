package new

import (
	"testing"

	"github.com/rhysmah/CLI-Note-App/models"
)

func TestValidateNoteFIle(t *testing.T) {

	tests := []struct {
		name      string
		noteTitle string
		wantErr   bool
	}{
		{
			name:      "Empty Title",
			noteTitle: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := models.Note{Title: tt.noteTitle}
			err := validateNoteTitle(note)

			// Check if the error result matches our expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("validateNoteTitle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
