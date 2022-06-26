package webapi

// UsersListResponse is the response of the users.list API.
type UsersListResponse struct {
	OK       bool   `json:"ok"`
	Error    string `json:"error"`
	Needed   string `json:"needed"`
	Provided string `json:"provided"`
	Members  []User `json:"members"`
}

// User represents the Slack user.
type User struct {
	ID       string `json:"id,omitempty"`
	TeamID   string `json:"team_id,omitempty"`
	Name     string `json:"name,omitempty"`
	RealName string `json:"real_name,omitempty"`
	IsBot    bool   `json:"is_bot,omitempty"`
}
