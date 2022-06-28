package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v44/github"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"golang.org/x/oauth2"
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
	Name     string `json:"Name"`
	Status   string `json:"Status"`
	Time     string `json:"Created Time"`
	Runtime  string `json:"Runtime"`
	Provider string `json:"Provider"`
	Url      string `json:"URL"`
}

const (
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
	serveType := flag.String("type", "serve", "Use this flag for serving api server or getting static report")
	flag.Parse()
	if *serveType == "serve" {
		handleRequests()
	}
}

func (ghClient *GitHub) getResult(date string) (Result, error) {
	res := &Result{
		Date:      fmt.Sprintf("%v", timeNow.Format(time.RFC822)),
		Status:    "",
		BuildType: "daily",
		Providers: []Providers{},
	}

	res.Status = "success"
	res.Providers = mygh.getProviderResults(date)
	res.Plugins = mygh.getPluginResults(date)
	log.Infof("%v", res)
	return *res, nil
}

//Todo Remove Redundant code for provide and plugin
func (ghClient *GitHub) getProviderResults(date string) []Providers {
	workflows := strings.Split(os.Getenv(workflowsEnv), ",")
	log.Infof("%s", workflows)
	w := []Providers{}
	for _, workflow := range workflows {
		log.Infof("checking for workflow id %s", workflow)
		w_id, _ := strconv.ParseInt(workflow, 10, 64)
		run, err := ghClient.listWorkflowFromId(w_id, date)
		if err != nil {
			log.Errorf("Error %s ", err)
		}
		y := run.WorkflowRuns[0]

		log.Infof("%s %s ", *y.Conclusion, *y.Name)
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
func (ghClient *GitHub) getPluginResults(date string) []Plugins {
	workflows := strings.Split(os.Getenv(pluginEnv), ",")
	log.Infof("%s", workflows)
	w := []Plugins{}
	for _, workflow := range workflows {
		log.Infof("checking for workflow id %s", workflow)
		w_id, _ := strconv.ParseInt(workflow, 10, 64)
		run, err := ghClient.listWorkflowFromId(w_id, date)
		if err != nil {
			log.Errorf("Error %s ", err)
		}
		y := run.WorkflowRuns[0]
		log.Infof("%s %s ", *y.Conclusion, *y.Name)
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

func (ghClient *GitHub) listWorkflowFromId(workflowId int64, date string) (*github.WorkflowRuns, error) {
	opts := &github.ListWorkflowRunsOptions{}
	if date != "" {
		opts = &github.ListWorkflowRunsOptions{Created: fmt.Sprintf("=%s", date)}
	}
	runs, _, err := ghClient.client.Actions.ListWorkflowRunsByID(ctx, owner, repo, workflowId, opts)
	//log.Infof("Getting from Date %s \n %v", date, opts)
	if err != nil {
		log.Errorf("Workflows Listing failed. Err: %v\n", err)
		return nil, err
	}

	return runs, nil
}

func handleRequests() {
	port := "8080"
	log.Infof("Starting HTTP api server on localhost:%s", port)
	http.HandleFunc("/", home)
	http.HandleFunc("/viewjson", jsonview)
	http.HandleFunc("/viewhtml", viewhtml)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
func home(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
}

func jsonview(w http.ResponseWriter, r *http.Request) {
	log.Infof("Trying to get Result for runs ")
	date := r.URL.Query().Get("date")
	res, _ := mygh.getResult(date)
	results, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Errorf("Error %s ", err)
	}
	w.Write(results)
}
func viewhtml(w http.ResponseWriter, r *http.Request) {
	log.Infof("Rendering HTML View")
	date := r.URL.Query().Get("date")
	results, _ := mygh.getResult(date)
	t, _ := template.ParseFiles("templates/report.html")
	err := t.Execute(w, &results)
	if err != nil {
		log.Errorf("Something went wrong")
		return
	}

}
