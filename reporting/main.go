package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v44/github"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"golang.org/x/oauth2"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var ctx = context.Background()

type GitHub struct {
	client *github.Client
}
type Result struct {
	Date      string      `json:"Date"`
	Status    string      `json:"Status"`
	BuildType string      `json:"BuildType"`
	Providers []Providers `json:"Providers"`
	Plugins   []Plugins   `json:"Plugins"`
	Packages  []Packages  `json:"Packages"`
}
type Providers struct {
	Name    string `json:"Name"`
	Status  string `json:"Status"`
	Time    string `json:"Created Time"`
	Runtime string `json:"Runtime"`
	Url     string `json:"URL"`
}
type Plugins struct {
	Name    string `json:"Name"`
	Status  string `json:"Status"`
	Time    string `json:"Created Time"`
	Runtime string `json:"Runtime"`
	Url     string `json:"URL"`
}
type Packages struct {
	Name    string `json:"Name"`
	Status  string `json:"Status"`
	Time    string `json:"Created Time"`
	Runtime string `json:"Runtime"`
	Url     string `json:"URL"`
}

const (
	lookup       = "E2E"
	owner        = "vmware-tanzu"
	repo         = "community-edition"
	workflowsEnv = "WORKFLOWS" //To be set as csv values of workflows id
	pluginEnv    = "PLUGIN"    //To be set as csv values of workflows id
	token        = "GH_TOKEN"  //To be set as csv values of workflows id
)

var timeNow = time.Now()
var ts = oauth2.StaticTokenSource(
	&oauth2.Token{AccessToken: os.Getenv(token)},
)
var tc = oauth2.NewClient(ctx, ts)

var client = github.NewClient(tc)
var mygh = &GitHub{
	client: client,
}

func main() {
	log.InitLogger("reporting")
	handleRequests()
}

func (ghClient *GitHub) getResults() ([]byte, error) {
	res := &Result{
		Date:      fmt.Sprintf("%v", timeNow.Format(time.RFC822)),
		Status:    "",
		BuildType: "daily",
		Providers: []Providers{},
	}

	res.Status = "success"
	res.Providers = mygh.getProviderResults()
	res.Plugins = mygh.getPluginResults()
	ParseJson, err := json.Marshal(res)
	if err != nil {
		log.Errorf("Error %s ", err)
	}
	log.Infof("%s", string(ParseJson))
	return json.MarshalIndent(res, "", "  ")
}

func (ghClient *GitHub) getResultsStruct() (Result, error) {
	res := &Result{
		Date:      fmt.Sprintf("%v", timeNow.Format(time.RFC822)),
		Status:    "",
		BuildType: "daily",
		Providers: []Providers{},
	}

	res.Status = "success"
	res.Providers = mygh.getProviderResults()
	res.Plugins = mygh.getPluginResults()

	return *res, nil
}

//Todo Remove Redundant code for provide and plugin
func (ghClient *GitHub) getProviderResults() []Providers {
	workflows := strings.Split(os.Getenv(workflowsEnv), ",")
	log.Infof("%s", workflows)
	w := []Providers{}
	for _, workflow := range workflows {
		log.Infof("checking for workflow id %s", workflow)
		w_id, _ := strconv.ParseInt(workflow, 10, 64)
		run, err := ghClient.listWorkflowFromId(w_id)
		if err != nil {
			log.Errorf("Error %s ", err)
		}
		y := run.WorkflowRuns[1]
		log.Infof("%s %s %s %s", *y.Conclusion, *y.Name, *y.CreatedAt, *y.HTMLURL)
		runtime := y.UpdatedAt.Time.Sub(y.CreatedAt.Time).Minutes()
		x := Providers{
			Name:    *y.Name,
			Status:  *y.Conclusion,
			Time:    fmt.Sprintf("%s", *y.CreatedAt),
			Runtime: fmt.Sprintf("%.1f Minutes", runtime),
			Url:     *y.HTMLURL,
		}
		w = append(w, x)
	}
	return w

}

//Todo Remove Redundant code for provide and plugin
func (ghClient *GitHub) getPluginResults() []Plugins {
	workflows := strings.Split(os.Getenv(pluginEnv), ",")
	log.Infof("%s", workflows)
	w := []Plugins{}
	for _, workflow := range workflows {
		log.Infof("checking for workflow id %s", workflow)
		w_id, _ := strconv.ParseInt(workflow, 10, 64)
		run, err := ghClient.listWorkflowFromId(w_id)
		if err != nil {
			log.Errorf("Error %s ", err)
		}
		y := run.WorkflowRuns[1]
		log.Infof("%s %s %s %s", *y.Conclusion, *y.Name, *y.CreatedAt, *y.HTMLURL)
		runtime := y.UpdatedAt.Time.Sub(y.CreatedAt.Time).Minutes()
		x := Plugins{
			Name:    *y.Name,
			Status:  *y.Conclusion,
			Time:    fmt.Sprintf("%s", *y.CreatedAt),
			Runtime: fmt.Sprintf("%.1f Minutes", runtime),
			Url:     *y.HTMLURL,
		}
		w = append(w, x)
	}
	return w
}

func (ghClient *GitHub) listWorkflowFromId(workflowId int64) (*github.WorkflowRuns, error) {
	opts := &github.ListWorkflowRunsOptions{
		Actor:       "",
		Branch:      "",
		Event:       "",
		Status:      "",
		Created:     "",
		ListOptions: github.ListOptions{},
	}
	runs, _, err := ghClient.client.Actions.ListWorkflowRunsByID(ctx, "vmware-tanzu", "community-edition", workflowId, opts)
	if err != nil {
		log.Errorf("Workflows Listing failed. Err: %v\n", err)
		return nil, err
	}

	return runs, nil
}

func (ghClient *GitHub) listWorkflows() ([]*github.WorkflowRuns, error) {
	days := timeNow.AddDate(0, 0, -1).Format("2006-01-02")
	opts := &github.ListWorkflowRunsOptions{Status: "failure", Created: fmt.Sprint(">", days)}
	//opts := &github.ListWorkflowRunsOptions{Status: "failure", Created: ">2022-05-19"}
	opt := &github.ListOptions{}
	workflows, _, err := ghClient.client.Actions.ListWorkflows(ctx, owner, repo, opt)
	if err != nil {
		log.Errorf("Workflows Listing failed. Err: %v\n", err)
		return nil, err
	}
	TceWorkflow := []*github.WorkflowRuns{}
	for _, workflow := range workflows.Workflows {
		if strings.Contains(*workflow.Name, lookup) {
			log.Infof("%s %d", *workflow.Name, *workflow.ID)
			runs, _, err := ghClient.client.Actions.ListWorkflowRunsByID(ctx, owner, repo, *workflow.ID, opts)
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

func handleRequests() {
	port := "8080"
	log.Infof("Starting HTTP api server on localhost:%s", port)
	http.HandleFunc("/", home)
	http.HandleFunc("/viewjson", jsonview)
	http.HandleFunc("/viewhtml", view)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
}

func jsonview(w http.ResponseWriter, r *http.Request) {
	log.Infof("Trying to get Result for runs ")
	res, _ := mygh.getResults()
	w.Write(res)
}
func view(w http.ResponseWriter, r *http.Request) {
	log.Infof("Rendering HTML View")
	x, _ := mygh.getResultsStruct()
	t, _ := template.ParseFiles("templates/report.html")
	err := t.Execute(w, &x)
	if err != nil {
		log.Errorf("Something went wrong")
		return
	}
	log.Infof("%v", &x)

}
