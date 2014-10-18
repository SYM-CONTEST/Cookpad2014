package crawler

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"log"
	"strings"
)

func failIfNeeded(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func printResults(ts []anaconda.Tweet) {
	for _, t := range ts {
		printResult(t)
	}
}

func printResult(t anaconda.Tweet) {
	fmt.Println("InReplyToStatusID: ", t.InReplyToStatusID)
	p := Parser{}
	tt := p.parseNouns(t.Text)
	wc := wordCount(tt)
	for v, k := range wc {
		fmt.Println(k, v)
	}
}

func wordCount(ss []string) map[string]int {
	res := make(map[string]int)
	for _, str := range ss {
		res[strings.ToLower(str)]++
	}
	return res
}
