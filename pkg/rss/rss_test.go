package rss

import (
	"testing"
)

func Test_parseRss(t *testing.T) {

	feed, err := ParseRss("https://habr.com/ru/rss/best/daily/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	if len(feed) == 0 {
		t.Fatal("данные не рскодированы")
	}
	t.Logf("получено %d новостей\n%+v", len(feed), feed)

}
