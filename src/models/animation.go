package models

type Animation struct {
	Id              int64
	Tier            int
	Image           string
	Frames          float64
	Framespersecond float64
	Category        string
	Version         int
	Filesize        float64
}
