package cmd

import (
	"fmt"
	jk "jk/src"
	"strconv"
	"strings"
)

func argsToStr(args []string) string {
	var str string
	for _, arg := range args {
		str += arg
	}
	return str
}

func argsToIds(args []string) (ids []int, err error) {
	ids = make([]int, 0)
	var id int
	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		if arg == "" {
			continue
		}
		id, err = strconv.Atoi(arg)
		if err != nil {
			return
		}
		ids = append(ids, id)
	}
	return
}

func argsToIdAndStr(args []string) (id int, str string, err error) {
	if len(args) != 2 {
		err = &jk.InvalidInputError{
			Message: fmt.Sprintf("expect: <id> <str>, got: %v", args),
		}
		return
	}
	id, err = strconv.Atoi(args[0])
	if err != nil {
		return
	}
	str = strings.TrimSpace(args[1])
	return
}

func argsToIdAndOptionalStr(args []string) (id int, str string, err error) {
	if len(args) != 1 && len(args) != 2 {
		err = &jk.InvalidInputError{
			Message: fmt.Sprintf("expect: <id> [str], got: %v", args),
		}
		return
	}

	id, err = strconv.Atoi(args[0])
	if err != nil {
		return
	}
	if len(args) == 1 {
		return
	}

	str = strings.TrimSpace(args[1])
	return
}
