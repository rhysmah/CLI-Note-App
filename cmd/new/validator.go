package new

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rhysmah/CLI-Note-App/models"
	"github.com/rhysmah/CLI-Note-App/validator"
)

const (
	illegalChars     string = "\\/:*?\"<>|: ."
	noteNameMaxLimit int    = 20
	noteNameMinLimit int    = 1
	dateTimeFormat   string = "2006_01_02_15_04"
)

// newValidator creates and returns a new validator for Note objects
// with predefined validation rules.
func newValidator() *validator.Validator[models.Note] {
	return &validator.Validator[models.Note]{
		Rules: []validator.ValidationRule[models.Note]{
			validateNoteTitleLength,
			validateNoteTitleCharacters,
		},
	}
}

// validateNoteTitle checks if the note's title meets the required criteria:
// - Does not exceed the character limit
// - Does not contain illegal characters
// Returns an error if validation fails.
func validateNoteTitleLength(note models.Note) error {
	noteNameTrimmed := strings.TrimSpace(note.Title)

	if len(noteNameTrimmed) < noteNameMinLimit {
		errMsg := fmt.Sprintf("note name %q must be greater than %d character", noteNameTrimmed, noteNameMinLimit)
		return errors.New(errMsg)
	}
	if len(noteNameTrimmed) > noteNameMaxLimit {
		errMsg := fmt.Sprintf("note name %q must be less than %d characters", noteNameTrimmed, noteNameMaxLimit)
		return errors.New(errMsg)
	}
	return nil
}

// checkForIllegalCharacters verifies that the note name doesn't contain any
// forbidden characters defined in illegalChars. Returns an error listing any
// illegal characters found.
func validateNoteTitleCharacters(note models.Note) error {
	var illegalCharsFound []rune

	for _, char := range note.Title {
		if strings.ContainsRune(illegalChars, char) {
			illegalCharsFound = append(illegalCharsFound, char)
		}
	}
	if len(illegalCharsFound) > 0 {
		errMsg := fmt.Sprintf("name contains illegal characters: %q", string(illegalCharsFound))
		return errors.New(errMsg)
	}
	return nil
}
