package yaml

import (
	"testing"
)

type TestCase struct {
	Name     string
	Template string
	Params   string
	Filled   string
}

var testCases = []TestCase{
	TestCase{
		Name:     "top-level number",
		Template: `a: ${NUMBER}`,
		Params:   `NUMBER: 10`,
		Filled: `a: 10
`,
	},
	TestCase{
		Name:     "top-level hole inside string",
		Template: `a: as${STRING}jk`,
		Params:   `STRING: dfgh`,
		Filled: `a: asdfghjk
`,
	},
	TestCase{
		Name:     "top-level hole filled by object",
		Template: `a: ${OBJ}`,
		Params: `
OBJ:
  aa: 0`,
		Filled: `a:
  aa: 0
`,
	},
	TestCase{
		Name: "nested and multiple holes",
		Template: `a:
  aa:
  - ${AA}
  - aaa: ${AAA}
b: ${B}`,
		Params: `AA:
  zzz: 0
AAA: 1
B: 2`,
		Filled: `a:
  aa:
  - zzz: 0
  - aaa: 1
b: 2
`,
	},
}

func TestConvert(t *testing.T) {
	for _, tc := range testCases {
		result, err := Fill(
			[]byte(tc.Template),
			[]byte(tc.Params))
		if err != nil {
			t.Error(err)
		}

		if string(result) != tc.Filled {
			t.Errorf("Unexpected result for test case %s\n%s\n%s\n%s\n%s\n%v\n%v",
				tc.Name,
				tc.Template,
				tc.Params,
				tc.Filled,
				result,
				len(tc.Filled),
				len(result))
		}
	}
}
