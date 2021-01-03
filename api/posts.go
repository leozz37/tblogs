package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ezeoleaf/tblogs/cfg"
	"github.com/ezeoleaf/tblogs/models"
)

var posts map[int]models.PostCache

const defaultTime = 4

// GetPostsByBlog returns a list of Posts for a single blog
func GetPostsByBlog(blogID int, cfg *cfg.Config) models.Posts {

	if len(posts[blogID].Posts.Posts) > 0 {
		d := time.Now()
		diff := d.Sub(posts[blogID].DateUpdated).Hours()

		if diff < defaultTime {
			return posts[blogID].Posts
		}
	}

	if len(posts) == 0 {
		posts = make(map[int]models.PostCache)
	}

	pr := models.PostRequest{Blogs: []int{blogID}}

	postsResp := fetchPosts(pr, cfg)

	pc := models.PostCache{Posts: postsResp, DateUpdated: time.Now()}

	posts[blogID] = pc

	return postsResp
}

// GetPosts returns a list of Posts for a list of Blogs using the Blogs ids
func GetPosts(blogs []int, cfg *cfg.Config) models.Posts {
	pr := models.PostRequest{Blogs: blogs}

	postsResp := fetchPosts(pr, cfg)

	return postsResp
}

func fetchPosts(reqPost models.PostRequest, cfg *cfg.Config) models.Posts {
	rJSON, err := json.Marshal(reqPost)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	payload := strings.NewReader(string(rJSON))

	req, err := http.NewRequest("GET", cfg.API.Host+"/posts", payload)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("BLOGIO-KEY", cfg.API.Key)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	posts := models.Posts{}
	err = json.Unmarshal(body, &posts)
	if err != nil {
		panic(err)
	}
	return posts
}
