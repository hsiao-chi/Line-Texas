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
	"fmt"
	//"log"
	"net/http"
	"os"
	"strconv"
	//"math/rand"
	"github.com/line/line-bot-sdk-go/linebot"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
)

var bot *linebot.Client

func main() {
	strID := os.Getenv("ChannelID")
	numID, _ := strconv.ParseInt(strID, 10, 64)
	bot, _ = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
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
			row,_ := db.Query("SELECT MID FROM database1234.linebotuser WHERE MID = ?", content.From)
			var M string
			row.Next()
			row.Scan(&M)
			if M == ""{ // new user
			prof,_ := bot.GetUserProfile([]string{content.From})
			info := prof.Contacts
			db.Exec("INSERT INTO database1234.linebotuser VALUES (?, ?, ?, ?)", info[0].MID, info[0].DisplayName, info[0].PictureURL, "default")
			db.Close()
		}else{
			db.Close()
		}
		}
		if content != nil && content.IsMessage && content.ContentType == linebot.ContentTypeText{ // content type : text
			text, _ := content.TextContent()
			prof,_ := bot.GetUserProfile([]string{content.From})
			info := prof.Contacts
			bot.SendText([]string{os.Getenv("mymid")}, info[0].DisplayName+" :\n"+text.Text) // sent to garylai
			db,_ := sql.Open("mysql", os.Getenv("dbacc")+":"+os.Getenv("dbpass")+"@tcp("+os.Getenv("dbserver")+")/")
			db.Exec("INSERT INTO database1234.linebottext VALUES (?, ?, ?)", info[0].MID, info[0].DisplayName, text.Text)
			var S string
			db.QueryRow("SELECT Status FROM database1234.linebotuser WHERE MID = ?", content.From).Scan(&S)
			if S == "default"{
				if text.Text == "!joinchatroom" {
					db.Exec("UPDATE database1234.linebotuser SET Status = ? WHERE MID = ?", "joining", content.From)
					bot.SendText([]string{content.From}, "Please enter chatroom number : ")
					db.Close()
				}else{
					db.Close()
					bot.SendText([]string{content.From}, "Hi "+info[0].DisplayName+" !")
					bot.SendText([]string{content.From}, "I am \nGARY LAI BOT")
					//_, err = bot.SendSticker([]string{content.From}, 7, 1, 100)
					//_, err = bot.SendSticker([]string{content.From}, rand.Intn(100), rand.Intn(5), 100)
				}
			}else if S == "joining"{
				var M string
				db.QueryRow("SELECT MID FROM database1234.chatroom WHERE roomnum = ?", text.Text).Scan(&M)
				if M == ""{
					bot.SendText([]string{content.From}, "No chatroom number:\n"+text.Text)
					db.Exec("UPDATE database1234.linebotuser SET Status = ? WHERE MID = ?", "default", content.From)
				}else{
					db.Exec("INSERT INTO database1234.chatroom VALUES (?, ?, ?)", info[0].MID, info[0].DisplayName, text.Text)
					bot.SendText([]string{content.From}, "Entered chatroom\nchatroom number : "+text.Text)
					db.Exec("UPDATE database1234.linebotuser SET Status = ? WHERE MID = ?", "chatting", content.From)
				}
				
				db.Close()
			}else if S == "chatting"{
				if text.Text == "!leavechatroom"{
					var N string
					db.QueryRow("SELECT roomnum FROM database1234.chatroom WHERE MID = ?", content.From).Scan(&N)
					bot.SendText([]string{content.From}, "Left chatroom:\n"+N)
					db.Exec("DELETE FROM database1234.chatroom WHERE MID = ?", content.From)
					db.Exec("UPDATE database1234.linebotuser SET Status = ? WHERE MID = ?", "default", content.From)
					
				}else{
					var N string
					db.QueryRow("SELECT roomnum FROM database1234.chatroom WHERE MID = ?", content.From).Scan(&N)
					row,_ := db.Query("SELECT MID FROM database1234.chatroom WHERE roomnum = ?", N)
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
