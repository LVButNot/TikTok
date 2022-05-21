package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	np := md5.Sum([]byte("364365706" + "123123"))
	tok := fmt.Sprintf("%X", np)
	fmt.Println(tok)
	pas := md5.Sum([]byte("123123"))
	pasmd5 := fmt.Sprintf("%X", pas)
	fmt.Println(pasmd5)
}
