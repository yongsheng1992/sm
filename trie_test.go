package sm

import (
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
	trie.BFS(func(key []byte, node *Node) {
		if node.IsKey {
			if node.Value != string(key) {
				t.Error("bfs failed")
			}
		}
	})
}

func TestTrie_SeekBefore(t *testing.T) {
	urls := ExtractUrl("http://www.sina.com.cn")
	trie := NewTrie()
	for _, url := range urls {
		trie.Insert([]byte(url), url)
	}

	it := trie.SeekAfter([]byte("https://finance.sina.com.cn"))
	for it.HasNext() {
		key, node := it.Next()
		if node.IsKey {
			if string(key) != node.Value {
				t.Error("failed")
			}
		}
	}
}
