package list

import (
	"sort"

	"github.com/rhysmah/CLI-Note-App/models"
)

// sortFiles sorts a slice of files based on the specified field and order.
// It uses the compareFiles function to determine the ordering between any two files.
func sortNotes(notes []models.Note, field SortBy, order SortOrder) {
	sort.Slice(notes, func(a, b int) bool {
		return compareNotes(notes[a], notes[b], field, order)
	})
}

// compareNotes compares two files based on the specified sort field and order.
// It returns true if file 'a' should come before file 'b' in the sorted result.
func compareNotes(a, b models.Note, field SortBy, order SortOrder) bool {
	ascending := order == SortOrderAscending // Default order = ascending

	switch field {

	case SortByTitle:
		if ascending {
			return a.Title < b.Title // A - Z
		}
		return a.Title > b.Title // Z - A

	case SortByCreated:
		if ascending {
			return a.CreatedAt.Before(b.CreatedAt)
		}
		return a.CreatedAt.After(b.CreatedAt)

	case SortByModified:
		if ascending {
			return a.ModifiedAt.Before(b.ModifiedAt)
		}
		return a.ModifiedAt.After(b.ModifiedAt)

	default:
		return a.Title < b.Title
	}
}