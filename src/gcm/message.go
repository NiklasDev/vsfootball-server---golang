package gcm

type Message struct {
	RegistrationIDs []string `json:"registration_ids"`
	//   RegistrationIDs string
	CollapseKey           string            `json:"collapse_key,omitempty"`
	Data                  map[string]string `json:"data,omitempty"`
	DelayWhileIdle        bool              `json:"delay_while_idle,omitempty"`
	TimeToLive            int               `json:"time_to_live,omitempty"`
	RestrictedPackageName string            `json:"restricted_package_name,omitempty"`
	DryRun                bool              `json:"dry_run,omitempty"`
}

// NewMessage returns a new Message with the specified payload
// and registration ids.
func NewMessage(data map[string]string, regIds ...string) *Message {
	return &Message{RegistrationIDs: regIds, Data: data}
}
