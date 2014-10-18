package crawler

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"net/url"
)

type Crawler struct {
	Api *anaconda.TwitterApi
}

func (c Crawler) AnalyzeAnniversary() {
	vs := url.Values{}
	ts, e := c.Api.GetMentionsTimeline(vs)
	failIfNeeded(e)
	r := make([]anaconda.Tweet, 0, len(ts))
	for _, t := range ts {
		r = c.getReplyRecursively(t, r)
	}
	for _, rt := range r {
		fmt.Println(rt.Text)
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
	rid := t.InReplyToStatusID
	if rid == 0 {
		return replies
	}
	nextT, e := c.Api.GetTweet(t.InReplyToStatusID, nil)
	failIfNeeded(e)
	fmt.Println("count: ", len(replies))
	replies = append(replies, nextT)
	return c.getReplyRecursively(nextT, replies)
}
