package main

import (
	_"reflect"
	"github.com/ChimeraCoder/anaconda"
	"github.com/SYM-CONTEST/Cookpad2014/kinennbi/crawler"
)

func main() {
	anaconda.SetConsumerKey("qercs5E6GZeMwNaDOSocGnMQ6")
	anaconda.SetConsumerSecret("gh3j3EtQU3Kr8PtEf8GhoDJokkYHwCnFvrt2muu8W75bXXFi72")
	api := anaconda.NewTwitterApi("35741880-LCcGJm0iTTj84yn0Ch1tmpuw0ujCCl1JkdwErkPfj", "V7Vpxbu1l40Qo2q8tzosF5JwEhryNuRsP6CU94ixSo4VX")

	//	analyzeHoge(api)
	crawler.AnalyzeAnniversary(api)
}
