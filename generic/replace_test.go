package generic

import (
	"reflect"
	"testing"

	"github.com/kr/pretty"
)

type TestCase struct {
	Name     string
	Template interface{}
	Params   map[string]interface{}
	Filled   interface{}
}

var testCases = []TestCase{
	TestCase{
		Name:     "top-level number",
		Template: map[string]interface{}{"a": "${NUMBER}"},
		Params:   map[string]interface{}{"NUMBER": 10},
		Filled:   map[string]interface{}{"a": 10},
	},
	TestCase{
		Name:     "top-level hole inside string",
		Template: map[string]interface{}{"a": "as${STRING}jk"},
		Params:   map[string]interface{}{"STRING": "dfgh"},
		Filled:   map[string]interface{}{"a": "asdfghjk"},
	},
	TestCase{
		Name:     "top-level multiple holes inside string",
		Template: map[string]interface{}{"a": "as${STRING}:${INT}:${FLOAT}"},
		Params:   map[string]interface{}{"STRING": "df", "INT": 12, "FLOAT": 3.45},
		Filled:   map[string]interface{}{"a": "asdf:12:3.45"},
	},
	TestCase{
		Name:     "top-level hole filled by object",
		Template: map[string]interface{}{"a": "${OBJ}"},
		Params:   map[string]interface{}{"OBJ": map[string]interface{}{"aa": 0}},
		Filled:   map[string]interface{}{"a": map[string]interface{}{"aa": 0}},
	},
	TestCase{
		Name: "nested and multiple holes",
		Template: map[string]interface{}{
			"a": map[string]interface{}{"aa": []interface{}{
				"${AA}",
				map[string]interface{}{"aaa": "${AAA}"}}},
			"b": "${B}"},
		Params: map[string]interface{}{
			"AA":  map[string]interface{}{"zzz": 0},
			"AAA": 1,
			"B":   2},
		Filled: map[string]interface{}{
			"a": map[string]interface{}{"aa": []interface{}{
				map[string]interface{}{"zzz": 0},
				map[string]interface{}{"aaa": 1}}},
			"b": 2},
	},
}

func TestConvert(t *testing.T) {
	for _, tc := range testCases {
		result := Fill(tc.Template, tc.Params)

		if !reflect.DeepEqual(result, tc.Filled) {
			t.Error(pretty.Sprintf(
				"Unexpected result for test case %s\n%# v\n%# v\n%# v\n%# v",
				tc.Name,
				tc.Template,
				tc.Params,
				tc.Filled,
				result))
		}
	}
}
