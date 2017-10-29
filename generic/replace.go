package generic

import (
	"regexp"

	"github.com/golang/glog"
	"github.com/kr/pretty"
)

/*

"Generic" means format-agnostic.
This package operates on an AST-like object obtained by parsing
a yaml or json (or other?) document.

Template "holes" are represented as the string "${NAME}".

If a "hole" is part of (but not all of) a string, then only string/number values are supported.
This behavior is defined in `generic.fillString`.

Parameter values may also have document structure and will retain
this structure when inserted into the template.

*/

// Fill a template using the provided parameter values.
func Fill(template interface{}, params map[string]interface{}) interface{} {
	return replaceAny(template, params)
}

func replaceAny(template interface{}, params map[string]interface{}) interface{} {
	switch template := template.(type) {
	case string:
		return replaceString(template, params)
	case []interface{}:
		return replaceSlice(template, params)
	case map[string]interface{}:
		return replaceMap(template, params)
	default:
		// No template parameters in other data types.
	}

	return template
}

func replaceMap(template, params map[string]interface{}) map[string]interface{} {
	for key, val := range template {
		template[key] = replaceAny(val, params)
	}

	return template
}

func replaceSlice(template []interface{}, params map[string]interface{}) []interface{} {
	for ix, val := range template {
		template[ix] = replaceAny(val, params)
	}

	return template
}

func replaceString(template string, params map[string]interface{}) interface{} {
	// Find all template holes and replace them with param values.
	expanded, modified := expandString(template, params)
	if modified {
		return expanded
	}

	return fillString(template, params)
}

// Returns true if it expanded the template.
func expandString(template string, params map[string]interface{}) (interface{}, bool) {
	re := regexp.MustCompile("^\\$\\{([^\\{\\}]*)\\}$")
	matches := re.FindStringSubmatch(template)
	if len(matches) == 0 {
		return template, false
	}

	key := matches[1]
	if val, ok := params[key]; ok {
		return val, true
	}

	glog.Warning(pretty.Sprintf(
		"no value for (%s) in (%# v)", key, params))
	return template, false
}

func fillString(template string, params map[string]interface{}) string {
	re := regexp.MustCompile("\\$\\{[^\\{\\}]*\\}")
	result := re.ReplaceAllFunc([]byte(template), func(match []byte) []byte {
		key := match[2 : len(match)-1]
		if val, ok := params[string(key)]; ok {
			switch val := val.(type) {
			case string:
				return []byte(val)
			case float64:
				return []byte(pretty.Sprintf("%v", val))
			case int:
				return []byte(pretty.Sprintf("%v", val))
			}

			glog.Warning(pretty.Sprintf(
				"value for (%s) is not a string (%# v)",
				key, val))
		}

		glog.Warning(pretty.Sprintf(
			"no value for (%s) in (%# v)", key, params))

		return match
	})

	return string(result)
}
