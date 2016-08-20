// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main 

import (
	//"fmt"
	//"log"
	//"math/rand" rand.Intn(100)
	"net/http"
	"os"
	"strconv"
	"github.com/line/line-bot-sdk-go/linebot"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"github.com/DB"
)

var bot *linebot.Client

func main() {
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64) // string to integer
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	received, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, result := range received.Results {
		content := result.Content()
		db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
		prof,_ := bot.GetUserProfile([]string{content.From})
		info := prof.Contacts
		if content != nil {
			var M string
			db.QueryRow("SELECT MID FROM sql6131889.User WHERE MID = ?", content.From).Scan(&M)
			if M == ""{ // new user
			bot.SendText([]string{content.From}, "歡迎光臨LineBot遊戲機器人!") // put user profile into database
			db.Exec("INSERT INTO sql6131889.User (MID, UserName, UserStatus, UserTitle, UserPicture) VALUES (?, ?, ?, ?, ?)", info[0].MID, info[0].DisplayName, 10, "菜鳥", info[0].PictureURL)
			bot.SendText([]string{content.From}, "請設定您的暱稱:")
			db.Exec("UPDATE sql6131889.User SET UserStatus = ? WHERE MID = ?", 400, content.From)
			bot.SendText([]string{content.From}, "請善用指令提示:\n!指令")

			}
			if content.ContentType == linebot.ContentTypeText{ // content type : text
				text, _ := content.TextContent()
				bot.SendText([]string{os.Getenv("mymid")}, info[0].DisplayName+" :\n"+text.Text) // sent to tester
				var nn string
				db.QueryRow("SELECT UserNickName FROM sql6131889.User WHERE MID = ?", content.From).Scan(&nn)
				db.Exec("INSERT INTO sql6131889.text (MID, Text)VALUES (?, ?)", info[0].MID, text.Text)
				var S int
				db.QueryRow("SELECT UserStatus FROM sql6131889.User WHERE MID = ?", content.From).Scan(&S) // get user status
				if S == 10{
					if text.Text == "!加入房間" { // cheak if enter commands
						db.Exec("UPDATE sql6131889.User SET UserStatus = ? WHERE MID = ?", 11, content.From)
						bot.SendText([]string{content.From}, "請輸入房間名稱:")
					}else if text.Text == "!建立房間" {
						db.Exec("UPDATE sql6131889.User SET UserStatus = ? WHERE MID = ?", 12, content.From)
						bot.SendText([]string{content.From}, "請輸入房間名稱:")
					}else if text.Text == "!更改暱稱"{
						db.Exec("UPDATE sql6131889.User SET UserStatus = ? WHERE MID = ?", 400, content.From)
						bot.SendText([]string{content.From}, "請輸入新暱稱:")
					}else if text.Text == "!指令"{
						bot.SendText([]string{content.From}, "哈囉! "+nn+"!\n您現在位於 大廳\n可用指令為:\n!建立房間\n!加入房間\n!更改暱稱")
					}else{
						bot.SendText([]string{content.From}, "請善用指令提示:\n!指令")
					}
				}else if S == 12{
					var rn string
					db.QueryRow("SELECT RoomName FROM sql6131889.Room WHERE RoomName = ?", text.Text).Scan(&rn)
					if rn != ""{
						bot.SendText([]string{content.From}, "房間名稱已重複")
						db.Exec("UPDATE sql6131889.User SET UserStatus = ? WHERE MID = ?", 10, content.From)
					}else{
						db.Exec("INSERT INTO sql6131889.Room (RoomName, RoomPass) VALUES (?, ?)", text.Text, content.From)
						bot.SendText([]string{content.From}, "請設定房間密碼:")
						db.Exec("UPDATE sql6131889.User SET UserStatus = ? WHERE MID = ?", 13, content.From)
					}
				}else if S == 13{
					db.Exec("UPDATE sql6131889.Room SET RoomPass = ? WHERE RoomPass = ?", text.Text, content.From)
					db.Exec("UPDATE sql6131889.User SET UserStatus = ? WHERE MID = ?", 10, content.From)
					var rn string
					db.QueryRow("SELECT RoomName FROM sql6131889.Room WHERE RoomPass = ?", text.Text).Scan(&rn)
					bot.SendText([]string{content.From}, "房間 "+rn+" 已建立")
					db.Exec("UPDATE sql6131889.User SET UserRoom = ? WHERE MID = ?", rn, content.From)
					db.Exec("UPDATE sql6131889.User SET UserStatus = ? WHERE MID = ?", 1000, content.From)
					bot.SendText([]string{content.From}, "您已進入房間 "+rn)
				}else if S == 11{
					var pw string
					db.QueryRow("SELECT RoomPass FROM sql6131889.Room WHERE RoomName = ?", text.Text).Scan(&pw)
					if pw == ""{
						bot.SendText([]string{content.From}, "房間 "+text.Text+" 不存在")
						db.Exec("UPDATE sql6131889.User SET UserStatus = ? WHERE MID = ?", 10, content.From)
					}else{
						db.Exec("UPDATE sql6131889.User SET UserStatus = ? WHERE MID = ?", 14, content.From)
						db.Exec("UPDATE sql6131889.User SET UserRoom = ? WHERE MID = ?", text.Text, content.From)
						bot.SendText([]string{content.From}, "請輸入房間密碼:")
					}
				}else if S == 14{
					var rp string
					var rn string
					db.QueryRow("SELECT UserRoom FROM sql6131889.User WHERE MID = ?", content.From).Scan(&rn)
					db.QueryRow("SELECT RoomPass FROM sql6131889.Room WHERE RoomName = ?", rn).Scan(&rp)
					if text.Text == rp{ // correct password
						bot.SendText([]string{content.From}, "進入房間 "+rn)
						db.Exec("UPDATE sql6131889.User SET UserStatus = ? WHERE MID = ?", 1000, content.From)
					}else{
						bot.SendText([]string{content.From}, "密碼錯誤")
						db.Exec("UPDATE sql6131889.User SET UserStatus = ? WHERE MID = ?", 10, content.From)
					}
				}else if S == 1000{
					var cards [2]int
					cards = DB.GetTwoCards(content.From)
					if text.Text == "!手牌"{
						c1 := DB.GetCardName(cards[0])
						c2 := DB.GetCardName(cards[1])
						bot.SendText([]string{content.From}, "您的手牌為：\n" + c1 + "\n" + c2)
					}
					if text.Text == "!離開房間"{
						var R string
						db.QueryRow("SELECT UserRoom FROM sql6131889.User WHERE MID = ?", content.From).Scan(&R)
						var playerInGame string
						db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE MID = ?", content.From).Scan(&playerInGame)
						if playerInGame != "" {
							DB.CancelGameAction(content.From)
							DB.CancelGame(content.From)
							bot.SendText([]string{content.From}, "結束遊戲...")
						}
						bot.SendText([]string{content.From}, "已離開房間: "+R)
						db.Exec("UPDATE sql6131889.User SET UserStatus = ? WHERE MID = ?", 10, content.From)
						db.Exec("UPDATE sql6131889.User SET UserRoom = ? WHERE MID = ?", 1000, content.From)
						}else if text.Text == "!指令"{
							DB.InRoomInst(content.From)
						}else if text.Text == "!新遊戲"{
							DB.InRoomNewGame(content.From)
						}else if text.Text == "!加入遊戲"{
							DB.InRoomJoinGame(content.From)
						}else if text.Text == "!開始遊戲"{
							DB.InRoomStartGame(content.From)
						}else if text.Text == "!結束遊戲"{
							var playerInGame string
							db.QueryRow("SELECT MID FROM sql6131889.GameAction WHERE MID = ?", content.From).Scan(&playerInGame)
							if playerInGame != "" {
								DB.CancelGameAction(content.From)
								DB.CancelGame(content.From)
							}else{
								bot.SendText([]string{content.From}, "您的狀態為 閒置中")
							}
						}else{
							var R string
							db.QueryRow("SELECT UserRoom FROM sql6131889.User WHERE MID = ?", content.From).Scan(&R)
							row,_ := db.Query("SELECT MID FROM sql6131889.User WHERE UserRoom = ? AND UserStatus = ?", R, 1000)
							for row.Next() {
								var mid1 string
								row.Scan(&mid1)
								if mid1 != content.From{
									bot.SendText([]string{mid1}, nn+":\n"+text.Text)
								}
							}
						}
				}else if S == 400{
					db.Exec("UPDATE sql6131889.User SET UserNickName = ? WHERE MID = ?", text.Text, content.From)
					var temp string
					db.QueryRow("SELECT UserNickName FROM sql6131889.User WHERE MID = ?", content.From).Scan(&temp)
					bot.SendText([]string{content.From}, "您的暱稱更新為 "+temp)
					db.Exec("UPDATE sql6131889.User SET UserStatus = ? WHERE MID = ?", 10, content.From)
				}
			}else if content.ContentType == linebot.ContentTypeSticker{ // content type : sticker
			sticker, _ := content.StickerContent()
			bot.SendSticker([]string{content.From}, 7, 1, 100)
			bot.SendText([]string{os.Getenv("mymid")}, info[0].DisplayName+" sent a sticker") // sent to tester
			db.Exec("INSERT INTO sql6131889.Stiker (MID, PackageID, StickerID, Version)VALUES (?, ?, ?, ?)", info[0].MID, sticker.PackageID, sticker.ID, sticker.Version)
			}
		}
		db.Close()
	}
}
