package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"strconv"
)

var db *sql.DB

func main() {
	var err error

	connStr := "postgres://postgres:example@localhost:5432/jsonplaceholder?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	done := make(chan struct{})
	go MakePosts(done)
	<-done
}

type Post struct {
	UserID int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Comments struct {
	PostID int    `json:"postId"`
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func AddPosts(p Post) (err error) {
	query := `INSERT INTO public_posts(user_id, id, title, body) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(query, p.UserID, p.Id, p.Title, p.Body)
	return err
}

func MakePosts(done chan struct{}) {
	defer func() {
		done <- struct{}{}
	}()
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts?userId=7")
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
	posts := []Post{}
	err = json.Unmarshal(body, &posts)
	if err != nil {
		fmt.Println(err)
		return
	}
	doneCom := make(chan struct{}, len(posts))
	for _, post := range posts {
		err = AddPosts(post)
		if err != nil {
			fmt.Println(err)
			return
		}
		go MakeComments(strconv.Itoa(post.Id), doneCom)
	}
	for j := 1; j <= len(posts); j++ {
		<-doneCom
	}
}

func AddComments(c Comments) (err error) {
	query := `INSERT INTO public_comments(post_id, id, name, email, body) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(query, c.PostID, c.Id, c.Name, c.Email, c.Body)
	return err
}

func MakeComments(postID string, com chan struct{}) {
	defer func() {
		com <- struct{}{}
	}()
	resp, err := http.Get("https://jsonplaceholder.typicode.com/comments?postId=" + postID)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	comments := []Comments{}
	err = json.Unmarshal(body, &comments)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, comment := range comments {
		err = AddComments(comment)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
