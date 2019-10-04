package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	// ExitCodeOK exits normally
	ExitCodeOK = iota
	// ExitCodeParseFlagError exits with flag parse error
	ExitCodeParseFlagError
)

// GoJiraOption describes custom-command's option of go-jira/jira
type GoJiraOption string

// String fomrats GoJiraOption as string
// Option is <no value> if it is not set
func (g GoJiraOption) String() string {
	if g == "<no value>" {
		return ""
	}
	return string(g)
}

var (
	project    string
	issueType  string
	component  string
	state      string
	assignee   string
	resolution string
)

// CLI has streams
type CLI struct {
	Out io.Writer
	Err io.Writer
}

// Run runs Command
func (c *CLI) Run(args []string) int {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.SetOutput(c.Err)
	flags.StringVar(&issueType, "i", "", "Issue type to search for")
	flags.StringVar(&component, "c", "", "Component to search for")
	flags.StringVar(&state, "S", "", "Filter on issue status")
	flags.StringVar(&assignee, "a", "", "User assigned the issue")
	flags.StringVar(&resolution, "r", "", "Filter on issue resolution")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagError
	}

	queryMap := make(map[string]string)
	queryMap["project"] = GoJiraOption(project).String()
	queryMap["type"] = GoJiraOption(issueType).String()
	queryMap["component"] = GoJiraOption(component).String()
	queryMap["state"] = GoJiraOption(state).String()
	queryMap["assignee"] = GoJiraOption(assignee).String()
	queryMap["resolution"] = GoJiraOption(resolution).String()

	query := buildQuery(queryMap)

	fmt.Fprintf(c.Out, query)
	return ExitCodeOK
}

func buildQuery(queryMap map[string]string) string {
	query := make([]byte, 0, 2048)
	for k, v := range queryMap {
		if v != "" {
			if len(query) != 0 {
				query = append(query, []byte(" AND ")...)
			}
			query = append(query, []byte(k)...)
			query = append(query, []byte("=")...)
			query = append(query, []byte(v)...)
		}
	}
	return string(query)
}

func main() {
	c := &CLI{
		Out: os.Stdout,
		Err: os.Stderr,
	}

	project = os.Getenv("JIRA_PROJECT")
	os.Exit(c.Run(os.Args))
}
