package main

import (
	"github.com/SYM-CONTEST/Cookpad2014/kinennbi/crawler"
	"log"
)

func main() {
	// 認証した人のcrawlerを生成
	c := crawler.NewCrawler("35741880-nEBx773oAZqU4aMq246CrO3o3YguX5fpvjSSFue0V", "GYPNfRUHKOmkIEqSf72Ilh5aqAn37L78asXevSfFopGkc")
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
	second, statusId := a.CreateSecondMessage()
	// TODO: secondは希に空文字がありうるので何とかしたい
	log.Println(second)
	// 埋め込み用
	embed := c.GetOEmbed(statusId, a.Tweets)
	log.Println("embed HTML: ", embed.Html)
	log.Println("embed original url: ", embed.Url)
	// Tweet用メッセージ
	full := a.CreateFullMessage(first, second)
	log.Println(full)
	// Tweetする時はこれで
	//	c.PostByAniv(full)
}
