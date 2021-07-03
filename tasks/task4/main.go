package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func main() {
	t := make(chan struct{}, 100)
	for i := 1; i <= 100; i++ {
		go MakeRequest(strconv.Itoa(i), t)
	}
	for j := 1; j <= 100; j++ {
		<-t
	}
}

type Post struct {
	UserID int    `json:"UserId"`
	Id     int    `json:"Id"`
	Title  string `json:"Title"`
	Body   string `json:"Body"`
}

func MakeRequest(id string, t chan struct{}) {
	defer func() {
		t <- struct{}{}
	}()
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/" + id)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	post := Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = MakeFiles(id, post)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func MakeFiles(id string, post Post) error {
	err := ioutil.WriteFile("./storage/posts/"+id+".txt",
		[]byte(fmt.Sprintf("%v", post)), 0644)
	if err != nil {
		return err
	}
	return nil
}
