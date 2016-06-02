package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	// "strconv"
	// "strings"
	// "time"
	// "crypto/sha1"
	// "encoding/base64"
)

func main() {
	// fmt.Println("signup test")

	// fmt.Println(time.Now())
	// for i := 0; i < 100000; i++ {

	// 	// fmt.Println(time.Now())
	// 	values := make(url.Values)
	// 	values.Set("email", "ravi@galleonlabs.com")
	// 	// values.Set("firstname","rpmd30@gmail.com")
	// 	// values.Set("lastname","rpmd30@gmail.com")
	// 	values.Set("password", "password")
	// 	values.Set("platform", "iOS")
	// 	// myclient , err := http.PostForm("http://localhost:8080/user/signup",values)
	// fmt.Println(time.Now())
	// values := make(url.Values)
	// values.Set("email", "ravi@galleonlabs.com")
	// // values.Set("firstname","rpmd30@gmail.com")
	// // values.Set("lastname","rpmd30@gmail.com")
	// values.Set("password", "password")
	// values.Set("platform", "iOS")
	// myclient, _ := http.PostForm("http://localhost:8080/user/signup", values)
	// // 	func() {
	// // 		// http.PostForm("http://localhost:8080/user/signup", values)
	// // 		myclient, err := http.PostForm("http://localhost:8080/user/signup", values)
	// // 		if err != nil {
	// // 			fmt.Println(string(err.Error()))
	// // 		} else {
	// // 			ioutil.ReadAll(myclient.Body)
	// // 			myclient.Body.Close()
	// ioutil.ReadAll(myclient.Body)
	// myclient.Body.Close()
	// 			// if err != nil {
	// 			// 	fmt.Println(string(err.Error()))
	// 			// } else {
	// 			// 	fmt.Println(string(body))
	// 			// }
	// 		}
	// 	}()
	// }
	// fmt.Println(time.Now())
	// fmt.Println("login test")

	// values2 := make(url.Values)
	// values2.Set("username","rpmd30@gmail.com")
	//             // hasher := sha1.New()
	//            // hasher.Write([]byte("ZW5nYWdlbW9iaWxlMTIz"))
	//            // sha1output := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	// values2.Set("password","rpmd30@gmail.com")
	// myclient2, loginErr := http.PostForm("http://localhost:8080/user/login",values2)
	// if(loginErr != nil){
	// 	fmt.Println(string(loginErr.Error()))
	// } else{
	// 	body, err := ioutil.ReadAll(myclient2.Body)
	// 	myclient2.Body.Close()
	// 	if(err != nil){
	// 		fmt.Println(string(loginErr.Error()))
	// 	} else {
	// 		fmt.Println(string(body))
	// 	}
	// }

	// fmt.Println("verify test")
	// query := "http://localhost:8080/32a18f665cfd4f28614d1bdcda6ae6c6/verify"
	// myclient3, err := http.Get(query)
	// fmt.Println(err)
	// body, _ := ioutil.ReadAll(myclient3.Body)
	// myclient3.Body.Close()
	// fmt.Println(body)

	// fmt.Println("resend email test")
	// query2 := "http://localhost:8080/user/sendemailnotification"
	// values4 := make(url.Values)
	// values4.Set("email","rpmd30@gmail.com")
	// myclient4,_ := http.PostForm(query2,values4)
	// body2 ,_ := ioutil.ReadAll(myclient4.Body)
	// myclient4.Body.Close()
	// fmt.Println(body2)

	// fmt.Println("forgot email test")
	// query2 := "http://localhost:8080/user/forgotpassword"
	// values4 := make(url.Values)
	// values4.Set("email", "rpmd30@gmail.com")
	// myclient4, _ := http.PostForm(query2, values4)
	// body2, _ := ioutil.ReadAll(myclient4.Body)
	// myclient4.Body.Close()
	// fmt.Println(body2)

	// query6 := "http://localhost:8080/8991812860d941ea6f4c230b8fe2531e/feedback"
	// values6 := make(url.Values)
	// values6.Set("comment", "my cool comments")
	// values6.Set("gameid", "2")
	// myclient6, _ := http.PostForm(query6, values6)
	// body6, _ := ioutil.ReadAll(myclient6.Body)
	// myclient6.Body.Close()
	// fmt.Println(body6)
	// for i := 0; i < 100000; i++ {
	// 	query7 := "http://localhost:8080/8991812860d941ea6f4c230b8fe2531e/games"
	// 	values7 := make(url.Values)
	// 	values7.Set("inviteEmail", "ravi@galleonlabs.com")
	// 	values7.Set("possession", "5")
	// 	values7.Set("teamName", "ravi's team")
	// 	values7.Set("playIdSelected", "37")
	// 	myclient7, _ := http.PostForm(query7, values7)
	// 	body7, _ := ioutil.ReadAll(myclient7.Body)
	// 	myclient7.Body.Close()
	// 	fmt.Println(string(body7))
	// }

	// query8 :="http://localhost:8080/32a18f665cfd4f28614d1bdcda6ae6c6/games"
	// myclient8,_ := http.Get(query8)
	// body8,_ :=ioutil.ReadAll(myclient8.Body)
	// myclient8.Body.Close()
	// fmt.Println(string(body8))

	// query9 := "http://localhost:8080/8991812860d941ea6f4c230b8fe2531e/games/resign"
	// values9 := make(url.Values)
	// values9.Set("gameid", "2")
	// myclient9, _ := http.PostForm(query9, values9)
	// body9, _ := ioutil.ReadAll(myclient9.Body)
	// myclient9.Body.Close()
	// fmt.Println(string(body9))

	// values10 := make(url.Values)
	// values10.Set("email","ravi@galleonlabs.com")
	// values10.Set("firstname","rpmd30@gmail.com")
	// values10.Set("lastname","rpmd30@gmail.com")
	// values10.Set("accounttype","facebook:rpmd30@gmail.com")
	// values10.Set("accesstoken","facebook:rpmd30@gmail.com")
	// values10.Set("tokenexpiration","facebook:rpmd30@gmail.com")
	// values10.Set("platform","iOS")
	// // myclient , err := http.PostForm("http://localhost:8080/user/signup",values)
	// myclient10 , _ := http.PostForm("http://localhost:8080/user/signup/facebook",values10)
	// 	body10, _ := ioutil.ReadAll(myclient10.Body)
	// 	myclient10.Body.Close()
	// 		fmt.Println(string(body10))

	// values11 := make(url.Values)
	// values11.Set("gameid","2")
	// values11.Set("teamname","das2")
	// myclient11,_ := http.PostForm("http://localhost:8080/b702995f863a4ace5bb81a304aecff14/games/confirm",values11)
	// body11,_ := ioutil.ReadAll(myclient11.Body)
	// myclient11.Body.Close()
	// fmt.Println(string(body11))
	// fmt.Println(time.Now())
	// for i:= 4 ; i<250	 ; i++{
	// fmt.Println(strings.Replace("Results: player2 7-0", "player2", "ravi", -1))
	// fmt.Println(strings.Contains("Results: player2 7-0", "player2"))
	// i := 143

	// // // myclient31, _ := http.Get("http://localhost:8080/8991812860d941ea6f4c230b8fe2531e/games/2/turn/1")
	// // // body31, _ := ioutil.ReadAll(myclient31.Body)
	// // // myclient31.Body.Close()
	// // // fmt.Println(string(body31))

	// values30 := make(url.Values)
	// values30.Set("playid", "20")
	// myclient30, _ := http.PostForm("http://localhost:8080/8991812860d941ea6f4c230b8fe2531e/games/turn/"+strconv.Itoa(i)+"/play", values30)
	// body30, _ := ioutil.ReadAll(myclient30.Body)
	// myclient30.Body.Close()
	// fmt.Println(string(body30))

	// values12 := make(url.Values)
	// values12.Set("playid", "28")
	// myclient12, _ := http.PostForm("http://localhost:8080/b702995f863a4ace5bb81a304aecff14/games/turn/"+strconv.Itoa(i)+"/play", values12)
	// body12, _ := ioutil.ReadAll(myclient12.Body)
	// myclient12.Body.Close()
	// fmt.Println(string(body12))

	// }
	// fmt.Println(time.Now())

	// values18 := make(url.Values)
	// values18.Set("playid", "15")
	// values18.Set("receipt", "15")

	// myclient18, _ := http.PostForm("http://localhost:8080/8991812860d941ea6f4c230b8fe2531e/playbook/itunes", values18)
	// body18, _ := ioutil.ReadAll(myclient18.Body)
	// myclient18.Body.Close()
	// fmt.Println(string(body18))
	// values18 := make(url.Values)
	// values18.Set("playid", "15")
	// values18.Set("receipt", "12354324265436543654365436543")

	// myclient18, _ := http.PostForm("http://vsf001.engagemobile.com/ef524eb0c0db43346363e78429a9cfcd/playbook/itunes", values18)
	// body18, _ := ioutil.ReadAll(myclient18.Body)
	// myclient18.Body.Close()
	// fmt.Println(string(body18))
	// now := time.Now().Unix()
	// fmt.Println(now)
	// // 1377207608
	// var past int64 = 1376870026 //time.Unix(1376869552593828557, 0)

	// fmt.Println(now - past)
	// fmt.Println(past - now)
	// values18 := make(url.Values)
	// values18.Set("registrationid", "APA91bFL_cyKOlXcPwE3iNp4J0oB9A_vpdXowd7wu5j6kkx7AQQSwS_9VGU5DfnIZjIKAPU3XQrl4YeXDomnNC1VSBSj06Vzdzrw79D9fhTlunB7xaUI3j3m9ip1_ZpxbYWuJrmKbtZnSuVOznYKTtRtE6R4vYjQa1SuxUw7Rrp4DHOkpxlZia0")
	// myclient18, _ := http.PostForm("http://localhost:8080/8991812860d941ea6f4c230b8fe2531e/android/register", values18)
	// body18, _ := ioutil.ReadAll(myclient18.Body)
	// myclient18.Body.Close()
	// fmt.Println(string(body18))
	// values31 := make(url.Values)
	// values31.Set("inviteeFacebookId", "1283188515")
	// values31.Set("possession", "4")
	// values31.Set("playId", "37")
	// values31.Set("teamName", "Ravi")
	// myclient31, _ := http.PostForm("http://vsf001.engagemobile.com/8f23fc6f10d6428a6534b0d5862aabe1/facebook/game/invite", values31)
	// body31, _ := ioutil.ReadAll(myclient31.Body)
	// myclient31.Body.Close()
	// fmt.Println(string(body31))

	// values33 := make(url.Values)
	// values33.Set("message", "testing")
	// values33.Set("deviceid", "APA91bGhJROyjSWhWl8eTM2tru3y6pJp15BDCliuCK7WGIZqL3LUtyfAZJmjOKf-PccLu0SroX8DgQIV93foIxyYmAD33kSkoIgLqSuu1ZHNFJEFBby-qs08rgbwW1YKyBqHN4e9Yh-MMrFW0FiYio-A7gbVd93CRjgCHc2tphGFONhjAIRqtJQ")
	// myclient33, _ := http.PostForm("http://vsf001.engagemobile.com/test/gcm/echo", values33)
	// body33, _ := ioutil.ReadAll(myclient33.Body)
	// myclient33.Body.Close()
	// fmt.Println(string(body33))
	// values33 := make(url.Values)
	// // values33.Set("message", "testing")
	// values33.Set("devicetoken", "e2884292a311e3898acbde0086356dad45d3baabbf48c7265ea43b9b316a56e5")
	// myclient33, _ := http.PostForm("http://localhost:8080/8991812860d941ea6f4c230b8fe2531e/ios/register", values33)
	// body33, _ := ioutil.ReadAll(myclient33.Body)
	// myclient33.Body.Close()
	// fmt.Println(string(body33))

	// values12 := make(url.Values)
	// values12.Set("message", "team name")
	// myclient12, err := http.PostForm("http://localhost:8080/8991812860d941ea6f4c230b8fe2531e/game/1/chat", values12)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// body12, _ := ioutil.ReadAll(myclient12.Body)
	// myclient12.Body.Close()
	// fmt.Println(string(body12))

	// myclient12, err := http.Get("http://localhost:8080/d6419873f2b847c04053a5cd4e70e884/game/1/chat")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// body12, _ := ioutil.ReadAll(myclient12.Body)
	// myclient12.Body.Close()
	// fmt.Println(string(body12))
	query7 := "http://vsfdev.engagemobile.com/8991812860d941ea6f4c230b8fe2531e/game/random"
	values7 := make(url.Values)
	myclient7, _ := http.PostForm(query7, values7)
	body7, _ := ioutil.ReadAll(myclient7.Body)
	myclient7.Body.Close()
	fmt.Println(string(body7))

	query7 = "http://vsfdev.engagemobile.com/123/game/random"
	values7 = make(url.Values)
	myclient7, _ = http.PostForm(query7, values7)
	body7, _ = ioutil.ReadAll(myclient7.Body)
	myclient7.Body.Close()
	fmt.Println(string(body7))

}
