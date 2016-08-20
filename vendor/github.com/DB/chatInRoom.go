package DB
import(
	"os"
	"strconv"
	"github.com/line/line-bot-sdk-go/linebot"
	"database/sql"
	_"github.com/go-sql-driver/mysql"

)

var bot *linebot.Client

func chatInRoom(mID string,gID int,t string) {
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
		if callToken1(mID,text){
			S++
		}
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

//第一輪加注
func callToken1(mID string, text string) bool{
	// every function needs to open db again
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	var uR string//在的房間name
	db.QueryRow("SELECT UserRoom FROM sql6131889.User WHERE MID = ?",mID).Scan(&uR)
	var rID int//在的房間ID
	db.QueryRow("SELECT ID FROM sql6131889.room WHERE RoomName = ?",uR).Scan(&rID)
	var gID int//輸入者在玩的GAMEID
	db.QueryRow("SELECT GameID FROM sql6131889.Game WHERE RoomId = ?",rID).Scan(&gID)
	var tN int//GAME的狀態turn
	db.QueryRow("SELECT Turn FROM sql6131889.Game WHERE ID = ?",gID).Scan(&tN)
	var money int = 5//money 小盲柱
	var P int//輸入者的身分
	db.QueryRow("SELECT PlayerX FROME sql6131889.GameAction WHERE MID?",mID).Scan(&P)
	//row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ?", gID)
	var mT int//最高投注金額
	db.QueryRow("SELECT MaxToken FROM sql6131889.Game WHERE ID = ?",gID).Scan(&mT)
	var pN int//遊戲人數
	db.QueryRow("SELECT PlayerNum FROM sql6131889.Game WHERE ID = ?",gID).Scan(&pN)

	if P == tN{
		runOne(mID,text,gID,rID,mT,(tN+1)%pN)
	}else{
		chatInRoom(mID,gID,text)
	}
	var tmp int = 0
	row,_ := db.Query("SELECT Action FROM sql6131889.GameAction WHERE GameID = ?", gID)
	for row.Next() {
		var act int
		row.Scan(&act)
		if act == mT || act == -1{
			tmp++
		}
	}
	return tmp == pN
}


func runOne (mID string,text string,gID int,rID int,mT int,nextS int) {
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
		if text == "!Call"{

			AddPlayerToken(mID,(-1)*mT)
			AddGameToken(rID,mT)

			db.Exec("UPDATE sql6131889.Game SET Turn = ? WHERE RoomId = ?",nextS,gID)
			db.Exec("UPDATE sql6131889.GameAction SET Action = ? WHERE MID = ?",mT,mID)
			row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ?", gID)
			for row.Next() {
				var mid1 string
				row.Scan(&mid1)
				if mid1 != mID{
					var n string
					db.QueryRow("SELECT UserName FROM sql6131889.GameAction WHERE MID = ?",mID).Scan(&n)
					bot.SendText([]string{mid1}, n+": 跟注")
				}
			}
			var mid2 string
			db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE PlayerX = ?",nextS).Scan(&mid2)
			bot.SendText([]string{mid2}, "系統: 跟注金額"+strconv.Itoa(mT)+" 請選擇指令\n!Call\n!Fold\n!allin")
		}else if text == "!Fold"{
				bot.SendText([]string{content.From},"系統: \nFold")
				db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE RoomId = ?",nextS,gID)
				db.Exec("UPDATE sql6131889.GameAction SET Action = ? WHERE MID = ?",-1,mID)
				row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ?", gID)
				for row.Next() {
					var mid1 string
					row.Scan(&mid1)
					if mid1 != mID{
						var n string
						db.QueryRow("SELECT UserName FROM sql6131889.GameAction WHERE MID = ?",mID).Scan(&n)
						bot.SendText([]string{mid1}, n+": Fold")
					}
				}
			var mid2 string
			db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE PlayerX = ?",nextS).Scan(&mid2)
			bot.SendText([]string{mid2}, "系統: 跟注金額"+strconv.Itoa(mT)+" 請選擇指令\n!Call\n!Fold\n!allin")
		}else if text == "!allin"{

			
		}else{//聊天
			ChatInRoom(content.From,gID,text.Text)
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
