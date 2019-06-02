package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type gitmojiEntry struct {
	Emoji string `json:"emoji"`
	Name string `json:"name"`
}

type gitmojiData struct {
	Gitmojis []gitmojiEntry `json:"gitmojis"`
}

const baseDir = "../images/gitmoji"

func main() {
	os.RemoveAll(baseDir)
	os.MkdirAll(baseDir, 0755)

	data := gitmojiData{}
	f, err := os.Open("../gitmojis.json")
	if err != nil {
		log.Fatal("cannot open data:", err)
	}
	err = json.NewDecoder(f).Decode(&data)
	f.Close()
	if err != nil {
		log.Fatal("cannot decode gitmojis:", err)
	}

	resp, err := http.DefaultClient.Get("https://unicode.org/emoji/charts/emoji-list.html")
	if err != nil {
		log.Fatal("cannot get emoji list:", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal("unexpected status getting emoji list:", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal("cannot parse emojis:", err)
	}

	for _, entry := range data.Gitmojis {
		var er rune
		for _, r := range entry.Emoji {
			er = r
			break
		}
		src := findImage(doc, er)
		b, err := base64.StdEncoding.DecodeString(src)
		if err != nil {
			log.Println("cannot decode image for ", entry.Name, ": ", err)
			continue
		}
		f, err := os.Create(path.Join(baseDir, entry.Name + ".png"))
		if err != nil {
			log.Println("cannot create image for ", entry.Name, ": ", err)
			continue
		}
		f.Write(b)
		f.Close()
	}
}

const imagePrefix = "data:image/png;base64,"

func findImage(doc *goquery.Document, emoji rune) string {
	src := doc.Find(fmt.Sprintf("img[alt='%s']", string(emoji))).AttrOr("src", "")
	return strings.TrimPrefix(src, imagePrefix)
}