package linux_helper

import (
	"fmt"

	"github.com/cjlapao/common-go/commands"
	"github.com/cjlapao/common-go/guard"
	"github.com/cjlapao/common-go/helper"
)

func ChangeOwner(path string, username string, groupName string, recursive bool) error {
	if err := guard.EmptyOrNil(path); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(username); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(groupName); err != nil {
		return err
	}

	createParameters := make([]string, 0)
	if recursive {
		createParameters = append(createParameters, "-R")
	}

	createParameters = append(createParameters, fmt.Sprintf("%v:%v", username, groupName))
	createParameters = append(createParameters, helper.ToOsPath(path))

	_, err := commands.Execute("chown", createParameters...)

	if err != nil {
		return err
	}

	return nil
}
