package crawler

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"log"
	"strings"
)

func failIfNeeded(e error) {
	if e != nil {
		log.Println(e)
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

func containsString(ss []string, target string) bool {
	for _, s := range ss {
		if s == target {
			return true
		}
	}
	return false
}
func containsNearlyString(ss []string, target string) bool {
	for _, s := range ss {
		if strings.Contains(target, s) {
			return true
		}
	}
	return false
}

func OutputAniversarries(aniversaries []Anniversary) {
	log.Println("============SUMMARY============")
	log.Println("aniversary count: ", len(aniversaries))
	for _, aniv := range aniversaries {
		first := aniv.CreateFirstMessage()
		second, _ := aniv.CreateSecondMessage()
		log.Println("message: ", aniv.CreateFullMessage(first, second))
	}
}
