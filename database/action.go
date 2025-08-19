package database

const (
	ActionBlockRequest = "block_request" // blocks the request
	ActionBlockIP      = "block_ip"
	ActionRedirect     = "redirect"
)

type Action struct {
	Type string `json:"type"`           // e.g "block_request"
	Data any    `json:"data,omitempty"` // additional data for the action, e.g. redirect URL
}
