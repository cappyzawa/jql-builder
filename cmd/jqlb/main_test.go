package main

import (
	"bytes"
	"testing"
)

func TestGojiraOptionString(t *testing.T) {
	cases := []struct {
		name   string
		opt    GoJiraOption
		expect string
	}{
		{name: "gojira default string", opt: "<no value>", expect: ""},
		{name: "with option", opt: "Task", expect: "Task"},
	}

	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.opt.String()
			if actual != c.expect {
				t.Errorf("actual should be %s, but it is %s", c.expect, actual)
			}
		})
	}
}

func TestBuildQuery(t *testing.T) {
	expect := "project=PROJECT"
	queryMap := make(map[string]string)
	queryMap["project"] = "PROJECT"
	actual := buildQuery(queryMap)
	if actual != expect {
		t.Errorf("actual should be %s, but it is %s", expect, actual)
	}
}

func TestCLIRun(t *testing.T) {
	project = "PROJECT"
	cases := []struct {
		name   string
		args   []string
		expect int
	}{
		{
			name:   "with no option",
			args:   []string{"test"},
			expect: ExitCodeOK,
		},
		{
			name:   "with issueType option",
			args:   []string{"test", "-i", "Task"},
			expect: ExitCodeOK,
		},
		{
			name:   "with issuetype and component option",
			args:   []string{"test", "-i", "Task", "-c", "Component"},
			expect: ExitCodeOK,
		},
		{
			name:   "with invalid option",
			args:   []string{"test", "--invalid"},
			expect: ExitCodeParseFlagError,
		},
	}

	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			outS := new(bytes.Buffer)
			errS := new(bytes.Buffer)
			cli := &CLI{
				Out: outS,
				Err: errS,
			}
			actual := cli.Run(c.args)
			if actual != c.expect {
				t.Errorf("exitCode should be %d, but actual is %d", c.expect, actual)
			}
		})
	}
}
