package DB
import(
	"os"
	"strconv"
	"github.com/line/line-bot-sdk-go/linebot"
	"database/sql"
	_"github.com/go-sql-driver/mysql"

)

//var bot *linebot.Client

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
			db.QueryRow("SELECT UserName FROM sql6131889.User WHERE MID = ?",mID).Scan(&n)
			bot.SendText([]string{mid1}, n+":\n"+t)
		}
	}
	db.Close()
}


func Management(mID string, text string) { // if playing call this func 
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")

	var uR string
	db.QueryRow("SELECT UserRoom FROM sql6131889.User WHERE MID = ?",mID).Scan(&uR)
	var rid int
	db.QueryRow("SELECT ID FROM sql6131889.Room WHERE RoomName = ?", uR).Scan(&rid)
	var S int
	db.QueryRow("SELECT GameStatus FROM sql6131889.Game WHERE RoomId = ?",rid).Scan(&S)
	var gID int//輸入者在玩的GAMEID
	db.QueryRow("SELECT GameID FROM sql6131889.GameAction WHERE MID = ?",mID).Scan(&gID)
	if S == 1{//等人

	}else if S == 2{//開始Game
		row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ?", gID)
		for row.Next() {
			var mid1 string
			row.Scan(&mid1)
			bot.SendText([]string{mid1}, "現在開始遊戲-Texas")
		}
		S=3
	}
	if S == 3{//發牌=一人2張
	
		row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ?", gID)
		for row.Next() {
			var mid1 string
			row.Scan(&mid1)
			var cards [2]int
			cards = GetTwoCards(mid1)
			c1 := GetCardName(cards[0])
			c2 := GetCardName(cards[1])
			bot.SendText([]string{mid1}, "您的手牌為mid：\n" + c1 + "\n" + c2)
		}
		var p1 string
		db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE PlayerX = ?AND GameID = ?",1,gID).Scan(&p1)
		bot.SendText([]string{p1}, "系統: 跟注金額 5$\n請選擇指令 !Call")
		S=4
	}else if S == 4{//第一輪下注
		if callToken1(mID,text,S){
			S=5
		}
	}else if S == 5{//發牌=檯面3張

	}else if S == 6{//第二輪下注
		if callToken1(mID,text,S){
			S=7
		}
	}else if S == 7{//發牌=檯面4張

	}else if S == 8{//第三輪下注
		if callToken1(mID,text,S){
			S=9
		}
	}else if S == 9{//發牌=檯面5張

	}else if S == 10{//第四輪下注
		if callToken1(mID,text,S){
			S=11
		}
	}else if S == 11{//輸贏+分錢

	}
	db.Close()
}

//第一輪加注
func callToken1(mID string, text string,S int) bool{

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
	mT = money
	if P == tN{
		if S == 4{
			runOne(mID,text,gID,rID,mT,(tN+1)%pN)
		}else if S>4{
			runTwo(mID,text,gID,rID,mT,(tN+1)%pN)
		}
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
	//db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
		if text == "!Call"{
			runCall(mID,text,gID,rID,mT,nextS)
		
		}else if text == "!Fold"{
			runFold(mID,text,gID,mT,nextS)
				
		}else if text == "!Raise"{
			runRaise(mID,text,gID,rID,mT,nextS)
			
		}else{//聊天
			chatInRoom(mID,gID,text)
		}
		
}
func runTwo (mID string,text string,gID int,rID int,mT int,nextS int) {
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	if text == "!Call"{
		runCall(mID,text,gID,rID,mT,nextS)
		
	}else if text == "!Fold"{
		runFold(mID,text,gID,mT,nextS)
	}else if text == "!Raise"{
		runRaise(mID,text,gID,rID,mT,nextS)
		
	}else if text == "!Pass"{
		if mT == 0{
			bot.SendText([]string{mID},"系統: \nPass")
			db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE RoomId = ?",nextS,gID)
			db.Exec("UPDATE sql6131889.GameAction SET Action = ? WHERE MID = ?",0,mID)
			row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ?", gID)
			for row.Next() {
				var mid1 string
				row.Scan(&mid1)
				if mid1 != mID{
					var n string
					db.QueryRow("SELECT UserName FROM sql6131889.GameAction WHERE MID = ?",mID).Scan(&n)
					bot.SendText([]string{mid1}, n+": Pass")
				}
			}
			var mid2 string
			db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE PlayerX = ?",nextS).Scan(&mid2)
			bot.SendText([]string{mid2}, "系統: 跟注金額"+strconv.Itoa(mT)+" 請選擇指令\n!Call\n!Fold\n!Raise")
		}else{
			bot.SendText([]string{mID}, "你不能pass 只能\n!Call\n!Fold\n!Raise")
		}
	}else{//聊天
		chatInRoom(mID,gID,text)
	}
		
}

//跟注
func runCall(mID string,text string,gID int,rID int,mT int,nextS int) {
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	
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
	bot.SendText([]string{mid2}, "系統: 跟注金額"+strconv.Itoa(mT)+" 請選擇指令\n!Call\n!Fold\n!Raise")
}
//棄牌
func runFold(mID string,text string,gID int,mT int,nextS int){
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	bot.SendText([]string{mID},"系統: \nFold")
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
	bot.SendText([]string{mid2}, "系統: 跟注金額"+strconv.Itoa(mT)+" 請選擇指令\n!Call\n!Fold\n!Raise")
}
//加注
func runRaise(mID string,text string,gID int,rID int,mT int,nextS int) {
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	mT*=2
	AddPlayerToken(mID,(-1)*mT)
	AddGameToken(rID,mT)
	db.Exec("UPDATE sql6131889.Game SET MaxToken = ? WHERE RoomId = ?",mT,gID)
	db.Exec("UPDATE sql6131889.Game SET Turn = ? WHERE RoomId = ?",nextS,gID)
	db.Exec("UPDATE sql6131889.GameAction SET Action = ? WHERE MID = ?",mT,mID)
	row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ?", gID)
	for row.Next() {
		var mid1 string
		row.Scan(&mid1)
		if mid1 != mID{
			var n string
			db.QueryRow("SELECT UserName FROM sql6131889.GameAction WHERE MID = ?",mID).Scan(&n)
			bot.SendText([]string{mid1}, n+": 加注")
		}
	}
	var mid2 string
	db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE PlayerX = ?",nextS).Scan(&mid2)
	bot.SendText([]string{mid2}, "系統: 跟注金額"+strconv.Itoa(mT)+" 請選擇指令\n!Call\n!Fold\n!Raise")
}


