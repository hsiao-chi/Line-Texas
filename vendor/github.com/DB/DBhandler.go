/*
使用者是否在遊戲中
他手上有哪兩張牌
*/
package DB
import (
	"math/rand"
	"time"
	"os"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
)
/*
傳入使用者MID
回傳使用者是否正在遊戲
*/
func UserGamming(MID string) bool{
	var GameID int
	GameID = 0;
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	db.QueryRow("SELECT GameID FROM sql6131889.GameAction WHERE MID = ? and Cancel = 0", MID ).Scan(&GameID)
	db.Close()
	if GameID == 0{
		return false
	}else{
		return true
	}
}

//假的假的 發牌會發一堆重複的
func GetTwoCards(MID string) [2]int{
	var GameID int
	GameID = 0;
	var card1 int
	var card2 int
	cards := [2]int{-1, -1}
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	db.QueryRow("SELECT GameID, PlayerCard1, PlayerCard2 FROM sql6131889.GameAction WHERE MID = ? and Cancel = 0", MID ).Scan(&GameID, &card1, &card2)
	if card1 != 0{
		//db.QueryRow("名字 FROM GameAction, 撲克牌參照表 WHERE MID = ? and PlayerCard1 = 編號", MID ).Scan(&card1name)
		//db.QueryRow("名字 FROM GameAction, 撲克牌參照表 WHERE MID = ? and PlayerCard2 = 編號", MID ).Scan(&card2name)
		cards = [2]int{card1, card2}
	}else{
		rand.Seed(time.Now().UTC().UnixNano())
		card1 = 1 + rand.Intn(52)
		card2 = 1 + rand.Intn(52)
		cards = [2]int{card1, card2}
		db.Exec("UPDATE sql6131889.GameAction SET PlayerCard1 = ? , PlayerCard2 = ? WHERE MID = ? AND Cancel = ?", card1, card2, MID, 0) 
	}
	return cards
}

////拿新的兩張手牌
func NewTwoCards(MID string) [2]int{
	var GameID int
	GameID = 0;
	var card1 int
	var card2 int
	cards := [2]int{-1, -1}
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	db.QueryRow("SELECT GameID, PlayerCard1, PlayerCard2 FROM sql6131889.GameAction WHERE MID = ? and Cancel = 0", MID ).Scan(&GameID, &card1, &card2)
	if card1 != 0{
		//db.QueryRow("名字 FROM GameAction, 撲克牌參照表 WHERE MID = ? and PlayerCard1 = 編號", MID ).Scan(&card1name)
		//db.QueryRow("名字 FROM GameAction, 撲克牌參照表 WHERE MID = ? and PlayerCard2 = 編號", MID ).Scan(&card2name)
		rand.Seed(time.Now().UTC().UnixNano())
		card1 = 1 + rand.Intn(52)
		card2 = 1 + rand.Intn(52)
		cards = [2]int{card1, card2}
		db.Exec("UPDATE sql6131889.GameAction SET PlayerCard1 = ?, PlayerCard2 = ?  WHERE MID = ?", card1, card2, MID)
		cards = [2]int{card1, card2}
	}else{
		rand.Seed(time.Now().UTC().UnixNano())
		card1 = 1 + rand.Intn(52)
		card2 = 1 + rand.Intn(52)
		cards = [2]int{card1, card2}
		db.Exec("INSERT INTO sql6131889.GameAction (ID, MID, GameID, PlayerX, Action, PlayerCard1, Cancel, PlayerCard2)VALUES (?, ?, ?, ?, ?, ?, ?, ?)", 123, MID, 321, 11, 2, card1, 0, card2)
	}
	return cards
}


//拿新的五張牌並回傳
func NewFiveCards(GameID int) [5]int{
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	rand.Seed(time.Now().UTC().UnixNano())
	tablecards := [5]int{-1, -1, -1, -1, -1}
	i := 0
	for i < 5{
		tablecards[i] = 1 + rand.Intn(52)
		i = i + 1
	}
	db.Exec("UPDATE sql6131889.Game SET Card1 = ?, Card2 = ?, Card3 = ?, Card4 = ?, Card5 = ? WHERE ID = GameID", tablecards[0], tablecards[1], tablecards[2], tablecards[3], tablecards[4])
	return tablecards
}

//回傳目前牌桌的牌
func GetFiveCards(GameID int) [5]int{
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	tablecards := [5]int{-1, -1, -1, -1, -1}
	db.QueryRow("SELECT Card1, Card2, Card3, Card4, Card5 FROM sql6131889.Game WHERE ID = ?", GameID ).Scan(&tablecards[0], &tablecards[1], &tablecards[2], &tablecards[3], &tablecards[4])
	return tablecards
}

//丟卡片編號數字進來就給你卡片名子，如 Spade A
func GetCardName(cardNo int) string{
	var cardname string
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	db.QueryRow("select 名字 FROM sql6131889.撲克牌參照表 WHERE 編號 = ?", cardNo ).Scan(&cardname)
	return cardname
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

//洗牌啦，不回傳值
func Shuffle(){
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	rand.Seed(time.Now().UTC().UnixNano())
	i := 0
	a := 0
	b := 0
	
	cards := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
		14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
		25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35,
		36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46,
		47, 48, 49, 50, 51, 52}
	
	for i < 52{
		a = rand.Intn(52)
		b = rand.Intn(52) 
		cards[a], cards[b] = Swap(cards[a], cards[b])
		i = i +1 
	}
	i = 0
	for i < 52{
		db.QueryRow("UPDATE sql6131889.洗過的牌堆 SET 牌堆順序 = ? WHERE 編號 = ?", cards[i], i + 1)
		i= i + 1
	}
}

func Swap(x, y int) (int, int) {
	return y, x
}
