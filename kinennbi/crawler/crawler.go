package crawler

import (
	"github.com/ChimeraCoder/anaconda"
	"log"
	"net/url"
)

type Crawler struct {
	Api *anaconda.TwitterApi
}

func (c Crawler) AnalyzeAnniversary() {
	vs := url.Values{}
	ts, e := c.Api.GetMentionsTimeline(vs)
	failIfNeeded(e)
	aniversaries := make([]Aniversary, 0, len(ts))
	for _, t := range ts {
		if containsInAnivs(aniversaries, t) {
			continue
		}
		r := make([]anaconda.Tweet, 0, len(ts))
		r = c.getReplyRecursively(t, r)
		aniv := Aniversary{Tweets: r}
		//		log.Println("r length: ", len(r))
		if len(aniv.names()) > 1 {
			aniversaries = append(aniversaries, aniv)
		}
	}

	log.Println("============SUMMARY============")
	log.Println("aniversary count: ", len(aniversaries))

	for _, aniv := range aniversaries {
		//		log.Println("mention count: ", len(aniv.Tweets))
		log.Println("message: ", aniv.createMessage())
		//		for _, t := range aniv.Tweets {
		//			log.Println(t.CreatedAt, t.Text)
		//		}
	}
}

func (c Crawler) AnalyzeMentions() {
	vs := url.Values{}
	ts, e := c.Api.GetMentionsTimeline(vs)
	failIfNeeded(e)
	printResults(ts)
}

func (c Crawler) AnalyzeTimeline() {
	vs := url.Values{}
	vs.Set("count", "200")
	ts, e := c.Api.GetUserTimeline(vs)
	failIfNeeded(e)
	printResults(ts)
}

func (c Crawler) getReplyRecursively(t anaconda.Tweet, replies []anaconda.Tweet) []anaconda.Tweet {
	//	log.Println("count: ", len(replies))
	replies = append(replies, t)
	log.Println("text: ", t.Text)
	rid := t.InReplyToStatusID
	if rid == 0 {
		return replies
	}
	nextT, e := c.Api.GetTweet(t.InReplyToStatusID, nil)
	failIfNeeded(e)
	return c.getReplyRecursively(nextT, replies)
}
