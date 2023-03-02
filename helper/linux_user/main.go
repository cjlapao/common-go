package linux_user

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cjlapao/common-go/commands"
	"github.com/cjlapao/common-go/guard"
)

type LinuxUser struct {
	ID            int
	Name          string
	GroupId       int
	HomeDirectory string
	Shell         string
}

type LinuxUserCreateOptions struct {
	Key   string
	Value string
}

func Marshal(value string) (*LinuxUser, error) {
	if value == "" {
		return nil, errors.New("value is empty")
	}

	parts := strings.Split(value, ":")
	if len(parts) != 7 {
		return nil, errors.New("wrong format")
	}

	userId, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, errors.New("the user id is invalid")
	}

	groupId, err := strconv.Atoi(parts[3])
	if err != nil {
		return nil, errors.New("the group id is invalid")
	}

	result := LinuxUser{
		ID:            userId,
		Name:          parts[0],
		GroupId:       groupId,
		HomeDirectory: parts[5],
		Shell:         parts[6],
	}

	return &result, nil
}

func Exists(userName string) bool {
	_, err := Get(userName)

	return err == nil
}

func Get(userName string) (*LinuxUser, error) {
	output, err := commands.Execute("getent", "passwd", userName)
	if err != nil {
		return nil, errors.New("user does not exist")
	}

	if output.GetAllOutputs() == "" {
		return nil, errors.New("user does not exist")
	}

	user, err := Marshal(output.GetAllOutputs())

	if err != nil {
		return nil, err
	}

	return user, nil
}

func UserGroupCreateOption(groupName string) LinuxUserCreateOptions {
	return LinuxUserCreateOptions{
		Key:   "-g",
		Value: groupName,
	}
}

func UserHomeDirectoryCreateOption(path string) LinuxUserCreateOptions {
	return LinuxUserCreateOptions{
		Key:   "-d",
		Value: path,
	}
}

func Create(userName string, userId int, options ...LinuxUserCreateOptions) error {
	if userName == "" {
		return errors.New("group name cannot be empty")
	}

	if userId <= 1 {
		return errors.New("group Id needs to be greater than 0")
	}

	exists := Exists(userName)
	if exists {
		return fmt.Errorf("user %v already exists", userName)
	}

	createParameters := make([]string, 0)
	createParameters = append(createParameters, "-u")
	createParameters = append(createParameters, fmt.Sprintf("%v", userId))
	createParameters = append(createParameters, userName)

	if len(options) > 0 {
		for _, option := range options {
			createParameters = append(createParameters, option.Key)
			createParameters = append(createParameters, option.Value)
		}
	}

	output, err := commands.Execute("useradd", createParameters...)

	if err != nil {
		return fmt.Errorf("there was an error creating user %v with id %v, err %v", userName, userId, err.Error())
	}

	if strings.ContainsAny(output.GetAllOutputs(), "already exists") {
		return fmt.Errorf("there was an error creating user %v with id %v, user already exists", userName, userId)
	}

	return nil
}

func AddToGroup(userName string, groupName string) error {
	if err := guard.EmptyOrNil(userName, "username"); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(groupName, "group name"); err != nil {
		return err
	}

	_, err := commands.Execute("adduser", userName, groupName)

	if err != nil {
		return err
	}

	return nil
}
