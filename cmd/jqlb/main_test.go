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
	cases := []struct {
		name                   string
		expect                 string
		queryMapFunc           func() map[string]string
		enableResolvedOrClosed bool
	}{
		{
			name:   "enableResolvedOrClosed is false",
			expect: "project=PROJECT AND status!=Closed AND status!=Resoved",
			queryMapFunc: func() map[string]string {
				queryMap := make(map[string]string)
				queryMap["project"] = "PROJECT"
				return queryMap
			},
			enableResolvedOrClosed: false,
		},
		{
			name:   "enableResolvedOrClosed is true",
			expect: "project=PROJECT",
			queryMapFunc: func() map[string]string {
				queryMap := make(map[string]string)
				queryMap["project"] = "PROJECT"
				return queryMap
			},
			enableResolvedOrClosed: true,
		},
		{
			name:   "enableResolvedOrClosed is false but status is specified Closed",
			expect: "project=PROJECT AND status=Closed",
			queryMapFunc: func() map[string]string {
				queryMap := make(map[string]string)
				queryMap["project"] = "PROJECT"
				queryMap["status"] = "Closed"
				return queryMap
			},
			enableResolvedOrClosed: false,
		},
	}
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			enableResolvedOrClosed = c.enableResolvedOrClosed
			queryMap := c.queryMapFunc()
			actual := buildQuery(queryMap)
			if actual != c.expect {
				t.Errorf("actual should be %s, but it is %s", c.expect, actual)
			}
		})
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
