package query

// VideoSearchQuery defines the parameters for searching videos.
// example: limit, offset, sort order, etc.
type VideoSearchQuery struct {
	Limit int
	// add more fields as needed
}

type VideoRangeQuery struct {
	Start int64
	End   *int64 // nil means until the end of the file
}
