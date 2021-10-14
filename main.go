package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	txtrecords, _ := net.LookupTXT("claime-dev.tk")
	t := "claime-ownership-claim="
	for _, txt := range txtrecords {
		if strings.HasPrefix(txt, t) {
			address := strings.ReplaceAll(txt, t, "")
			fmt.Println(address)
		}
	}
}
