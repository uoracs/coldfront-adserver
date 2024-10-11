package coldfront_adserver

import (
	"slices"
	"strings"
)

// DiffLists takes a two lists, a list of new strings and a list of existing
// strings, and returns two lists, a list of strings to add, and a list of strings
// to delete.
func DiffLists(newUsers, existingUsers []string) ([]string, []string) {
	var addUsers []string
	var delUsers []string
	// only record the users that don't already exist
	for _, nu := range newUsers {
		if !slices.Contains(existingUsers, nu) {
			addUsers = append(addUsers, nu)
		}
	}
	// any existing user that's not present in newusers,
	// create a delete operation
	for _, eu := range existingUsers {
		if !slices.Contains(newUsers, eu) {
			delUsers = append(delUsers, eu)
		}
	}
	return addUsers, delUsers
}

// CleanList takes a block of text and returns all non-empty lines
// with all the leading/trailing spaces cleaned
func CleanList(text string) []string {
	var out []string
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l != "" {
			out = append(out, l)
		}
	}
	return out
}
