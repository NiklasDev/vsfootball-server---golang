package jsonOutputs

type StatsOutput struct {
	Success    string
	Message    string
	Statuscode int
	Wins       int
	Losses     int
	Ties       int
}

type OpponentsOutput struct {
	Success    string
	Message    string
	Statuscode int
	Opponents  []PlayerGuidEmailOutput
}

type PlayerGuidEmailOutput struct {
	Guid  string
	Email string
}
