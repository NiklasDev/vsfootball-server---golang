package models

type Turn struct {
	Id                  int64
	Gameid              int64
	Player1id           string
	Player2id           string
	Previousturn        int64
	Yardline            int
	Down                int
	Downdistance        int
	Player1playselected int64
	Player2playselected int64
	Player1role         int
	Player2role         int
	Results             string
	Playtime            int
	Timeelapsedingame   int
	Currentplayer1score int
	Currentplayer2score int
	Player1playdatetime int64
	Player2playdatetime int64
	Player1playflipped  int
	Player2playflipped  int
	Animationforturn    string
}
