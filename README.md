# spoce

> soc proof of concept engine 

## Support

- [x] xray
- [ ] nuclei

## Download

```bash
go get github.com/yanmengfei/spoce
```

## Example

```go
package main

import (
    "log"

    "github.com/yanmengfei/spoce/build"
    "github.com/yanmengfei/spoce/library/http"
    "github.com/yanmengfei/spoce/scanner"
)

var urls = []string{
    "http://117.161.6.2:8180",
    "https://27.221.68.244:443",
    "http://13.75.117.202:3000",
    "https://113.108.174.45:443",
}

var pocYamlStr = `name: poc-yaml-weblogic-console

rules:
  - method: GET
    path: /console/login/LoginForm.jsp
    headers:
      User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:52.0) Gecko/20100101 Firefox/52.0
    expression: response.status==200
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

```
