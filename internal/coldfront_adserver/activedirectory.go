package coldfront_adserver

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func GetCurrentProjectOwner(ctx context.Context, projectName string) (string, error) {
	slog.Debug("getting current project owner", "project", projectName)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Get-PirgPIUsername -Pirg %s", projectName)
	output, err := ex.Execute(command)
	if err != nil {
		return "", fmt.Errorf("failed to get project users: %v", err)
	}
	fmt.Println(output)
	names := strings.Split(output, "\n")
	if len(names) != 1 {
		return "", fmt.Errorf("more than one existing owner for project '%s', fix manually!", projectName)
	}
	return names[0], nil
}

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

func GetCurrentProjectGroupNames(ctx context.Context, projectName string) ([]string, error) {
	slog.Debug("getting current project groups", "project", projectName)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Get-PirgGroupNames -Pirg %s", projectName)
	output, err := ex.Execute(command)
	if err != nil {
		return nil, fmt.Errorf("failed to get project groups: %v", err)
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
	slog.Info("adding user to project", "project", projectName, "user", username)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Add-PirgUser -Pirg %s -User %s", projectName, username)
	_, err := ex.Execute(command)
	if err != nil {
		return fmt.Errorf("failed to add user '%s' to project '%s': %v", username, projectName, err)
	}
	return nil
}

func DeleteUserFromProject(ctx context.Context, projectName, username string) error {
	slog.Info("removing user from project", "project", projectName, "user", username)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Remove-PirgUser -Pirg %s -User %s", projectName, username)
	_, err := ex.Execute(command)
	if err != nil {
		return fmt.Errorf("failed to remove user '%s' from project '%s': %v", username, projectName, err)
	}
	return nil
}

func AddAdminUserToProject(ctx context.Context, projectName, username string) error {
	slog.Info("adding admin user to project", "project", projectName, "user", username)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Add-PirgAdmin -Pirg %s -User %s", projectName, username)
	_, err := ex.Execute(command)
	if err != nil {
		return fmt.Errorf("failed to add admin '%s' to project '%s': %v", username, projectName, err)
	}
	return nil
}

func DeleteAdminUserFromProject(ctx context.Context, projectName, username string) error {
	slog.Info("removing admin user from project", "project", projectName, "user", username)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Remove-PirgAdmin -Pirg %s -User %s", projectName, username)
	_, err := ex.Execute(command)
	if err != nil {
		return fmt.Errorf("failed to remove admin '%s' from project '%s': %v", username, projectName, err)
	}
	return nil
}

func AddGroupToProject(ctx context.Context, projectName, groupName string) error {
	slog.Info("adding group to project", "project", projectName, "group", groupName)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("New-PirgGroup -Pirg %s -Name %s", projectName, groupName)
	_, err := ex.Execute(command)
	if err != nil {
		return fmt.Errorf("failed to add group '%s' to project '%s': %v", groupName, projectName, err)
	}
	return nil
}

func DeleteGroupFromProject(ctx context.Context, projectName, groupName string) error {
	slog.Info("removing user from project group", "project", projectName, "group", groupName)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Remove-PirgGroup -Pirg %s -Name %s", projectName, groupName)
	_, err := ex.Execute(command)
	if err != nil {
		return fmt.Errorf("failed to remove group '%s' from project '%s': %v", groupName, projectName, err)
	}
	return nil
}

func AddGroupUserToProject(ctx context.Context, projectName, groupName, username string) error {
	slog.Info("adding user to project group", "project", projectName, "group", groupName, "user", username)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Add-PirgGroupUser -Pirg %s -Group %s -User %s", projectName, groupName, username)
	_, err := ex.Execute(command)
	if err != nil {
		return fmt.Errorf("failed to add group '%s' user '%s' to project '%s': %v", groupName, username, projectName, err)
	}
	return nil
}

func DeleteGroupUserFromProject(ctx context.Context, projectName, groupName, username string) error {
	slog.Info("removing user from project group", "project", projectName, "group", groupName, "user", username)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Remove-PirgGroupUser -Pirg %s -Group %s -User %s", projectName, groupName, username)
	_, err := ex.Execute(command)
	if err != nil {
		return fmt.Errorf("failed to remove group '%s' user '%s' from project '%s': %v", groupName, username, projectName, err)
	}
	return nil
}

func SetProjectOwner(ctx context.Context, projectName, username string) error {
	slog.Info("setting project owner", "project", projectName, "owner", username)
	ex := ctx.Value(ExecutorKey).(Executor)
	command := fmt.Sprintf("Set-PirgPI -Pirg %s -User %s", projectName, username)
	_, err := ex.Execute(command)
	if err != nil {
		return fmt.Errorf("failed to set owner '%s' on project '%s': %v", username, projectName, err)
	}
	return nil
}
