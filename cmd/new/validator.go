package new

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rhysmah/CLI-Note-App/models"
	"github.com/rhysmah/CLI-Note-App/validator"
)

const (
	illegalChars      string = "\\/:*?\"<>|: ."
	noteNameCharLimit int    = 50
	dateTimeFormat    string = "2006_01_02_15_04"
)

func newValidator() *validator.Validator[models.Note] {
	return &validator.Validator[models.Note]{
		Rules: []validator.ValidationRule[models.Note]{
			validateNoteTitle,
		},
	}
}

func validateNoteTitle(note models.Note) error {
	noteNameTrimmed := strings.TrimSpace(note.Title)

	if len(noteNameTrimmed) > noteNameCharLimit {
		errMsg := fmt.Sprintf("name exceeds %d character limit", noteNameCharLimit)
		return errors.New(errMsg)
	}
	if err := checkForIllegalCharacters(noteNameTrimmed); err != nil {
		return fmt.Errorf("invalid characters in note name: %w", err)
	}
	return nil
}

func checkForIllegalCharacters(noteName string) error {
	var illegalCharsFound []rune

	for _, char := range noteName {
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
