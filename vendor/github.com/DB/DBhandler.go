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
		db.Exec("INSERT INTO sql6131889.GameAction (ID, MID, GameID, PlayerX, Action, PlayerCard1, Cancel, PlayerCard2)VALUES (?, ?, ?, ?, ?, ?, ?, ?)", 123, MID, 321, 11, 2, card1, 0, card2)
	}
	return cards
}

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
