package socketmode

// Option represents the client's option.
type Option func(*Client) error

// Debug is the debug option.
func Debug() Option {
	return func(c *Client) error {
		c.debug = true
		return nil
	}
}
