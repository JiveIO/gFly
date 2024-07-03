package response

// ServerInfo struct to describe System info response.
type ServerInfo struct {
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
	Server string `json:"server"`
}
