//go:generate go-enum --marshal
package httpx

// ENUM(success, error)
type Status string

type APIResponse struct {
	Status Status `json:"status"`
	Data   any    `json:"data,omitempty"`
	Error  string `json:"error,omitempty"`
}
