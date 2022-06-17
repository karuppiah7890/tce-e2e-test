package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v44/github"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/slack-go/slack"
	"os"
	"strings"
	"time"
)

var ctx = context.Background()

type GitHub struct {
	client *github.Client
}

const (
	lookup    = "E2E"
	channelID = "SLACK_CHANNEL_ID"
	token     = "SLACK_TOKEN"
	owner     = "vmware-tanzu"
	repo      = "community-edition"
)

var timeNow = time.Now()

func main() {
	log.InitLogger("slacker")
	client := github.NewClient(nil)

	mygh := &GitHub{
		client: client,
	}
	result, _ := mygh.listWorkflows()
	log.Infof("%d", result)
	for _, runs := range result {
		for _, y := range runs.WorkflowRuns {
			if "failure" == *y.Conclusion {
				runtime := y.UpdatedAt.Time.Sub(y.CreatedAt.Time).Minutes()
				log.Infof("%s %s %s %s %f", *y.Conclusion, *y.Name, *y.CreatedAt, *y.HTMLURL, runtime)
				sendSlack(fmt.Sprintf("%s", *y.Name), fmt.Sprintf("Commit Message :%s \n Commit from : %s", *y.HeadCommit.Message, *y.HeadCommit.Author.Name), fmt.Sprintf("%s", *y.Conclusion), fmt.Sprintf("%s", *y.HTMLURL), fmt.Sprintf("%f", runtime), fmt.Sprintf("%s", *y.CreatedAt))

			}
		}
	}

}

//func (g *GitHub) listWorkflowFromId(workflowId int64) (*github.WorkflowRuns, error) {
//	opts := &github.ListWorkflowRunsOptions{}
//
//	runs, _, err := g.client.Actions.ListWorkflowRunsByID(ctx, "vmware-tanzu", "community-edition", workflowId, opts)
//	if err != nil {
//		log.Errorf("Workflows Listing failed. Err: %v\n", err)
//		return nil, err
//	}
//
//	return runs, nil
//}
//func (g *GitHub) getRunStatus() error {
//	succeeded := false
//	for i := 0; i < defaultNumOfTimesToPoll; i++ {
//		runner, err := g.listWorkflowFromId(42524244)
//		if err != nil {
//			log.Errorf("Workflows Listing failed. Err: %v\n", err)
//			return err
//		}
//		if *runner.WorkflowRuns[0].Status != "failure" {
//			succeeded = true
//			break
//		}
//		log.Infof("Attempt poll %d... sleeping\n", i)
//		time.Sleep(time.Duration(defaultSleepBetweenPoll) * time.Second)
//	}
//	return nil
//}

func sendSlack(heading, details, status, url, runtime, created string) {
	token := os.Getenv(token)
	channelID := os.Getenv(channelID)
	client := slack.New(token, slack.OptionDebug(false))
	attachment := slack.Attachment{
		Pretext: heading,
		Text:    details,
		Color:   "#ffff00",
		// Fields are Optional extra data!
		Fields: []slack.AttachmentField{
			{
				Title: "Status",
				Value: status,
			},
			{
				Title: "Github URL",
				Value: url,
			},
			{
				Title: "Runtime",
				Value: fmt.Sprint(runtime, " Minutes"),
			}, {
				Title: "Created On",
				Value: created,
			},
		},
	}
	_, timestamp, err := client.PostMessage(
		channelID,
		// uncomment the item below to add a extra Header to the message, try it out :)
		//slack.MsgOptionText("New message from bot", false),
		slack.MsgOptionAttachments(attachment),
	)
	if err != nil {
		panic(err)
	}
	log.Infof("Message sent at %s", timestamp)
}

func (g *GitHub) listWorkflows() ([]*github.WorkflowRuns, error) {
	days := timeNow.AddDate(0, 0, -1).Format("2006-01-02")
	opts := &github.ListWorkflowRunsOptions{Created: fmt.Sprint(">", days)}
	//opts := &github.ListWorkflowRunsOptions{Status: "failure", Created: ">2022-05-19"}
	opt := &github.ListOptions{}
	workflows, _, err := g.client.Actions.ListWorkflows(ctx, owner, repo, opt)
	if err != nil {
		log.Errorf("Workflows Listing failed. Err: %v\n", err)
		return nil, err
	}
	TceWorkflow := []*github.WorkflowRuns{}
	for _, workflow := range workflows.Workflows {
		if strings.Contains(*workflow.Name, lookup) {
			log.Infof("%s %d", *workflow.Name, *workflow.ID)
			runs, _, err := g.client.Actions.ListWorkflowRunsByID(ctx, owner, repo, *workflow.ID, opts)
			if err != nil {
				log.Errorf("Workflows Listing failed. Err: %v\n", err)
				return nil, err
			}
			TceWorkflow = append(TceWorkflow, runs)
			//return runs, nil
			log.Infof("%d", *runs.TotalCount)
		}
	}
	return TceWorkflow, nil
}
