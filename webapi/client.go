package webapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultTimeout represents the time to wait for a response from slack.
	DefaultTimeout = 10 * time.Second
)

const (
	postMessageEndpoint = "https://slack.com/api/chat.postMessage"
	filesUploadEndpoint = "https://slack.com/api/files.upload"
	usersListEndpoint   = "https://slack.com/api/users.list"
)

// Client represents a Slack client for Web API.
type Client struct {
	mux        sync.Mutex
	token      string
	httpclient *http.Client
	usersCache map[string]User
	debug      bool
}

// New creates a client with a bot token.
func New(token string, opts ...Option) (*Client, error) {
	ret := Client{
		token: token,
		httpclient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
	for _, opt := range opts {
		if err := opt(&ret); err != nil {
			return nil, err
		}
	}
	return &ret, nil
}

// PostMessage sends a message to the Slack channel.
// see. https://api.slack.com/methods/chat.postMessage
func (c Client) PostMessage(ctx context.Context, channelID string, msg string) (*MessageResponse, error) {
	body := url.Values{
		"token":   {c.token},
		"channel": {channelID},
		"text":    {msg},
	}.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postMessageEndpoint, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpclient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("slack chat.postMessage failed: %w", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response body read error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("slack chat.postMessage failed: %v, %q", resp.Status, body)
	}
	var ret MessageResponse
	if err := json.Unmarshal(b, &ret); err != nil {
		return nil, fmt.Errorf("response body unmarshal error: body=%q, %w", string(b), err)
	}
	return &ret, nil
}

// UploadImage uploads an image by files.upload API.
// see. https://api.slack.com/methods/files.upload
func (c Client) UploadImage(channels []string, title, fileName, fileType, comment string, img io.Reader) error {
	if c.token == "" {
		return fmt.Errorf("slack token is empty")
	}
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	part, err := mw.CreateFormFile("file", fileName)
	if err != nil {
		return fmt.Errorf("multipart create from file error, %v, %v", title, err)
	}
	if _, err := io.Copy(part, img); err != nil {
		return fmt.Errorf("file copy error, %v, %v", title, err)
	}
	// for slack settings
	settings := map[string]string{
		"token":           c.token,
		"channels":        strings.Join(channels, ","),
		"filetype":        fileType,
		"title":           title,
		"initial_comment": comment,
	}
	for k, v := range settings {
		if err := mw.WriteField(k, v); err != nil {
			return fmt.Errorf("write field error, %v:%v, %v", k, v, err)
		}
	}
	if err := mw.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", filesUploadEndpoint, &buf)
	if err != nil {
		return fmt.Errorf("slack files.uplad new request error, %v", err)
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	cl := &http.Client{Timeout: 10 * time.Second}
	resp, err := cl.Do(req)
	if err != nil {
		return fmt.Errorf("slack files.upload error, %v", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("response body read error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack files.upload status error: %s, %q", resp.Status, string(b))
	}
	return nil
}

// UsersList lists all users in a Slack team.
// see. https://api.slack.com/methods/users.list
func (c Client) UsersList() ([]User, error) {
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, usersListEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpclient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("slack chat.postMessage failed: %w", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response body read error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("slack chat.postMessage failed: %v", resp.Status)
	}
	var ul UsersListResponse
	if err := json.Unmarshal(b, &ul); err != nil {
		return nil, fmt.Errorf("response body unmarshal error: body=%q, %w", string(b), err)
	}
	if !ul.OK {
		if ul.Error == "missing_scope" {
			return nil, fmt.Errorf("%s: needed: %q, provided: %q", ul.Error, ul.Needed, ul.Provided)
		}
		return nil, errors.New(ul.Error)
	}
	return ul.Members, nil
}

// Users lists all users in a Slack team and returns it's userID map.
func (c Client) Users() (map[string]User, error) {
	list, err := c.UsersList()
	if err != nil {
		return nil, err
	}
	ret := map[string]User{}
	for _, v := range list {
		ret[v.ID] = v
	}
	return ret, nil
}

// RefreshUsersCache updates the client's cached user map.
func (c *Client) RefreshUsersCache() error {
	us, err := c.Users()
	if err != nil {
		return err
	}
	defer c.mux.Unlock()
	c.mux.Lock()
	c.usersCache = us
	return nil
}

// User returns the user corresponding to user ID from the client's user cache.
func (c Client) User(id string) (User, bool) {
	u, ok := c.usersCache[id]
	if ok {
		return u, true
	}
	return u, false
}
