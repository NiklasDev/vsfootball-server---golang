package jsonOutputs

type Animation struct {
	Id       int64
	Tier     int
	Image    string
	Duration float64
	Version  int
	Filesize float64
}

type AnimationOutput struct {
	Success      string
	Message      string
	Statuscode   int
	AnimationObj Animation
}

type AnimationIdsOutput struct {
	Success      string
	Message      string
	Statuscode   int
	Animationids []int64
}
