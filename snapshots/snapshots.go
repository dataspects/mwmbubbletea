package snapshots

import (
	"fmt"
	"strings"
)

func ManageSnapshotsInterface() string {
	ifce := []string{
		"Manage snapshots",
		fmt.Sprintf("%s or %s", "dasd", "asd"),
	}
	return strings.Join(ifce, "\n")
}
