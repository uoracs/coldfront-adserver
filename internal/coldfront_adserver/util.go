package coldfront_adserver

import (
	"context"
	"slices"
)

// DiffUserLists takes a two lists, a list of new users and a list of existing
// users, and returns two lists, a list of users to add, and a list of users
// to delete.
func DiffUserLists(ctx context.Context, newUsers, existingUsers []string) ([]string, []string) {
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
