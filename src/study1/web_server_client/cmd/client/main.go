package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var complete chan int

func modify(chl chan int) {
	host := "http://localhost:8080"

	resp, err := http.PostForm(host+"/test/modify", url.Values{"key": {"test1"}, "data": {"哈哈。。。"}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	result, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(result))

	chl <- 1
}
func main() {
	t1 := time.Now()
	time.Sleep(time.Second)
	num := 1000
	count := 100

for k := 0; k < num; k++ {
	chls := make([]chan int,count)
	for i := 0; i < count; i++ {

		chls[i] = make(chan int)
		go modify(chls[i])

	}

	//排空信道
	for j,_:= range chls {
		<-chls[j]
	}

}
	t2 := time.Now()
	d := t2.Sub(t1)
	fmt.Println(d)
}
