package crawler

import (
	"bitbucket.org/rerofumi/mecab"
	"log"
	"strings"
)

type Parser struct {
}

func (parser Parser) parseNouns(s string) []string {
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
