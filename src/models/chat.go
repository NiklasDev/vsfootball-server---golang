package models

// Put chat database functionality here

type Chatmessage struct {
	Id             int64
	Sender         string
	Senderteamname string
	Gameid         int64
	Created        int64
	Message        string
}
