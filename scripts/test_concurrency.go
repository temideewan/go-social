package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type UpdatePostPayload struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

func updatePost(postID int, p UpdatePostPayload, wg *sync.WaitGroup) {
	defer wg.Done()

	// construct the URL for update endpoint
	url := fmt.Sprintf("http://localhost:4000/v1/posts/%d", postID)

	// create the json payload
	b, _ := json.Marshal(p)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Error creating request", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Update response status:", resp.Status)
}

func main() {
	var wg sync.WaitGroup
	postID := 4

	// simulate user A and userB updating the same post
	wg.Add(2)
	content := "NEW CONTENT FROM USER B"
	title := "NEW TITLE FROM USER A"

	go updatePost(postID, UpdatePostPayload{Title: &title}, &wg)
	go updatePost(postID, UpdatePostPayload{Content: &content}, &wg)
	wg.Wait()
}
