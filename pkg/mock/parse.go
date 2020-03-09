package mock

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	pathRegexp = regexp.MustCompile(`~~(?P<key>\w+)~(?P<value>\w+)~~`)
)

// parsePath takes a path and matches it to an API endpoint. This works by
// iterating through the path and searching for path components that match
// the pattern:
//   ~~KEY~VALUE~~
// These will be replaced with :KEY, and a new substitution variable will be
// added that replaces KEY with VALUE. For instance:
//   /users/~~user~rae~~/repos
// is replaced by:
//   /users/:user/repos
// and a new substitution for:
//   user => rae
// will be created for that request.
func (ms *MockServer) parsePath(path string) (string, []*VariableSubstitution) {
	subs := []*VariableSubstitution{}
	matches := pathRegexp.FindAllStringSubmatch(path, -1)
	for _, match := range matches {
		result := make(map[string]string)
		for i, name := range pathRegexp.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}
		subs = append(subs, &VariableSubstitution{
			key:   result["key"],
			value: result["value"],
		})

		path = strings.Replace(path, match[0], fmt.Sprintf(":%s", result["key"]), 1)
	}

	return path, subs
}
