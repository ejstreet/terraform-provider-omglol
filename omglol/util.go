package omglol

import (
	"regexp"
)

var errorNotFoundRegexp = regexp.MustCompile("status: 404|NOT_FOUND|does not exist|No records match the filter criteria")

func isNotFoundError(err error) bool {
	return errorNotFoundRegexp.MatchString(err.Error())
}