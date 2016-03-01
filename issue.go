package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Issue struct {
	Title       string
	Description string
	AssigneeId  string
	Labels      string
}

type GitlabResponse struct {
	ProjectID int         `json:"project_id"`
	ID        int         `json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	Iid       int         `json:"iid"`
	Title     string      `json:"title"`
	State     string      `json:"state"`
	Assignee  interface{} `json:"assignee"`
	Labels    []string    `json:"labels"`
	Author    struct {
		Name      string      `json:"name"`
		AvatarURL interface{} `json:"avatar_url"`
		State     string      `json:"state"`
		WebURL    string      `json:"web_url"`
		ID        int         `json:"id"`
		Username  string      `json:"username"`
	} `json:"author"`
	Description interface{} `json:"description"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Milestone   interface{} `json:"milestone"`
}

func addIssue(project string) {
	reader := bufio.NewReader(os.Stdin)

	issue := Issue{}

	fmt.Print("Enter the title of an issue: ")

	issue.Title, _ = reader.ReadString('\n')

	fmt.Print("Enter the description of an issue: ")

	issue.Description = readWithVim()

	fmt.Printf("Enter the ID of a user to assign issue: (default %v) ", Config.DefaultAssigneeId)

	issue.AssigneeId, _ = reader.ReadString('\n')

	fmt.Print("Enter comma-separated label names for an issue: ")

	issue.Labels, _ = reader.ReadString('\n')

	confirmIssue(issue, project)
}

func confirmIssue(issue Issue, project string) {

	fmt.Println("Do you really want add this issue?\n")
	fmt.Printf("Title:\n %v", issue.Title)
	fmt.Printf("Description:\n %v", issue.Description)
	fmt.Printf("Labels:\n %v", issue.Labels)

	c := confirm("Store this issue? Y/N ")
	if c == true {
		storeIssue(issue, project)
	}
}

func storeIssue(issue Issue, project string) {
	appUrl := Config.Host + apiVerions + "projects/" + project + "/issues"

	form := url.Values{}

	form.Add("title", issue.Title)
	form.Add("description", issue.Description)
	form.Add("assignee_id", issue.AssigneeId)
	form.Add("labels", issue.Labels)

	req, err := http.NewRequest("POST", appUrl, strings.NewReader(form.Encode()))

	if err != nil {
		logger.Critical(err)
	}

	req.Header.Set("PRIVATE-TOKEN", Config.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Critical(err)
	}
	defer resp.Body.Close()

	var response GitlabResponse

	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &response)

	logger.Infof("Added issue with idd: %d\n", response.Iid)
}

func readWithVim() string {
	vi := "vim"
	tmpDir := os.TempDir()
	tmpFile, tmpFileErr := ioutil.TempFile(tmpDir, "tempFilePrefix")
	if tmpFileErr != nil {

		logger.Criticalf("Error %s while creating tempFile", tmpFileErr)
	}
	path, err := exec.LookPath(vi)
	if err != nil {

		logger.Criticalf("Error %s while looking up for %s!!", path, vi)
	}

	logger.Infof("%s is available at %s, Calling it with file %s \n", vi, path, tmpFile.Name())

	cmd := exec.Command(path, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		logger.Criticalf("Start failed: %s", err)
	}

	err = cmd.Wait()

	if err != nil {
		logger.Criticalf("Command finished with error: %v", err)
	}

	dat, _ := ioutil.ReadFile(tmpFile.Name())

	os.Remove(tmpFile.Name())

	return string(dat)
}
