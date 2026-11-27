// Package timeline builds pages of home timelines.
package timeline

import "slices"

// Page is one page of a home timeline.
type Page struct {
	// PostIDs holds the post IDs on the page, newest first.
	PostIDs []int64
	// NextCursor is the cursor for the next page. It is 0 on the last page.
	NextCursor int64
}

// Merge builds one page of a home timeline from several sources, such as
// the account's inbox, the recent posts of each celebrity account it
// follows, and its own recent posts.
//
// Each source lists post IDs newest first. Merge keeps only the IDs that
// are older than cursor; a cursor of 0 means the first page. An ID that
// appears in more than one source appears only once on the page.
//
// To tell whether there is a next page, each source must hold up to
// limit+1 IDs older than cursor, not just limit. Merge panics if limit
// is less than 1.
func Merge(sources [][]int64, cursor int64, limit int) Page {
	if limit < 1 {
		panic("timeline: limit must be at least 1")
	}

	var ids []int64
	for _, source := range sources {
		for _, id := range source {
			if cursor == 0 || id < cursor {
				ids = append(ids, id)
			}
		}
	}

	// IDs grow over time, so sorting them from largest to smallest puts the
	// newest post first. Sorting also puts duplicates next to each other,
	// which is what Compact needs to remove them.
	slices.Sort(ids)
	slices.Reverse(ids)
	ids = slices.Compact(ids)

	if len(ids) <= limit {
		return Page{PostIDs: ids}
	}
	page := ids[:limit]
	return Page{PostIDs: page, NextCursor: page[limit-1]}
}
