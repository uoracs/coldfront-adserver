package coldfront_adserver

import (
	"fmt"
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

func ProcessProjectUsers(project CFProject) error {
	existingProjectUsers, err := GetCurrentProjectUsers(project.Name)
	if err != nil {
		return fmt.Errorf("failed to get current project users: %v", err)
	}
	addList, delList := DiffUserLists(project.Users, existingProjectUsers)
	for _, u := range delList {
		err := DeleteUserFromProject(project.Name, u)
		if err != nil {
			return fmt.Errorf("failed to remove user from project: %v", err)
		}
	}
	for _, u := range addList {
		err := AddUserToProject(project.Name, u)
		if err != nil {
			return fmt.Errorf("failed to add user to project: %v", err)
		}
	}
	return nil
}

func ProcessProjectAdmins(project CFProject) error {
	// manage the admins of the project
	existingAdminUsers, err := GetCurrentProjectAdminUsers(project.Name)
	if err != nil {
		return fmt.Errorf("failed to get current project admin users: %v", err)
	}
	addList, delList := DiffUserLists(project.Admins, existingAdminUsers)
	for _, u := range delList {
		err := DeleteAdminUserFromProject(project.Name, u)
		if err != nil {
			return fmt.Errorf("failed to remove admin user from project: %v", err)
		}
	}
	for _, u := range addList {
		err := AddAdminUserToProject(project.Name, u)
		if err != nil {
			return fmt.Errorf("failed to add admin user to project: %v", err)
		}
	}
	return nil
}

func ProcessProjectGroup(project CFProject, group CFGroup) error {
	existingGroupUsers, err := GetCurrentProjectGroupUsers(project.Name, group.Name)
	if err != nil {
		return fmt.Errorf("failed to get current project group users: %v", err)
	}
	addList, delList := DiffUserLists(group.Users, existingGroupUsers)
	for _, u := range delList {
		err := DeleteGroupUserFromProject(project.Name, group.Name, u)
		if err != nil {
			return fmt.Errorf("failed to remove group user from project: %v", err)
		}
	}
	for _, u := range addList {
		err := AddGroupUserToProject(project.Name, group.Name, u)
		if err != nil {
			return fmt.Errorf("failed to add group user to project: %v", err)
		}
	}
	return nil
}

func ProcessProject(project CFProject) error {
	err := ProcessProjectUsers(project)
	if err != nil {
		return fmt.Errorf("failed to process project users: %v", err)
	}
	err = ProcessProjectAdmins(project)
	if err != nil {
		return fmt.Errorf("failed to process project admins: %v", err)
	}
	for _, group := range project.Groups {
		err = ProcessProjectGroup(project, group)
		if err != nil {
			return fmt.Errorf("failed to process project group: %v", err)
		}
	}
	return nil
}
