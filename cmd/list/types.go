package list

type SortBy string
type SortOrder string

const (
	SortByModified SortBy = "modified"
	SortByCreated  SortBy = "created"
	SortByTitle    SortBy = "title"

	SortOrderAscending    SortOrder = "ascending"
	SortOrderDescending   SortOrder = "descending"
)

