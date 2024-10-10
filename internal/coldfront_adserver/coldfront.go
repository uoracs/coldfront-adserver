package coldfront_adserver

import (
	"context"
	"fmt"
	"log/slog"
)

type CFGroup struct {
	Name  string   `json:"name,omitempty"`
	Users []string `json:"users,omitempty"`
}

type CFProject struct {
	Name   string    `json:"name,omitempty"`
	Owner  string    `json:"owner,omitempty"`
	Admins []string  `json:"admins,omitempty"`
	Users  []string  `json:"users,omitempty"`
	Groups []CFGroup `json:"groups,omitempty"`
}

type CFProjectsRequest struct {
	Projects []CFProject
}

func ProcessProjectUsers(ctx context.Context, project CFProject) error {
	existingProjectUsers, err := GetCurrentProjectUsers(ctx, project.Name)
	if err != nil {
		return fmt.Errorf("failed to get current project users: %v", err)
	}
	addList, delList := DiffLists(project.Users, existingProjectUsers)
	for _, u := range delList {
		err := DeleteUserFromProject(ctx, project.Name, u)
		if err != nil {
			return fmt.Errorf("failed to remove user from project: %v", err)
		}
	}
	for _, u := range addList {
		err := AddUserToProject(ctx, project.Name, u)
		if err != nil {
			return fmt.Errorf("failed to add user to project: %v", err)
		}
	}
	return nil
}

func ProcessProjectAdmins(ctx context.Context, project CFProject) error {
	// manage the admins of the project
	existingAdminUsers, err := GetCurrentProjectAdminUsers(ctx, project.Name)
	if err != nil {
		return fmt.Errorf("failed to get current project admin users: %v", err)
	}
	addList, delList := DiffLists(project.Admins, existingAdminUsers)
	for _, u := range delList {
		err := DeleteAdminUserFromProject(ctx, project.Name, u)
		if err != nil {
			return fmt.Errorf("failed to remove admin user from project: %v", err)
		}
	}
	for _, u := range addList {
		err := AddAdminUserToProject(ctx, project.Name, u)
		if err != nil {
			return fmt.Errorf("failed to add admin user to project: %v", err)
		}
	}
	return nil
}

func ProcessProjectGroup(ctx context.Context, project CFProject, group CFGroup) error {
	existingGroupUsers, err := GetCurrentProjectGroupUsers(ctx, project.Name, group.Name)
	if err != nil {
		return fmt.Errorf("failed to get current project group users: %v", err)
	}
	addList, delList := DiffLists(group.Users, existingGroupUsers)
	for _, u := range delList {
		err := DeleteGroupUserFromProject(ctx, project.Name, group.Name, u)
		if err != nil {
			return fmt.Errorf("failed to remove group user from project: %v", err)
		}
	}
	for _, u := range addList {
		err := AddGroupUserToProject(ctx, project.Name, group.Name, u)
		if err != nil {
			return fmt.Errorf("failed to add group user to project: %v", err)
		}
	}
	return nil
}

func ProcessProjectGroups(ctx context.Context, project CFProject) error {
	existingGroupNames, err := GetCurrentProjectGroupNames(ctx, project.Name)
	if err != nil {
		return fmt.Errorf("failed to get current project groups: %v", err)
	}
	newGroupNames := GetIncomingGroupNames(project.Groups)
	addList, delList := DiffLists(newGroupNames, existingGroupNames)
	for _, g := range delList {
		err := DeleteGroupFromProject(ctx, project.Name, g)
		if err != nil {
			return fmt.Errorf("failed to remove group from project: %v", err)
		}
	}
	for _, g := range addList {
		err := AddGroupToProject(ctx, project.Name, g)
		if err != nil {
			return fmt.Errorf("failed to add group to project: %v", err)
		}
	}
	return nil
}

func ProcessProject(ctx context.Context, project CFProject) error {
	slog.Debug("processing project", "project", project.Name)
	err := ProcessProjectUsers(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to process project users: %v", err)
	}
	err = ProcessProjectAdmins(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to process project admins: %v", err)
	}
	// this handles things like making sure groups exist or get deleted
	err = ProcessProjectGroups(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to process project groups: %v", err)
	}
	for _, group := range project.Groups {
		// this one handles group membership of groups
		err = ProcessProjectGroup(ctx, project, group)
		if err != nil {
			return fmt.Errorf("failed to process project group: %v", err)
		}
	}
	return nil
}

func GetIncomingGroupNames(groups []CFGroup) []string {
	var names []string
	for _, g := range groups {
		names = append(names, g.Name)
	}
	return names
}
