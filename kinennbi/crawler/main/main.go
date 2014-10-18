package main

import (
	"github.com/SYM-CONTEST/Cookpad2014/kinennbi/crawler"
	"log"
)

func main() {
	// 認証した人のcrawlerを生成
	c := crawler.NewCrawler("298482612-AzpnvM6K8TfLw1kbOVnJTlwlQjEWGEGPgXdd7Viz", "HMvID4dg5K3WF6jo3urHYmsVk2MkAmY4V43kLBRE190DH")
	// 認証した人のメンションを分析してそれっぽい記念日群を抽出
	as := c.AnalyzeAnniversary()
	// ただの確認出力なので不要
	//	crawler.OutputAniversarries(as)
	// 記念日群からもっともよさげなものを選定 (今は取り急ぎ最古のもの)
	a := crawler.ChooseBestAniversary(as)
	// サイト用メッセージ1
	first := a.CreateFirstMessage()
	log.Println(first)
	// サイト用メッセージ2
	second := a.CreateSecondMessage()
	log.Println(second)
	// Tweet用メッセージ
	full := a.CreateFullMessage(first, second)
	log.Println(full)
	// Tweetする時はこれで
	//	c.PostByAniv(full)
}
