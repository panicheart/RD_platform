package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GiteaClient provides a client for Gitea API
type GiteaClient struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewGiteaClient creates a new Gitea API client
func NewGiteaClient(baseURL, token string) *GiteaClient {
	return &GiteaClient{
		baseURL: baseURL,
		 token:   token,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// CreateRepoRequest represents a request to create a repository
type CreateRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	AutoInit    bool   `json:"auto_init"`
}

// CreateRepo creates a new repository
func (c *GiteaClient) CreateRepo(owner string, req CreateRepoRequest) (*Repository, error) {
	url := fmt.Sprintf("%s/api/v1/user/repos", c.baseURL)
	return c.createRepo(url, req)
}

// CreateOrgRepo creates a new repository in an organization
func (c *GiteaClient) CreateOrgRepo(org string, req CreateRepoRequest) (*Repository, error) {
	url := fmt.Sprintf("%s/api/v1/orgs/%s/repos", c.baseURL, org)
	return c.createRepo(url, req)
}

func (c *GiteaClient) createRepo(url string, req CreateRepoRequest) (*Repository, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create repo: %s", resp.Status)
	}

	var repo Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, err
	}

	return &repo, nil
}

// GetRepo retrieves a repository
func (c *GiteaClient) GetRepo(owner, repo string) (*Repository, error) {
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s", c.baseURL, owner, repo)

	resp, err := c.doRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get repo: %s", resp.Status)
	}

	var repository Repository
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return nil, err
	}

	return &repository, nil
}

// DeleteRepo deletes a repository
func (c *GiteaClient) DeleteRepo(owner, repo string) error {
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s", c.baseURL, owner, repo)

	resp, err := c.doRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete repo: %s", resp.Status)
	}

	return nil
}

// ListRepos lists repositories for a user
func (c *GiteaClient) ListRepos(username string, page, pageSize int) ([]Repository, error) {
	url := fmt.Sprintf("%s/api/v1/users/%s/repos?page=%d&limit=%d", c.baseURL, username, page, pageSize)

	resp, err := c.doRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list repos: %s", resp.Status)
	}

	var repos []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}

// CreateFile creates a file in a repository
func (c *GiteaClient) CreateFile(owner, repo, filepath string, content []byte, message string) error {
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/contents/%s", c.baseURL, owner, repo, filepath)

	reqBody := map[string]string{
		"content":  base64Encode(content),
		"message":  message,
	}

	body, _ := json.Marshal(reqBody)
	resp, err := c.doRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create file: %s", resp.Status)
	}

	return nil
}

// GetCommits retrieves commits for a repository
func (c *GiteaClient) GetCommits(owner, repo, branch string, page, pageSize int) ([]Commit, error) {
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/commits?sha=%s&page=%d&limit=%d", 
		c.baseURL, owner, repo, branch, page, pageSize)

	resp, err := c.doRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get commits: %s", resp.Status)
	}

	var commits []Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, err
	}

	return commits, nil
}

// GetDiff retrieves diff between two commits
func (c *GiteaClient) GetDiff(owner, repo, base, head string) (string, error) {
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/compare/%s...%s", c.baseURL, owner, repo, base, head)

	resp, err := c.doRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get diff: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// doRequest performs an HTTP request with authentication
func (c *GiteaClient) doRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Content-Type", "application/json")

	return c.client.Do(req)
}

func base64Encode(data []byte) string {
	const base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	result := make([]byte, 0, len(data)*4/3+4)
	
	for i := 0; i < len(data); i += 3 {
		b1, b2, b3 := data[i], byte(0), byte(0)
		if i+1 < len(data) {
			b2 = data[i+1]
		}
		if i+2 < len(data) {
			b3 = data[i+2]
		}

		result = append(result, base64Chars[b1>>2])
		result = append(result, base64Chars[((b1&0x03)<<4)|(b2>>4)])
		if i+1 < len(data) {
			result = append(result, base64Chars[((b2&0x0f)<<2)|(b3>>6)])
		} else {
			result = append(result, '=')
		}
		if i+2 < len(data) {
			result = append(result, base64Chars[b3&0x3f])
		} else {
			result = append(result, '=')
		}
	}
	
	return string(result)
}

// Repository represents a Gitea repository
type Repository struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	HTMLURL     string `json:"html_url"`
	CloneURL    string `json:"clone_url"`
	DefaultBranch string `json:"default_branch"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Commit represents a Git commit
type Commit struct {
	SHA       string    `json:"sha"`
	Message   string    `json:"message"`
	Author    *Author  `json:"author"`
	Committer *Author  `json:"committer"`
	HTMLURL   string    `json:"html_url"`
}

// Author represents a Git author
type Author struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}
