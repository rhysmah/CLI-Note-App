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

func getLongestFileName(notes []models.Note) int {
	longestName := 0

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

func calculateRowLineLength(notes []models.Note) int {
	fileNameWidth := max(len(headerFileName), getLongestFileName(notes))
	return fileNameWidth + (dateTimeWidth * 2) + (len(separator) * 3)
}

func printHeader(sort SortBy, order SortOrder, rowLineLength int) {
	sortName := getSortString(sort)
	orderName := getOrderString(sort, order)

	rowLine := strings.Repeat(lineSymbol, rowLineLength)
	header := fmt.Sprintf("Notes sorted by %s (%s)", sortName, orderName)

	fmt.Printf("%s\n%s\n%s\n", rowLine, header, rowLine)
}

func DisplayNotes(notes []models.Note, sort SortBy, order SortOrder) {
	rowLineLength := calculateRowLineLength(notes)
	longestFileNameLength := getLongestFileName(notes)

	printHeader(sort, order, rowLineLength)
	printNotesTable(notes, rowLineLength, longestFileNameLength)
}

func printNotesTable(notes []models.Note, rowLineLength, longestFileNameLength int) {
	// Create the divider line
	rowLine := strings.Repeat(lineSymbol, rowLineLength)

	// Print the header row with consistent formatting
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
