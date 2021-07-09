package main

import (
	"container/list"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func checkDataInSlice(data interface{}, mylist []string) int {
	for _, lItem := range mylist {
		if data == lItem {
			return 1
		}
	}
	return 0
}

func htmlPublishData(essay *list.List) {
	htmlTags := []string{"a", "p", "div", "title", "span", "br", "table", "img", "blockquote", "li", "script", "tbody", "b"}

	ignoreContent := []string{"சினிமா ", "விளையாட்டு", "கோயில்கள்", "என்.ஆர்.ஐ", "வர்த்தகம்", "கல்விமலர்", "புத்தகங்கள் ", "iPaper",
		" Newspaper Subscription ", "போட்டோக்கள்! ", " 2021", "          செய்திகள்-2020", " பக்கங்கள்", "தினமலர் டெலிகிராம் சேனலில் பார்க்கலாம்",
		"  Telegram Channel for FREE",
		" ", "  ", "		", "	", "\n", "		", "	"}
	//ignoreContent := []string{"சினிமா ", "விளையாட்டு"}
	skipDataList := append(htmlTags, ignoreContent...)
	// Iterate through list and print its contents.
	for e := essay.Front(); e != nil; e = e.Next() {
		if checkDataInSlice(e.Value, skipDataList) == 0 {
			fmt.Println(e.Value)
			continue
		}
	}
}

func html_page_parser(tokenizer *html.Tokenizer, essay *list.List) {
	isPrintData := 0
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				//end of the file, break out of the loop
				break
			}
			fmt.Printf("error tokenizing HTML: %v", tokenizer.Err())
		}
		token := tokenizer.Token()
		switch tokenType {
		case html.ErrorToken:
			fmt.Printf("Case: ERROR Token")
			break
		case html.TextToken:
			if isPrintData == 1 {
				essay.PushBack(token.Data)
			}
		case html.StartTagToken, html.SelfClosingTagToken:
			if token.Data == "p" || token.Data == "title" || token.Data == "h1" || token.Data == "h2" || token.Data == "h2" || token.Data == "h4" {
				isPrintData = 1
			}
		case html.EndTagToken:
			if token.Data == "p" || token.Data == "title" || token.Data == "h1" || token.Data == "h2" || token.Data == "h2" || token.Data == "h4" {
				isPrintData = 0
			}
		default:
			continue
		}
	}
}

func main() {
	resp, err := http.Get("https://tamil.oneindia.com/news/chennai/tamilnadu-politics-why-alternative-parties-join-dmk-and-aiadmk-426525.html")

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	//check response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("response status code was %d\n", resp.StatusCode)
	}

	//check response content type
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		fmt.Printf("response content type was %s not text/html\n", ctype)
	}

	essay := list.New()
	tokenizer := html.NewTokenizer(resp.Body)
	html_page_parser(tokenizer, essay)
	htmlPublishData(essay)
}
