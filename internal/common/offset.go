package common

// GetPageOffset calculates the offset number based on the limit and page number.
//
// It takes two parameters:
// - limit: the maximum number of items per page.
// - pageNum: the current page number, where pages are 0-indexed.
// If pageNum is less than or equal to 1, it returns 0 as the offset.
// Otherwise, it calculates the offset by multiplying the limit with (pageNum - 1).
//
// This formula assumes that page numbers start from 0.
func GetPageOffset(limit, pageNum int) int {
	if pageNum <= 1 {
		return 0
	}
	return limit * (pageNum - 1)
}
