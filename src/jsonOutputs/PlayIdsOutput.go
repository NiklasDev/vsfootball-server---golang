package jsonOutputs

type PlayIdsOutput struct {
	Success    string
	Message    string
	Statuscode int
	Playids    []int64
}
