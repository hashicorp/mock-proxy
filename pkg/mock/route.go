package mock

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
)

// The Route struct represents a single mocked route.
type Route struct {
	Host string `hcl:"host"`
	Path string `hcl:"path"`
	Type string `hcl:"type"`
}

// RouteConfig is a type alias for many Routes.
type RouteConfig []*Route

// RouteConfigHCL is used for converting HCL Blocks to RouteConfig.
type RouteConfigHCL struct {
	RouteConfig RouteConfig `hcl:"route,block"`
}

// ParseRoutes parses an input Routes file, using HCL2, into RouteConfig.
func ParseRoutes(inFile string) (RouteConfig, error) {
	input, err := os.Open(inFile)
	if err != nil {
		return []*Route{}, fmt.Errorf(
			"error in ParseRoutes opening config file: %w", err,
		)
	}
	defer input.Close()

	src, err := ioutil.ReadAll(input)
	if err != nil {
		return []*Route{}, fmt.Errorf(
			"error in ParseRoutes reading input `%s`: %w", inFile, err,
		)
	}

	parser := hclparse.NewParser()
	srcHCL, diag := parser.ParseHCL(src, inFile)
	if diag.HasErrors() {
		return []*Route{}, fmt.Errorf(
			"error in ParseRoutes parsing HCL: %w", diag,
		)
	}

	rc := &RouteConfigHCL{}
	if diag := gohcl.DecodeBody(srcHCL.Body, nil, rc); diag.HasErrors() {
		return []*Route{}, fmt.Errorf(
			"error in ParseRoutes decoding HCL configuration: %w", diag,
		)
	}

	// Return an instantiated RouteConfig instead of a nil pointer.
	if rc.RouteConfig == nil {
		rc.RouteConfig = RouteConfig{}
	}

	return rc.RouteConfig, nil
}

// ParseURL is used by a single Route to convert that route to a filepath, a
// list of transforms created by dynamic URLs, and an error. This should only
// by used on a URL that the Route Matches, as determined below.
func (r *Route) ParseURL(in *url.URL) (string, []Transformer, error) {
	switch r.Type {
	case "http":
		// An early escape for empty paths
		if r.Path == "" || r.Path == "/" {
			return fmt.Sprintf("%s/index.mock", r.Host), []Transformer{}, nil
		}

		subs, err := findSubstitutions(r.Path, in.EscapedPath())
		if err != nil {
			return "", []Transformer{},
				fmt.Errorf("error performing substitutions: %w", err)
		}

		return fmt.Sprintf("%s%s.mock", r.Host, r.Path), subs, nil
	case "git":
		// At this time, you can't template anything about git repos, because
		// of how references work.
		return filepath.Join("/git", r.Host, r.Path, ".git"), []Transformer{}, nil
	default:
		return "", []Transformer{}, fmt.Errorf("unknown route type %s", r.Type)
	}
}

// findSubstitutions is a helper function that abstracts out some pretty nasty
// regexp logic. In short, take a dynamic URL, and a templating Path from the
// Route, and convert the dynamic URL to a set of transformations with the
// :foo values turned into keys and the actual values as values.
//   Template: /mypath/:foo/bar/:baz
//   Input:    /mypath/1/bar/2
//   Output:   []VariableSubstitution{{key: foo, value: 1},{key: baz, value: 2}}
func findSubstitutions(tmplPath, inputPath string) ([]Transformer, error) {
	// First, generate a regexp with capture groups to find everywhere the
	// Route Path has a /:foo/ or /:foo value.
	pathSubRegexp := regexp.MustCompile(`(\/:\w+(?:\/|\z))+`)

	// An early exit here, if no matches, we can bail.
	pathMatches := pathSubRegexp.FindAllString(tmplPath, -1)
	if len(pathMatches) == 0 {
		return []Transformer{}, nil
	}

	// Next, use those captured segments to generate a new regexp with a
	// named capture group at each of those locations.
	captureRegexpString := fmt.Sprintf(`\A%s\z`, regexp.QuoteMeta(tmplPath))
	for _, pm := range pathMatches {
		var cg string
		if strings.HasSuffix(pm, "/") {
			cg = fmt.Sprintf(`\/(?P<%s>\S+)\/`, strings.Trim(pm, "/:"))
		} else {
			cg = fmt.Sprintf(`\/(?P<%s>\S+)`, strings.TrimLeft(pm, "/:"))
		}
		captureRegexpString = strings.Replace(captureRegexpString, regexp.QuoteMeta(pm), cg, 1)
	}
	captureRegexp, err := regexp.Compile(captureRegexpString)
	if err != nil {
		return []Transformer{}, fmt.Errorf("error generating capture group regexp: %w", err)
	}

	// Finally, generate transformers using the capture groups to create names.
	cgMatches := captureRegexp.FindStringSubmatch(inputPath)

	// If there aren't enough matches to fulfil the capture, error.
	if len(cgMatches) != len(captureRegexp.SubexpNames()) {
		return []Transformer{}, fmt.Errorf("insufficient capture groups detected")
	}

	transformers := []Transformer{}
	for i, name := range captureRegexp.SubexpNames() {
		if i != 0 && name != "" {
			val, _ := url.PathUnescape(cgMatches[i])
			transformers = append(transformers, &VariableSubstitution{
				key: name, value: val,
			})
		}
	}

	return transformers, nil
}

// match is a helper function that says if a single Route matches a single URL.
func (r *Route) match(in *url.URL) bool {
	// Easy case, if the hosts don't match, they don't match
	if r.Host != in.Host {
		return false
	}

	switch r.Type {
	case "http":
		// Another easy out, if the Paths already match, then true.
		if r.Path == in.Path || (r.Path == "" && in.Path == "/") {
			return true
		}

		// If this satisfies the input subsititution algorithm, go with it.
		subs, err := findSubstitutions(r.Path, in.EscapedPath())
		return err == nil && len(subs) != 0
	case "git":
		pathRequest := in.Path
		if len(in.RawQuery) != 0 {
			pathRequest = fmt.Sprintf("%s?%s", pathRequest, in.RawQuery)
		}
		switch pathRequest {
		case fmt.Sprintf("%s/info/refs?service=git-upload-pack", r.Path):
			return true
		case fmt.Sprintf("%s/git-upload-pack", r.Path):
			return true
		default:
			return false
		}
	default:
		// This is a bit oversimplified, but nice not to have to return an
		// error from this function.
		return false
	}
}

// MatchRoute returns the Route from a list of Routes that matches a given
// input URL.
func (rc RouteConfig) MatchRoute(in *url.URL) (*Route, error) {
	var match *Route
	var specificity int
	for _, route := range rc {
		if route.match(in) {
			// "specificity" is a measure of how many route components match
			// between the input and the matching route.
			currentSpecificity := len(strings.Split(route.Path, "/"))

			// An equally specific match is an error. Overlapping routes of
			// this type cannot be easily chosen between.
			if currentSpecificity == specificity {
				return nil, fmt.Errorf("multiple routes matched input: %s", in.String())
			}

			// A more specific match replaces the current match.
			if currentSpecificity > specificity {
				specificity = currentSpecificity
				match = route
			}

			// A less specific match is ignored.
		}
	}

	return match, nil
}
