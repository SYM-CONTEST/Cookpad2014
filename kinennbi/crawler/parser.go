package crawler

import (
	"bitbucket.org/rerofumi/mecab"
	"bytes"
	"log"
	"strings"
	"unicode/utf8"
)

var containsBlacklist = []string{
	"@",
	"_",
	"‿",
	"´",
	"｀",
	"｀",
	")",
	"(",
	"/",
	":",
	"bPYx",
	"LAt",
	"http",
	".com",
}
var matchBlacklist = []string{
	"さん",
	"これ",
	"at",
	"via",
	"co",
	"..",
	"どの",
	"そう",
	"よう",
	"lt",
	"gt",
	"ww",
	"fi",
	"なに",
}

type Parser struct {
}

func (parser Parser) ParseToNouns(ss []string) []string {
	var b bytes.Buffer
	for _, s := range ss {
		b.WriteString(s)
	}
	s := b.String()
	return parser.parseToNouns(s)
}

func (parser Parser) filterNoise(src []string, a Anniversary, minLength int) []string {
	r := make([]string, 0, len(src))
	for _, s := range src {
		if a.containsNealyName(s) {
			continue
		}
		if utf8.RuneCountInString(s) < minLength {
			continue
		}
		if containsString(matchBlacklist, s) {
			continue
		}
		if containsNearlyString(containsBlacklist, s) {
			continue
		}
		r = append(r, s)
	}
	return r
}

func (parser Parser) parseToNouns(s string) []string {
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
