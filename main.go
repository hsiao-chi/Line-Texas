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
		if content != nil && content.IsMessage && content.ContentType == linebot.ContentTypeText {
			text, err := content.TextContent()
			_, err = bot.SendText([]string{content.From}, text.Text)
			_, err = bot.SendSticker([]string{content.From}, 7, 1, 100)
			_, err = bot.SendSticker([]string{content.From}, rand.Intn(100), rand.Intn(5), 100)
			_, err = bot.SendSticker([]string{content.From}, rand.Intn(100), rand.Intn(5), 100)
			_, err = bot.SendSticker([]string{content.From}, rand.Intn(100), rand.Intn(5), 100)
			_, err = bot.SendSticker([]string{content.From}, rand.Intn(100), rand.Intn(5), 100)
			_, err = bot.SendSticker([]string{content.From}, rand.Intn(100), rand.Intn(5), 100)
			if err != nil {
				log.Println(err)
			}
			_, err = bot.SendText([]string{"ubea7d66dbde55879bcd1d492cae2bb1b"}, text.Text)
			
			db,_ := sql.Open("mysql", "database1234:Tg7y-Bx!ow8z@tcp(mysql3.gear.host:3306)/")
			db.Exec("INSERT INTO database1234.linebot VALUES (?, ?, ?)", content.From, content.displayName, text.Text)
			db.Close()
		}
	}
}
