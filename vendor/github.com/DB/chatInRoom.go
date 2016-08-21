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
	row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ? AND Cancel = ?", gID, 0)
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
	db.QueryRow("SELECT ID FROM sql6131889.Room WHERE RoomName = ? AND Cancel = ?", uR, 0).Scan(&rid)
	var S int
	db.QueryRow("SELECT GameStatus FROM sql6131889.Game WHERE RoomID = ? AND Cancel = ?",rid, 0).Scan(&S)
	var gID int//輸入者在玩的GAMEID
	db.QueryRow("SELECT GameID FROM sql6131889.GameAction WHERE MID = ? AND Cancel = ?",mID, 0).Scan(&gID)
	var cardOnTable [5]int
	cardOnTable = GetFiveCards(gID)
	if S == 1{//等人
		chatInRoom(mID,gID,text)
	}else if S == 2{//開始Game
		
	}else if S == 3{//發牌=一人2張

	}
	if S == 4{//第一輪下注
		if callToken1(mID,text,S){
			db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE RoomID = ? AND Cancel = ?",5,gID, 0)
		}
	}
	if S == 5{//發牌=檯面3張
		row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ? AND Cancel = ?", gID, 0)
		for row.Next() {
			var mid1 string
			row.Scan(&mid1)
			bot.SendText([]string{mid1}, "開牌3張")
			c1 := GetCardName(cardOnTable[0])
			c2 := GetCardName(cardOnTable[1])
			c3 := GetCardName(cardOnTable[2])
			bot.SendText([]string{mid1}, "牌桌上：\n" + c1 + "\n" + c2+ "\n" + c3)
		}
				
		db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE RoomID = ? AND Cancel = ?",6,gID, 0)
	}
	if S == 6{//第二輪下注
		if callToken1(mID,text,S){
			db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE RoomID = ? AND Cancel = ?",7,gID, 0)
		}
	}
	if S == 7{//發牌=檯面4張
		row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ? AND Cancel = ?", gID, 0)
		for row.Next() {
			var mid1 string
			row.Scan(&mid1)
			bot.SendText([]string{mid1}, "開牌1張")
			c1 := GetCardName(cardOnTable[0])
			c2 := GetCardName(cardOnTable[1])
			c3 := GetCardName(cardOnTable[2])
			c4 := GetCardName(cardOnTable[3])
			bot.SendText([]string{mid1}, "牌桌上：\n" + c1 + "\n" + c2+ "\n" + c3+ "\n" + c4)
		}
		db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE RoomID = ? AND Cancel = ?",8,gID, 0)
	}
	if S == 8{//第三輪下注
		if callToken1(mID,text,S){
			db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE RoomID = ? AND Cancel = ?",9,gID, 0)
		}
	}
	if S == 9{//發牌=檯面5張
		row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ? AND Cancel = ?", gID, 0)
		for row.Next() {
			var mid1 string
			row.Scan(&mid1)
			bot.SendText([]string{mid1}, "開牌1張")
			c1 := GetCardName(cardOnTable[0])
			c2 := GetCardName(cardOnTable[1])
			c3 := GetCardName(cardOnTable[2])
			c4 := GetCardName(cardOnTable[3])
			c5 := GetCardName(cardOnTable[4])
			bot.SendText([]string{mid1}, "牌桌上：\n" + c1 + "\n" + c2+ "\n" + c3+ "\n" + c4+ "\n" + c5)
		}
		db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE RoomID = ? AND Cancel = ?",10,gID, 0)
	}
	if S == 10{//第四輪下注
		if callToken1(mID,text,S){
		db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE RoomID = ? AND Cancel = ?",11,gID, 0)
		}
	}
	var great [11][6]int
	if S == 11{//輸贏+分錢
		

		row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ? AND Cancel = ?", gID, 0)
		for row.Next() {
			var mid1 string
			row.Scan(&mid1)
			var color [4]int
			var point [13]int
			color,point=BigCombi (mid1)
			var pX int
			db.QueryRow("SELECT PlayerX FROM sql6131889.GameAction WHERE GameID = ? AND Cancel = ? AND MID = ?", gID, 0,mid1).Scan(&pX)
			great[pX]= CalculatePoint (color,point)
			
		}
		db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE RoomID = ? AND Cancel = ?",12,gID, 0)
	}
	if S == 12{
		win:=checkwinner(great[1],great[2],great[3],great[4],great[5],great[6],great[7],great[8],great[9],great[10])
		var winid int
		db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE GameID = ? AND Cancel = ? AND PlayerX = ?", gID, 0,win).Scan(&winid)
		var winname string
		db.QueryRow("SELECT UserName FROM sql6131889.User WHERE MID", winid).Scan(&winname)
		row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ? AND Cancel = ?", gID, 0)
		for row.Next() {
			var mid1 string
			row.Scan(&mid1)
			bot.SendText([]string{mid1}, "遊戲結束\nWnner: "+winname)
		}

	}
	if S == 200{
		var md string
		db.QueryRow("SELECT Template1 FROM sql6131889.Game WHERE ID = ? AND Cancel = ?",gID, 0).Scan(&md)
		if md == mID {
			setbetprize(mID,text, gID)
		}else{
			chatInRoom(mID,gID,text)
		}
	}
	db.Close()
}

//第一輪加注
func callToken1(mID string, text string,S int) bool{
	// every function needs to open db again
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))

	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	var uR string//在的房間name
	db.QueryRow("SELECT UserRoom FROM sql6131889.User WHERE MID = ?",mID).Scan(&uR)
	var rID int//在的房間ID
	db.QueryRow("SELECT ID FROM sql6131889.Room WHERE RoomName = ? AND Cancel = ?",uR, 0).Scan(&rID)
	var gID int//輸入者在玩的GAMEID
	db.QueryRow("SELECT ID FROM sql6131889.Game WHERE RoomID = ? AND Cancel = ?",rID, 0).Scan(&gID)
	var tN int//GAME的狀態turn
	db.QueryRow("SELECT Turn FROM sql6131889.Game WHERE ID = ? AND Cancel = ?",gID, 0).Scan(&tN)
	//var money int = 5//money 小盲柱
	var P int//輸入者的身分
	db.QueryRow("SELECT PlayerX FROM sql6131889.GameAction WHERE MID = ? AND Cancel = ?",mID, 0).Scan(&P)
	//row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ?", gID)
	var mT int//最高投注金額
	db.QueryRow("SELECT MaxToken FROM sql6131889.Game WHERE ID = ? AND Cancel = ?",gID, 0).Scan(&mT)
	var pN int//遊戲人數
	db.QueryRow("SELECT PlayerNum FROM sql6131889.Game WHERE ID = ? AND Cancel = ?",gID, 0).Scan(&pN)
	if P == tN{
		if S == 4{
			runOne(mID,text,gID,rID,mT,tN)
		}else if S>4{
			runTwo(mID,text,gID,rID,mT,tN)
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
	bot.SendText([]string{mID}, "test\ntmp= "+strconv.Itoa(tmp)+"\npN= "+strconv.Itoa(pN)+"\nmT= "+strconv.Itoa(mT))//////////////test
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
			
		}else if text == "!Bet"{
			runBet(mID,text,gID,rID,mT,nextS)
		}else if text == "!See"{
			See(mID, gID)
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
			db.Exec("UPDATE sql6131889.Game SET Turn = ? WHERE RoomID = ?",nextS,gID)
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
	var act int
	db.QueryRow("SELECT Action FROM sql6131889.GameAction WHERE MID = ? AND Cancel = ?",mID, 0).Scan(&act)
	AddPlayerToken(mID,(-1)*(mT - act))
	AddGameToken(rID,mT - act)
	db.Exec("UPDATE sql6131889.GameAction SET Action = ? WHERE MID = ? AND Cancel = ?",mT,mID, 0)

	
	row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ? AND Cancel = ?", gID, 0)
	for row.Next() {
		var mid1 string
		row.Scan(&mid1)
		if mid1 != mID{
			var n string
			db.QueryRow("SELECT UserName FROM sql6131889.User WHERE MID = ?",mID).Scan(&n)
			bot.SendText([]string{mid1}, n+": 跟注")
		}
	}
	var pN int//遊戲人數
	db.QueryRow("SELECT PlayerNum FROM sql6131889.Game WHERE ID = ? AND Cancel = ?",gID, 0).Scan(&pN)
	for i := 0;i < pN;i++ {
		nextS += 1
		if nextS > pN {
			nextS = 1
		}
		var act int
		db.QueryRow("SELECT Action FROM sql6131889.GameAction WHERE GameID = ? AND PlayerX = ? AND Cancel = ?",gID, nextS, 0).Scan(&act)
		if act != -1 {
			break
		}
	}
	db.Exec("UPDATE sql6131889.Game SET Turn = ? WHERE ID = ? AND Cancel = ?",nextS,gID, 0)
	var mid2 string
	db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE PlayerX = ? AND Cancel = ? AND GameID = ?",nextS, 0, gID).Scan(&mid2)
	bot.SendText([]string{mid2}, "目前下注金額 $"+strconv.Itoa(mT)+" 請選擇指令\n!Bet\n!Check\n!Call\n!Raise\n!Fold")
}
//棄牌
func runFold(mID string,text string,gID int,mT int,nextS int){

	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	bot.SendText([]string{mID},"You Fold")
	db.Exec("UPDATE sql6131889.GameAction SET Action = ? WHERE MID = ? AND Cancel = ?",-1,mID, 0)
	row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ?", gID)
	for row.Next() {
		var mid1 string
		row.Scan(&mid1)
		if mid1 != mID{
			var n string
			db.QueryRow("SELECT UserName FROM sql6131889.User WHERE MID = ?",mID).Scan(&n)
			bot.SendText([]string{mid1}, n+": Fold")
		}
	}
	counts := 0
	var pN int//遊戲人數
	db.QueryRow("SELECT PlayerNum FROM sql6131889.Game WHERE ID = ? AND Cancel = ?",gID, 0).Scan(&pN)
	var winner int
	for i := 1;i <= pN;i++ {
		var act int
		db.QueryRow("SELECT Action FROM sql6131889.GameAction WHERE GameID = ? AND PlayerX = ? AND Cancel = ?",gID, i, 0).Scan(&act)
		if act == -1 {
			counts += 1
		}else{
			winner = i
		}
	}
	if counts == pN - 1 { // one player left
		var pN int//遊戲人數
		db.QueryRow("SELECT PlayerNum FROM sql6131889.Game WHERE ID = ? AND Cancel = ?",gID, 0).Scan(&pN)
		db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE ID = ? AND Cancel = ?",12, gID, 0)
		var md string // winner mid
		db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE PlayerX= ? AND Cancel = ?",winner, 0).Scan(&md)
		var n string
		db.QueryRow("SELECT UserName FROM sql6131889.User WHERE MID = ?",md).Scan(&n)
		row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ? AND Cancel = ?", gID, 0)
		for row.Next() {
			var mid1 string
			row.Scan(&mid1)
				if md == mid1{
					bot.SendText([]string{mid1}, "YOU WIN")
				}else{
					bot.SendText([]string{mid1}, n+" WIN")
				}
		}
	}else{
		var pN int//遊戲人數
		db.QueryRow("SELECT PlayerNum FROM sql6131889.Game WHERE ID = ? AND Cancel = ?",gID, 0).Scan(&pN)
		for i := 0;i < pN;i++ {
			nextS += 1
			if nextS > pN {
				nextS = 1
			}
			var act int
			db.QueryRow("SELECT Action FROM sql6131889.GameAction WHERE GameID = ? AND PlayerX = ? AND Cancel = ?",gID, nextS, 0).Scan(&act)
			if act != -1 {
				break
			}
		}
		db.Exec("UPDATE sql6131889.Game SET Turn = ? WHERE ID = ? AND Cancel = ?",nextS,gID, 0)
		var mid2 string
		db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE PlayerX = ? AND Cancel = ? AND GameID = ?",nextS, 0, gID).Scan(&mid2)
		bot.SendText([]string{mid2}, "目前下注金額 $"+strconv.Itoa(mT)+" 請選擇指令\n!See\n!Bet\n!Check\n!Call\n!Raise\n!Fold")
	}

	
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
	db.Exec("UPDATE sql6131889.Game SET MaxToken = ? WHERE RoomID = ?",mT,gID)
	db.Exec("UPDATE sql6131889.Game SET Turn = ? WHERE RoomID = ?",nextS,gID)
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

func runBet(mID string,text string,gID int,rID int,mT int,nextS int) {
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	var S int
	db.QueryRow("SELECT GameStatus FROM sql6131889.Game WHERE ID = ? AND Cancel = ?",gID, 0).Scan(&S)
	db.Exec("UPDATE sql6131889.Game SET Template2 = ? WHERE ID = ? AND Cancel = ?", S, gID, 0)
	db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE ID = ? AND Cancel = ?",200, gID, 0) // GameStatus = 200
	db.Exec("UPDATE sql6131889.Game SET Template1 = ? WHERE ID = ? AND Cancel = ?", mID, gID, 0)
	bot.SendText([]string{mID}, "請輸入金額:")
	db.Close()
}

func setbetprize(mID string,text string, gID int){
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	prize , _ := strconv.Atoi(text)
	db.Exec("UPDATE sql6131889.GameAction SET Action = ? WHERE MID = ? AND Cancel = ?", prize, mID, 0)
	db.Exec("UPDATE sql6131889.Game SET MaxToken = ? WHERE ID = ? AND Cancel = ?", prize, gID, 0)
	var S int
	db.QueryRow("SELECT Template2 FROM sql6131889.Game WHERE ID = ? AND Cancel = ?",gID, 0).Scan(&S)
	db.Exec("UPDATE sql6131889.Game SET GameStatus = ? WHERE ID = ? AND Cancel = ?", S, gID, 0)
	row,_ := db.Query("SELECT MID FROM sql6131889.GameAction WHERE GameID = ? AND Cancel = ?", gID, 0)
	for row.Next() {
		var mid1 string
		row.Scan(&mid1)
		if mid1 != mID{
			var n string
			db.QueryRow("SELECT UserName FROM sql6131889.User WHERE MID = ?",mID).Scan(&n)
			bot.SendText([]string{mid1}, n+" 下注 $"+text)
		}
	}
	bot.SendText([]string{mID}, "你下注 $"+text)
	var pN int//遊戲人數
	db.QueryRow("SELECT PlayerNum FROM sql6131889.Game WHERE ID = ? AND Cancel = ?",gID, 0).Scan(&pN)
	var nextS int
	db.QueryRow("SELECT PlayerX FROM sql6131889.GameAction WHERE MID = ? AND Cancel = ?",mID, 0).Scan(&nextS)
	for i := 0;i < pN;i++ {
		nextS += 1
		if nextS > pN {
			nextS = 1
		}
		var act int
		db.QueryRow("SELECT Action FROM sql6131889.GameAction WHERE GameID = ? AND PlayerX = ? AND Cancel = ?",gID, nextS, 0).Scan(&act)
		if act != -1 {
			break
		}
	}
	db.Exec("UPDATE sql6131889.Game SET Turn = ? WHERE ID = ? AND Cancel = ?",nextS,gID, 0)
	var mid2 string
	db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE PlayerX = ? AND Cancel = ? AND GameID = ?",nextS, 0, gID).Scan(&mid2)
	bot.SendText([]string{mid2}, "目前下注金額 $"+text+" 請選擇指令\n!See\n!Bet\n!Check\n!Call\n!Raise\n!Fold")
	db.Close()
}

func See(mID string, gID int){
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	var card1 int
	var card2 int
	var card3 int
	var card4 int
	var card5 int
	db.QueryRow("SELECT Card1, Card2, Card3, Card4, Card5 FROM sql6131889.Game WHERE ID = ? AND Cancel = ?",gID, 0).Scan(&card1, &card2, &card3, &card4, &card5)
	cn1 := GetCardName(card1)
	cn2 := GetCardName(card2)
	cn3 := GetCardName(card3)
	cn4 := GetCardName(card4)
	cn5 := GetCardName(card5)
	bot.SendText([]string{mID}, "牌桌上的牌:\n"+cn1+"\n"+cn2+"\n"+cn3+"\n"+cn4+"\n"+cn5)
	db.QueryRow("SELECT PlayerCard1, PlayerCard2 FROM sql6131889.GameAction WHERE MID = ? AND Cancel = ?",mID, 0).Scan(&card1, &card2)
	cn1 = GetCardName(card1)
	cn2 = GetCardName(card2)
	bot.SendText([]string{mID}, "您的手牌:\n"+cn1+"\n"+cn2)
	var mT int//最高投注金額
	db.QueryRow("SELECT MaxToken FROM sql6131889.Game WHERE ID = ? AND Cancel = ?",gID, 0).Scan(&mT)
	prize := strconv.Itoa(mT)
	bot.SendText([]string{mID}, "目前下注金額 $"+prize+" 請選擇指令\n!See\n!Bet\n!Check\n!Call\n!Raise\n!Fold")
	db.Close()
}

