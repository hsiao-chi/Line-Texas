package DB
import(
	"os"
	"strconv"
	"github.com/line/line-bot-sdk-go/linebot"
	"database/sql"
	_"github.com/go-sql-driver/mysql"

)

var bot *linebot.Client

func ChatInRoom(mID string,gID int,t string) {
	//
	
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ?", gID)
	for row.Next() {
		var mid1 string
		row.Scan(&mid1)
		if mid1 != mID{
			var n string
			db.QueryRow("SELECT UserName FROM sql6131889.GameAction WHERE MID = ?",mID).Scan(&n)
			bot.SendText([]string{mid1}, n+":\n"+t)
		}
	}
	db.Close()
}


func Management(mID string, text string) { // if playing call this func 
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	var uR int
	db.QueryRow("SELECT UserRoom FROM sql6131889.User WHERE MID = ?",mID).Scan(&uR)
	var S int
	db.QueryRow("SELECT GameStatus FROM sql6131889.Game WHERE RoomId = ?",uR).Scan(&S)
	var gID int//輸入者在玩的GAMEID
	db.QueryRow("SELECT GameID FROM sql6131889.GameAction WHERE MID = ?",mID).Scan(&gID)
	if S == 1{//等人
		//yu-chi
	}else if S == 2{//開始Game

	}else if S == 3{//發牌=一人2張

	}else if S == 3{//第一輪下注

	}else if S == 4{//發牌=檯面3張

	}else if S == 5{//第二輪下注

	}else if S == 6{//發牌=檯面4張

	}else if S == 7{//第三輪下注

	}else if S == 8{//發牌=檯面5張

	}else if S == 9{//第四輪下注

	}else if S == 10{//輸贏+分錢

	}
	db.Close()
}


func CallToken(mID string, text string) {
	// every function needs to open db again
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	var gID int//輸入者在玩的GAMEID
	db.QueryRow("SELECT GameID FROM sql6131889.GameAction WHERE MID = ?",mID).Scan(&gID)
	var tN int//GAME的狀態
	db.QueryRow("SELECT Turn FROM sql6131889.Game WHERE ID = ?",gID).Scan(&tN)
	var money int = 5//money 小盲柱
	var P int//輸入者的身分
	db.QueryRow("SELECT PlayerX FROME sql6131889.GameAction WHERE MID?",mID).Scan(&P)
	//row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ?", gID)
	var mT int//最高投注金額
	db.QueryRow("SELECT MaxToken FROM sql6131889.Game WHERE ID = ?",gID).Scan(&mT)
	if P == tN{
		if text == "!call"{ //跟注
			Call()
			//var pmoney int
			//db.QueryRow("SELECT UserToken FROM sql6131889.User WHERE MID = ?",content.From).Scan(&pmoney)
			//db.Exec("UPDATE sql6131889.User SET UserToken = ? WHERE MID = ?",pmoney-money,content.From)
			//bot.SendText([]string{content.From},"系統: \nflow: "+money+"\n餘額: "+pmoney-money)
			//db.Exec("UPDATE sql6131889.GameAction SET Action = ? WHERE MID = ?",31,content.From)
			//var rID int
			//db.Exec("UPDATE sql6131889.Game SET GameTokens = ? WHERE RoomId",)

			// db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE RoomId = ?",21,gID)
			// db.Exec("UPDATE sql6131889.GameAction SET Action = ? WHERE MID = ?",mT,content.From)
			// row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ?", gID)
			// for row.Next() {
			// 	var mid1 string
			// 	row.Scan(&mid1)
			// 	if mid1 != content.From{
			// 		bot.SendText([]string{mid1}, "player1: 跟注")
			// 	}
			// }
		}else if text == "!pass"{
			// bot.SendText([]string{content.From},"系統: \npass")
			// db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE RoomId = ?",21,gID)
			// db.Exec("UPDATE sql6131889.GameAction SET Action = ? WHERE MID = ?",0,content.From)
			// row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ?", gID)
			// for row.Next() {
			// 	var mid1 string
			// 	row.Scan(&mid1)
			// 	if mid1 != content.From{
			// 		bot.SendText([]string{mid1}, "player1: pass")
			// 	}
			// }
		}else if text == "!allin"{

		}else{//聊天
				ChatInRoom(mID,gID,text)
		}
	}else{
		ChatInRoom(mID,gID,text)
	}
}
//Call WHEN PlayerToken ADD OR SUB
func AddPlayerToken(MID string,addtoken int){
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	db.QueryRow("UPDATE sql6131889.User SET UserToken=UserToken+? WHERE MID =?",addtoken,MID)
	db.Close()
}
//Call WHEN GAMETOKEN ADD OR SUB
func AddGameToken(RoomId int,addtoken int){
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	db.QueryRow("UPDATE sql6131889.Game SET GameToken=GameToken+? WHERE RoomID =?",addtoken,RoomId)
	db.Close()
}
func RunOne (mID string,nowS int,st int,text string,gID int,mT int) {
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	if nowS == st{
		if text == "!flow"{
				
				//var pmoney int
				//db.QueryRow("SELECT UserToken FROM sql6131889.User WHERE MID = ?",content.From).Scan(&pmoney)
				//db.Exec("UPDATE sql6131889.User SET UserToken = ? WHERE MID = ?",pmoney-money,content.From)
				//bot.SendText([]string{content.From},"系統: \nflow: "+money+"\n餘額: "+pmoney-money)
				//db.Exec("UPDATE sql6131889.GameAction SET Action = ? WHERE MID = ?",31,content.From)

				//var rID int

				//db.Exec("UPDATE sql6131889.Game SET GameTokens = ? WHERE RoomId",)

			db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE RoomId = ?",21,gID)
			db.Exec("UPDATE sql6131889.GameAction SET Action = ? WHERE MID = ?",mT,mID)
			row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ?", gID)
			for row.Next() {
				var mid1 string
				row.Scan(&mid1)
				if mid1 != content.From{
					bot.SendText([]string{mid1}, "player1: 跟注")
				}
			}
		}else if text.Text == "!pass"{
				bot.SendText([]string{content.From},"系統: \npass")
				db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE RoomId = ?",21,gID)
				db.Exec("UPDATE sql6131889.GameAction SET Action = ? WHERE MID = ?",0,content.From)
				row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ?", gID)
				for row.Next() {
					var mid1 string
					row.Scan(&mid1)
					if mid1 != content.From{
						bot.SendText([]string{mid1}, "player1: pass")
					}
				}
			}else{//聊天
				ChatInRoom(content.From,gID,text.Text)
			}
		}else{
			ChatInRoom(content.From,gID,text.Text)
		}
}