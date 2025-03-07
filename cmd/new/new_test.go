package new

import "testing"

const (
	testNoteTitle        = "new_note"
	testInvalidNoteTitle = "new:note"
)

func TestNewNote(t *testing.T) {
	note, err := createNote(testNoteTitle)
	if err != nil {
		t.Errorf("Couldn't create note: %v", err)
	}

	// Verify all fields
	if note.ID == "" {
		t.Errorf("Note ID not properly set; got %s", note.ID)
	}

	if note.Title != testNoteTitle {
		t.Errorf("Note title not correct; got %s", note.Title)
	}

	if note.Content != "" {
		t.Errorf("Note content not empty; got %s", note.Content)
	}

	if note.CreatedAt.IsZero() {
		t.Errorf("Created timestamp not set")
	}

	if note.ModifiedAt.IsZero() {
		t.Errorf("Modified timestamp not set")
	}

	if len(note.Tags) != 0 {
		t.Errorf("Should be no tags; got %v", note.Tags)
	}
}

func TestInvalidNote(t *testing.T) {
	_, err := createNote(testInvalidNoteTitle)
	if err == nil {
		t.Errorf("Note title invalid; should have thrown error; got nil")
	}
}
