package list

import (
	"fmt"
	"strings"
	"time"

	"github.com/rhysmah/CLI-Note-App/models"
)

const (
	dateTimeFormat = "Jan 2, 2006 15:04"
	headerFileName = "File Name"
	headerCreated  = "Created Date"
	headerModified = "Modified Date"
)

const (
	dateTimeWidth   = 18
	paddingWidth    = 2
	separatorSymbol = "|"
	lineSymbol      = "-"
	separator       = "  |  "
)

// DisplayNotes renders a formatted table of notes.
// It displays notes based on the specified sort criteria and order..
//
// Parameters:
//   - notes: The slice of notes to display
//   - sort: The field by which to sort the notes (e.g., by name, date)
//   - order: The order in which to sort (ascending or descending)
func DisplayNotes(notes []models.Note, sort SortBy, order SortOrder) {
	rowLineLength := calculateRowLineLength(notes)
	longestFileNameLength := getLongestFileName(notes)

	printHeader(sort, order, rowLineLength)
	printNotesTable(notes, rowLineLength, longestFileNameLength)
}

// Calculates and returns the length of the longest file name,
// either the header -- "File Name" -- or the name of an actual
// file. This is used for table-formatting purposes.
func getLongestFileName(notes []models.Note) int {
	longestName := len(headerFileName)

	if len(notes) == 0 {
		return longestName
	}
	for _, note := range notes {
		if len(note.Title) > longestName {
			longestName = len(note.Title)
		}
	}
	return longestName
}

// Calculates the length of the table row dashes for table-formatting purposes.
func calculateRowLineLength(notes []models.Note) int {
	fileNameWidth := max(len(headerFileName), getLongestFileName(notes))
	return fileNameWidth + (dateTimeWidth * 2) + (len(separator) * 2)
}


// printHeader prints a formatted header for the notes list display.
// It shows how the notes are sorted (by date, title, etc.) and the sort order (ascending/descending).
//
// Parameters:
//   - sort: The field by which notes are sorted
//   - order: The direction of the sort (ascending/descending)
//   - rowLineLength: Length of the decorative lines surrounding the header
func printHeader(sort SortBy, order SortOrder, rowLineLength int) {
	sortName := getSortString(sort)
	orderName := getOrderString(sort, order)

	rowLine := strings.Repeat(lineSymbol, rowLineLength)
	header := fmt.Sprintf("Notes sorted by %s (%s)", sortName, orderName)

	fmt.Printf("%s\n%s\n%s\n", rowLine, header, rowLine)
}

// printNotesTable displays a formatted table of notes with columns for 
// filename, creation date, and modification date.
func printNotesTable(notes []models.Note, rowLineLength, longestFileNameLength int) {
	// Create the divider line
	rowLine := strings.Repeat(lineSymbol, rowLineLength)

	// Print header row
	fmt.Printf("%-*s%s%-*s%s%-*s\n",
		longestFileNameLength, headerFileName,
		separator, dateTimeWidth, headerCreated,
		separator, dateTimeWidth, headerModified)

	// Print divider after header
	fmt.Println(rowLine)

	// Print each note with the same formatting as the header
	for _, note := range notes {
		fmt.Printf("%-*s%s%-*s%s%-*s\n",
			longestFileNameLength, note.Title,
			separator, dateTimeWidth, formatDateTime(note.CreatedAt),
			separator, dateTimeWidth, formatDateTime(note.ModifiedAt))
	}
}

func formatDateTime(dt time.Time) string {
	return dt.Format(dateTimeFormat)
}

// Returns a string describing how the data is ordered.
// i.e., A - Z (if by title), newest to olders (if by a date)
func getOrderString(sort SortBy, order SortOrder) string {
	if sort == SortByTitle {
		if order == SortOrderAscending {
			return "A - Z"
		} else {
			return "Z - A"
		}
	} else {
		if order == SortOrderAscending {
			return "Newest to Oldest"
		} else {
			return "Oldest to Newest"
		}
	}
}

// Returns a string describing in what way the data is ordered
// i.e., title, creation date, modified date.
func getSortString(sort SortBy) string {
	switch sort {
	case SortByTitle:
		return "Title"
	case SortByCreated:
		return "Creation Date"
	case SortByModified:
		return "Modified Date"
	default:
		return ""
	}
}
