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
	"math/rand"
	"github.com/line/line-bot-sdk-go/linebot"
	//_ "github.com/go-sql-driver/mysql"
    "database/sql"
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
	db, err := sql.Open("mysql","database1234:Tg7y-Bx!ow8z@tcp(mysql3.gear.host)/linebot")
	a := "test1"
	db.Query("INSERT INTO linebot (request, awnser) VALUES (a, a)")
	db.Close()
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
		if content != nil && content.IsMessage && content.ContentType == linebot.ContentTypeText {
			text, err := content.TextContent()
			_, err = bot.SendText([]string{content.From}, text.Text+"呦")
			_, err = bot.SendSticker([]string{content.From}, 7, 1, 100)
			_, err = bot.SendSticker([]string{content.From}, rand.Intn(100), rand.Intn(5), 100)
			_, err = bot.SendSticker([]string{content.From}, rand.Intn(100), rand.Intn(5), 100)
			_, err = bot.SendSticker([]string{content.From}, rand.Intn(100), rand.Intn(5), 100)
			_, err = bot.SendSticker([]string{content.From}, rand.Intn(100), rand.Intn(5), 100)
			_, err = bot.SendSticker([]string{content.From}, rand.Intn(100), rand.Intn(5), 100)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
