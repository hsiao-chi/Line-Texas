package DB

import (
	"os"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
)

func BigCombi (mid string) ([4]int, [13]int){ // Biggest Combination
	db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
	// create an 4 * 13 array
	deck := [4][13]int{}
	for i := 0;i < 4;i++{
		for j:= 0;j < 13;j++{
			deck[i][j] = 0
		}
	}
	// create suits[4] 
	suits := [4]int{}
	for i := 0;i < 4;i++{
		suits[i] = 0
	}
	// create nums[13]
	nums := [13]int{}
	for i := 0;i < 13;i++{
		nums[i] = 0
	}
	// create cards[7]
	cards := [7]int{}
	// get players cards * 2
	db.QueryRow("SELECT PlayerCard1 PlayerCard2 FROM sql6131889.GameAction WHERE MID = ?", mid).Scan(&cards[0], &cards[1])
	// get desktop cards * 5
	var gid string
	db.QueryRow("SELECT GameID FROM sql6131889.GameAction WHERE MID = ?", mid).Scan(&gid)
	db.QueryRow("SELECT Card1 Card2 Card3 Card4 Card5 FROM sql6131889.Game WHERE ID = ?", gid).Scan(&cards[2], &cards[3], &cards[4], &cards[5], &cards[6])
	// value into arrays
	for i := 0;i < 7;i++{ // cards into deck
		cards[i] -= 1
		deck[cards[i] / 13][cards[i] % 13] += 1
	}
	for i := 0;i < 4;i++{ // deck into suits and nums
		for j:= 0;j < 13;j++{
			suits[i] += deck[i][j]
			nums[j] += deck[i][j]
		}
	}
	db.Close()
	temp := nums[0]
	for i := 0;i < 12;i++{
		nums[i] = nums[i+1]
	}
	nums[12] = temp
	return suits, nums
}