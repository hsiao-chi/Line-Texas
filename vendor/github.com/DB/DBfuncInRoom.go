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
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	var nn string
	db.QueryRow("SELECT UserNickName FROM sql6131889.User WHERE MID = ?", MID).Scan(&nn)
	var R string
	db.QueryRow("SELECT UserRoom FROM sql6131889.User WHERE MID = ?", MID).Scan(&R)
	bot.SendText([]string{MID}, "哈囉! "+nn+"!\n您現在位於房間 "+R+"\n可用指令為:\n!離開房間\n!加入遊戲\n!開始遊戲\n!新遊戲")
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
	var gameCancel int
	row,_ := db.Query("SELECT Cancel FROM sql6131889.Game WHERE RoomID = ?", RID)
	for row.Next() { 
		row.Scan(&gameCancel)
	}
	if haveGame == "" || gameCancel == 1 {
		db.Exec("INSERT INTO sql6131889.Game (GameName, RoomID, GameStatus, GameTokens, GamePlayer1, GameMaster, PlayerNum, Cancel) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", "TexasPoker", RID, 1, 0, MID, "0", 1, 0)	
		db.QueryRow("SELECT ID FROM sql6131889.Game WHERE RoomID = ?", RID).Scan(&GID)
		db.Exec("INSERT INTO sql6131889.GameAction (MID, GameID, PlayerX, Action, Cancel) VALUES (?, ?, ?, ?, ?)", MID, GID, 1, 0, 0)
		db.Exec("UPDATE sql6131889.Room SET RoomStatus = ? WHERE RoomName = ?", 101, R)
		bot.SendText([]string{MID}, "您建立了一個新遊戲")
		bot.SendText([]string{MID}, "您是 Player1")
	}else{
		bot.SendText([]string{MID}, "這間房間已有另一場遊戲")
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
		bot.SendText([]string{MID}, "請用 !新遊戲 建立遊戲")
	}else{
		var playerInGame string
		db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE MID = ?", MID).Scan(&playerInGame)
		var nextPlayer int
		var gameActionCancel int
		db.QueryRow("SELECT Cancel FROM sql6131889.GameAction WHERE MID = ?", MID).Scan(&gameActionCancel)
		row,_ := db.Query("SELECT Cancel FROM sql6131889.GameAction WHERE MID = ?", MID)
		for row.Next() { 
			row.Scan(&playerInGame)
		}
		if playerInGame == "" || gameActionCancel == 0{
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
				bot.SendText([]string{MID}, "您是 Player1")
			}else if nextPlayer == 2 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer2 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "您是 Player2")
			}else if nextPlayer == 3 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer3 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "您是 Player3")
			}else if nextPlayer == 4 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer4 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "您是 Player4")
			}else if nextPlayer == 5 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer5 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "您是 Player5")
			}else if nextPlayer == 6 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer6 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "您是 Player6")
			}else if nextPlayer == 7 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer7 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "您是 Player7")
			}else if nextPlayer == 8 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer8 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "您是 Player8")
			}else if nextPlayer == 9 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer9 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "您是 Player9")
			}else if nextPlayer == 10 {
				db.Exec("UPDATE sql6131889.Game SET GamePlayer10 = ? WHERE ID = ?", MID, GID)
				bot.SendText([]string{MID}, "您是 Player10")
			}
			db.Exec("UPDATE sql6131889.Game SET PlayerNum = ? WHERE ID = ?", nextPlayer, GID)
		}else if nextPlayer == 50 {
			bot.SendText([]string{MID}, "您正在進行遊戲中!!")
		}else{
			bot.SendText([]string{MID}, "此房間玩家已滿!!")
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
							bot.SendText([]string{MID}, "必須至少兩個玩家以上才能開始遊戲")
						}
					}else{
						bot.SendText([]string{MID}, "遊戲已開始")
					}
			}else{
				bot.SendText([]string{MID}, "您不在玩家名單內")
			}
		}else{
			bot.SendText([]string{MID}, "您不在玩家名單內")
		}
	}else{
		bot.SendText([]string{MID}, "這個房間尚未存在遊戲")
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
