package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	jira "github.com/andygrunwald/go-jira"
)

func main() {

	summary := os.Args[1]

	issue, err := createIssue(summary)
	if err != nil {
		fmt.Printf("Error creating issue: %s\n", err)
		return
	}

	branchName, err := generateIssueBranchName(issue.Key, summary)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", branchName)

}

func generateIssueBranchName(issueKey string, summary string) (string, error) {
	var result strings.Builder
	var lastRune rune
	for _, r := range summary {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
			lastRune = r
		} else if (unicode.IsSpace(r) || r == '-') && lastRune != '-' {
			result.WriteRune('-')
			lastRune = '-'
		}
	}

	ticketName := issueKey + "-" + strings.Trim(result.String(), "-")

	return ticketName, nil
}

func createIssue(summary string) (*jira.Issue, error) {

	email := os.Getenv("GIT_JIRA_EMAIL_ADDRESS")
	apiToken := os.Getenv("GIT_JIRA_API_TOKEN")
	jiraURL := os.Getenv("GIT_JIRA_URL")

	tp := jira.BasicAuthTransport{
		Username: email,
		Password: apiToken,
	}

	client, err := jira.NewClient(tp.Client(), jiraURL)
	if err != nil {
		panic(err)
	}

	project, _, err := client.Project.Get("10007")
	if err != nil {
		return nil, fmt.Errorf("error getting project: %w", err)
	}

	issue, _, err := client.Issue.Get("IND-1906", nil)

	if err != nil {
		return nil, fmt.Errorf("error getting issue type: %w", err)
	}

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Priority: issue.Fields.Priority,
			Type:     issue.Fields.Type,
			Project: jira.Project{
				Key: project.Key,
			},
			Summary: summary,
		},
	}

	newIssue, _, err := client.Issue.Create(&i)
	if err != nil {
		return nil, fmt.Errorf("error creating issue: %w", err)
	}

	user, _, err := client.User.GetSelf()
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	_, err = client.Issue.UpdateAssignee(newIssue.ID, user)
	if err != nil {
		return nil, fmt.Errorf("error updating assignee: %w", err)
	}

	return newIssue, nil
}
