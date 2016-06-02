package main

import (
	"controllers"
	"fmt"
	"github.com/gorilla/mux"
	"models"
	"net/http"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("ENV says GOMAXPROCS: :%s: \n", os.Getenv("GOMAXPROCS"))
	fmt.Printf("runtime says MAXPROCS = %d \n", runtime.NumCPU())
	models.Init()
	r := mux.NewRouter()
	r.HandleFunc("/", controllers.HandleHome).Methods("GET")
	r.HandleFunc("/user/signup", controllers.HandleSignUp).Methods("POST")
	r.HandleFunc("/user/login", controllers.HandleLogin).Methods("POST")
	r.HandleFunc("/user/sendemailnotification", controllers.HandleEmailNotification).Methods("POST")
	r.HandleFunc("/user/forgotpassword", controllers.HandleForgotPassword).Methods("POST")
	r.HandleFunc("/{guid}/verify", controllers.HandleVerify).Methods("GET")
	r.HandleFunc("/{guid}/feedback", controllers.HandleFeedback).Methods("POST")
	r.HandleFunc("/{guid}/games", controllers.HandleCreateGame).Methods("POST")
	r.HandleFunc("/{guid}/games", controllers.HandleGamesList).Methods("GET")
	r.HandleFunc("/{guid}/games/resign", controllers.HandleGameResign).Methods("POST")
	r.HandleFunc("/user/signup/facebook", controllers.HandleFacebookSignup).Methods("POST")
	r.HandleFunc("/{guid}/invite/email", controllers.HandleInviteEmail).Methods("POST")
	r.HandleFunc("/{guid}/games/confirm", controllers.HandleConfirmGame).Methods("POST")
	r.HandleFunc("/{guid}/games/turn/{turnid}/play", controllers.HandleSelectPlay).Methods("POST")
	r.HandleFunc("/{guid}/playbook", controllers.HandlePlaybook).Methods("GET")
	r.HandleFunc("/{guid}/game/{gameid}/chat", controllers.HandleChatSubmit).Methods("POST")
	r.HandleFunc("/{guid}/game/{gameid}/chat", controllers.HandleChatRetrieve).Methods("GET")
	r.HandleFunc("/{guid}/games/id", controllers.HandleGetGameIds).Methods("GET")
	r.HandleFunc("/{guid}/games/{gameid}/turn/{returnlimit}", controllers.HandleTurnsByGame).Methods("GET")
	r.HandleFunc("/oddstable", controllers.HandleOddsTableUpdate).Methods("POST")
	r.HandleFunc("/playbook/ids", controllers.HandleGetAllPlayIds).Methods("GET")
	r.HandleFunc("/playbook/{playid}", controllers.HandleGetPlayForPlayId).Methods("GET")
	r.HandleFunc("/{guid}/playbook/owned", controllers.HandleGetOwnedPlays).Methods("GET")
	r.HandleFunc("/{guid}/playbook/itunes", controllers.HandleBuyItunesConnectPlay).Methods("POST")
	r.HandleFunc("/{guid}/playbook/playstore", controllers.HandlePlayStorePlayPurchase).Methods("POST")
	r.HandleFunc("/animations/ids", controllers.HandleGetAnimationIds).Methods("GET")
	r.HandleFunc("/animations/{animationid}", controllers.HandleGetAnimationForAnimationId).Methods("GET")
	r.HandleFunc("/{guid}/android/register", controllers.HandleAndroidRegister).Methods("POST")
	r.HandleFunc("/{guid}/ios/register", controllers.HandleIosRegister).Methods("POST")
	r.HandleFunc("/{guid}/facebook/game/invite", controllers.HandleFacebookGameInvite).Methods("POST")
	r.HandleFunc("/{guid}/game/rematch", controllers.HandleRematch).Methods("POST")
	r.HandleFunc("/test/gcm/echo", controllers.HandleGcmTest).Methods("POST")
	r.HandleFunc("/turnfailcheck", controllers.HandleTurnFailCheck).Methods("POST")
	r.HandleFunc("/{guid}/password", controllers.HandleChangePassword).Methods("POST")
	r.HandleFunc("/viralcoefficient", controllers.HandleViralCoefficient).Methods("GET")
	r.HandleFunc("/{guid}/game/random", controllers.HandleRandomOpponent).Methods("POST")
	r.HandleFunc("/login", controllers.HandlePushHome).Methods("GET", "POST")
	r.HandleFunc("/push", controllers.HandlePush).Methods("POST")
	r.HandleFunc("/push/devices", controllers.HandleGetAllDevices).Methods("GET")
	r.HandleFunc("/push/device/{type}", controllers.HandlePush).Methods("POST")
	r.HandleFunc("/{guid}/game/stats/{opponent_guid}", controllers.HandleStatsTotal).Methods("GET")
	r.HandleFunc("/{guid}/game/opponents", controllers.HandleStatsForOpponent).Methods("GET")
	r.HandleFunc("/{guid}/whatsontheline/preset", controllers.HandleWhatsOnTheLinePreset).Methods("GET")
	r.HandleFunc("/{guid}/whatsontheline/{gameid}", controllers.HandleWhatsOnTheLine).Methods("POST")
	r.HandleFunc("/{guid}/whatsontheline/{gameid}", controllers.HandleWhatsOnTheLineForAGame).Methods("GET")
	r.HandleFunc("/{guid}/whatsontheline/{gameid}/confirm/{whatsonthelineid}", controllers.HandleWhatsOnTheLineConfirm).Methods("POST")

	// r.HandleFunc("/test/apns/echo", controllers.HandleAPNSTest).Methods("POST")
	// r.HandleFunc("/{guid}/products", controllers.HandleProductsList).Methods("GET")

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
