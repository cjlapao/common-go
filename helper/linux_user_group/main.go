package linux_user_group

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cjlapao/common-go/commands"
)

type LinuxUserGroup struct {
	ID    int
	Name  string
	Users []string
}

func Marshal(value string) (*LinuxUserGroup, error) {
	if value == "" {
		return nil, errors.New("value is empty")
	}

	parts := strings.Split(value, ":")
	if len(parts) != 4 {
		return nil, errors.New("wrong format")
	}

	groupId, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, errors.New("the group id is invalid")
	}

	result := LinuxUserGroup{
		ID:   groupId,
		Name: parts[0],
	}

	if len(parts[3]) == 0 {
		result.Users = make([]string, 0)
	} else {
		users := strings.Split(parts[3], ",")
		result.Users = append(result.Users, users...)
	}

	return &result, nil
}

func Exists(groupName string) bool {
	_, err := Get(groupName)

	return err == nil
}

func Get(groupName string) (*LinuxUserGroup, error) {
	output, err := commands.Execute("getent", "group", groupName)
	if err != nil {
		return nil, errors.New("group does not exists")
	}

	if output == "" {
		return nil, errors.New("group does not exists")
	}

	group, err := Marshal(output)

	if err != nil {
		return nil, err
	}

	return group, nil
}

func Create(groupName string, groupId int) error {
	if groupName == "" {
		return errors.New("group name cannot be empty")
	}

	if groupId <= 1 {
		return errors.New("group Id needs to be greater than 0")
	}

	exists := Exists(groupName)
	if exists {
		return fmt.Errorf("group %v already exists", groupName)
	}

	output, err := commands.Execute("groupadd", "-g", fmt.Sprintf("%v", groupId), groupName)

	if err != nil {
		return fmt.Errorf("there was an error creating group %v with id %v, err %v", groupName, groupId, err.Error())
	}

	if strings.ContainsAny(output, "already exists") {
		return fmt.Errorf("there was an error creating group %v with id %v, group already exists", groupName, groupId)
	}

	return nil
}
