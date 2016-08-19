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
		if content != nil { // put user profile into database
			db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
			var M string
			db.QueryRow("SELECT MID FROM sql6131889.User WHERE MID = ?", content.From).Scan(&M)
			if M == ""{ // new user
			prof,_ := bot.GetUserProfile([]string{content.From})
			info := prof.Contacts
			bot.SendText([]string{content.From}, "歡迎!")
			bot.SendText([]string{content.From}, "請輸入您的暱稱")
			db.Exec("INSERT INTO sql6131889.User (MID, UserName, UserStatus, UserTitle, UserPicture) VALUES (?, ?, ?, ?, ?)", info[0].MID, info[0].DisplayName, 1, "菜鳥", info[0].PictureURL)
			}
		}
		if content != nil && content.IsMessage && content.ContentType == linebot.ContentTypeText{ // content type : text
			text, _ := content.TextContent()
			prof,_ := bot.GetUserProfile([]string{content.From})
			info := prof.Contacts
			bot.SendText([]string{os.Getenv("mymid")}, "測試\n"+info[0].DisplayName+" :\n"+text.Text) // sent to garylai
			db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
			db.Exec("INSERT INTO database1234.linebottext VALUES (?, ?, ?)", info[0].MID, info[0].DisplayName, text.Text)
			var S string
			db.QueryRow("SELECT Status FROM database1234.linebotuser WHERE MID = ?", content.From).Scan(&S) // get user status
			if S == "default"{
				if text.Text == "!進入房間" { // cheak if enter commands
					db.Exec("UPDATE database1234.linebotuser SET Status = ? WHERE MID = ?", "joining", content.From)
					bot.SendText([]string{content.From}, "請輸入房間號碼:")
					db.Close()
				}else if text.Text == "!開新房間" {
					db.Exec("UPDATE database1234.linebotuser SET Status = ? WHERE MID = ?", "creating", content.From)
					bot.SendText([]string{content.From}, "請輸入新房間名字(純數字):")
				}else if text.Text == "!提示" {
					bot.SendText([]string{content.From}, "哈囉! "+info[0].DisplayName+"!\n您目前位於大廳\n系統指令提示:\n!開新房間\n!進入房間\n")
				}else{
					db.Close()
					bot.SendText([]string{content.From}, "請善用系統指令:\n!提示")
				}
			}else if S == "creating"{
				var rn string
				db.QueryRow("SELECT roomnum FROM database1234.chatroom WHERE roomnum = ?", text.Text).Scan(&rn)
				if rn != ""{
					bot.SendText([]string{content.From}, "已有此房間")
					db.Exec("UPDATE database1234.linebotuser SET Status = ? WHERE MID = ?", "default", content.From)
				}else{
					db.Exec("INSERT INTO database1234.chatroom VALUES (?, ?)", text.Text, content.From)
					bot.SendText([]string{content.From}, "請輸入密碼:")
					db.Exec("UPDATE database1234.linebotuser SET Status = ? WHERE MID = ?", "creatingpw", content.From)
				}
				db.Close()
			}else if S == "creatingpw"{
				db.Exec("UPDATE database1234.chatroom SET roompw = ? WHERE roompw = ?", text.Text, content.From)
				db.Exec("UPDATE database1234.linebotuser SET Status = ? WHERE MID = ?", "default", content.From)
				var rn string
				db.QueryRow("SELECT roomnum FROM database1234.chatroom WHERE roompw = ?", text.Text).Scan(&rn)
				bot.SendText([]string{content.From}, "房間: "+rn+"\n已建立")
				db.Close()
			}else if S == "joining"{
				var pw string
				db.QueryRow("SELECT roompw FROM database1234.chatroom WHERE roomnum = ?", text.Text).Scan(&pw)
				if pw == ""{
					bot.SendText([]string{content.From}, "房間 : "+text.Text+"\n不存在")
					db.Exec("UPDATE database1234.linebotuser SET Status = ? WHERE MID = ?", "default", content.From)
				}else{
					db.Exec("INSERT INTO database1234.chatroomuser VALUES (?, ?, ?)", info[0].MID+"q", info[0].DisplayName, text.Text)
					bot.SendText([]string{content.From}, "請輸入房間密碼:")
					db.Exec("UPDATE database1234.linebotuser SET Status = ? WHERE MID = ?", "enterpw", content.From)
				}
				db.Close()
			}else if S == "enterpw"{
				var rp string
				var rn string
				db.QueryRow("SELECT roomnum FROM database1234.chatroomuser WHERE MID = ?", content.From+"q").Scan(&rn)
				db.QueryRow("SELECT roompw FROM database1234.chatroom WHERE roomnum = ?", rn).Scan(&rp)
				if text.Text == rp{ // correct password
					bot.SendText([]string{content.From}, "進入房間:\n"+rn)
					db.Exec("UPDATE database1234.chatroomuser SET MID = ? WHERE MID = ?", content.From, content.From+"q")
					db.Exec("UPDATE database1234.linebotuser SET Status = ? WHERE MID = ?", "chatting", content.From)
				}else{
					bot.SendText([]string{content.From}, "密碼錯誤")
					db.Exec("DELETE FROM database1234.chatroomuser WHERE MID = ?", content.From+"q")
					db.Exec("UPDATE database1234.linebotuser SET Status = ? WHERE MID = ?", "default", content.From)
				}
				db.Close()
			}else if S == "chatting"{
				if text.Text == "!離開房間"{
					var N string
					db.QueryRow("SELECT roomnum FROM database1234.chatroomuser WHERE MID = ?", content.From).Scan(&N)
					bot.SendText([]string{content.From}, "已離開房間:\n"+N)
					db.Exec("DELETE FROM database1234.chatroomuser WHERE MID = ?", content.From)
					db.Exec("UPDATE database1234.linebotuser SET Status = ? WHERE MID = ?", "default", content.From)
				}else if text.Text == "!提示"{
					var rn string
					db.QueryRow("SELECT roomnum FROM database1234.chatroomuser WHERE MID = ?", content.From).Scan(&rn)
					bot.SendText([]string{content.From}, "哈囉! "+info[0].DisplayName+"!\n您目前位於房間: "+rn+"\n系統指令提示:\n!建立新牌局\n!進入牌局\n!離開牌局\n!離開房間")
				}else if text.Text == "!建立新牌局"{
					var N string
					db.QueryRow("SELECT roomnum FROM database1234.chatroomuser WHERE MID = ?", content.From).Scan(&N)
					row,_ := db.Query("SELECT MID FROM database1234.chatroomuser WHERE roomnum = ?", N)
					for row.Next() {
						var mid1 string
						row.Scan(&mid1)
						bot.SendText([]string{mid1}, "玩家: "+info[0].DisplayName+" 建立新牌局")
						bot.SendText([]string{mid1}, "玩家: "+info[0].DisplayName+" 進入牌局")
					}
					//把房間state改成遊戲中
					//把玩家state改成playing //S == "playing"
				}else if text.Text == "!進入牌局"{
					var N string
					db.QueryRow("SELECT roomnum FROM database1234.chatroomuser WHERE MID = ?", content.From).Scan(&N)
					row,_ := db.Query("SELECT MID FROM database1234.chatroomuser WHERE roomnum = ?", N)
					for row.Next() {
						var mid1 string
						row.Scan(&mid1)
						bot.SendText([]string{mid1}, "玩家: "+info[0].DisplayName+" 進入牌局")
					}
					//把玩家state改成playing //S == "playing"
				}else{
					var N string
					db.QueryRow("SELECT roomnum FROM database1234.chatroomuser WHERE MID = ?", content.From).Scan(&N)
					row,_ := db.Query("SELECT MID FROM database1234.chatroomuser WHERE roomnum = ?", N)
					for row.Next() {
						var mid1 string
						row.Scan(&mid1)
						if mid1 != content.From{
							bot.SendText([]string{mid1}, info[0].DisplayName+":\n"+text.Text)
						}
					}
				}
				db.Close()
			}
		}else if content != nil && content.ContentType == linebot.ContentTypeSticker{ // content type : sticker
			sticker, _ := content.StickerContent()
			prof,_ := bot.GetUserProfile([]string{content.From})
			info := prof.Contacts
			bot.SendSticker([]string{content.From}, 7, 1, 100)
			bot.SendText([]string{os.Getenv("mymid")}, info[0].DisplayName+" sent a sticker") // sent to garylai
			db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
			db.Exec("INSERT INTO database1234.linebotsticker VALUES (?, ?, ?, ?, ?)", info[0].MID, info[0].DisplayName, sticker.PackageID, sticker.ID, sticker.Version)
			db.Close()
		}
		
		
	}
}
