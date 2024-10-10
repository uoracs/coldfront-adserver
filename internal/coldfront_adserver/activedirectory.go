package coldfront_adserver

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func GetCurrentProjectUsers(ctx context.Context, projectName string) ([]string, error) {
	slog.Debug("getting current project users", "project", projectName)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Get-PirgUsernames -Pirg %s", projectName)
	output, err := ex.Execute(command)
	if err != nil {
		return nil, fmt.Errorf("failed to get project users: %v", err)
	}
	return strings.Split(output, "\n"), nil
}

func GetCurrentProjectAdminUsers(ctx context.Context, projectName string) ([]string, error) {
	slog.Debug("getting current project admins", "project", projectName)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Get-PirgAdminUsernames -Pirg %s", projectName)
	output, err := ex.Execute(command)
	if err != nil {
		return nil, fmt.Errorf("failed to get project admin users: %v", err)
	}
	return strings.Split(output, "\n"), nil
}

func GetCurrentProjectGroupUsers(ctx context.Context, projectName, groupName string) ([]string, error) {
	slog.Debug("getting current project group users", "project", projectName, "group", groupName)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Get-PirgGroupUsernames -Pirg %s -Group %s", projectName, groupName)
	output, err := ex.Execute(command)
	if err != nil {
		return nil, fmt.Errorf("failed to get project group users: %v", err)
	}
	return strings.Split(output, "\n"), nil
}

func AddUserToProject(ctx context.Context, projectName, username string) error {
	slog.Debug("adding user to project", "project", projectName, "user", username)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Add-PirgUser -Pirg %s -User %s", projectName, username)
	_, err := ex.Execute(command)
	if err != nil {
		return fmt.Errorf("failed to add user '%s' to project '%s': %v", username, projectName, err)
	}
	return nil
}

func DeleteUserFromProject(ctx context.Context, projectName, username string) error {
	slog.Debug("removing user from project", "project", projectName, "user", username)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Remove-PirgUser -Pirg %s -User %s", projectName, username)
	_, err := ex.Execute(command)
	if err != nil {
		return fmt.Errorf("failed to remove user '%s' from project '%s': %v", username, projectName, err)
	}
	return nil
}

func AddAdminUserToProject(ctx context.Context, projectName, username string) error {
	slog.Debug("adding admin user to project", "project", projectName, "user", username)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Add-PirgAdmin -Pirg %s -User %s", projectName, username)
	_, err := ex.Execute(command)
	if err != nil {
		return fmt.Errorf("failed to add admin '%s' to project '%s': %v", username, projectName, err)
	}
	return nil
}

func DeleteAdminUserFromProject(ctx context.Context, projectName, username string) error {
	slog.Debug("removing admin user from project", "project", projectName, "user", username)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Remove-PirgAdmin -Pirg %s -User %s", projectName, username)
	_, err := ex.Execute(command)
	if err != nil {
		return fmt.Errorf("failed to remove admin '%s' from project '%s': %v", username, projectName, err)
	}
	return nil
}

func AddGroupUserToProject(ctx context.Context, projectName, groupName, username string) error {
	slog.Debug("adding user to project group", "project", projectName, "group", groupName, "user", username)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Add-PirgGroupUser -Pirg %s -Group %s -User %s", projectName, groupName, username)
	_, err := ex.Execute(command)
	if err != nil {
		return fmt.Errorf("failed to add group '%s' user '%s' to project '%s': %v", groupName, username, projectName, err)
	}
	return nil
}

func DeleteGroupUserFromProject(ctx context.Context, projectName, groupName, username string) error {
	slog.Debug("removing user from project group", "project", projectName, "group", groupName, "user", username)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Remove-PirgGroupUser -Pirg %s -Group %s -User %s", projectName, groupName, username)
	_, err := ex.Execute(command)
	if err != nil {
		return fmt.Errorf("failed to remote group '%s' user '%s' from project '%s': %v", groupName, username, projectName, err)
	}
	return nil
}
