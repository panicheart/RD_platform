package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	zoteroAPIBaseURL = "https://api.zotero.org"
	defaultTimeout   = 30 * time.Second
)

// ZoteroClient provides access to Zotero Web API v3
type ZoteroClient struct {
	apiKey     string
	userID     string
	httpClient *http.Client
	baseURL    string
}

// ZoteroConfig holds configuration for Zotero client
type ZoteroConfig struct {
	APIKey string
	UserID string
}

// ZoteroItem represents a Zotero library item from the API
type ZoteroItem struct {
	Key          string                 `json:"key"`
	Version      int                    `json:"version"`
	ItemType     string                 `json:"itemType"`
	Title        string                 `json:"title"`
	Creators     []ZoteroCreator        `json:"creators,omitempty"`
	AbstractNote string                 `json:"abstractNote,omitempty"`
	Publication  string                 `json:"publicationTitle,omitempty"`
	Volume       string                 `json:"volume,omitempty"`
	Issue        string                 `json:"issue,omitempty"`
	Pages        string                 `json:"pages,omitempty"`
	Date         string                 `json:"date,omitempty"`
	DOI          string                 `json:"DOI,omitempty"`
	URL          string                 `json:"url,omitempty"`
	Tags         []ZoteroTag            `json:"tags,omitempty"`
	Collections  []string               `json:"collections,omitempty"`
	Relations    map[string]interface{} `json:"relations,omitempty"`
}

// ZoteroCreator represents an author or contributor
type ZoteroCreator struct {
	CreatorType string `json:"creatorType"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName"`
	Name        string `json:"name,omitempty"` // For institutional authors
}

// ZoteroTag represents a tag attached to an item
type ZoteroTag struct {
	Tag  string `json:"tag"`
	Type int    `json:"type,omitempty"`
}

// ZoteroAttachment represents a file attachment
type ZoteroAttachment struct {
	Key         string `json:"key"`
	ItemType    string `json:"itemType"`
	Title       string `json:"title"`
	Filename    string `json:"filename,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	URL         string `json:"url,omitempty"`
}

// ZoteroCollection represents a Zotero collection/folder
type ZoteroCollection struct {
	Key              string `json:"key"`
	Version          int    `json:"version"`
	Name             string `json:"name"`
	ParentCollection string `json:"parentCollection,omitempty"`
}

// NewZoteroClient creates a new Zotero API client
func NewZoteroClient(config ZoteroConfig) (*ZoteroClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("Zotero API key is required")
	}
	if config.UserID == "" {
		return nil, fmt.Errorf("Zotero user ID is required")
	}

	return &ZoteroClient{
		apiKey: config.APIKey,
		userID: config.UserID,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: zoteroAPIBaseURL,
	}, nil
}

// SetBaseURL allows overriding the base URL for testing
func (c *ZoteroClient) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// TestConnection validates the API credentials
func (c *ZoteroClient) TestConnection() error {
	req, err := c.newRequest("GET", fmt.Sprintf("/users/%s/items?limit=1", c.userID), nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Zotero API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Zotero API returned status %d", resp.StatusCode)
	}

	return nil
}

// GetItems retrieves items from the user's library
func (c *ZoteroClient) GetItems(options *ItemQueryOptions) ([]ZoteroItem, int, error) {
	if options == nil {
		options = &ItemQueryOptions{}
	}

	path := fmt.Sprintf("/users/%s/items", c.userID)
	if options.CollectionKey != "" {
		path = fmt.Sprintf("/users/%s/collections/%s/items", c.userID, options.CollectionKey)
	}

	query := options.toQueryParams()
	req, err := c.newRequest("GET", path+"?"+query.Encode(), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch items: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("Zotero API returned status %d", resp.StatusCode)
	}

	// Parse total results from header
	totalResults := 0
	if totalHeader := resp.Header.Get("Total-Results"); totalHeader != "" {
		fmt.Sscanf(totalHeader, "%d", &totalResults)
	}

	var items []ZoteroItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return items, totalResults, nil
}

// GetItem retrieves a specific item by key
func (c *ZoteroClient) GetItem(itemKey string) (*ZoteroItem, error) {
	path := fmt.Sprintf("/users/%s/items/%s", c.userID, itemKey)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch item: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("item not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Zotero API returned status %d", resp.StatusCode)
	}

	var item ZoteroItem
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &item, nil
}

// GetItemChildren retrieves child items (attachments, notes) of an item
func (c *ZoteroClient) GetItemChildren(itemKey string) ([]ZoteroAttachment, error) {
	path := fmt.Sprintf("/users/%s/items/%s/children", c.userID, itemKey)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch children: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Zotero API returned status %d", resp.StatusCode)
	}

	var children []ZoteroAttachment
	if err := json.NewDecoder(resp.Body).Decode(&children); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return children, nil
}

// GetCollections retrieves all collections from the user's library
func (c *ZoteroClient) GetCollections() ([]ZoteroCollection, error) {
	path := fmt.Sprintf("/users/%s/collections", c.userID)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch collections: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Zotero API returned status %d", resp.StatusCode)
	}

	var collections []ZoteroCollection
	if err := json.NewDecoder(resp.Body).Decode(&collections); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return collections, nil
}

// GetTags retrieves all tags from the user's library
func (c *ZoteroClient) GetTags() ([]ZoteroTag, error) {
	path := fmt.Sprintf("/users/%s/tags", c.userID)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Zotero API returned status %d", resp.StatusCode)
	}

	var tags []ZoteroTag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return tags, nil
}

// GetPDFURL returns a signed URL for downloading a PDF attachment
func (c *ZoteroClient) GetPDFURL(attachmentKey string) (string, error) {
	// Zotero API doesn't provide direct download URLs for attachments
	// We return a path that can be used with the API key
	return fmt.Sprintf("%s/users/%s/items/%s/file/view", c.baseURL, c.userID, attachmentKey), nil
}

// DownloadAttachment downloads an attachment file
func (c *ZoteroClient) DownloadAttachment(attachmentKey string) ([]byte, string, error) {
	path := fmt.Sprintf("/users/%s/items/%s/file", c.userID, attachmentKey)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download attachment: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, "", fmt.Errorf("attachment is a linked URL, not a file")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("Zotero API returned status %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read attachment data: %w", err)
	}

	return data, contentType, nil
}

// ItemQueryOptions represents query parameters for fetching items
type ItemQueryOptions struct {
	Limit         int
	Start         int
	Sort          string
	Direction     string
	CollectionKey string
	Tag           string
	Q             string // Full-text search query
	ItemType      string
	Since         int64 // Version timestamp for sync
}

func (o *ItemQueryOptions) toQueryParams() url.Values {
	q := url.Values{}

	if o.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", o.Limit))
	} else {
		q.Set("limit", "25") // Default limit
	}

	if o.Start > 0 {
		q.Set("start", fmt.Sprintf("%d", o.Start))
	}

	if o.Sort != "" {
		q.Set("sort", o.Sort)
	}

	if o.Direction != "" {
		q.Set("direction", o.Direction)
	}

	if o.Tag != "" {
		q.Set("tag", o.Tag)
	}

	if o.Q != "" {
		q.Set("q", o.Q)
	}

	if o.ItemType != "" {
		q.Set("itemType", o.ItemType)
	}

	if o.Since > 0 {
		q.Set("since", fmt.Sprintf("%d", o.Since))
	}

	return q
}

// newRequest creates a new HTTP request with authentication
func (c *ZoteroClient) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	url := c.baseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Zotero-API-Version", "3")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// GetUserID returns the configured user ID
func (c *ZoteroClient) GetUserID() string {
	return c.userID
}

// GetAPIKey returns a masked version of the API key (for display)
func (c *ZoteroClient) GetAPIKey() string {
	if len(c.apiKey) <= 8 {
		return "***"
	}
	return c.apiKey[:4] + "..." + c.apiKey[len(c.apiKey)-4:]
}
