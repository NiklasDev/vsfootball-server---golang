package models

import (
	//"crypto/sha1"
	"database/sql"
	// "encoding/base64"
	// "apns"
	"encoding/json"
	"fmt"
	"gcm"
	"github.com/coopernurse/gorp"
	// logger "github.com/llimllib/loglevel"
	"github.com/nu7hatch/gouuid"
	"github.com/stathat/amzses"
	"github.com/stathat/jconfig"
	_ "github.com/ziutek/mymysql/godrv"
	"io/ioutil"
	"jsonOutputs"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

/**
This function creates a connection to the database. It shouldn't have to know anything
about the pool, It will be called N times where N is the size of the requested pool.
*/
func initDbConnections() (*gorp.DbMap, error) {
	db, connectionError := sql.Open("mymysql", "tcp:"+dbaddress+":"+port+"*"+"vsfootball/"+dbusername+"/"+dbpassword)
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "LATIN_SWEDISH_CI"}}
	fmt.Println(connectionError)
	if isDbQueryError(connectionError) {
		return nil, connectionError
	}
	dbmap.AddTableWithName(User{}, "user").SetKeys(true, "Id")
	dbmap.AddTableWithName(Feedback{}, "feedback").SetKeys(true, "Id")
	dbmap.AddTableWithName(Game{}, "game").SetKeys(true, "Id")
	dbmap.AddTableWithName(Turn{}, "turn").SetKeys(true, "Id")
	dbmap.AddTableWithName(Play{}, "play").SetKeys(true, "Id")
	dbmap.AddTableWithName(Animation{}, "animation").SetKeys(true, "Id")
	dbmap.AddTableWithName(Androiddevice{}, "androiddevice").SetKeys(true, "Id")
	dbmap.AddTableWithName(Iosdevice{}, "iosdevice").SetKeys(true, "Id")
	dbmap.AddTableWithName(Storepurchase{}, "storepurchase").SetKeys(true, "Id")
	dbmap.AddTableWithName(Chatmessage{}, "chatmessage").SetKeys(true, "Id")
	dbmap.AddTableWithName(Usergamelookup{}, "usergamelookup").SetKeys(true, "Id")
	//Add Preset
	dbmap.AddTableWithName(Preset{}, "preset").SetKeys(true, "Id")
	dbmap.CreateTables()
	dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
	// log.SetPriorityString("info")
	// log.Info("doing stuff")
	return dbmap, nil
}

var (
	environment    = os.Getenv("environment")
	config         = jconfig.LoadConfig("/etc/aws.conf")
	serverlocation = config.GetString("serverlocation")
	port           = config.GetString("dbport")
	dbusername     = config.GetString("dbusername")
	dbpassword     = config.GetString("dbpassword")
	dbaddress      = config.GetString("dbaddress")
	nodeName       = config.GetString("node_name")
	db, _          = sql.Open("mymysql", "tcp:"+dbaddress+":"+port+"*"+"vsfootball/"+dbusername+"/"+dbpassword)
	plays          map[int64]*Play
	playIds        []int64
	animations     map[string][]*Animation
	animationsById map[int64]*Animation
	animationIds   []int64
	ImageLocation  = serverlocation + "static/"
	pem            string
	dbPool         = &ConnectionPoolWrapper{}
)

func CloseAll() {
	con := dbPool.GetConnection()
	i := 0
	for ; con != nil; con = dbPool.GetConnection() {
		fmt.Println(con)
		fmt.Println(i)
		con.Db.Close()
		i += 1
	}
}
func Init() {
	err := dbPool.InitPool(100, initDbConnections)
	isDbQueryError(err)
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)

	var playsFromDatabase []*Play
	gatherAllPlayDataFromDatabase := "select * from play"
	_, playDataGatherError := dbmap.Select(&playsFromDatabase, gatherAllPlayDataFromDatabase)
	if isDbQueryError(playDataGatherError) {
		fmt.Println(playDataGatherError)
	} else {
		plays = make(map[int64]*Play)
	}
	for play := range playsFromDatabase {
		plays[playsFromDatabase[play].Id] = playsFromDatabase[play]
		playIds = append(playIds, playsFromDatabase[play].Id)
	}
	gatherAllAnimationsFromDatabase := "select * from animation"
	var animationsFromDatabase []*Animation
	_, animationGatherError := dbmap.Select(&animationsFromDatabase, gatherAllAnimationsFromDatabase)
	if isDbQueryError(animationGatherError) {
		fmt.Println(animationGatherError)
	} else {
		animations = make(map[string][]*Animation)
		animationsById = make(map[int64]*Animation)
	}
	for index := range animationsFromDatabase {
		fmt.Println(animationsFromDatabase[index].Category)
		animations[animationsFromDatabase[index].Category] = append(animations[animationsFromDatabase[index].Category], animationsFromDatabase[index])
		animationIds = append(animationIds, animationsFromDatabase[index].Id)
		animationsById[animationsFromDatabase[index].Id] = animationsFromDatabase[index]
	}
	fmt.Println(playIds)
	fmt.Println(animationIds)

	fmt.Println("database initialized.")
}

func SendEmailVerification(email string) (string, string, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	query := "select user.Guid from user where Username =? limit 1"
	var results []*User
	_, err := dbmap.Select(&results, query, email)
	if isDbQueryError(err) {
		return "false", "We've experienced a problem on our end.", 501
	} else {
		if len(results) > 0 {
			message := "Thank you for creating an account with Vs. Football. Please verify your email by clicking the link below\n" + serverlocation + results[0].Guid + "/verify" + ":, and then you're ready to play some football. \nVs. Football is the best football strategy game on the App Store and Google Play where you play football versus your friends. \nDownload the free Vs. Football app for iPhone or Android.  Let's play some football!\n\nLink to App on iTunes\nLink to App on Google Play Store\nwww.PlayVsFootball.com"
			amzses.SendMail("GameOn@PlayVsFootball.com", email, "Verification of your Vsfootball account.", message)
			// amzses.SendMail("GameOn@PlayVsFootball.com", email, "Verifcation of your Vsfootball account.", "Please verify your Vsfootball account at "+serverlocation+results[0].Guid+"/verify")
			return "true", "Email verification has been sent.", 200
		} else {
			return "false", "Account could not be found.", 401
		}
	}
}

func ForgotPassword(email string) (string, string, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	query := "select user.Guid from user where Username =? limit 1"
	var results []*User
	_, userSearchErr := dbmap.Select(&results, query, email)

	if isDbQueryError(userSearchErr) {
		return "false", "Database error", 501
	} else {
		if len(results) > 0 {
			rawGuid, guidErr := uuid.NewV4()
			if guidErr != nil {
				fmt.Println(guidErr.Error())
			}
			guid := strings.Split(rawGuid.String(), "-")
			updateQuery := "update user set user.Password= SHA1(?) , user.Updated= ? where Guid=?"
			_, err := dbmap.Exec(updateQuery, guid[0], time.Now().Unix(), results[0].Guid)
			if isDbQueryError(err) {
				return "false", "Database error", 501
			} else {
				message := "Your password is below:\n" + guid[0] + " Thanks for playing Vs. Football. \n\nVs. Football is the best football strategy game on the App Store and Google Play where you play football versus your friends. \n\nDownload the free Vs. Football app for iPhone or Android.  Let's play some football!\n\nLink to App on iTunes\nLink to App on Google Play Store  \nwww.PlayVsFootball.com"
				amzses.SendMail("GameOn@PlayVsFootball.com", email, "Vs. Football Password", message)
				return "true", "Email verification has been sent.", 200
			}
		} else {
			return "false", "Account does not exist.", 401
		}
	}
}

func AddFeedback(feedback *Feedback) (string, string, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	err := dbmap.Insert(feedback)
	if isDbQueryError(err) {
		return "false", "Database error.", 501
	} else {
		go getAndSendFeedbackEmails(feedback)
		return "true", "Successfully added feedback.", 200
	}
}

func getAndSendFeedbackEmails(feedback *Feedback) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	fmt.Println("inside feedback email.")
	userSelectQuery := "select * from user where user.Guid=? limit 1"
	var user []*User
	_, userSelectError := dbmap.Select(&user, userSelectQuery, feedback.Userid)
	if !isDbQueryError(userSelectError) && len(user) > 0 {
		message := "Feedback from " + user[0].Username + ": \n" + feedback.Comment + " on game id " + feedback.Gameid + "\n on OS:" + feedback.Os
		go amzses.SendMail("GameOn@PlayVsFootball.com", "Dclawson@engagemobile.com", "FeedBack:"+serverlocation, message)
		go amzses.SendMail("GameOn@PlayVsFootball.com", "zhenzhelu@hengtiansoft.com", "FeedBack:"+serverlocation, message)
		go amzses.SendMail("GameOn@PlayVsFootball.com", "ravi@engagemobile.com", "FeedBack:"+serverlocation, message)
		go amzses.SendMail("GameOn@PlayVsFootball.com", "stimperley@engagemobile.com", "FeedBack:"+serverlocation, message)
		go amzses.SendMail("GameOn@PlayVsFootball.com", "chris@playvsfootball.com", "FeedBack:"+serverlocation, message)
		go amzses.SendMail("GameOn@PlayVsFootball.com", "mbarksdale@engagemobile.com", "FeedBack:"+serverlocation, message)
		go amzses.SendMail("GameOn@PlayVsFootball.com", "ksamuelson@playvsfootball.com", "FeedBack:"+serverlocation, message)
	}

}

//Messages need to come from here not controller
func CreateAccount(user *User) (string, string, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	query := "select * from user where user.Username =? limit 1"
	var results []*User
	_, userSearchErr := dbmap.Select(&results, query, user.Username)
	fmt.Println(results)
	fmt.Println(userSearchErr)
	if isDbQueryError(userSearchErr) {
		return "false", "Database error", 501
	}
	if len(results) > 0 && results[0].Verified == 2 {
		updateQuery := "update user set user.Firstname=?, user.Lastname=?,user.Platform=?,user.Updated=?,user.Password=? where user.Username=?"
		_, updateErr := dbmap.Exec(updateQuery, user.Firstname, user.Lastname, user.Platform, time.Now().Unix(), user.Password, user.Username)
		if isDbQueryError(updateErr) {
			return "false", "Database error", 501
		} else {
			return "true", "Account was successfully created.", 200
		}
	} else {
		if len(results) == 0 {
			fmt.Println("in create account.", user)
			userInsertError := dbmap.Insert(user)
			if isDbQueryError(userInsertError) {
				return "false", "Database error", 501
			}
			message := "Thank you for creating an account with Vs. Football. Please verify your email by clicking the link below:\n " + serverlocation + user.Guid + "/verify" + " , and then you're ready to play some football. \nVs. Football is the best football strategy game on the App Store and Google Play where you play football versus your friends. \nDownload the free Vs. Football app for iPhone or Android.  Let's play some football!\n\nLink to App on iTunes\nLink to App on Google Play Store\nwww.PlayVsFootball.com"
			go amzses.SendMail("GameOn@PlayVsFootball.com", user.Username, "Verification of your Vsfootball account.", message)
			return "true", "Account was successfully created", 200
		} else {
			return "false", "Account already exists.", 403
		}
	}
}

//Messages need to come from here not controller
func AccountLogin(username, password string) (string, string, string, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	query := "select user.Guid , user.Password, user.Verified from user where Username =? limit 1"
	var results []*User
	_, err := dbmap.Select(&results, query, username)
	if isDbQueryError(err) {
		return "", "false", "Database error", 501
	} else {
		if len(results) > 0 {
			if password == results[0].Password {
				if results[0].Verified == 0 {
					return results[0].Guid, "false", "Account has not been verified", 402
				} else {
					return results[0].Guid, "true", "Successful login", 200
				}
			} else {
				return "", "false", "Invalid Username/Password or Account Does Not Exist.", 400
			}
		} else {
			return "", "false", "Account does not exist.", 401
		}
	}
}

func FacebookSignup(user *User) (string, string, string, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	query := "select * from user where user.Accounttype=? limit 1"
	var results []*User
	_, err := dbmap.Select(&results, query, user.Accounttype)
	if isDbQueryError(err) {
		return "false", "Database error.", "", 501
	}
	if len(results) > 0 {
		return "true", "Account Exists", results[0].Guid, 200
	} else {
		if user.Username == "" || user.Username == "<null>" || len(user.Accounttype) < 10 {
			return "false", "Looks like your Facebook email is private, so we can't use your email to create an account.  Please try creating an account using your email address.", "", 405
		}
		createError := dbmap.Insert(user)
		if isDbQueryError(createError) {
			if len(user.Accounttype) > 9 {
				updateQuery := "update user set user.Accounttype=? where user.Username=?"
				_, updateFailure := dbmap.Exec(updateQuery, user.Accounttype, user.Username)
				if isDbQueryError(updateFailure) {
					return "false", "Database failure", "", 501
				} else {
					return "true", "Successful login.", results[0].Guid, 200
				}
			}
		}
		return "true", "Account created.", user.Guid, 201
	}
}

func VerifyAccount(guid string) (bool, string) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	query := "update user set user.Verified=1, user.Updated=? where Guid = ?"
	_, err := dbmap.Exec(query, time.Now().Unix(), guid)
	isDbQueryError(err)
	return true, "/static/VsFootballAccountVerified.html"
}
func GetPlaybook(guid string) (string, string, []jsonOutputs.PlayInPlaybook, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	playsQuery := "select * from play"

	var playbook []*Play
	_, selectPlaybookError := dbmap.Select(&playbook, playsQuery)
	var outputPlaybook []jsonOutputs.PlayInPlaybook
	for play := range playbook {
		image := ImageLocation + playbook[play].Image
		if playbook[play].Image == "" {
			image = ""
		}
		outputPlaybook = append(outputPlaybook, jsonOutputs.PlayInPlaybook{
			Id:         playbook[play].Id,
			Title:      playbook[play].Title,
			Image:      image,
			Type:       playbook[play].Type,
			Possession: playbook[play].Possession,
			Premium:    playbook[play].Premium})
	}
	if selectPlaybookError == nil {
		return "true", "Successfully gathered playbook data.", outputPlaybook, 200
	} else {
		return "false", "Failure to gather playbook data.", nil, 501
	}

}

func FacebookInviteToGame(guid, facebookId, possession, teamName, playId, flippedPlay string, whatsOnTheLine string) (string, string, int, int64) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	if possession != "4" && possession != "5" {
		return "false", "Possession input not valid.", 408, -1
	}
	// if flippedPlay != "0" && flippedPlay != "1" {
	// 	return "false", "FlippedPlay input not valid.", 408, -1
	// }
	inviteeQuery := "select * from user where user.Accounttype=? limit 1"
	var inviteeResults []*User
	_, inviteeSelectError := dbmap.Select(&inviteeResults, inviteeQuery, "facebook:"+facebookId)
	if isDbQueryError(inviteeSelectError) {
		return "false", "Database failure.", 501, -1
	}
	if len(inviteeResults) <= 0 {
		return "false", "User does not have a Vs Football account connected with Facebook.", 421, -1
	}
	var player2Possession int
	if possession == "4" {
		player2Possession = 5
	} else {
		player2Possession = 4
	}
	var inviterResults []*User
	inviterQuery := "select * from user where user.Guid=? limit 1"
	_, inviterSelectError := dbmap.Select(&inviterResults, inviterQuery, guid)
	if isDbQueryError(inviterSelectError) {
		return "false", "Database failure.", 501, -1
	}
	if len(inviterResults) <= 0 {
		return "false", "Could not find user.", 401, -1
	}
	invitee := inviteeResults[0]
	inviter := inviterResults[0]
	if invitee.Guid == inviter.Guid {
		return "false", "Can't play with one's self.", 421, -1
	}
	player2teamname := invitee.Username
	if len(player2teamname) > 10 {
		player2teamname = player2teamname[0:10]
	}
	var flippedPlayInt, _ = strconv.Atoi(flippedPlay)
	player1Possession, _ := strconv.Atoi(possession)
	if inviter.Verified != 0 {
		playIdInt, _ := strconv.ParseInt(playId, 10, 64)
		success, message, gameId := innerCreateGameTurn(guid, invitee.Guid, inviter.Username, invitee.Username, teamName, player2teamname, player1Possession, player2Possession, flippedPlayInt, playIdInt,whatsOnTheLine)
		if success {
			pushmessage := "Invitation from " + inviter.Username + " to play a game."
			go sendAndroidMessage(invitee.Guid, pushmessage, "")
			go sendIosMessage(invitee.Guid, pushmessage, "", "")
			return "true", message, 200, gameId
		} else {
			return "false", message, 501, -1
		}
	} else {
		return "false", "Inviter user not verified.", 402, -1
	}
	return "false", "Database failure.", 501, -1
}

func Rematch(guid, invitedGuid, possession, teamName, playIdSelected, flippedPlay string, whatsOnTheLine string) (string, string, int, int64) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	if possession != "4" && possession != "5" {
		return "false", "Possession input not valid.", 408, -1
	}
	inviter := GetUserByGuid(guid)
	if inviter == nil || inviter.Verified == 0 {
		return "false", "Inviter guid is not valid or user does not exist.", 401, -1
	}
	var player2Possession, player1Possession int
	if possession == "4" {
		player1Possession = 4
		player2Possession = 5
	} else {
		player1Possession = 5
		player2Possession = 4
	}
	var flippedPlayInt, _ = strconv.Atoi(flippedPlay)
	results := GetUserByGuid(invitedGuid)
	if results != nil && results.Verified != 0 {
		inviteeShortName := results.Username
		if len(inviteeShortName) > 10 {
			inviteeShortName = inviteeShortName[0:10]
		}
		playIdInt, _ := strconv.ParseInt(playIdSelected, 10, 64)
		success, message, gameId := innerCreateGameTurn(guid, invitedGuid, inviter.Username, results.Username, teamName, inviteeShortName, player1Possession, player2Possession, flippedPlayInt, playIdInt, whatsOnTheLine)
		if success {
			pushmessage := "Invitation from " + inviter.Username + " to play a game."
			go sendAndroidMessage(results.Guid, pushmessage, "")
			go sendIosMessage(results.Guid, pushmessage, "", "")
			return "true", message, 200, gameId
		} else {
			return "false", message, 501, -1
		}
	} else {
		return "false", "Invitee guid is not valid or user does not exist.", 402, -1
	}

	return "false", "Error occured during game creation", 501, -1
}

func CreateGame(guid, inviteEmail, possession, teamName, playIdSelected, flippedPlay string, whatsOnTheLine string) (string, string, int64, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)

	if possession != "4" && possession != "5" {
		return "Possession input not valid.", "false", -1, 408
	}
	//valid inviteEmail
	//matched result: [[kun.he@ht.com.cn kun.he ht.com.cn ht com.cn]]
	emailName := inviteEmail
	if len(emailName) > 10 {
		emailName = emailName[0:10]
	}
	inviter := GetUserByGuid(guid)
	if inviter == nil {
		return "Inviter guid is not valid or user does not exist.", "false", -1, 401
	}
	var flippedPlayInt, _ = strconv.Atoi(flippedPlay)
	var player2Possession, player1Possession int
	if possession == "4" {
		player1Possession = 4
		player2Possession = 5
	} else {
		player1Possession = 5
		player2Possession = 4
	}
	//todo://sql inject
	query := "select user.Guid , user.Verified, user.Username from user where Username =? limit 1"
	var results []*User
	_, userSearchErr := dbmap.Select(&results, query, inviteEmail)
	if isDbQueryError(userSearchErr) {
		return "Database error.", "false", -1, 501
	}
	playIdInt, _ := strconv.ParseInt(playIdSelected, 10, 64)
	if len(results) > 0 {
		fmt.Println(results[0])
		if results[0].Guid == inviter.Guid {
			return "Can't play with one's self.", "false", -1, 401
		}

		if results[0].Verified != 0 {
			success, message, gameId := innerCreateGameTurn(guid, results[0].Guid, inviter.Username, results[0].Username, teamName, emailName, player1Possession, player2Possession, flippedPlayInt, playIdInt,whatsOnTheLine)
			if success {
				emailMessage := "You’ve been challenged to a game of Vs. Football by " + inviter.Username + ". Click a link below or head to PlayVsFootball.com to download the free app and start playing. \n\nVs. Football is the best football strategy game on the App Store and Google Play where you play football versus your friends. \n\nDownload the free Vs. Football app for iPhone or Android.  Let's play some football!\n\nApp on iTunes:\nhttps://itunes.apple.com/us/app/vs-football/id700513212?ls=1&mt=8 \nPlay Store:\nhttp://play.google.com/store/apps/details?id=com.engagemobile.vsfootball\n/www.PlayVsFootball.com"
				go amzses.SendMail("GameOn@PlayVsFootball.com", inviteEmail, "Invitation from "+inviter.Username+" to join Vsfootball.", emailMessage)
				pushmessage := "Invitation from " + inviter.Username + " to play a game."
				go sendAndroidMessage(results[0].Guid, pushmessage, "")
				go sendIosMessage(results[0].Guid, pushmessage, "", "")
				return message, "true", gameId, 200
			} else {
				return message, "false", -1, 501
			}

		} else {
			return "User is not verified", "false", -1, 402
		}
	} else {
		player2Guid, guidErr := uuid.NewV4()
		if guidErr != nil {
			fmt.Println(guidErr.Error())
		}
		user := &User{
			Created:         time.Now().Unix(),
			Updated:         -1,
			Firstname:       "",
			Lastname:        "",
			Guid:            strings.Replace(player2Guid.String(), "-", "", -1),
			Password:        "",
			Accesstoken:     "",
			Accounttype:     "",
			Tokenexpiration: "",
			Username:        inviteEmail,
			Verified:        2,
			Playsowned:      "",
			Platform:        ""}
		createAccountError := dbmap.Insert(user)
		if !isDbQueryError(createAccountError) {
			success, message, gameId := innerCreateGameTurn(guid, user.Guid, inviter.Username, user.Username, teamName, emailName, player1Possession, player2Possession, flippedPlayInt, playIdInt,whatsOnTheLine)
			if success {
				emailMessage := "You’ve been challenged to a game of Vs. Football by " + inviter.Username + ". Click a link below or head to PlayVsFootball.com to download the free app and start playing. \n\nVs. Football is the best football strategy game on the App Store and Google Play where you play football versus your friends. \n\nDownload the free Vs. Football app for iPhone or Android.  Let's play some football!\n\nApp on iTunes:\nhttps://itunes.apple.com/us/app/vs-football/id700513212?ls=1&mt=8 \nPlay Store:\nhttp://play.google.com/store/apps/details?id=com.engagemobile.vsfootball\n/www.PlayVsFootball.com"
				go amzses.SendMail("GameOn@PlayVsFootball.com", inviteEmail, "Invitation from "+inviter.Username+" to join Vsfootball.", emailMessage)
				return message, "true", gameId, 200
			} else {
				return message, "false", -1, 501
			}
		}
	}
	return "Error occured during game creation", "false", -1, 501
}
func innerCreateGameTurn(player1, player2, player1Email, player2Email, player1TeamName, player2TeamName string, player1Possession, player2Possession, flippedPlay int, playId int64, whatsOnTheLine string) (bool, string, int64) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	game := &Game{
		Created:             time.Now().Unix(),
		Player1:             player1,
		Player2:             player2,
		Player1teamname:     player1TeamName,
		Player2teamname:     player2TeamName,
		Outcome:             "",
		Player1accepteddate: time.Now().Unix(),
		Player2accepteddate: -1,
		Enddate:             -1,
		Inviteaccepted:      0,
		Lastupdated:         time.Now().Unix(),
		Player1prefix:       strings.Split(player1Email, "@")[0],
		Player2prefix:       strings.Split(player2Email, "@")[0],
		WhatsOnTheLine:      whatsOnTheLine,
	}
	gameSaveErr := dbmap.Insert(game)
	if !isDbQueryError(gameSaveErr) {
		turn := &Turn{
			Gameid:              game.Id,
			Player1id:           game.Player1,
			Player2id:           game.Player2,
			Previousturn:        -1,
			Yardline:            35,
			Down:                1,
			Downdistance:        10,
			Player1playselected: playId,
			Player2playselected: -1,
			Player1role:         player1Possession,
			Player2role:         player2Possession,
			Results:             "",
			Timeelapsedingame:   0,
			Currentplayer1score: 0,
			Currentplayer2score: 0,
			Player1playdatetime: time.Now().Unix(),
			Player2playdatetime: 0,
			Player1playflipped:  flippedPlay}
		turnSaveError := dbmap.Insert(turn)
		if !isDbQueryError(turnSaveError) {
			userGameLookupPlayer1 := &Usergamelookup{
				Userid: game.Player1,
				Gameid: game.Id}
			userGameLookupPlayer2 := &Usergamelookup{
				Userid: game.Player2,
				Gameid: game.Id}
			lookupPlayer1SaveError := dbmap.Insert(userGameLookupPlayer1)
			lookupPlayer2SaveError := dbmap.Insert(userGameLookupPlayer2)
			if isDbQueryError(lookupPlayer1SaveError) || isDbQueryError(lookupPlayer2SaveError) {
				return false, "Failed to save game data.", -1
			} else {
				return true, "Awaiting confirmation from Player 2", game.Id
			}
		} else {
			return false, "Failed to save turn data.", -1
		}
	} else {
		return false, "Failed to save game data.", -1
	}
}

func GamesList(guid string) (bool, string, []jsonOutputs.GameInList) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	query := "select * from game where game.player1 = '" + guid + "' or game.player2= '" + guid + "'"
	var gameList []*Game
	_, gamesListError := dbmap.Select(&gameList, query)
	var gameListJson []jsonOutputs.GameInList
	if gamesListError == nil {
		for game := range gameList {
			query := "select * from turn where turn.Gameid=" + strconv.FormatInt(gameList[game].Id, 10) + " order by turn.Id desc limit 1"
			var turnList []*Turn
			_, turnsErr := dbmap.Select(&turnList, query)
			fmt.Println(turnList[0])
			if turnsErr == nil {
				var turnsInGame []jsonOutputs.TurnInGame
				for turn := range turnList {
					turnsInGame = append(turnsInGame, jsonOutputs.TurnInGame{
						Id:                  turnList[turn].Id,
						Player1id:           turnList[turn].Player1id,
						Player2id:           turnList[turn].Player2id,
						Previousturn:        turnList[turn].Previousturn,
						Yardline:            turnList[turn].Yardline,
						Down:                turnList[turn].Down,
						Togo:                turnList[turn].Downdistance,
						Player1playselected: turnList[turn].Player1playselected,
						Player2playselected: turnList[turn].Player2playselected,
						Player1role:         turnList[turn].Player1role,
						Player2role:         turnList[turn].Player2role,
						Results:             turnList[turn].Results,
						Playtime:            turnList[turn].Playtime,
						Timeelapsedingame:   turnList[turn].Timeelapsedingame,
						Currentplayer1score: turnList[turn].Currentplayer1score,
						Currentplayer2score: turnList[turn].Currentplayer2score,
						Player1playdatetime: turnList[turn].Player1playdatetime,
						Player2playdatetime: turnList[turn].Player2playdatetime})
				}
				gameListJson = append(gameListJson, jsonOutputs.GameInList{
					Inviteaccepted:  gameList[game].Inviteaccepted,
					Gameid:          gameList[game].Id,
					Player1:         gameList[game].Player1,
					Player2:         gameList[game].Player2,
					Player1teamname: gameList[game].Player1teamname,
					Player2teamname: gameList[game].Player2teamname,
					Outcome:         gameList[game].Outcome,
					Lastupdated:     gameList[game].Lastupdated,
					Turns:           turnsInGame})
			} else {
				fmt.Println(turnsErr)
				return false, "Error gathering turn data", nil
			}
		}
	} else {
		fmt.Println(gamesListError)
		return false, "Error gathering games.", nil
	}
	return true, "Successfully gathered game/turn data.", gameListJson
}

func GameIdsForUser(guid string) (string, string, []int64, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	gameIdsForUserQuery := "select usergamelookup.Gameid from usergamelookup where usergamelookup.Userid=?"
	var gameList []*Usergamelookup
	_, gameIdsError := dbmap.Select(&gameList, gameIdsForUserQuery, guid)
	if isDbQueryError(gameIdsError) {
		return "false", "Error while gathering game list.", nil, 401
	} else {
		var outputArray []int64
		for _, game := range gameList {
			outputArray = append(outputArray, game.Gameid)
		}
		return "true", "Successfully gathered game IDs.", outputArray, 200
	}
}

func TurnsForGame(guid, gameId, returnLimit string) (string, string, jsonOutputs.GameInList, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	gameQuery := "select * from game where game.Id=" + gameId + " limit 1"
	var gameList []*Game
	_, gameSelectError := dbmap.Select(&gameList, gameQuery)
	if isDbQueryError(gameSelectError) {
		return "false", "Database error.", jsonOutputs.GameInList{}, 501
	}
	if len(gameList) == 0 {
		return "false", "Could not find game.", jsonOutputs.GameInList{}, 406
	} else {
		var turns []*Turn
		turnQuery := "select * from turn where turn.Gameid='" + gameId + "' order by turn.Id desc"
		if returnLimit != "0" {
			turnQuery += " limit " + returnLimit
		}

		_, turnSelectError := dbmap.Select(&turns, turnQuery)
		if isDbQueryError(turnSelectError) {
			return "false", "Error while retrieving turn data.", jsonOutputs.GameInList{}, 501
		} else {
			var turnDataList []jsonOutputs.TurnInGame
			for _, turn := range turns {
				player1AllowedPlays := ""
				if turn.Player1role == 0 || turn.Player1role == 2 {
					player1AllowedPlays = "0,2"
				} else {
					if turn.Player1role == 1 || turn.Player1role == 3 {
						player1AllowedPlays = "1,3"
					} else {
						if turn.Player1role == 4 {
							player1AllowedPlays = "4"
						} else {
							player1AllowedPlays = "5"
						}
					}

				}
				player2AllowedPlays := ""
				if turn.Player2role == 0 || turn.Player2role == 2 {
					player2AllowedPlays = "0,2"
				} else {
					if turn.Player2role == 1 || turn.Player2role == 3 {
						player2AllowedPlays = "1,3"
					} else {
						if turn.Player2role == 4 {
							player2AllowedPlays = "4"
						} else {
							player2AllowedPlays = "5"
						}
					}
				}
				lineOfScrimage := turn.Yardline
				if turn.Yardline > 50 {
					lineOfScrimage = 100 - turn.Yardline
				}

				quarter := turn.Timeelapsedingame / 720
				maxQuarterTime := 720 * (quarter + 1)
				timeDelta := maxQuarterTime - turn.Timeelapsedingame
				min := strconv.Itoa(timeDelta / 60)
				sec := strconv.Itoa(timeDelta % 60)

				// timeRemainingInQuarter := turn.Timeelapsedingame - (720 * (turn.Timeelapsedingame / 720))
				// min := strconv.Itoa(12 - (timeRemainingInQuarter / 60))
				// sec := "0"
				// if timeRemainingInQuarter%60 != 0 {
				// sec = strconv.Itoa(60 - (timeRemainingInQuarter % 60))
				// }

				fmt.Println(len(sec))
				if len(sec) < 2 {
					sec = "0" + sec
				}
				offensiveDirection := "Left"
				if turn.Player1role == 0 || turn.Player1role == 2 || turn.Player1role == 4 {
					offensiveDirection = "Right"
				}
				firstDownYardLine := 0
				if offensiveDirection == "Left" {
					firstDownYardLine = turn.Yardline - turn.Downdistance
				} else {
					firstDownYardLine = turn.Yardline + turn.Downdistance
				}
				sideOfField := "their"
				if turn.Player1id == guid && turn.Yardline < 50 {
					sideOfField = "your"
				} else if turn.Player2id == guid && turn.Yardline >= 50 {
					sideOfField = "your"
				}
				if quarter == 4 {
					quarter = 3
					min = "0"
					sec = "0"
				}
				var animationForTurn []int64
				animationSplit := strings.Split(turn.Animationforturn, ",")
				for _, element := range animationSplit {
					intAnimation, _ := strconv.ParseInt(element, 10, 64)
					animationForTurn = append(animationForTurn, intAnimation)
				}
				turnData := jsonOutputs.TurnInGame{
					Id:                     turn.Id,
					Player1id:              turn.Player1id,
					Player2id:              turn.Player2id,
					Previousturn:           turn.Previousturn,
					Yardline:               turn.Yardline,
					Down:                   turn.Down,
					Togo:                   turn.Downdistance,
					Player1playselected:    turn.Player1playselected,
					Player2playselected:    turn.Player2playselected,
					Player1role:            turn.Player1role,
					Player2role:            turn.Player2role,
					Results:                turn.Results,
					Playtime:               turn.Playtime,
					Timeelapsedingame:      turn.Timeelapsedingame,
					Currentplayer1score:    turn.Currentplayer1score,
					Currentplayer2score:    turn.Currentplayer2score,
					Player1playdatetime:    turn.Player1playdatetime,
					Player2playdatetime:    turn.Player2playdatetime,
					Player1allowedplays:    player1AllowedPlays,
					Player2allowedplays:    player2AllowedPlays,
					Ballon:                 lineOfScrimage,
					Quarter:                quarter + 1,
					Timeremaininginquarter: min + ":" + sec,
					Offensivedirection:     offensiveDirection,
					Player1timeoutsleft:    3,
					Player2timeoutsleft:    3,
					Firstdownyardline:      firstDownYardLine,
					SideOfField:            sideOfField,
					Player1playflipped:     turn.Player1playflipped,
					Player2playflipped:     turn.Player2playflipped,
					Animationforturn:       animationForTurn}
				turnDataList = append(turnDataList, turnData)
			}
			timeDifference := time.Now().Unix() - gameList[0].Lastupdated
			timeBehind := ""
			fmt.Println(timeDifference)
			if timeDifference/86400 > 0 {
				days := strconv.FormatInt(timeDifference/86400, 10)
				timeBehind = days + " Days ago"
			} else if timeDifference/2400 > 0 {
				hours := strconv.FormatInt(timeDifference/2400, 10)
				timeBehind = hours + " Hours ago"
			} else if timeDifference/60 > 0 {
				minutes := strconv.FormatInt(timeDifference/60, 10)
				timeBehind = minutes + " Minutes ago"
			} else {
				seconds := strconv.FormatInt(timeDifference, 10)
				timeBehind = seconds + " Seconds ago"
			}
			winner := ""
			if gameList[0].Outcome != "" {
				// player1Winner := "Results: " + gameInArray.Player1teamname + " won"
				player2Winner := "Results: " + gameList[0].Player2teamname + " won"
				player1Resigned := "Retired:" + gameList[0].Player1teamname + " Lost."
				// player2Resigned := "Retired:" + gameInArray.Player2teamname + " Lost."
				fmt.Println(turnDataList[0].Currentplayer1score)
				fmt.Println(turnDataList[0].Currentplayer2score)
				if strings.Contains(gameList[0].Outcome, "Tie") {
					winner = "3"
				} else if player1Resigned == gameList[0].Outcome || player2Winner == gameList[0].Outcome {
					winner = "2"
				} else {
					winner = "1"
				}
			}
			gameInList := jsonOutputs.GameInList{
				Inviteaccepted:  gameList[0].Inviteaccepted,
				Gameid:          gameList[0].Id,
				Player1:         gameList[0].Player1,
				Player2:         gameList[0].Player2,
				Player1teamname: gameList[0].Player1teamname,
				Player2teamname: gameList[0].Player2teamname,
				Outcome:         gameList[0].Outcome,
				Lastupdated:     gameList[0].Lastupdated,
				Waitingtime:     timeBehind,
				Turns:           turnDataList,
				Player1prefix:   gameList[0].Player1prefix,
				Player2prefix:   gameList[0].Player2prefix,
				Winner:          winner}
			return "true", "Successfully gathered turn data.", gameInList, 200
		}
	}
}

func ResignGame(guid, gameid string) (string, string, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	query := "select * from game where game.Id=? limit 1"
	var gameList []*Game
	_, gameErr := dbmap.Select(&gameList, query, gameid)
	if isDbQueryError(gameErr) {
		return "false", "Failed to update game.", 501
	}
	if len(gameList) > 0 {
		if gameList[0].Player1 == guid || gameList[0].Player2 == guid {
			update := "update game set game.Outcome=? where game.Id=?"
			teamNameReplace := ""
			if gameList[0].Player1 == guid {
				teamNameReplace = gameList[0].Player1teamname
			} else {
				teamNameReplace = gameList[0].Player2teamname
			}
			_, updateErr := dbmap.Exec(update, "Retired:"+teamNameReplace+" Lost.", gameid)
			if isDbQueryError(updateErr) {
				return "false", "Failed to update game.", 501
			} else {
				return "true", "Game was successfully resigned.", 200
			}
		} else {
			return "false", "Invalid game id.", 411
		}
	} else {
		return "false", "Could not find game id.", 406
	}
}

func ConfirmGame(guid, gameid, teamname string) (string, string, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	query := "update game set game.Inviteaccepted=1,game.Player2accepteddate=?,game.Player2teamname=? where game.Id=? and game.Player2=?"
	_, updateErr := dbmap.Exec(query, time.Now().Unix(), teamname, gameid, guid)
	if isDbQueryError(updateErr) {
		return "false", "Failed to confirm game.", 501
	} else {
		return "true", "Game confirmed.", 200
	}
}

func SelectPlay(guid, turnid, playid, flip string) (string, string, int, []int64) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	query := "select * from turn where turn.Id=? limit 1"
	var turns []*Turn

	_, selectTurnError := dbmap.Select(&turns, query, turnid)

	updateOccurred := false
	updatedPlayer := 0
	if isDbQueryError(selectTurnError) {
		return "false", "Error has occured during turn update.", 501, nil
	}
	if len(turns) > 0 {
		gameQuery := "select * from game where game.Id=? limit 1"
		var game []*Game
		_, gameSelectError := dbmap.Select(&game, gameQuery, turns[0].Gameid)
		fmt.Println(gameSelectError)
		if gameSelectError == nil && game[0].Outcome == "" {

			fmt.Println(turns[0])
			turn := turns[0]
			// Need to check if play select is correct possession for player1
			playIdInt, playIdConversionError := strconv.ParseInt(playid, 10, 64)
			if playIdConversionError != nil {
				return "false", "Play ID is invalid", 410, nil
			}
			currentPlay, playExists := plays[playIdInt]
			if turn.Player1id == guid && turn.Player1playselected == -1 {
				if playExists {
					playPossession := currentPlay.Possession
					playerRole := turn.Player1role
					playAllowed := false
					if playPossession == playerRole {
						playAllowed = true
					} else {
						if (playPossession == 0 || playPossession == 2) && (playerRole == 0 || playerRole == 2) {
							playAllowed = true
						} else {
							playAllowed = (playPossession == 1 || playPossession == 3) && (playerRole == 1 || playerRole == 3)
						}
					}
					fmt.Println("Check play")
					fmt.Println(playAllowed)
					if playAllowed {
						updatedPlayer = 1
						turn.Player1playselected = playIdInt
						turn.Player1playdatetime = time.Now().Unix()
						updateOccurred = true
					} else {
						return "false", "Play is not viable for Possession.", 411, nil
					}
				} else {
					return "false", "Play does not exist", 410, nil
				}
			}
			if turn.Player2id == guid && turn.Player2playselected == -1 {
				if playExists {
					playPossession := currentPlay.Possession
					playerRole := turn.Player2role
					playAllowed := false
					if playPossession == playerRole {
						playAllowed = true
					} else {
						if (playPossession == 0 || playPossession == 2) && (playerRole == 0 || playerRole == 2) {
							playAllowed = true
						} else {
							playAllowed = (playPossession == 1 || playPossession == 3) && (playerRole == 1 || playerRole == 3)
						}
					}
					fmt.Println("checkplay")
					fmt.Println(playAllowed)
					if playAllowed {
						updatedPlayer = 2
						turn.Player2playselected = playIdInt
						turn.Player2playdatetime = time.Now().Unix()
						updateOccurred = true
					} else {
						return "false", "Play is not viable for Possession.", 411, nil
					}
				} else {
					return "false", "Play does not exist.", 410, nil
				}
			}
			if turn.Player1playselected != -1 && turn.Player2playselected == -1 {
				player2 := turn.Player2id

				//push notification Android
				go sendAndroidMessage(player2, "It is your turn to call a play.", "")

				//push notification IOS
				go sendIosMessage(player2, "It is your turn to call a play.", "", "")
			}
			if turn.Player2playselected != -1 && turn.Player1playselected == -1 {
				player1 := turn.Player1id
				//push notification Android
				go sendAndroidMessage(player1, "It is your turn to call a play.", "")

				//push notification IOS
				go sendIosMessage(player1, "It is your turn to call a play.", "", "")
			}
			if turn.Player2playselected != -1 && turn.Player1playselected != -1 && updateOccurred {
				fmt.Println("Result needs to be computed")
				result, animation, completedGame := GatherResults(turn, guid)
				result = strings.Replace(result, "player1", game[0].Player1teamname, -1)
				result = strings.Replace(result, "player2", game[0].Player2teamname, -1)
				completedGame = strings.Replace(completedGame, "player1", game[0].Player1teamname, -1)
				completedGame = strings.Replace(completedGame, "player2", game[0].Player2teamname, -1)
				var animationString []string
				for _, element := range animation {
					elementString := strconv.FormatInt(element, 10)
					animationString = append(animationString, elementString)
				}
				// strings.Replace(result, "player1", 	, 0)
				if updatedPlayer == 1 {
					turnUpdate := "Update turn set turn.Player1playselected=?, turn.Player1playdatetime=?,turn.Playtime=?,turn.Results=?,turn.Player1playflipped=?,turn.Animationforturn=? where turn.Id=?"
					_, err := dbmap.Exec(turnUpdate, playid, time.Now().Unix(), turn.Playtime, result, flip, strings.Join(animationString, ","), turnid)
					if isDbQueryError(err) {
						return "false", "Error has occured during turn update.", 501, nil
					} else {
						gameUpdateQuery := "Update game set game.Lastupdated=?,game.Outcome=? where game.Id=?"
						_, updateError := dbmap.Exec(gameUpdateQuery, time.Now().Unix(), completedGame, turn.Gameid)
						isDbQueryError(updateError)
					}
				} else {
					if updatedPlayer == 2 {
						turnUpdate := "Update turn set turn.Player2playselected=?, turn.Player2playdatetime=?,turn.Playtime=?,turn.Results=?,turn.Player2playflipped=?,turn.Animationforturn=? where turn.Id=?"
						_, err := dbmap.Exec(turnUpdate, playid, time.Now().Unix(), turn.Playtime, result, flip, strings.Join(animationString, ","), turnid)
						if isDbQueryError(err) {
							return "false", "Error has occured during turn update.", 501, nil
						} else {
							gameUpdateQuery := "Update game set game.Lastupdated=?,game.Outcome=? where game.Id=?"
							_, updateError := dbmap.Exec(gameUpdateQuery, time.Now().Unix(), completedGame, turn.Gameid)
							isDbQueryError(updateError)
						}
					}
				}
				return "true", result, 200, animation
			} else {
				if updatedPlayer != 0 {
					if updatedPlayer == 1 {
						turnUpdate := "Update turn set turn.Player1playselected=?, turn.Player1playdatetime=?,turn.Player1playflipped=? where turn.Id=?"
						_, err := dbmap.Exec(turnUpdate, playid, time.Now().Unix(), flip, turnid)
						if isDbQueryError(err) {
							return "false", "Error has occured during turn update.", 501, nil
						} else {
							gameUpdateQuery := "Update game set game.Lastupdated=? where game.Id=?"
							_, updateError := dbmap.Exec(gameUpdateQuery, time.Now().Unix(), turn.Gameid)
							isDbQueryError(updateError)
						}
					} else {
						if updatedPlayer == 2 {
							turnUpdate := "Update turn set turn.Player2playselected=?, turn.Player2playdatetime=?,turn.Player2playflipped=? where turn.Id=?"
							_, err := dbmap.Exec(turnUpdate, playid, time.Now().Unix(), flip, turnid)
							if isDbQueryError(err) {
								return "false", "Error has occured during turn update.", 501, nil
							} else {
								gameUpdateQuery := "Update game set game.Lastupdated=? where game.Id=?"
								_, updateError := dbmap.Exec(gameUpdateQuery, time.Now().Unix(), turn.Gameid)
								isDbQueryError(updateError)
							}
						}
					}
				}
				return "true", "Awaiting oppositions play selection.", 200, nil
			}
		} else {
			return "false", "Game has already been completed.", 420, nil
		}
	}

	return "false", "Turn ID is not valid.", 413, nil
}
func sendAndroidMessage(guid, message, data string) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	androidQuery := "select * from user where user.Guid = ?"
	adeviceQuery := "select * from androiddevice where androiddevice.Id = ?"
	var androiduser []*User
	var androiddevice []*Androiddevice
	_, androidSelectError := dbmap.Select(&androiduser, androidQuery, guid)
	if isDbQueryError(androidSelectError) || len(androiduser) <= 0 {
		return
	}
	androiddeviceid := strings.Split(androiduser[0].Androiddevices, ",")
	anumberdevices := strings.Count(androiduser[0].Androiddevices, ",")
	for i := 1; i < (1 + anumberdevices); i++ {
		_, androiddeviceSelectError := dbmap.Select(&androiddevice, adeviceQuery, androiddeviceid[i])
		if !isDbQueryError(androiddeviceSelectError) && len(androiddevice) > 0 {
			SendGCMMessage(androiddevice[0].Deviceid, message, data)
		}
	}
}
func sendIosMessage(guid, message, key, value string) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	fmt.Println(message)
	iosQuery := "select * from user where user.Guid = ?"
	ideviceQuery := "select * from iosdevice where iosdevice.Id = ?"
	var iosuser []*User
	var iosdevice []*Iosdevice
	_, iosSelectError := dbmap.Select(&iosuser, iosQuery, guid)
	fmt.Println("IOS check ")
	fmt.Println(len(iosuser))
	if isDbQueryError(iosSelectError) || len(iosuser) <= 0 {
		return
	}
	iosdeviceid := strings.Split(iosuser[0].Iosdevices, ",")
	inumberdevices := strings.Count(iosuser[0].Iosdevices, ",")
	for i := 1; i < (1 + inumberdevices); i++ {
		_, iosdeviceSelectError := dbmap.Select(&iosdevice, ideviceQuery, iosdeviceid[i])
		if !isDbQueryError(iosdeviceSelectError) && len(iosdevice) > 0 {
			SendAPNSMessageInner(iosdevice[0].Devicetoken, message, key, value)
		}
	}
}
func SendAPNSMessageInner(tokenInput, message, key, value string) {
	// SendAPNSMessage(iosdevice[0].Devicetoken, "It is your turn to select a play")
	token := stripchars(tokenInput, "< >")
	fmt.Println(token)
	requestUrl := "http://localhost:9090/APNS"
	values := make(url.Values)
	values.Set("message", message)
	values.Set("devicetoken", token)
	values.Set("key", key)
	values.Set("value", value)

	response, err := http.PostForm(requestUrl, values)
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		defer response.Body.Close()
		message, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("error", err)
		}
		type Device struct {
			id          string
			devicetoken string
		}
		fmt.Println(string(message))
	}
}
func getOddsTableResults(offensivePlayId, defensivePlayId string) (map[string]interface{}, bool) {
	values := make(url.Values)
	values.Set("offensivePlayId", offensivePlayId)
	values.Set("defensivePlayId", defensivePlayId)
	client, requestError := http.PostForm("http://localhost:8082/odds", values)
	if requestError != nil {
		go amzses.SendMail("GameOn@PlayVsFootball.com", "dclawson@engagemobile.com", "Vsfootball Error", "Error has occured during odds table lookup for "+offensivePlayId+" "+defensivePlayId+"\n"+"Node Name: "+nodeName)
		go amzses.SendMail("GameOn@PlayVsFootball.com", "ravi@engagemobile.com", "Vsfootball Error", "Error has occured during odds table lookup for "+offensivePlayId+" "+defensivePlayId+"\n"+"Node Name: "+nodeName)
		go amzses.SendMail("GameOn@PlayVsFootball.com", "stimperley@engagemobile.com", "Vsfootball Error", "Error has occured during odds table lookup for "+offensivePlayId+" "+defensivePlayId+"\n"+"Node Name: "+nodeName)
		cmd := exec.Command("/home/ec2-user/oddstable/odds", "/home/ec2-user/oddstable/tables", "&")
		go cmd.Output()
		defer client.Body.Close()
		return nil, false
	} else {
		defer client.Body.Close()
	}

	body, _ := ioutil.ReadAll(client.Body)
	fmt.Println(string(body))
	var result interface{}
	json.Unmarshal(body, &result)
	return result.(map[string]interface{}), true
	// return result

}
func GatherResults(turn *Turn, guid string) (string, []int64, string) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	playtime := r.Intn(8) + 60

	gainPlayer1 := true
	if turn.Player2role == 0 || turn.Player2role == 4 {
		gainPlayer1 = false
	}

	outcomeAmount := r.Intn(40)
	var oddsTableResults map[string]interface{}
	positive := 1
	player1Play := strconv.FormatInt(turn.Player1playselected, 10)
	player2Play := strconv.FormatInt(turn.Player2playselected, 10)
	success := false
	changeofpossession := false
	yardLine := turn.Yardline
	down := turn.Down
	downdistance := turn.Downdistance
	var results string
	var animation []int64
	var completedgame string
	player1role := turn.Player1role
	player2role := turn.Player2role
	player1score := turn.Currentplayer1score
	player2score := turn.Currentplayer2score
	var newTimeElapsed int
	var outputyards int

	if gainPlayer1 {
		oddsTableResults, success = getOddsTableResults(player1Play, player2Play)
	} else {
		oddsTableResults, success = getOddsTableResults(player2Play, player1Play)
	}

	if success && oddsTableResults["status"].(float64) == 0 {
		outcomeAmount, positive = getYardageAndPositive(oddsTableResults)
		if oddsTableResults["isChangeOfPossession"].(float64) == 1 {
			changeofpossession = true
		}
	}
	playname := oddsTableResults["playName"].(string)

	if gainPlayer1 {
		if turn.Player1role == 4 {
			playtime, player1role, player2role, player1score, player2score, down, downdistance, outputyards, animation, results, newTimeElapsed, completedgame = GetKickoff(turn, guid, gainPlayer1, changeofpossession, outcomeAmount, playname)
		} else if plays[turn.Player1playselected].Type == "Run" || turn.Player1playselected == 16 {
			playtime, player1role, player2role, player1score, player2score, down, downdistance, outputyards, animation, results, newTimeElapsed, completedgame = GetRun(turn, guid, gainPlayer1, changeofpossession, playname, outcomeAmount, positive, yardLine, down, downdistance)
		} else if plays[turn.Player1playselected].Type == "Pass" {
			playtime, player1role, player2role, player1score, player2score, down, downdistance, outputyards, animation, results, newTimeElapsed, completedgame = GetPass(turn, guid, gainPlayer1, changeofpossession, playname, outcomeAmount, positive, yardLine, down, downdistance)
		} else if turn.Player1playselected == 15 {
			playtime, player1role, player2role, player1score, player2score, down, downdistance, outputyards, animation, results, newTimeElapsed, completedgame = GetPunt(turn, guid, gainPlayer1, changeofpossession, outcomeAmount, positive, yardLine)
		} else if turn.Player1playselected == 42 {
			message := "failed"
			var fieldgoalcheck bool
			fieldgoalcheck, message = FieldGoal(yardLine, gainPlayer1)
			playtime = r.Intn(8) + 1
			newTimeElapsed = turn.Timeelapsedingame + playtime
			if fieldgoalcheck {
				player1role, player2role, player1score, player2score = scoring(player1role, player2role, player1score, player2score, 3)
				down = 0
				downdistance = 0
				outputyards = 0
				results = message
			} else {
				player1role, player2role = cop(player1role, player2role)
				outputyards = yardLine - 7
				results = message
				down = 1
				downdistance = 10
			}
		}
	} else {
		if turn.Player2role == 4 {
			playtime, player1role, player2role, player1score, player2score, down, downdistance, outputyards, animation, results, newTimeElapsed, completedgame = GetKickoff(turn, guid, gainPlayer1, changeofpossession, outcomeAmount, playname)
		} else if plays[turn.Player2playselected].Type == "Run" || turn.Player2playselected == 16 {
			playtime, player1role, player2role, player1score, player2score, down, downdistance, outputyards, animation, results, newTimeElapsed, completedgame = GetRun(turn, guid, gainPlayer1, changeofpossession, playname, outcomeAmount, positive, yardLine, down, downdistance)
		} else if plays[turn.Player2playselected].Type == "Pass" {
			playtime, player1role, player2role, player1score, player2score, down, downdistance, outputyards, animation, results, newTimeElapsed, completedgame = GetPass(turn, guid, gainPlayer1, changeofpossession, playname, outcomeAmount, positive, yardLine, down, downdistance)
		} else if turn.Player2playselected == 15 {
			playtime, player1role, player2role, player1score, player2score, down, downdistance, outputyards, animation, results, newTimeElapsed, completedgame = GetPunt(turn, guid, gainPlayer1, changeofpossession, outcomeAmount, positive, yardLine)
		} else if turn.Player2playselected == 42 {
			message := "failed"
			var fieldgoalcheck bool
			playtime = r.Intn(8) + 1
			newTimeElapsed = turn.Timeelapsedingame + playtime
			fieldgoalcheck, message = FieldGoal(yardLine, gainPlayer1)
			if fieldgoalcheck {
				player1role, player2role, player1score, player2score = scoring(player1role, player2role, player1score, player2score, 3)
				results = message
				down = 0
				downdistance = 0
				outputyards = 0
			} else {
				player1role, player2role = cop(player1role, player2role)
				outputyards = yardLine + 7
				results = message
				down = 1
				downdistance = 10
			}
		}
	}
	newTurn := &Turn{
		Gameid:              turn.Gameid,
		Player1id:           turn.Player1id,
		Player2id:           turn.Player2id,
		Previousturn:        turn.Id,
		Yardline:            outputyards,
		Down:                down,
		Downdistance:        downdistance,
		Player1playselected: -1,
		Player2playselected: -1,
		Player1role:         player1role,
		Player2role:         player2role,
		Results:             "",
		Playtime:            0,
		Timeelapsedingame:   newTimeElapsed,
		Currentplayer1score: player1score,
		Currentplayer2score: player2score,
		Player1playdatetime: -1,
		Player2playdatetime: -1}

	isDbQueryError(dbmap.Insert(newTurn))
	turn.Playtime = playtime
	turn.Results = results
	return results, animation, completedgame
}

func getYardageAndPositive(results map[string]interface{}) (int, int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fmt.Println(strings.TrimSpace(results["yardageName"].(string)))
	fmt.Println(len(strings.TrimSpace(results["yardageName"].(string))))
	switch strings.TrimSpace(results["yardageName"].(string)) {
	case "Gain remaining yds":
		return 100, 1
	case "Gain of 61 yds or more":
		return r.Intn(40) + 61, 1
	case "Gain of 41 - 60 yds":
		return r.Intn(20) + 41, 1
	case "Gain of 21 - 40 yds":
		return r.Intn(20) + 21, 1
	case "Gain of 11 - 20 yds":
		return r.Intn(10) + 11, 1
	case "Gain of 6 - 10 yds":
		return r.Intn(4) + 6, 1
	case "Gain of 1 - 5 yds":
		return r.Intn(5) + 1, 1
	case "No Gain":
		return 0, 1
	case "Loss of 1 - 5 yds":
		return r.Intn(5) + 1, -1
	case "Loss of 6 - 10 yds":
		return r.Intn(5) + 6, -1
	case "Loss of 11 - 15 yds":
		return r.Intn(5) + 11, -1
	case "Loss of 16 yds or more":
		return r.Intn(85) + 16, -1
	default:
		fmt.Println("something messed up.")
	}
	return 0, 1
}

func getRandResultString(runResult string, gain int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var result string
	switch runResult {
	case "Run - Tackle":
		runTackled := r.Intn(3)
		switch runTackled {
		case 0:
			if gain == 0 {
				result = "No Gain!"
				return result
			}
			if gain < 0 {
				result = fmt.Sprintf("Loss of %v \n", gain*-1)
				return result
			}
			if gain > 0 {
				result = fmt.Sprintf("Gain of %v \n", gain)
				return result
			}
		case 1:
			if gain == 0 {
				result = "No Gain!"
				return result
			}
			if gain < 0 {
				result = fmt.Sprintf("Loss of %v on the play.\n", gain*-1)
				return result
			}
			if gain > 0 {
				result = fmt.Sprintf("Gain of %v on the play.\n", gain)
				return result
			}
		case 2:
			if gain == 0 {
				result = "No Gain!"
				return result
			}
			if gain < 0 {
				result = fmt.Sprintf("%v-yard loss on the play.\n", gain*-1)
				return result
			}
			if gain > 0 {
				result = fmt.Sprintf("%v-yard gain on the play.\n", gain)
				return result
			}
		}
	case "Run - Fumble, Kept":
		runFumbleKept := r.Intn(2)
		switch runFumbleKept {
		case 0:
			return "Fumble on the play. Recovered by the offense.\n"
		case 1:
			return "Fumble! Recovered by the runner.\n"
		}
	case "Run - Fumble, Lost - COP":
		runFumbleLost := r.Intn(2)
		switch runFumbleLost {
		case 0:
			return "There's a fumble on the play, recovered by the defense.\n"
		case 1:
			return "Fumble!\n"
		}
	case "Pass - Incomplete":
		passIncomplete := r.Intn(3)
		switch passIncomplete {
		case 0:
			return "He couldn't pull it in. Pass is incomplete.\n"
		case 1:
			return "Right through his hands. Pass incomplete.\n"
		case 2:
			return "The QB overthrew him. Pass is incomplete.\n"
		}
	case "Pass - Complete, Tackle":
		passCompleteTackle := r.Intn(3)
		switch passCompleteTackle {
		case 0:
			if gain == 0 {
				result = "No Gain!"
				return result
			}
			if gain < 0 {
				result = fmt.Sprintf("Loss of %v on the play.\n", gain*-1)
				return result
			}
			if gain > 0 {
				if gain > 15 {
					result = "Beautiful catch.\n"
				}
				result += fmt.Sprintf("Gain of %v on the play.\n", gain)
				return result
			}
		case 1:
			if gain == 0 {
				result = "No Gain!"
				return result
			}
			if gain < 0 {
				result = fmt.Sprintf("Loss of %v on the play.\n", gain*-1)
				return result
			}
			if gain > 0 {
				if gain > 15 {
					result = "What a catch!\n"
				}
				result += fmt.Sprintf("Gain of %v on the play.\n", gain)
				return result
			}
		case 2:
			if gain == 0 {
				result = "No Gain!"
				return result
			}
			if gain < 0 {
				result = fmt.Sprintf("Loss of %v on the play.\n", gain*-1)
				return result
			}
			if gain > 0 {
				if gain > 15 {
					result = "Great pitch and catch.\n"
				}
				result += fmt.Sprintf("Gain of %v on the play.\n", gain)
				return result
			}
		}
	case "Pass - QB Scramble, Tackle":
		qbScramble := r.Intn(2)
		if gain > 0 {
			result = "QB scramble!\n"
			result += fmt.Sprintf("Gain of %v on the play.\n", gain)
			return result
		}
		switch qbScramble {
		case 0:
			if gain == 0 {
				result = "QB sack!\n"
				result += "No Gain!"
				return result
			} else {
				result = "QB sack!\n"
				result += fmt.Sprintf("Loss of %v on the play.\n", gain*-1)
				return result
			}
		case 1:
			if gain == 0 {
				result = "QB brought down behind the line.\n"
				result += "No Gain!"
				return result
			} else {
				result = "QB brought down behind the line.\n"
				result += fmt.Sprintf("Loss of %v on the play.\n", gain*-1)
				return result
			}
		}
	case "Pass - Fumble, Kept":
		passFumbleKept := r.Intn(2)
		switch passFumbleKept {
		case 0:
			return "Great catch, then a fumble on the play.  Receiver fell on it.\n"
		case 1:
			return "Caught but loose ball. Offense recovers.\n"
		}
	case "Pass - Fumble, Lost - COP":
		passFumbleLoss := r.Intn(2)
		switch passFumbleLoss {
		case 0:
			return "Fumble!\n"
		case 1:
			return "Big hit and a fumble. Defense recovers.\n"
		}
	case "Pass - Interception, Tackle - COP":
		passInterceptionTackle := r.Intn(3)
		switch passInterceptionTackle {
		case 0:
			return "Interception!\n"
		case 1:
			return "Oskie!\n"
		case 2:
			return "I-N-T on the play.\n"
		}
	}

	return ""
}

func GetPlaybookIds() (string, string, int, []int64) {
	return "true", "Successfully gather play ids.", 200, playIds
}
func GetPlayForPlayId(playId string) (string, string, int, []jsonOutputs.PlayInPlaybook) {
	intId, _ := strconv.ParseInt(playId, 10, 64)
	play, found := plays[intId]
	if found {
		image := ImageLocation + play.Image
		if play.Image == "" {
			image = ""
		}
		description := ImageLocation + play.Description
		if play.Description == "" {
			description = ""
		}
		var playReturn []jsonOutputs.PlayInPlaybook
		playReturn = append(playReturn, jsonOutputs.PlayInPlaybook{
			Id:          play.Id,
			Title:       play.Title,
			Image:       image,
			Type:        play.Type,
			Possession:  play.Possession,
			Premium:     play.Premium,
			Price:       play.Price,
			Productid:   play.Productid,
			Producttype: play.Producttype,
			Version:     play.Version,
			Filesize:    play.Filesize,
			Canflip:     play.Canflip,
			Description: description})
		return "true", "Succesfully gathered play information", 200, playReturn
	} else {
		return "false", "Play was not found.", 410, nil
	}

}
func GetOwnedPlays(guid string) (string, string, int, string) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	query := "select user.Playsowned from user where user.Guid=? limit 1"
	var playIds []*User
	_, err := dbmap.Select(&playIds, query, guid)
	if isDbQueryError(err) {
		return "false", "Database failure", 501, ""
	} else {
		if len(playIds) > 0 {
			return "true", "Successfully gathered", 200, playIds[0].Playsowned
		} else {
			return "false", "User was not found.", 401, ""
		}

	}
}

func PurchaseItem(guid, playid, receipt, platform string) (string, string, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	query := "select * from user where user.Guid=? limit 1"
	var user []*User
	_, err := dbmap.Select(&user, query, guid)
	playIdInt, intConversionError := strconv.ParseInt(playid, 10, 64)
	if intConversionError != nil {
		return "false", "Playid is invalid", 405
	}
	if playIdInt < 0 {
		return "false", "Playid is invalid", 405
	}
	if isDbQueryError(err) {
		return "false", "Database error", 501
	}
	if len(user) > 0 {
		currentPlaysOwned := user[0].Playsowned
		playIdStrings := strings.Split(currentPlaysOwned, ",")
		playIdStrings = append(playIdStrings, playid)
		updateQuery := "update user set user.Playsowned=? where user.Guid=?"
		_, updateFailure := dbmap.Exec(updateQuery, strings.Join(playIdStrings, ","), guid)
		if isDbQueryError(updateFailure) {
			return "false", "Database failure.", 501
		} else {
			playIdInt, _ := strconv.ParseInt(playid, 10, 64)
			purchase := &Storepurchase{
				Userid:    guid,
				Productid: plays[playIdInt].Productid,
				Playid:    playIdInt,
				Receipt:   receipt,
				Os:        platform}
			insertError := dbmap.Insert(purchase)
			if isDbQueryError(insertError) {
				return "false", "Database failure", 501
			}
			return "true", "Successfully acquired play.", 200
		}
	} else {
		return "false", "Database failure.", 501
	}
}
func GetAnimationIds() (string, string, int, []int64) {
	return "true", "Successfully gather animation ids.", 200, animationIds
}
func GetAnimationForAnimationId(id string) (string, string, int, jsonOutputs.Animation) {
	animationId, _ := strconv.ParseInt(id, 10, 64)
	animation, found := animationsById[animationId]
	if found {
		image := ImageLocation + animation.Image
		if animation.Image == "" {
			image = ""
		}
		animationReturn := jsonOutputs.Animation{
			Id:       animation.Id,
			Image:    image,
			Tier:     animation.Tier,
			Duration: animation.Frames / animation.Framespersecond,
			Version:  animation.Version,
			Filesize: animation.Filesize}

		return "true", "Succesfully gathered animation information", 200, animationReturn
	} else {
		return "false", "Animation not found.", 415, jsonOutputs.Animation{}
	}

}
func RegisterAndroidDevice(guid, registrationid string) (string, string, int) {
	var register = DeviceRegister{Register: &AndroidDeviceRegister{}}
	return register.RegisteDevice(guid, registrationid)
}

func RegisterIosDevice(guid, deviceToken string) (string, string, int) {
	var register = DeviceRegister{Register: &IosDeviceRegister{}}
	return register.RegisteDevice(guid, deviceToken)
}

// func GetProductsList(guid string) (string, string, []jsonOutputs.ProductInList, int) {

// }
// func Rematch(guid, gameId string) (string, string, int) {

// }
func SendGCMMessage(deviceId, message, data string) {
	var regIds []string
	regIds = append(regIds, deviceId)
	payload := map[string]string{"message": message, "data": data}
	msg := gcm.NewMessage(payload, regIds...)
	//set GCM API key
	sender := gcm.New("AIzaSyD0PKSQH_4V9PYpkrIpBW_D7PK6Af5ELOY")

	//	 fmt.Println(sender)

	//	 fmt.Println(msg)

	// Send the message and receive the response after at most two retries.
	response, err := sender.Send(msg, 2)
	if err != nil {
		fmt.Println("Failed to send message: " + err.Error())
		return
	} else {
		fmt.Println(response)
	}
}

// func SendAPNSMessage(deviceToken, message string) {
// 	certificate := pem
// 	key := "Engage!" //keys
// 	token := stripchars(deviceToken, "< >")

// 	endpoint := "gateway.sandbox.push.apple.com:2195"
// 	Client, err := apns.NewClient(endpoint, certificate, key)
// 	//preparing binary payload from Json format

// 	fmt.Println(err)

// 	data := map[string]string{"It is your turn to select a play": message}
// 	bdata, err := json.Marshal(data)

// 	response := Client.SendPayloadString(token, bdata, 1)

// 	if err != nil {
// 		fmt.Println("Failed to send message: " + err.Error())
// 		return
// 	} else {
// 		fmt.Println(response)
// 	}
// }
func stripchars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}
func TurnFailCheck() {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	turncheckquery := "select * from turn where turn.Player1playselected > ? AND Player2playselected > ? AND Results = ?"
	var turncheck []*Turn
	_, turncheckError := dbmap.Select(&turncheck, turncheckquery, 0, 0, "")
	if turncheckError != nil {
		fmt.Println("Turn Check Failed")
	} else {
		updateturnQuery := "update turn set turn.Player2playselected = ? where turn.Id=?"
		_, err := dbmap.Exec(updateturnQuery, turncheck[0], -1, turncheck[0].Id)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Turn Check Successful")
		}
	}
}

func FieldGoal(yardline int, gainPlayer1 bool) (bool, string) {
	var success int
	var failleft int
	var failright int
	var blocked int
	var failshort int
	var distance int

	if gainPlayer1 {
		distance = 100 - yardline + 7
	} else {
		distance = yardline + 7
	}
	// fmt.Println("distance = ", distance)
	if distance <= 10 {
		success = 98
		failleft = success + 1
		failright = failleft + 1
		failshort = failright
		blocked = failshort
	}
	if distance <= 20 && 10 < distance {
		success = 95
		failleft = success + 2
		failright = failleft + 2
		failshort = failright
		blocked = failshort + 1
	}
	if distance <= 30 && 20 < distance {
		success = 92
		failleft = success + 3
		failright = failleft + 3
		failshort = failright
		blocked = failshort + 2
	}
	if distance <= 40 && 30 < distance {
		success = 90
		failleft = success + 4
		failright = failleft + 4
		failshort = failright
		blocked = failshort + 2
	}
	if distance <= 50 && 40 < distance {
		success = 84
		failleft = success + 6
		failright = failleft + 6
		failshort = failright + 1
		blocked = failshort + 3
	}
	if distance <= 55 && 50 < distance {
		success = 50
		failleft = success + 20
		failright = failleft + 20
		failshort = failright + 5
		blocked = failshort + 5
	}
	if distance <= 60 && 55 < distance {
		success = 40
		failleft = success + 15
		failright = failleft + 15
		failshort = failright + 20
		blocked = failshort + 10
	}
	if distance <= 63 && 60 < distance {
		success = 20
		failleft = success + 10
		failright = failleft + 10
		failshort = failright + 40
		blocked = failshort + 20
	}
	if distance <= 100 && 63 < distance {
		success = 0
		failleft = success + 20
		failright = failleft + 20
		failshort = failright + 50
		blocked = failshort + 10
	}
	// fmt.Println("success = ", success)
	// fmt.Println("failleft = ", failleft)
	// fmt.Println("failright = ", failright)
	// fmt.Println("failshort = ", failshort)
	// fmt.Println("blocked = ", blocked)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	checkfieldgoal := r.Intn(100)
	// fmt.Println("Random number = ", checkfieldgoal)
	if checkfieldgoal <= success {
		// fmt.Println("Success")
		return true, "The kick is good!\n"
	}
	if checkfieldgoal <= failleft && success < checkfieldgoal {
		return false, "The kick missed left!\nNo good!\n"
		// fmt.Println("Change of Possesion Fail Left")
	}
	if checkfieldgoal <= failright && failleft < checkfieldgoal {
		return false, "The kick missed right!\nNo good!\n"
		// fmt.Println("Change of Possesion Fail Right")
	}
	if checkfieldgoal <= failshort && failright < checkfieldgoal {
		return false, "The kick is short!\nNo good!\n"
		// fmt.Println("Change of Possesion Fail Short")
	}
	if checkfieldgoal <= blocked && failshort < checkfieldgoal {
		return false, "The kick is blocked!\nNo good!\n"
		// fmt.Println("Change of Possesion Blocked")
	}
	return false, "The kick was blocked!\nNo good!\n"
}

func ChangePassword(guid, oldPassword, newPassword string) (string, string, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	user := GetUserByGuid(guid)
	if user.Password == oldPassword {
		updatePassword := "update user set user.Password=? where user.Guid=?"
		_, updateError := dbmap.Exec(updatePassword, newPassword, guid)
		if isDbQueryError(updateError) {
			return "false", "Database error.", 501
		} else {
			return "true", "Password change successful.", 200
		}
	} else {
		return "false", "Old password did not match password on account.", 400
	}
}

func SubmitChat(sender string, gameid string, text string) (string, string, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)

	gameID, err0 := strconv.ParseInt(gameid, 10, 64)
	if err0 != nil {
		fmt.Printf(err0.Error())
	}

	//return success, response, status

	var results []*Game

	query := "select * from game where game.Id = ? limit 1"
	_, err := dbmap.Select(&results, query, gameid)
	if isDbQueryError(err) {
		return "false", "Database error", 501
	}

	if len(results) == 0 {
		return "false", "Invalid userid/gameid combination", 406
	}
	game := results[0]
	if game.Player1 != sender && game.Player2 != sender {
		return "false", "Invalid userid/gameid combination", 406
	}

	var teamName string
	var senderGuid string
	isFirstPlayerSender := false
	if game.Player1 == sender {
		isFirstPlayerSender = true
		teamName = game.Player1teamname
		senderGuid = game.Player1
	} else {
		teamName = game.Player2teamname
		senderGuid = game.Player2
	}
	fmt.Printf(teamName)
	chat := &Chatmessage{
		Sender:         senderGuid,
		Gameid:         gameID,
		Created:        time.Now().Unix(),
		Message:        text,
		Senderteamname: teamName}

	err2 := dbmap.Insert(chat)
	if isDbQueryError(err2) {
		fmt.Printf(err2.Error())
		return "false", "Database error", 501
	}
	data := map[string]string{"gameid": gameid}
	payload, _ := json.Marshal(data)
	if isFirstPlayerSender {
		pushmessage := "New message from " + game.Player1teamname
		go sendAndroidMessage(game.Player2, pushmessage, string(payload))
		go sendIosMessage(game.Player2, pushmessage, "gameid", gameid)
	} else {
		pushmessage := "New message from " + game.Player2teamname
		go sendAndroidMessage(game.Player1, pushmessage, string(payload))
		go sendIosMessage(game.Player1, pushmessage, "gameid", gameid)
	}
	return "true", "Chat message stored", 200
}

func RetrieveChat(gameid string, guid string) (string, string, []jsonOutputs.ChatSummary, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	gameID, err0 := strconv.ParseInt(gameid, 10, 64)
	if err0 != nil {
		fmt.Printf(err0.Error())
	}
	var games []*Game
	var chats []*Chatmessage

	// return success, message, chatList, status

	query0 := "select * from game where game.Id = ? limit 1"
	_, err1 := dbmap.Select(&games, query0, gameID)
	if isDbQueryError(err1) {
		return "false", "Database error", nil, 501
	}

	if len(games) == 0 {
		return "false", "Invalid userid/gameid combination", nil, 406
	}
	game := games[0]
	if game.Player1 != guid && game.Player2 != guid {
		return "false", "Invalid userid/gameid combination", nil, 406
	}
	query1 := "select * from chatmessage where gameid = ? order by Created desc limit 10"
	_, err2 := dbmap.Select(&chats, query1, gameID)
	if isDbQueryError(err2) {
		return "false", "Database error", nil, 501
	}

	// No errors if we got here

	var summaries []jsonOutputs.ChatSummary
	for index := range chats {
		timeDifference := time.Now().Unix() - chats[index].Created
		timeBehind := ""
		fmt.Println(timeDifference)
		if timeDifference/86400 > 0 {
			days := strconv.FormatInt(timeDifference/86400, 10)
			timeBehind = days + " Days ago"
		} else if timeDifference/2400 > 0 {
			hours := strconv.FormatInt(timeDifference/2400, 10)
			timeBehind = hours + " Hours ago"
		} else if timeDifference/60 > 0 {
			minutes := strconv.FormatInt(timeDifference/60, 10)
			timeBehind = minutes + " Minutes ago"
		} else {
			seconds := strconv.FormatInt(timeDifference, 10)
			timeBehind = seconds + " Seconds ago"
		}
		summaries = append(summaries, jsonOutputs.ChatSummary{
			Teamname:   chats[index].Sender,
			Timepassed: timeBehind,
			Text:       chats[index].Message,
			Timestamp:  chats[index].Created})
	}
	return "true", "Chat messages gathered successfully.", summaries, 200
}
func ViralCoefficient() (string, string, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)

	yesterday := time.Now().Unix() - 86400
	usersquery := "select count(*) from user where user.Created > ?"
	users, userserr := dbmap.SelectInt(usersquery, yesterday)
	if isDbQueryError(userserr) {
		return "false", "Database error", 501
	}

	invitequery := "select count(*) from user where user.Created > ? and user.Verified = ?"
	invitedusers, invitederr := dbmap.SelectInt(invitequery, yesterday, 2)
	if isDbQueryError(invitederr) {
		return "false", "Database error", 501
	}

	acceptquery := "select count(*) from user where user.Created > ? and user.Verified = ? and user.Updated != ?"
	acceptedusers, acceptederr := dbmap.SelectInt(acceptquery, yesterday, 2, -1)
	if isDbQueryError(acceptederr) {
		return "false", "Database error", 501
	}

	facebookquery := "select count(*) from user where user.Created > ? and user.Accounttype = ?"
	facebookusers, facebookerr := dbmap.SelectInt(facebookquery, yesterday, "")
	if isDbQueryError(facebookerr) {
		return "false", "Database error", 501
	}
	gamequery := "select count(*) from game where game.Outcome != ?"
	games, gameerr := dbmap.SelectInt(gamequery, "")
	if isDbQueryError(gameerr) {
		return "false", "Database error", 501
	}
	acceptedgamequery := "select count(*) from game where game.Outcome != ? and game.Inviteaccepted = ?"
	acceptedgames, acceptedgameerr := dbmap.SelectInt(acceptedgamequery, "", 1)
	if isDbQueryError(acceptedgameerr) {
		return "false", "Database error", 501
	}
	var inviterate float64
	var acceptrate float64
	var viral float64

	inviterate = (float64(invitedusers) / float64(users))
	acceptrate = (float64(acceptedusers) / float64(invitedusers))
	viral = (inviterate * acceptrate) * 100
	inviterate = inviterate * 100
	acceptrate = acceptrate * 100
	test := fmt.Sprintf("# of New Users: %v \n# of Email Invites: %v \n# of Email Accepts: %v \nEmail Invitation Rate: %.2f \nEmail Acceptance Rate: %.2f \nViral Coefficient: %.2f \n# of Facebook Users: %v \n# of Active Games: %v \n# of Accepted Games: %v", users, invitedusers, acceptedusers, inviterate, acceptrate, viral, facebookusers, games, acceptedgames)
	fmt.Println(test)
	go amzses.SendMail("GameOn@PlayVsFootball.com", "stimperley@engagemobile.com", "Viral Coefficient", fmt.Sprintf("# of New Users: %v \n# of Email Invites: %v \n# of Email Accepts: %v \n Email Invitation Rate: %.2f \n Email Acceptance Rate: %.2f \nViral Coefficient: %.2f \n# of Facebook Users: %v \n# of Active Games: %v \n# of Accepted Games: %v", users, invitedusers, acceptedusers, inviterate, acceptrate, viral, facebookusers, games, acceptedgames))
	go amzses.SendMail("GameOn@PlayVsFootball.com", "dclawson@engagemobile.com", "Viral Coefficient", fmt.Sprintf("# of New Users: %v \n# of Email Invites: %v \n# of Email Accepts: %v \nEmail Invitation Rate: %.2f \nEmail Acceptance Rate: %.2f \nViral Coefficient: %.2f \n# of Facebook Users: %v \n# of Active Games: %v \n# of Accepted Games: %v", users, invitedusers, acceptedusers, inviterate, acceptrate, viral, facebookusers, games, acceptedgames))
	go amzses.SendMail("GameOn@PlayVsFootball.com", "mbarksdale@engagemobile.com", "Viral Coefficient", fmt.Sprintf("# of New Users: %v \n# of Email Invites: %v \n# of Email Accepts: %v \nEmail Invitation Rate: %.2f \nEmail Acceptance Rate: %.2f \nViral Coefficient: %.2f \n# of Facebook Users: %v \n# of Active Games: %v \n# of Accepted Games: %v", users, invitedusers, acceptedusers, inviterate, acceptrate, viral, facebookusers, games, acceptedgames))
	go amzses.SendMail("GameOn@PlayVsFootball.com", "austin37@gmail.com", "Viral Coefficient", fmt.Sprintf("# of New Users: %v \n# of Email Invites: %v \n# of Email Accepts: %v \nEmail Invitation Rate: %.2f \nEmail Acceptance Rate: %.2f \nViral Coefficient: %.2f \n# of Facebook Users: %v \n# of Active Games: %v \n# of Accepted Games: %v", users, invitedusers, acceptedusers, inviterate, acceptrate, viral, facebookusers, games, acceptedgames))
	go amzses.SendMail("GameOn@PlayVsFootball.com", "ksamuelson@playvsfootball.com", "Viral Coefficient", fmt.Sprintf("# of New Users: %v \n# of Email Invites: %v \n# of Email Accepts: %v \nEmail Invitation Rate: %.2f \nEmail Acceptance Rate: %.2f \nViral Coefficient: %.2f \n# of Facebook Users: %v \n# of Active Games: %v \n# of Accepted Games: %v", users, invitedusers, acceptedusers, inviterate, acceptrate, viral, facebookusers, games, acceptedgames))
	return "true", "Viral Coefficient Sent", 200
}
func cop(player1role int, player2role int) (int, int) {
	if player1role == 0 {
		player1role = 1
		player2role = 0
		return player1role, player2role
	} else {
		player1role = 0
		player2role = 1
		return player1role, player2role
	}
	return player1role, player2role

}
func scoring(player1role int, player2role int, player1score int, player2score int, scoreamount int) (int, int, int, int) {
	if player1role == 0 {
		player1role = 4
		player2role = 5
		player1score += scoreamount
		return player1role, player2role, player1score, player2score
	}
	if player2role == 0 {
		player1role = 5
		player2role = 4
		player2score += scoreamount
		return player1role, player2role, player1score, player2score
	}
	return player1role, player2role, player1score, player2score
}

func isDbQueryError(err error) bool {
	if err != nil {
		go amzses.SendMail("GameOn@PlayVsFootball.com", "ravi@engagemobile.com", "DatabaseError", "Database Error: "+err.Error()+"\n"+"Node Name: "+nodeName)
		go amzses.SendMail("GameOn@PlayVsFootball.com", "dclawson@engagemobile.com", "DatabaseError", "Database Error: "+err.Error()+"\n"+"Node Name: "+nodeName)
		return true
	} else {
		return false
	}
}

func DbTest() {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	fmt.Print(dbmap)
	guid, guidErr := uuid.NewV4()
	if guidErr != nil {
		fmt.Errorf(guidErr.Error())
	}

	user := &User{
		Created:         time.Now().Unix(),
		Updated:         time.Now().Unix(),
		Firstname:       "Some",
		Lastname:        "CoolGuy",
		Guid:            strings.Replace(guid.String(), "-", "", -1),
		Password:        "mycoolpw",
		Accesstoken:     "token",
		Accounttype:     "",
		Tokenexpiration: ""}

	dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
	err2 := dbmap.Insert(user)
	// fmt.Print(user)

	if err2 != nil {
		fmt.Printf(err2.Error())
	}

	// fmt.Printf(err2.Error())

	var list []*User
	// obj,err3 := dbmap.Get(Weather{},0)
	_, err3 := dbmap.Select(&list, "select * from user")
	if err3 != nil {
		fmt.Printf(err3.Error())
	}

	for i := range list {
		fmt.Println(list[i])
	}
	// output := obj.(*Weather)

	// fmt.Printf(output.City)

}
func GetKickoff(turn *Turn, guid string, gainplayer1 bool, changeofpossession bool, outcomeamount int, playname string) (int, int, int, int, int, int, int, int, []int64, string, int, string) {
	var playtime int
	var score bool
	var results string
	var animation []int64
	var outputyards int
	var down int
	var downdistance int
	var lineOfScrimmage int
	var completedgame string

	player1role := turn.Player1role
	player2role := turn.Player2role
	player1playselected := turn.Player1playselected
	player2playselected := turn.Player2playselected
	player1score := turn.Currentplayer1score
	player2score := turn.Currentplayer2score

	score = false
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	playtime = r.Intn(8) + 1

	turn.Playtime = playtime
	newTimeElapsed := playtime + turn.Timeelapsedingame
	halfTime := false
	endOfRegulation := false
	if turn.Timeelapsedingame < 720 && 720 <= newTimeElapsed {
		newTimeElapsed = 720
	} else if turn.Timeelapsedingame < 1440 && 1440 <= newTimeElapsed {
		halfTime = true
		newTimeElapsed = 1440
	} else if turn.Timeelapsedingame < 2160 && 2160 <= newTimeElapsed {
		newTimeElapsed = 2160
	} else if 2880 <= newTimeElapsed {
		newTimeElapsed = 2880
		endOfRegulation = true
	}
	animation = append(animation, animations["kickoff"][0].Id)
	if playname == "Kick or Punt - Touchback - COP" {
		if gainplayer1 {
			outcomeamount = 0
			outputyards = 80
			player1role = 0
			player2role = 1
		} else {
			outcomeamount = 0
			outputyards = 20
			player1role = 1
			player2role = 0
		}
		lineOfScrimmage = 20
		playtime = 0
		results = "Touchback!\nBall on the 20 yardline.\n\n"
		down = 1
		downdistance = 10
	} else {
		if gainplayer1 {
			outputyards = 100 - outcomeamount
			player1role = 0
			player2role = 1
		} else {
			outputyards = outcomeamount
			player1role = 1
			player2role = 0
		}
		lineOfScrimmage = outputyards
		if outputyards > 50 {
			lineOfScrimmage = 100 - outputyards
		}
		if player1playselected == 14 || player2playselected == 14 {
			if changeofpossession {
				results = fmt.Sprintf("Ball recovered by the Receiving Team on the %v.\n\n", lineOfScrimmage)
				down = 1
				downdistance = 10
			} else {
				results = fmt.Sprintf("Ball recovered by the Kicking Team on the %v.\n\n", lineOfScrimmage)
				down = 1
				downdistance = 10
			}
		} else {
			if outcomeamount > 15 {
				animation = append(animation, animations["long-run"][0].Id)
			}
			results = fmt.Sprintf("Kickoff returned for %v yards.\n\n", outcomeamount)
			down = 1
			downdistance = 10
		}
	}
	if changeofpossession {
		player1role, player2role = cop(player1role, player2role)
	} else {
		animation = append(animation, animations["fumble"][0].Id)
		results = fmt.Sprintf("Fumble on the play\nBall recovered by the Kicking Team!\n\n")
		down = 1
		downdistance = 10
	}
	if outputyards > 99 && player1role == 0 || outputyards < 1 && player2role == 0 {
		score = true
	}
	if score {
		down = 0
		downdistance = 0
		if gainplayer1 {
			outputyards = 35
		} else {
			outputyards = 65
		}
		player1role, player2role, player1score, player2score = scoring(player1role, player2role, player1score, player2score, 7)
		animation = append(animation, animations["touchdowna"][0].Id)
		animation = append(animation, animations["touchdownb"][0].Id)
		animation = append(animation, animations["touchdownc"][0].Id)
		results += "Touchdown!!\n\nThe kick is up...and the extra point is good!"
	} else {
		if changeofpossession {
			animation = append(animation, animations["tackle"][0].Id)
		}
		if !gainplayer1 && downdistance > lineOfScrimmage && outputyards > 90 {
			downdistance = lineOfScrimmage
		}
		if gainplayer1 && downdistance > lineOfScrimmage && outputyards < 10 {
			downdistance = lineOfScrimmage
		}
		results += fmt.Sprintf("1st down and %v to go on the %v", downdistance, lineOfScrimmage)
	}
	fmt.Println("playtime = ", playtime)
	fmt.Println("player1role = ", player1role)
	fmt.Println("player2role = ", player2role)
	fmt.Println("player1score = ", player1score)
	fmt.Println("player2score = ", player2score)
	// fmt.Println("outcomeamount = ", outcomeamount)
	// fmt.Println("outputyards = ", outputyards)
	// fmt.Println("gainplayer1 = ", gainplayer1)
	// fmt.Println("score = ", score)
	// fmt.Println("changeofpossession = ", changeofpossession)
	fmt.Println("animations = ", animations)
	fmt.Println("results = ", results)

	completedgame = ""
	if halfTime {
		results += "\nEnd of half."
		player1role, player2role, player1score, player2score = scoring(player1role, player2role, player1score, player2score, 0)
	}
	if endOfRegulation {
		animation = append(animation, animations["gameovera"][0].Id)
		animation = append(animation, animations["gameoverb"][0].Id)
		tie := player1score == player2score
		if tie {
			results += "\nEnd of Regulation. Resulting in tie."
			completedgame = "Results: Tie score of " + strconv.Itoa(player1score)
		} else if player1score > player2score {
			results += "\nEnd of Regulation. Resulting in player1 winning game."
			completedgame = "Results: player1 won"
		} else if player2score > player1score {
			results += "\nEnd of Regulation. Resulting in player2 winning game."
			completedgame = "Results: player2 won"
		}
		player2role = 6
		player1role = 6
	}
	return playtime, player1role, player2role, player1score, player2score, down, downdistance, outputyards, animation, results, newTimeElapsed, completedgame
}
func GetRun(turn *Turn, guid string, gainplayer1 bool, changeofpossession bool, playname string, outcomeAmount int, positive int, yardLine int, down int, downdistance int) (int, int, int, int, int, int, int, int, []int64, string, int, string) {
	var playtime int
	var outputyards int
	var results string
	var animation []int64
	var lineOfScrimmage int
	var endOfRegulation bool
	var halfTime bool

	score := false
	player1role := turn.Player1role
	player2role := turn.Player2role
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	playtime = r.Intn(8) + 60
	player1score := turn.Currentplayer1score
	player2score := turn.Currentplayer2score

	turn.Playtime = playtime
	newTimeElapsed := playtime + turn.Timeelapsedingame
	halfTime = false
	endOfRegulation = false
	if turn.Timeelapsedingame < 720 && 720 <= newTimeElapsed {
		newTimeElapsed = 720
	} else if turn.Timeelapsedingame < 1440 && 1440 <= newTimeElapsed {
		halfTime = true
		newTimeElapsed = 1440
	} else if turn.Timeelapsedingame < 2160 && 2160 <= newTimeElapsed {
		newTimeElapsed = 2160
	} else if 2880 <= newTimeElapsed {
		newTimeElapsed = 2880
		endOfRegulation = true
	}

	animation = append(animation, animations["snap"][0].Id)
	animation = append(animation, animations["run"][0].Id)
	if gainplayer1 {
		outputyards = yardLine + positive*outcomeAmount
	} else {
		outputyards = yardLine - positive*outcomeAmount
	}
	lineOfScrimmage = outputyards
	if outputyards > 50 {
		lineOfScrimmage = 100 - outputyards
	}

	downdistance = downdistance - (outcomeAmount * positive)
	if downdistance <= 0 {
		downdistance = 10
		down = 1
		if gainplayer1 && downdistance > lineOfScrimmage && outputyards > 90 {
			downdistance = lineOfScrimmage
		}
		if !gainplayer1 && downdistance > lineOfScrimmage && outputyards < 10 {
			downdistance = lineOfScrimmage
		}
	} else {
		down += 1
	}

	turnOverOnDowns := false
	if down > 4 {
		turnOverOnDowns = true
	}
	results = getRandResultString(playname, outcomeAmount*positive) + "\n\n"

	if changeofpossession {
		player1role, player2role = cop(player1role, player2role)
		animation = append(animation, animations["fumble"][0].Id)
		if outputyards > 99 && player2role == 0 || outputyards < 1 && player1role == 0 {
			lineOfScrimmage = 20
			if gainplayer1 {
				outputyards = 80
			} else {
				outputyards = 20
			}
		}
		results += fmt.Sprintf("Ball on the %v yardline \n\n", lineOfScrimmage)

		down = 1
		downdistance = 10
	} else {
		if turnOverOnDowns {
			results += fmt.Sprintf("Turnover on downs ball on the %v yardline \n\n", lineOfScrimmage)
			player1role, player2role = cop(player1role, player2role)
			down = 1
			downdistance = 10
		}
		if outcomeAmount*positive > 15 {
			animation = append(animation, animations["long-run"][0].Id)
		}
		if strings.Contains(playname, "Fumble, Kept") {
			animation = append(animation, animations["fumble"][0].Id)
			results += fmt.Sprintf("Ball on the %v yardline \n\n", lineOfScrimmage)
		}
	}
	if outputyards > 99 && player1role == 0 || outputyards < 1 && player2role == 0 {
		down = 0
		downdistance = 0
		if gainplayer1 {
			outputyards = 35
		} else {
			outputyards = 65
		}
		player1role, player2role, player1score, player2score = scoring(player1role, player2role, player1score, player2score, 7)
		animation = append(animation, animations["touchdowna"][0].Id)
		animation = append(animation, animations["touchdownb"][0].Id)
		animation = append(animation, animations["touchdownc"][0].Id)
		if !changeofpossession {
			if gainplayer1 {
				gain := 100 - yardLine
				results = fmt.Sprintf("Gain of %v \n\n", gain)
			} else {
				gain := yardLine
				results = fmt.Sprintf("Gain of %v \n\n", gain)
			}
		} else {
			results = fmt.Sprintf("Fumble! \n\nRan back for a ")
		}
		lineOfScrimmage = 35
		results += "Touchdown!!\n\nThe kick is up...and the extra point is good!"
	} else {
		animation = append(animation, animations["tackle"][0].Id)
		downString := ""
		switch down {
		case 1:
			downString = "1st"
		case 2:
			downString = "2nd"
		case 3:
			downString = "3rd"
		case 4:
			downString = "4th"
		}
		if !endOfRegulation && !halfTime {
			results += downString + " Down and " + strconv.Itoa(downdistance) + " to go on the " + strconv.Itoa(lineOfScrimmage)
		}
	}

	fmt.Println("playtime = ", playtime)
	fmt.Println("player1role = ", player1role)
	fmt.Println("player2role = ", player2role)
	fmt.Println("player1score = ", player1score)
	fmt.Println("player2score = ", player2score)
	fmt.Println("outcomeamount = ", outcomeAmount)
	fmt.Println("outputyards = ", outputyards)
	fmt.Println("gainplayer1 = ", gainplayer1)
	fmt.Println("score = ", score)
	fmt.Println("changeofpossession = ", changeofpossession)
	fmt.Println("animations = ", animations)
	fmt.Println("results = ", results)

	completedgame := ""
	if halfTime {
		results += "\nEnd of half."
		player1role, player2role, player1score, player2score = scoring(player1role, player2role, player1score, player2score, 0)
		outputyards = 0
		down = 0
		downdistance = 0
	}
	if endOfRegulation {
		animation = append(animation, animations["gameovera"][0].Id)
		animation = append(animation, animations["gameoverb"][0].Id)
		tie := player1score == player2score
		if tie {
			results += "\nEnd of Regulation. Resulting in tie."
			completedgame = "Results: Tie score of " + strconv.Itoa(player1score)
		} else if player1score > player2score {
			results += "\nEnd of Regulation. Resulting in player1 winning game."
			completedgame = "Results: player1 won"
		} else if player2score > player1score {
			results += "\nEnd of Regulation. Resulting in player2 winning game."
			completedgame = "Results: player2 won"
		}
		player2role = 6
		player1role = 6
		outputyards = 0
		down = 0
		downdistance = 0
	}
	return playtime, player1role, player2role, player1score, player2score, down, downdistance, outputyards, animation, results, newTimeElapsed, completedgame
}
func GetPass(turn *Turn, guid string, gainplayer1 bool, changeofpossesion bool, playname string, outcomeAmount int, positive int, yardLine int, down int, downdistance int) (int, int, int, int, int, int, int, int, []int64, string, int, string) {
	var playtime int
	var outputyards int
	var results string
	var animation []int64
	var lineOfScrimmage int

	score := false
	player1role := turn.Player1role
	player2role := turn.Player2role
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	playtime = r.Intn(8) + 60
	player1score := turn.Currentplayer1score
	player2score := turn.Currentplayer2score

	animation = append(animation, animations["snap"][0].Id)
	if gainplayer1 {
		outputyards = yardLine + positive*outcomeAmount
	} else {
		outputyards = yardLine - positive*outcomeAmount
	}
	lineOfScrimmage = outputyards
	if outputyards > 50 {
		lineOfScrimmage = 100 - outputyards
	}
	downdistance = downdistance - (outcomeAmount * positive)
	if downdistance <= 0 {
		downdistance = 10
		down = 1
		if gainplayer1 && downdistance > lineOfScrimmage && outputyards > 90 {
			downdistance = lineOfScrimmage
		}
		if !gainplayer1 && downdistance > lineOfScrimmage && outputyards < 10 {
			downdistance = lineOfScrimmage
		}
	} else {
		down += 1
	}

	turnOverOnDowns := false
	if down > 4 {
		turnOverOnDowns = true
	}

	results = getRandResultString(playname, outcomeAmount*positive) + "\n"

	if changeofpossesion {
		player1role, player2role = cop(player1role, player2role)
		if strings.Contains(playname, "Interception") {
			animation = append(animation, animations["pass"][0].Id)
			animation = append(animation, animations["interception"][0].Id)
			if outputyards > 99 && player2role == 0 || outputyards < 1 && player1role == 0 {
				lineOfScrimmage = 20
				if gainplayer1 {
					outputyards = 80
				} else {
					outputyards = 20
				}
			}
			results += fmt.Sprintf("Ball on the %v yardline \n\n", lineOfScrimmage)
		}
		if strings.Contains(playname, "Fumble") {
			animation = append(animation, animations["pass"][0].Id)
			animation = append(animation, animations["complete"][0].Id)
			animation = append(animation, animations["fumble"][0].Id)
			if outputyards > 99 && player2role == 0 || outputyards < 1 && player1role == 0 {
				lineOfScrimmage = 20
				if gainplayer1 {
					outputyards = 80
				} else {
					outputyards = 20
				}
			}
			results += fmt.Sprintf("Ball on the %v yardline \n\n", lineOfScrimmage)
		}
		down = 1
		downdistance = 10
	} else {
		if turnOverOnDowns {
			results += fmt.Sprintf("Turnover on downs ball on the %v yardline \n\n", lineOfScrimmage)
			player1role, player2role = cop(player1role, player2role)
			down = 1
			downdistance = 10
		}
		if strings.Contains(playname, "Incomplete") {
			animation = append(animation, animations["pass"][0].Id)
			animation = append(animation, animations["incomplete"][0].Id)
			playtime = r.Intn(8) + 10
		}
		if strings.Contains(playname, "Complete") {
			animation = append(animation, animations["pass"][0].Id)
			if positive*outcomeAmount > 15 {
				animation = append(animation, animations["long-pass"][0].Id)
			}
			animation = append(animation, animations["complete"][0].Id)
		}
		if strings.Contains(playname, "QB Scramble") {
		}
		if strings.Contains(playname, "Fumble, Kept") {
			animation = append(animation, animations["pass"][0].Id)
			animation = append(animation, animations["complete"][0].Id)
			animation = append(animation, animations["fumble"][0].Id)
			results = fmt.Sprintf("Ball on the %v yardline \n\n", lineOfScrimmage)
		}
	}

	if outputyards > 99 && player1role == 0 || outputyards < 1 && player2role == 0 {
		score = true
		outputyards = 0
	}
	turn.Playtime = playtime
	newTimeElapsed := playtime + turn.Timeelapsedingame
	halfTime := false
	endOfRegulation := false
	if turn.Timeelapsedingame < 720 && 720 <= newTimeElapsed {
		newTimeElapsed = 720
	} else if turn.Timeelapsedingame < 1440 && 1440 <= newTimeElapsed {
		halfTime = true
		newTimeElapsed = 1440
	} else if turn.Timeelapsedingame < 2160 && 2160 <= newTimeElapsed {
		newTimeElapsed = 2160
	} else if 2880 <= newTimeElapsed {
		newTimeElapsed = 2880
		endOfRegulation = true
	}
	if score {
		down = 0
		downdistance = 0
		if gainplayer1 {
			outputyards = 35
		} else {
			outputyards = 65
		}
		player1role, player2role, player1score, player2score = scoring(player1role, player2role, player1score, player2score, 7)
		animation = append(animation, animations["touchdowna"][0].Id)
		animation = append(animation, animations["touchdownb"][0].Id)
		animation = append(animation, animations["touchdownc"][0].Id)
		if !changeofpossesion {
			if gainplayer1 {
				gain := 100 - yardLine
				if strings.Contains(playname, "QB Scramble") {
					results = fmt.Sprintf("QB scrambles for a gain of %v \n\n", gain)
				} else {
					results = fmt.Sprintf("Gain of %v \n\n", gain)
				}
			} else {
				gain := yardLine
				if strings.Contains(playname, "QB Scramble") {
					results = fmt.Sprintf("QB scrambles for a gain of %v \n\n", gain)
				} else {
					results = fmt.Sprintf("Gain of %v \n\n", gain)
				}
			}
		} else {
			if strings.Contains(playname, "Fumble") {
				results = "Fumble!\n\n"
			}
			if strings.Contains(playname, "Interception") {
				results = "Interception!\n\n"
			}
			results += fmt.Sprintf("Ran back for a ")
		}

		lineOfScrimmage = 35
		results += "Touchdown!!\n\nThe kick is up...and the extra point is good!"
	} else {
		if !strings.Contains(playname, "Incomplete") {
			animation = append(animation, animations["tackle"][0].Id)
		}
		downString := ""
		switch down {
		case 1:
			downString = "1st"
		case 2:
			downString = "2nd"
		case 3:
			downString = "3rd"
		case 4:
			downString = "4th"
		}
		if !endOfRegulation && !halfTime {
			results += downString + " Down and " + strconv.Itoa(downdistance) + " to go on the " + strconv.Itoa(lineOfScrimmage)
		}
	}

	fmt.Println("playname = ", playname)
	fmt.Println("player1role = ", player1role)
	fmt.Println("player2role = ", player2role)
	fmt.Println("player1score = ", player1score)
	fmt.Println("player2score = ", player2score)
	fmt.Println("outcomeamount = ", outcomeAmount)
	fmt.Println("outputyards = ", outputyards)
	fmt.Println("gainplayer1 = ", gainplayer1)
	fmt.Println("score = ", score)
	fmt.Println("changeofpossession = ", changeofpossesion)
	fmt.Println("animations = ", animations)
	fmt.Println("results = ", results)

	completedgame := ""
	if halfTime {
		results += "\nEnd of half."
		outputyards = 0
		down = 0
		downdistance = 0
		player1role, player2role, player1score, player2score = scoring(player1role, player2role, player1score, player2score, 0)
	}
	if endOfRegulation {
		animation = append(animation, animations["gameovera"][0].Id)
		animation = append(animation, animations["gameoverb"][0].Id)
		tie := player1score == player2score
		if tie {
			results += "\nEnd of Regulation. Resulting in tie."
			completedgame = "Results: Tie score of " + strconv.Itoa(player1score)
		} else if player1score > player2score {
			results += "\nEnd of Regulation. Resulting in player1 winning game."
			completedgame = "Results: player1 won"
		} else if player2score > player1score {
			results += "\nEnd of Regulation. Resulting in player2 winning game."
			completedgame = "Results: player2 won"
		}
		player2role = 6
		player1role = 6
		outputyards = 0
		down = 0
		downdistance = 0
	}
	return playtime, player1role, player2role, player1score, player2score, down, downdistance, outputyards, animation, results, newTimeElapsed, completedgame
}
func GetPunt(turn *Turn, guid string, gainplayer1 bool, changeofpossesion bool, outcomeamount int, positive int, yardLine int) (int, int, int, int, int, int, int, int, []int64, string, int, string) {
	var playtime int
	var score bool
	var outputyards int
	var results string
	var animation []int64
	var lineOfScrimmage int
	var down int
	var downdistance int

	score = false
	player1role := turn.Player1role
	player2role := turn.Player2role
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	playtime = r.Intn(8) + 1
	player1score := turn.Currentplayer1score
	player2score := turn.Currentplayer2score

	down = 1
	downdistance = 10
	turn.Playtime = playtime
	newTimeElapsed := playtime + turn.Timeelapsedingame
	halfTime := false
	endOfRegulation := false
	if turn.Timeelapsedingame < 720 && 720 <= newTimeElapsed {
		newTimeElapsed = 720
	} else if turn.Timeelapsedingame < 1440 && 1440 <= newTimeElapsed {
		halfTime = true
		newTimeElapsed = 1440
	} else if turn.Timeelapsedingame < 2160 && 2160 <= newTimeElapsed {
		newTimeElapsed = 2160
	} else if 2880 <= newTimeElapsed {
		newTimeElapsed = 2880
		endOfRegulation = true
	}

	animation = append(animation, animations["snap"][0].Id)
	if gainplayer1 {
		outputyards = yardLine + positive*outcomeamount
	} else {
		outputyards = yardLine - positive*outcomeamount
	}

	lineOfScrimmage = outputyards
	if outputyards > 50 {
		lineOfScrimmage = 100 - outputyards
	}
	if changeofpossesion {
		player1role, player2role = cop(player1role, player2role)
	} else {
		animation = append(animation, animations["fumble"][0].Id)
		downdistance = 10
		results = fmt.Sprintf("Fumble on the play\nBall recovered by the offense on the %v yardline \n\n1st down and %v to go on the %v", lineOfScrimmage, downdistance, lineOfScrimmage)
	}

	if outputyards > 99 && player1role == 0 || outputyards < 1 && player2role == 0 {
		score = true
	}
	touchback := false
	if outputyards > 99 && player2role == 0 || outputyards < 1 && player1role == 0 {
		playtime = 0
		results = "Touchback!\nBall on the 20 yardline\n\n1st Down and 10 to go on the 20"
		touchback = true
		lineOfScrimmage = 20
		if gainplayer1 {
			outputyards = 80
		} else {
			outputyards = 20
		}
	}
	if score {
		down = 0
		downdistance = 0
		if gainplayer1 {
			outputyards = 35
		} else {
			outputyards = 65
		}
		player1role, player2role, player1score, player2score = scoring(player1role, player2role, player1score, player2score, 7)
		animation = append(animation, animations["touchdowna"][0].Id)
		animation = append(animation, animations["touchdownb"][0].Id)
		animation = append(animation, animations["touchdownc"][0].Id)
		results = "Punt returned for a Touchdown!!\n\nThe kick is up...and the extra point is good!"
	} else {
		if changeofpossesion {
			if !touchback {
				downdistance = 10
				if !gainplayer1 && downdistance > lineOfScrimmage && outputyards > 90 {
					downdistance = lineOfScrimmage
				}
				if gainplayer1 && downdistance > lineOfScrimmage && outputyards < 10 {
					downdistance = lineOfScrimmage
				}
				results = fmt.Sprintf("Punt returned to the %v yardline \n\n1st down and %v to go on the %v", lineOfScrimmage, downdistance, lineOfScrimmage)
			}
		}
		animation = append(animation, animations["tackle"][0].Id)
	}

	fmt.Println("playtime = ", playtime)
	fmt.Println("player1role = ", player1role)
	fmt.Println("player2role = ", player2role)
	fmt.Println("player1score = ", player1score)
	fmt.Println("player2score = ", player2score)
	fmt.Println("outcomeamount = ", outcomeamount)
	fmt.Println("outputyards = ", outputyards)
	fmt.Println("gainplayer1 = ", gainplayer1)
	fmt.Println("score = ", score)
	fmt.Println("changeofpossession = ", changeofpossesion)
	fmt.Println("animations = ", animations)
	fmt.Println("results = ", results)

	completedgame := ""
	if halfTime {
		results += "\nEnd of half."
		player1role, player2role, player1score, player2score = scoring(player1role, player2role, player1score, player2score, 0)
	}
	if endOfRegulation {
		animation = append(animation, animations["gameovera"][0].Id)
		animation = append(animation, animations["gameoverb"][0].Id)
		tie := player1score == player2score
		if tie {
			results += "\nEnd of Regulation. Resulting in tie."
			completedgame = "Results: Tie score of " + strconv.Itoa(player1score)
		} else if player1score > player2score {
			results += "\nEnd of Regulation. Resulting in player1 winning game."
			completedgame = "Results: player1 won"
		} else if player2score > player1score {
			results += "\nEnd of Regulation. Resulting in player2 winning game."
			completedgame = "Results: player2 won"
		}
		player2role = 6
		player1role = 6
	}
	return playtime, player1role, player2role, player1score, player2score, down, downdistance, outputyards, animation, results, newTimeElapsed, completedgame
}

func RandomOpponentQueue(guid string) (string, string, int, string, []int64) {
	values := make(url.Values)
	values.Set("guid", guid)
	client, requestError := http.PostForm("http://"+dbaddress+"/random", values)
	// client, requestError := http.PostForm("http://localhost:8081"+"/random", values)
	if requestError != nil {
		go amzses.SendMail("GameOn@PlayVsFootball.com", "ravi@engagemobile.com", "Vsfootball Error", "Error has occured during random opponent info gathering"+requestError.Error())
		defer client.Body.Close()
		return "false", "Failure to queue user.", 502, "", nil
	}

	body, _ := ioutil.ReadAll(client.Body)
	fmt.Println(string(body))
	var result interface{}
	json.Unmarshal(body, &result)
	convertedJson := result.(map[string]interface{})
	fmt.Println(convertedJson)
	success := convertedJson["Success"].(string)
	matchedGuid := convertedJson["Matchedguid"].(string)
	message := convertedJson["Message"].(string)
	go sendRandomOpponentEmailsToAdmin(guid, string(body), matchedGuid)
	var animationIdsOutput []int64
	var status int
	statusFloat := convertedJson["Status"].(float64)
	if statusFloat == 200 {
		status = 200
		animationIdsOutput = append(animationIdsOutput, animations["randomopponent"][0].Id)
	} else if statusFloat == 201 {
		animationIdsOutput = nil
		status = 201
	}
	return success, message, status, matchedGuid, animationIdsOutput
}

func sendRandomOpponentEmailsToAdmin(inputGuid string, body string, matchedGuid string) {
	inputUser := GetUserByGuid(inputGuid)
	matchedUser := GetUserByGuid(matchedGuid)
	var inputUserEmail string
	var matchedUserEmail string
	if inputUser != nil {
		inputUserEmail = inputUser.Username
	}
	if matchedUser != nil {
		matchedUserEmail = matchedUser.Username
	}
	go amzses.SendMail("GameOn@PlayVsFootball.com", "service@playvsfootball.com", "Random Opponent Queue", body+"\nMatched User:"+matchedUserEmail+"\n Input User:"+inputUserEmail)
}

func GetIosDevices() []string {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	var iosDevices []*Iosdevice
	query := "select * from iosdevice"
	_, err := dbmap.Select(&iosDevices, query)
	fmt.Println(err)
	var devices []string
	for _, device := range iosDevices {
		devices = append(devices, device.Devicetoken)
	}
	return devices
}

func GetAndroidDevices() []string {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	var androidDevices []*Androiddevice
	query := "select * from androiddevice"
	_, err := dbmap.Select(&androidDevices, query)
	fmt.Println(err)
	var devices []string
	for _, device := range androidDevices {
		devices = append(devices, device.Deviceid)
	}
	return devices
}

func GetStats(guid, totalOrGuid string) (string, string, int, int, int, int) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	var query string
	var game []*Game
	if totalOrGuid == "total" {
		query = "select * from game where player1=? or player2=?"
		_, err := dbmap.Select(&game, query, guid, guid)
		if isDbQueryError(err) {
			return "false", "Database failure", 501, 0, 0, 0
		}
	} else {
		query = "select * from game where (player1=? and player2=?) or (player1=? and player2=?)"
		_, err := dbmap.Select(&game, query, guid, totalOrGuid, totalOrGuid, guid)
		if isDbQueryError(err) {
			return "false", "Database failure", 501, 0, 0, 0
		}
	}

	wins := 0
	losses := 0
	ties := 0
	for _, gameInArray := range game {
		var winnerString string
		var OtherPlayerResigned string
		if guid == gameInArray.Player1 {
			winnerString = "Results: " + gameInArray.Player1teamname + " won"
			OtherPlayerResigned = "Retired:" + gameInArray.Player2teamname + " Lost."
		} else {
			winnerString = "Results: " + gameInArray.Player2teamname + " won"
			OtherPlayerResigned = "Retired:" + gameInArray.Player1teamname + " Lost."
		}
		if gameInArray.Outcome == winnerString || gameInArray.Outcome == OtherPlayerResigned {
			wins += 1
		} else if strings.Contains(gameInArray.Outcome, "Tie") {
			ties += 1
		} else if gameInArray.Outcome != "" {
			losses += 1
		}
	}
	return "true", "Successfully gathered stats data.", 200, wins, losses, ties
}

func GetOpponentsForGuid(guid string) (string, string, int, []map[string]string) {
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	var games []*Game
	query := "select * from game where player1=? or player2=?"
	_, err := dbmap.Select(&games, query, guid, guid)
	if isDbQueryError(err) {
		return "false", "Database failure", 501, nil
	}
	var guidEmail []map[string]string
	guidsUsedSoFar := make(map[string]int)
	for _, game := range games {
		if game.Player1 == guid {
			if _, exists := guidsUsedSoFar[game.Player2]; !exists {
				opp := GetUserByGuid(game.Player2)
				if opp != nil {
					user := make(map[string]string)
					user["Email"] = opp.Username
					user["Guid"] = opp.Guid
					guidEmail = append(guidEmail, user)
					guidsUsedSoFar[game.Player2] = 0
				}
			}
		} else {
			if _, exists := guidsUsedSoFar[game.Player1]; !exists {
				opp := GetUserByGuid(game.Player1)
				if opp != nil {
					user := make(map[string]string)
					user["Email"] = opp.Username
					user["Guid"] = opp.Guid
					guidEmail = append(guidEmail, user)
				}
				guidsUsedSoFar[game.Player1] = 0
			}
		}
	}
	return "true", "Successfully found opponents", 200, guidEmail
}
const (
	getPresetQuery = `
		SELECT preset.* FROM preset;
	`
)
func GetPreset() (string,string,int,[]string){
	dbmap := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(dbmap)
	var presets []*Preset
	_,err:=dbmap.Select(&presets,getPresetQuery)
	if isDbQueryError(err) {
		return "false", "Database failure",501, nil
	}
	var presetString []string
	for _,preset := range presets{
		presetString = append(presetString,preset.Value)
	}
	return "true","Successfully gathered preset what's on the line items.",200,presetString
}

