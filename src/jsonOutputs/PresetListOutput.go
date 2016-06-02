package jsonOutputs

type WhatsOnTheLinePresetList struct {
	Success    string
	Message    string
	Statuscode int
	PresetList []string
}

type WhatsOnTheLineForAGame struct {
	Success          string
	Message          string
	Statuscode       int
	Confirmed        bool
	WhatsOnTheLine   string
	WhatsOnTheLineId int
}

type WhatsOnTheLineConfirmation struct {
	Success          string
	Message          string
	Statuscode       int
	WhatsOnTheLineId int
}
