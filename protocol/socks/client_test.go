package socks

import (
	"log"
	"net/http"
	"testing"
)

func TestClient_Dail(t *testing.T) {
	client , _ := NewClient(":888" , "aaa" , "aaa")
	tr := &http.Transport{Dial: client.Dail}
	hc := &http.Client{Transport: tr}
	log.Println(hc.Get("http://www.baidu.com"))
	//clientConn := httputil.NewClientConn(conn , bufio.NewReader(nil))
	//clientConn.Do()
}