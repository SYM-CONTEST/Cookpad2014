package crawler

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"log"
	_ "log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"encoding/json"
)

type Anniversary struct {
	Tweets []anaconda.Tweet
}

func (a Anniversary) contains(t anaconda.Tweet) bool {
	for _, myT := range a.Tweets {
		if myT.Id == t.Id {
			return true
		}
	}
	return false
}

func (a Anniversary) date() time.Time {
	t, e := a.Tweets[len(a.Tweets)-1].CreatedAtTime()
	failIfNeeded(e)
	return t
}

// ex: 「今日はお二人の」
func (a Anniversary) CreateFirstMessage() string {
	str := "今日は"
	names := a.Names()
	if len(names) <= 2 {
		str += "お二人の"
	} else {
		str += "皆さんの"
	}
	return str
}

// ex: 「りんご1周年記念日」
// ランダム要素が入っているので、毎回結果が異なります。
func (a Anniversary) CreateSecondMessage() (string, int64) {
	parser := Parser{}
	nouns := parser.ParseToNouns(a.tweetStrings())
	nouns = parser.filterNoise(nouns, a, 2)
	if len(nouns) < 1 {
		return "", -1
	}
	rand.Seed(time.Now().Unix())
	index := 0
	nlen := len(nouns) - 1
	if nlen > 0 {
		index = rand.Intn(len(nouns) - 1)
	}
	n := nouns[index]
	str := n
	str += a.createDateMessage()
	str += "記念日"
	return str, a.ownerTweet(n).Id
}

func (a Anniversary) ownerTweet(noun string) anaconda.Tweet {
	log.Println("ownerTweet")
	for i := len(a.Tweets) - 1; i >= 0; i-- {
		t := a.Tweets[i]
		log.Println("t: ", t)
		if strings.Contains(t.Text, noun) {
			return t
		}
	}
	failIfNeeded(nil)
	return a.Tweets[0]
}

// Tweet用のフルメッセージ(ex: 「@a @bさん、今日はお二人のりんご1周年記念日です！」
func (a Anniversary) CreateFullMessage(first string, second string) string {
	str := strings.Join(a.namesWithAtmark(), " ") + "さん、"
	str += first
	str += second
	str += "です！"
	return str
}

type EmbedResponse struct {
	Html string `json:"html"`
}

func (a Anniversary) GetEmbededHTML(statusId int64) string {
	url := fmt.Sprintf("https://api.twitter.com/1/statuses/oembed.json?id=%d&align=center", statusId)
	r, e := http.Get(url)
	failIfNeeded(e)
	body, e := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	var er EmbedResponse
	e2 := json.Unmarshal(body, &er)
	failIfNeeded(e2)
	return er.Html
}

func (a Anniversary) createDateMessage() string {
	d := a.date()
	now := time.Now()
	if d.Year() < now.Year() && d.Month() == now.Month() && d.Day() == now.Day() {
		return fmt.Sprintf("から%d周年", now.Year()-d.Year())
	}
	if d.Year() == now.Year() && d.Month() < now.Month() && d.Day() == now.Day() {
		return fmt.Sprintf("から%dヶ月", now.Month()-d.Month())
	}
	if d.Weekday() == now.Weekday() && d.Day() != now.Day() {
		return fmt.Sprintf("から%d週間", (now.YearDay()-d.YearDay())/7)
	}
	if d.Year() == now.Year() && d.Month() == now.Month() && d.Day() == now.Day() {
		return ""
	}
	return fmt.Sprintf("から%d日", now.YearDay()-d.YearDay())
}

func (a Anniversary) namesWithAtmark() []string {
	names := a.Names()
	r := make([]string, 0, len(names))
	for _, n := range names {
		r = append(r, "@"+n)
	}
	return r
}

func (a Anniversary) Names() []string {
	r := make([]string, 0, 2)
	for _, t := range a.Tweets {
		name := t.User.ScreenName
		if containsString(r, name) {
			continue
		}
		r = append(r, name)
	}
	for _, t := range a.Tweets {
		name := t.InReplyToScreenName
		if containsString(r, name) || name == "" {
			continue
		}
		r = append(r, name)
	}
	return r
}

func (a Anniversary) containsNealyName(target string) bool {
	for _, n := range a.Names() {
		if strings.Contains(n, target) {
			return true
		}
	}
	return false
}

func containsInAnivs(anivs []Anniversary, t anaconda.Tweet) bool {
	for _, a := range anivs {
		if a.contains(t) {
			return true
		}
	}
	return false
}

func (a Anniversary) tweetStrings() []string {
	r := make([]string, 0, len(a.Tweets))
	for _, t := range a.Tweets {
		r = append(r, t.Text)
	}
	return r
}
