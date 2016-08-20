package DB
import (
	"os"
	"strconv"
	"github.com/line/line-bot-sdk-go/linebot"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
)

var bot *linebot.Client

func InRoomInst(MID string){
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	/*
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	var haveGame string
	var RID string
	var R string
	db.QueryRow("SELECT UserRoom FROM sql6131889.User WHERE MID = ?", MID).Scan(&R)
	db.QueryRow("SELECT ID FROM sql6131889.User WHERE  RoomName = ?", R).Scan(&RID)
	db.QueryRow("SELECT RoomID FROM sql6131889.Game WHERE RoomID = ?", RID).Scan(&haveGame)
	if haveGame == ""{
		bot.SendText([]string{MID}, "You can use these instruction:\n!leavechatroom\n!newgame")
	}else{
		bot.SendText([]string{MID}, "You can use these instruction:\n!leavechatroom")
	}
	db.Close()*/
	bot.SendText([]string{MID}, "You can use these instruction:\n!leavechatroom\n!newgame\n!joingame")
}
func InRoomNewGame(MID string){
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	var haveGame string
	var RID int
	var R string
	var GID string
	db.QueryRow("SELECT UserRoom FROM sql6131889.User WHERE MID = ?", MID).Scan(&R)
	db.QueryRow("SELECT ID FROM sql6131889.Room WHERE  RoomName = ?", R).Scan(&RID)
	db.QueryRow("SELECT RoomID FROM sql6131889.Game WHERE RoomID = ?", RID).Scan(&haveGame)
	if haveGame == ""{
		db.Exec("INSERT INTO sql6131889.Game (GameName, RoomID, GameStatus, GameTokens, GamePlayer1, GameMaster, PlayerNum, Cancel) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", "TexasPoker", RID, 1, 0, MID, "0", 1, 0)	
		db.QueryRow("SELECT ID FROM sql6131889.Game WHERE RoomID = ?", RID).Scan(&GID)
		db.Exec("INSERT INTO sql6131889.GameAction (MID, GameID, PlayerX, Action, Cancel) VALUES (?, ?, ?, ?, ?)", MID, GID, 1, 0, 0)
		db.Exec("UPDATE sql6131889.Room SET RoomStatus = ? WHERE RoomName = ?", 101, R)
		bot.SendText([]string{MID}, "You created a new game")
		bot.SendText([]string{MID}, "You are Player1")
	}else{
		bot.SendText([]string{MID}, "There is already a game in this room!!")
	}
	db.Close()
}
func InRoomJoinGame(MID string){
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	var haveGame int
	var RID string
	var R string
	var GID string
	db.QueryRow("SELECT UserRoom FROM sql6131889.User WHERE MID = ?", MID).Scan(&R)
	db.QueryRow("SELECT ID FROM sql6131889.Room WHERE  RoomName = ?", R).Scan(&RID)
	db.QueryRow("SELECT ID FROM sql6131889.Game WHERE RoomID = ?", RID).Scan(&GID)
	db.QueryRow("SELECT RoomStatus FROM sql6131889.Room WHERE RoomName = ?", R).Scan(&haveGame)
	if haveGame == 100 {
		bot.SendText([]string{MID}, "Please create a new game use instruction:\n!newgame")
	}else{
		var playerInGame string
		db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE MID = ?", MID).Scan(&playerInGame)
		var nextPlayer int
		var gameActionCancel int
		db.QueryRow("SELECT Cancel FROM sql6131889.GameAction WHERE MID = ?", MID).Scan(&gameActionCancel)
		if playerInGame == "" || gameActionCancel == 1{
			db.QueryRow("SELECT PlayerNum FROM sql6131889.Game WHERE ID = ?", GID).Scan(&nextPlayer)
			nextPlayer = nextPlayer+1
		}else{
			nextPlayer = 50
		}
		if nextPlayer <= 10 {
			db.Exec("INSERT INTO sql6131889.GameAction (MID, GameID, PlayerX, Action, Cancel) VALUES (?, ?, ?, ?, ?)", MID, GID, nextPlayer, 0, 0)
			if nextPlayer == 1 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer1 = ? WHERE ID = ?", MID, GID)
				db.Exec("UPDATE sql6131889.GameAction SET PlayerX = ? WHERE MID = ?", 1, MID)
				bot.SendText([]string{MID}, "You are Player1")
			}else if nextPlayer == 2 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer2 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player2")
			}else if nextPlayer == 3 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer3 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player3")
			}else if nextPlayer == 4 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer4 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player4")
			}else if nextPlayer == 5 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer5 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player5")
			}else if nextPlayer == 6 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer6 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player6")
			}else if nextPlayer == 7 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer7 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player7")
			}else if nextPlayer == 8 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer8 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player8")
			}else if nextPlayer == 9 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer9 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player9")
			}else if nextPlayer == 10 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer10 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player10")
			}
			db.Exec("UPDATE sql6131889.Game SET PlayerNum = ? WHERE ID = ?", nextPlayer, GID)
		}else if nextPlayer == 50 {
			bot.SendText([]string{MID}, "You are already in this game!!")
		}else{
			bot.SendText([]string{MID}, "Full of player in this room!!")
		}
	}
	db.Close()
}
func InRoomStartGame(MID string){
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	var RID string
	var R string
	var GID string
	var haveGame int
	db.QueryRow("SELECT UserRoom FROM sql6131889.User WHERE MID = ?", MID).Scan(&R)
	db.QueryRow("SELECT ID FROM sql6131889.Room WHERE  RoomName = ?", R).Scan(&RID)
	db.QueryRow("SELECT ID FROM sql6131889.Game WHERE RoomID = ?", RID).Scan(&GID)
	var playerInGame string
	db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE MID = ?", MID).Scan(&playerInGame)
	db.QueryRow("SELECT RoomStatus FROM sql6131889.Room WHERE RoomName = ?", R).Scan(&haveGame)
	if haveGame == 101 {
		if playerInGame != ""{
			var gameActionCancel int
			db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE MID = ?", MID).Scan(&gameActionCancel)
			if gameActionCancel != 1{
					var waitingForStart int
					db.QueryRow("SELECT GameStatus FROM sql6131889.Game WHERE ID = ?", GID).Scan(&waitingForStart)
					if waitingForStart == 1 {
						var gamerNum int
						db.QueryRow("SELECT PlayerNum FROM sql6131889.Game WHERE ID = ?", GID).Scan(&gamerNum)
						if gamerNum > 1 {
							db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE ID = ?", 2, GID) //starting game now
							bot.SendText([]string{MID}, "== START THE GAME ==")
						}else{
							bot.SendText([]string{MID}, "The game can't be starting below 2 players!!")
						}
					}else{
						bot.SendText([]string{MID}, "The game is already started!!")
					}
			}else{
				bot.SendText([]string{MID}, "You are not in the game!!")
			}
		}else{
			bot.SendText([]string{MID}, "You are not in the game!!")
		}
	}else{
		bot.SendText([]string{MID}, "This room have no game!!")
	} 
}
func CancelGameAction(MID string){
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	var RID string
	var R string
	var GID string
	db.QueryRow("SELECT UserRoom FROM sql6131889.User WHERE MID = ?", MID).Scan(&R)
	db.QueryRow("SELECT ID FROM sql6131889.Room WHERE  RoomName = ?", R).Scan(&RID)
	db.QueryRow("SELECT ID FROM sql6131889.Game WHERE RoomID = ?", RID).Scan(&GID)
	db.Exec("UPDATE sql6131889.GameAction SET Cancel = ? WHERE GameID = ?", 1, GID)
}
func CancelGame(MID string){
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	var RID string
	var R string
	var GID string
	db.QueryRow("SELECT UserRoom FROM sql6131889.User WHERE MID = ?", MID).Scan(&R)
	db.QueryRow("SELECT ID FROM sql6131889.Room WHERE  RoomName = ?", R).Scan(&RID)
	db.QueryRow("SELECT ID FROM sql6131889.Game WHERE RoomID = ?", RID).Scan(&GID)
	db.Exec("UPDATE sql6131889.Game SET Cancel = ? WHERE RoomID = ?", 1, RID)
}
