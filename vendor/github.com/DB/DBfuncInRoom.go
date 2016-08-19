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
	/*strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
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
	var RID string
	var R string
	var GID string
	db.QueryRow("SELECT UserRoom FROM sql6131889.User WHERE MID = ?", MID).Scan(&R)
	db.QueryRow("SELECT ID FROM sql6131889.Room WHERE  RoomName = ?", R).Scan(&RID)
	db.QueryRow("SELECT ID FROM sql6131889.Game WHERE RoomID = ?", RID).Scan(&GID)
	db.QueryRow("SELECT RoomID FROM sql6131889.Game WHERE RoomID = ?", RID).Scan(&haveGame)
	if haveGame == ""{
		db.Exec("INSERT INTO sql6131889.Game (GameName, RoomID, GameStatus, GameTokens, GamePlayer1, GameMaster, Cancel) VALUES (?, ?, ?, ?, ?, ?, ?)", "TexasPoker", RID, 100, 0, MID, "0", 0)
		
		db.Exec("INSERT INTO sql6131889.GameAction (MID, GameID, PlayerX, Action, Cancel) VALUE (?, ?, ?, ?, ?)", MID, GID, 20, 0, 0)
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
	var haveGame string
	var RID string
	var R string
	var GID string
	db.QueryRow("SELECT UserRoom FROM sql6131889.User WHERE MID = ?", MID).Scan(&R)
	db.QueryRow("SELECT ID FROM sql6131889.Room WHERE  RoomName = ?", R).Scan(&RID)
	db.QueryRow("SELECT ID FROM sql6131889.Game WHERE RoomID = ?", RID).Scan(&GID)
	db.QueryRow("SELECT RoomID FROM sql6131889.Game WHERE RoomID = ?", RID).Scan(&haveGame)
	if haveGame == ""{
		bot.SendText([]string{MID}, "Please create a new game use instruction:\n!newgame")
	}else{
		var playerInGame string
		db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE MID = ?", MID).Scan(&playerInGame)
		var nextPlayer int
		if playerInGame == "" {
			row,_ := db.Query("SELECT PlayerX FROM sql6131889.GameAction WHERE GameID = ?", GID)
			for row.Next() {
				row.Scan(&nextPlayer)
			}
			nextPlayer = nextPlayer+1
		}else{
			nextPlayer = 50
		}
		if nextPlayer <= 29 {
			db.Exec("INSERT INTO sql6131889.GameAction (MID, GameID, PlayerX, Action, Cancel) VALUE (?, ?, ?, ?, ?)", MID, GID, nextPlayer, 0, 0)
			if nextPlayer == 1 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer1 = ? WHERE ID = ?", MID, GID)
				db.Exec("UPDATE sql6131889.GameAction SET PlayerX = ? WHERE MID = ?", 20, MID)
				bot.SendText([]string{MID}, "You are Player1")
			}else if nextPlayer == 21 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer2 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player2")
			}else if nextPlayer == 22 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer3 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player3")
			}else if nextPlayer == 23 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer4 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player4")
			}else if nextPlayer == 24 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer5 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player5")
			}else if nextPlayer == 25 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer6 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player6")
			}else if nextPlayer == 26 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer7 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player7")
			}else if nextPlayer == 27 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer8 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player8")
			}else if nextPlayer == 28 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer9 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player9")
			}else if nextPlayer == 29 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer10 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "You are Player10")
			}
		}else if nextPlayer == 50 {
			bot.SendText([]string{MID}, "You are already in this game!!")
		}else{
			bot.SendText([]string{MID}, "Full of player in this room!!")
		}
	}
	db.Close()
}
