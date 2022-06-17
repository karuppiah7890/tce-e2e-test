package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v44/github"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"os"
	"strconv"
	"strings"
	"time"
)

var ctx = context.Background()

type GitHub struct {
	client *github.Client
}
type result struct {
	date      string
	status    string
	buildType string
	jobs      job
}
type job struct {
	name       string
	status     string
	time       string
	runtime    string
	url        string
	trigrredBy string
}

const (
	lookup    = "E2E"
	owner     = "vmware-tanzu"
	repo      = "community-edition"
	workflows = "WORKFLOWS" //To be set as csv values of workflows id
)

var timeNow = time.Now()

func main() {
	log.InitLogger("reporting")
	client := github.NewClient(nil)
	//workflows := os.Getenv(workflows)
	mygh := &GitHub{
		client: client,
	}
	workflows := strings.Split(os.Getenv(workflows), ",")
	log.Infof("%s", workflows)
	for _, workflow := range workflows {
		log.Infof("checking for workflow id %s", workflow)
		w_id, _ := strconv.ParseInt(workflow, 10, 64)
		run, err := mygh.listWorkflowFromId(w_id)
		if err != nil {
			log.Errorf("Error %s ", err)
		}
		y := run.WorkflowRuns[1]
		log.Infof("%s %s %s %s", *y.Conclusion, *y.Name, *y.CreatedAt, *y.HTMLURL)
	}
	//result, _ := mygh.listWorkflows()
	//log.Infof("%d", result)
	//for _, runs := range result {
	//	for _, y := range runs.WorkflowRuns {
	//		runtime := y.UpdatedAt.Time.Sub(y.CreatedAt.Time).Minutes()
	//		log.Infof("%s %s %s %s %f", *y.Conclusion, *y.Name, *y.CreatedAt, *y.HTMLURL, runtime)
	//
	//	}
	//}

}

func (g *GitHub) listWorkflowFromId(workflowId int64) (*github.WorkflowRuns, error) {
	opts := &github.ListWorkflowRunsOptions{
		Actor:       "",
		Branch:      "",
		Event:       "",
		Status:      "",
		Created:     "",
		ListOptions: github.ListOptions{},
	}
	runs, _, err := g.client.Actions.ListWorkflowRunsByID(ctx, "vmware-tanzu", "community-edition", workflowId, opts)
	if err != nil {
		log.Errorf("Workflows Listing failed. Err: %v\n", err)
		return nil, err
	}

	return runs, nil
}

func (g *GitHub) listWorkflows() ([]*github.WorkflowRuns, error) {
	days := timeNow.AddDate(0, 0, -1).Format("2006-01-02")
	opts := &github.ListWorkflowRunsOptions{Status: "failure", Created: fmt.Sprint(">", days)}
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
