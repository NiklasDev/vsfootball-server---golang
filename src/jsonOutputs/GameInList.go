package jsonOutputs

type TurnInGame struct {
	Id                     int64
	Player1id              string
	Player2id              string
	Previousturn           int64
	Yardline               int
	Down                   int
	Togo                   int
	Player1playselected    int64
	Player2playselected    int64
	Player1role            int
	Player2role            int
	Results                string
	Playtime               int
	Timeelapsedingame      int
	Currentplayer1score    int
	Currentplayer2score    int
	Player1playdatetime    int64
	Player2playdatetime    int64
	Player1allowedplays    string
	Player2allowedplays    string
	Ballon                 int
	Quarter                int
	Timeremaininginquarter string
	Offensivedirection     string
	Player1timeoutsleft    int
	Player2timeoutsleft    int
	Firstdownyardline      int
	SideOfField            string
	Player1playflipped     int
	Player2playflipped     int
	Animationforturn       []int64
}

type GameInList struct {
	Gameid          int64
	Player1         string
	Player2         string
	Player1teamname string
	Player2teamname string
	Outcome         string
	Inviteaccepted  int64
	Turns           []TurnInGame
	Lastupdated     int64
	Player1prefix   string
	Player2prefix   string
	Waitingtime     string
	Winner          string
}
