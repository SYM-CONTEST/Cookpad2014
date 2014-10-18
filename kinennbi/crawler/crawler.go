package crawler

import (
	"bitbucket.org/rerofumi/mecab"
	"fmt"
	"strings"
	"github.com/ChimeraCoder/anaconda"
	"log"
	"net/url"
)

func AnalyzeAnniversary(api *anaconda.TwitterApi) {
	vs := url.Values{}
	ts, e := api.GetMentionsTimeline(vs)
	failIfNeeded(e)
	printResult(ts)
}

func AnalyzeHoge(api *anaconda.TwitterApi) {
	vs := url.Values{}
	vs.Set("count", "200")
	ts, e := api.GetUserTimeline(vs)
	failIfNeeded(e)
	printResult(ts)
}

func failIfNeeded(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func printResult(ts []anaconda.Tweet) {
	for _, t := range ts {
		fmt.Println(t.InReplyToStatusID)
		tt := parseNouns(t.Text)
		wc := wordCount(tt)

		for v, k := range wc {
			fmt.Println(k, v)
		}
	}
}

func parseNouns(s string) []string {
	p, e := mecab.Parse(s)
	if e != nil {
		log.Fatalln(e)
		return nil
	}
	r := make([]string, 0, len(p))
	for _, a := range p {
		v := strings.Fields(a)
		if strings.Split(v[1], ",")[0] == "名詞" {
			r = append(r, v[0])
		}
	}
	return r
}

func wordCount(ss []string) map[string]int {
	res := make(map[string]int)
	for _, str := range ss {
		res[strings.ToLower(str)]++
	}
	return res
}
