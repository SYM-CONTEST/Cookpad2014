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
	ns := p.ParseToNouns([]string{t.Text})
	wc := wordCount(ns)
	for noun, count := range wc {
		fmt.Println(count, noun)
	}
}

func wordCount(ss []string) map[string]int {
	res := make(map[string]int)
	for _, str := range ss {
		res[strings.ToLower(str)]++
	}
	return res
}

func containsString(ss[] string, target string) bool {
	for _, s := range ss {
		if s == target {
			return true
		}
	}
	return false
}
func containsNearlyString(ss[] string, target string) bool {
	for _, s := range ss {
		if strings.Contains(target, s) {
			return true
		}
	}
	return false
}
