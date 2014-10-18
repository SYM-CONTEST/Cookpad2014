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
	anaconda.SetConsumerKey("n567b7sH6HrPIWBZyhHM2QiaK")
	anaconda.SetConsumerSecret("ygYGJ7aXEh2UQgLI7pOOWU5cixK6o7pDWYVY4MmvRaerJjqLwT")
	c := Crawler{
		Api: anaconda.NewTwitterApi(accessToken, accessTokenSecret),
	}
	return c
}

func (c Crawler) PostByAniv(message string) {
	api := anaconda.NewTwitterApi("2862013525-LWE44BXKxKmfa2tDM0EOPxUkLZGm7labnskp6v7", "WX6AZIPsj5AVWPlaTYhhE8gOE3htiMUAzzQWgqzwqVtFZ")
	_, e := api.PostTweet(message, nil)
	failIfNeeded(e)
}

func (c Crawler) AnalyzeAnniversary() []Anniversary {
	vs := url.Values{}
	ts, e := c.Api.GetMentionsTimeline(vs)
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
