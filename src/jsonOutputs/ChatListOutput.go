package jsonOutputs

type ChatSummary struct {
	Teamname   string
	Timepassed string
	Text       string
	Timestamp  int64
}

type ChatListOutput struct {
	Success    string
	Message    string
	GameId     string
	Chatlist   []ChatSummary
	StatusCode int
}
