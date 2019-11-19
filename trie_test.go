package sm

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"testing"
)

func ExtractUrl(url string) (urls []string) {
	res, err := http.Get("http://metalsucks.net")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("aaaaa")
	fmt.Println("xxxx")
	doc.Find("/html/body/a").Each(func(i int, s *goquery.Selection) {
		href, exist := s.Attr("href")
		if exist {
			fmt.Println(href)
			urls = append(urls, href)
		}
	})
	return urls
}

func TestExtractUrl(t *testing.T) {
	ExtractUrl("http://www.sina.com.cn")
}

func TestNewTrie(t *testing.T) {
	urls := ExtractUrl("http://www.sina.com.cn")
	trie := NewTrie()
	for _, url := range urls {
		trie.Insert([]byte(url), true)
	}

}
