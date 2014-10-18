package crawler

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	_ "log"
	"math/rand"
	"strings"
	"time"
	"log"
)

type Aniversary struct {
	Tweets []anaconda.Tweet
}

func (a Aniversary) contains(t anaconda.Tweet) bool {
	for _, myT := range a.Tweets {
		if myT.Id == t.Id {
			return true
		}
	}
	return false
}

func (a Aniversary) date() time.Time {
	t, e := a.Tweets[len(a.Tweets)-1].CreatedAtTime()
	failIfNeeded(e)
	return t
}

func (a Aniversary) createMessage() string {
	str := strings.Join(a.namesWithAtmark(), " ") + "さん、"
	str += "今日は"
	names := a.names()
	if len(names) <= 2 {
		str += "お二人の"
	} else {
		str += "皆さんの"
	}
	parser := Parser{}
	nouns := parser.ParseToNouns(a.tweetStrings())
	nouns = parser.filterNoise(nouns, a, 2)
	rand.Seed(time.Now().Unix())
	index := 0
	nlen := len(nouns)-1
	if nlen > 0 {
		index = rand.Intn(len(nouns)-1)
	}
	str += nouns[index]
	str += a.createDateMessage()
	str += "記念日です！"
	return str
}

func (a Aniversary) createDateMessage() string {
	d := a.date()
	now := time.Now()
	if d.Year() < now.Year() && d.Month() == now.Month() && d.Day() == now.Day() {
		return fmt.Sprintf("から%d周年", now.Year()-d.Year())
	}
	if d.Year() == now.Year() && d.Month() < now.Month() && d.Day() == now.Day() {
		return fmt.Sprintf("から%dヶ月", now.Month()-d.Month())
	}
	if d.Weekday() == now.Weekday() && d.Day() != now.Day() {
		return fmt.Sprintf("から%d週間", (now.YearDay()-d.YearDay()) / 7)
	}
	if d.Year() == now.Year() && d.Month() == now.Month() && d.Day() == now.Day() {
		return ""
	}
	return fmt.Sprintf("から%d日", now.YearDay()-d.YearDay())
}

func (a Aniversary) namesWithAtmark() []string {
	names := a.names()
	r := make([]string, 0, len(names))
	for _, n := range names {
		r = append(r, "@"+n)
	}
	return r
}

func (a Aniversary) names() []string {
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

func (a Aniversary) containsNealyName(target string) bool {
	for _, n := range a.names() {
		if strings.Contains(n, target) {
			return true
		}
	}
	return false
}

func containsInAnivs(anivs []Aniversary, t anaconda.Tweet) bool {
	for _, a := range anivs {
		if a.contains(t) {
			return true
		}
	}
	return false
}

func (a Aniversary) tweetStrings() []string {
	r := make([]string, 0, len(a.Tweets))
	for _, t := range a.Tweets {
		r = append(r, t.Text)
	}
	return r
}
