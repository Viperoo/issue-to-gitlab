package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Projects []struct {
	ID              int           `json:"id"`
	Description     string        `json:"description"`
	DefaultBranch   string        `json:"default_branch"`
	TagList         []interface{} `json:"tag_list"`
	Public          bool          `json:"public"`
	Archived        bool          `json:"archived"`
	VisibilityLevel int           `json:"visibility_level"`
	SSHURLToRepo    string        `json:"ssh_url_to_repo"`
	HTTPURLToRepo   string        `json:"http_url_to_repo"`
	WebURL          string        `json:"web_url"`
	Owner           struct {
		Name      string `json:"name"`
		Username  string `json:"username"`
		ID        int    `json:"id"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"owner,omitempty"`
	Name                 string    `json:"name"`
	NameWithNamespace    string    `json:"name_with_namespace"`
	Path                 string    `json:"path"`
	PathWithNamespace    string    `json:"path_with_namespace"`
	IssuesEnabled        bool      `json:"issues_enabled"`
	MergeRequestsEnabled bool      `json:"merge_requests_enabled"`
	WikiEnabled          bool      `json:"wiki_enabled"`
	BuildsEnabled        bool      `json:"builds_enabled"`
	SnippetsEnabled      bool      `json:"snippets_enabled"`
	CreatedAt            time.Time `json:"created_at"`
	LastActivityAt       time.Time `json:"last_activity_at"`
	SharedRunnersEnabled bool      `json:"shared_runners_enabled"`
	CreatorID            int       `json:"creator_id"`
	Namespace            struct {
		ID          int         `json:"id"`
		Name        string      `json:"name"`
		Path        string      `json:"path"`
		OwnerID     int         `json:"owner_id"`
		CreatedAt   time.Time   `json:"created_at"`
		UpdatedAt   time.Time   `json:"updated_at"`
		Description string      `json:"description"`
		Avatar      interface{} `json:"avatar"`
	} `json:"namespace"`
	AvatarURL       interface{} `json:"avatar_url"`
	StarCount       int         `json:"star_count"`
	ForksCount      int         `json:"forks_count"`
	OpenIssuesCount int         `json:"open_issues_count"`
	PublicBuilds    bool        `json:"public_builds"`
	Permissions     struct {
		ProjectAccess struct {
			AccessLevel       int `json:"access_level"`
			NotificationLevel int `json:"notification_level"`
		} `json:"project_access"`
		GroupAccess interface{} `json:"group_access"`
	} `json:"permissions"`
	ForkedFromProject struct {
		ID                int    `json:"id"`
		Name              string `json:"name"`
		NameWithNamespace string `json:"name_with_namespace"`
		Path              string `json:"path"`
		PathWithNamespace string `json:"path_with_namespace"`
	} `json:"forked_from_project,omitempty"`
}

func listProjects() {

	projects := getProjects()

	for _, value := range projects {

		fmt.Printf("|%-6d|%8s\n", value.ID, value.NameWithNamespace)
	}
}

func getProjects() Projects {
	url := Config.Host + apiVerions + "projects?per_page=100"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(``)))

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

	var projects Projects

	body, _ := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(body, &projects)

	return projects
}
