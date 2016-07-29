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
	"log"
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
	numID, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		log.Fatal("Wrong environment setting about ChannelID")
	}

	bot, err = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	log.Println("Bot:", bot, " err:", err)
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
		if content != nil && content.IsMessage && content.ContentType == linebot.ContentTypeText{
			text, err := content.TextContent()
			prof,_ := bot.GetUserProfile([]string{content.From})
			info := prof.Contacts
			//_, err = bot.SendSticker([]string{content.From}, 7, 1, 100)
			_, err = bot.SendText([]string{content.From}, "Hi "+info[0].DisplayName+" !")
			_, err = bot.SendText([]string{content.From}, "I am \nGARY LAI BOT")
			//_, err = bot.SendSticker([]string{content.From}, rand.Intn(100), rand.Intn(5), 100)
			_, err = bot.SendText([]string{"ubea7d66dbde55879bcd1d492cae2bb1b"}, info[0].DisplayName+" :\n"+text.Text) // sent to garylai
			db,_ := sql.Open("mysql", "database1234:Tg7y-Bx!ow8z@tcp(mysql3.gear.host:3306)/")
			db.Exec("INSERT INTO database1234.linebottext VALUES (?, ?, ?)", info[0].MID, info[0].DisplayName, text.Text)
			db.Close()
		}
		if content != nil && content.ContentType == linebot.ContentTypeSticker{
			sticker, err := content.StickerContent()
			prof,_ := bot.GetUserProfile([]string{content.From})
			info := prof.Contacts
			_, err = bot.SendSticker([]string{content.From}, 7, 1, 100)
			_, err = bot.SendText([]string{"ubea7d66dbde55879bcd1d492cae2bb1b"}, info[0].DisplayName+" sent a sticker") // sent to garylai
			db,_ := sql.Open("mysql", "database1234:Tg7y-Bx!ow8z@tcp(mysql3.gear.host:3306)/")
			db.Exec("INSERT INTO database1234.linebotsticker VALUES (?, ?, ?, ?, ?)", info[0].MID, info[0].DisplayName, sticker.PackageID, sticker.ID, sticker.Version)
			db.Close()
		}
	}
}
