package sm

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	"testing"
)

func ExtractUrl(url string) (urls []string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal(err.Error())
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exist := s.Attr("href")
		href = strings.Trim(href, " ")
		href = strings.Replace(href, "http://", "https://", 1)

		if exist {
			if href != "" && href != "#" && href != "javascript:;" {
				//fmt.Println(href, s.Text())
				urls = append(urls, href)
			}
		}
	})
	return urls
}

func TestExtractUrl(t *testing.T) {
	ExtractUrl("http://www.sina.com.cn")
}

func TestNewTrie(t *testing.T) {
	urlMap := make(map[string]bool)
	urls := ExtractUrl("http://www.sina.com.cn")
	trie := NewTrie()
	for _, url := range urls {
		if _, ok := urlMap[url]; !ok {
			if strings.Contains(url, "https://finance.sina.com.cn") {
				//fmt.Println(url)
				trie.Insert([]byte(url), url)
			}
			urlMap[url] = true
		}
	}
	fmt.Println("=============================")
	trie.BFS(func(key []byte, node *Node) {
		if node.IsKey {
			//fmt.Println(string(key), node)
			fmt.Println(node.Value, string(key))
		}
	})

	//k := []byte("https://finance.sina.com.cn/money/forex/hq/USDCNY.shtml")
	//fmt.Println("===========================")
	//for _,idx := range trie.SeekBefore(k) {
	//	fmt.Println(string(k[0:idx]))
	//}
	_, val := trie.Find([]byte("https://finance.sina.com.cn/roll/2019-11-20/doc-iihnzahi2038825.shtml"))
	fmt.Println(val)
}

func TestAppend(t *testing.T) {
	s := make([]byte, 0)
	b := append(s, 56)
	c := append(s, 57)
	fmt.Println(string(s))
	fmt.Println(string(b))
	fmt.Println(string(c))
}
