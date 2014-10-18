package crawler

import (
	"github.com/ChimeraCoder/anaconda"
	"log"
	"net/url"
)

type Crawler struct {
	Api *anaconda.TwitterApi
}

func NewCrawler(accessToken string, accessTokenSecret string) Crawler {
	anaconda.SetConsumerKey("WKs1pXAfWbwat1MPimOWdmoBm")
	anaconda.SetConsumerSecret("0tvflEMDAz0yrjTEKUis3bueEVGjNhBy2tR25pNWRAZpKcAhrO")
	c := Crawler{
		Api: anaconda.NewTwitterApi(accessToken, accessTokenSecret),
	}
	return c
}

func (c Crawler) PostByAniv(message string) {
	api := anaconda.NewTwitterApi("2862013525-BGq8ZDZ4hxfhW1tbOwHjK63PlR6c2Sf6d9EqRgu", "QrEbfKfrCJF5jRQy7KBYbeXLeMB0W8zaBTK9CvwMQUjfi")
	_, e := api.PostTweet(message, nil)
	failIfNeeded(e)
}

func (c Crawler) GetOEmbed(statusId int64, candidates []anaconda.Tweet) anaconda.OEmbed {
	r, e := c.Api.GetOEmbedId(statusId, nil)
	if e != nil {
		return c.tryGetEmbedAll(candidates)
	}
	return r
}

func (c Crawler) tryGetEmbedAll(candidates []anaconda.Tweet) anaconda.OEmbed {
	for _, t := range candidates {
		r, e := c.Api.GetOEmbedId(t.Id, nil)
		if e != nil {
			continue
		}
		return r
	}
	return anaconda.OEmbed{}
}

func (c Crawler) AnalyzeAnniversary() []Anniversary {
	vs := url.Values{}
	vs.Set("count", "200")
	ts, e := c.Api.GetMentionsTimeline(vs)
	l := len(ts)
	start := l - 20
	if start < 0 {
		start = 0
	}
	ts = ts[start:]
	failIfNeeded(e)
	aniversaries := make([]Anniversary, 0, len(ts))
	for _, t := range ts {
		if containsInAnivs(aniversaries, t) {
			continue
		}
		r := make([]anaconda.Tweet, 0, len(ts))
		r = c.getReplyRecursively(t, r)
		aniv := Anniversary{Tweets: r}
		//		log.Println("r length: ", len(r))
		if len(aniv.Names()) > 1 {
			aniversaries = append(aniversaries, aniv)
		}
	}
	return aniversaries
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
