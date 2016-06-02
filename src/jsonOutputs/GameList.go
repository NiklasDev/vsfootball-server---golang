package jsonOutputs

type GameListOutput struct {
	Success    string
	Message    string
	Games      []GameInList
	Statuscode int
}

type GameTurnOutput struct {
	Success    string
	Message    string
	Game       GameInList
	Statuscode int
}
