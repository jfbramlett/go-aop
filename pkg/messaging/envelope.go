package messaging

type Envelope struct {
	Content interface{} `json:"content,omitempty"`
}
