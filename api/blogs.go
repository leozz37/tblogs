package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ezeoleaf/tblogs/models"

	"github.com/ezeoleaf/tblogs/cfg"
)

var blogs models.Blogs

// GetBlogs returns a list of Blogs from the Blogio API
func GetBlogs(cfg *cfg.Config) models.Blogs {

	if len(blogs.Blogs) > 0 {
		return blogs
	}

	client := &http.Client{}

	request, err := http.NewRequest("GET", cfg.API.Host+"/blogs", nil)

	if err != nil {
		panic(err)
	}

	request.Header.Add("BLOGIO-KEY", cfg.API.Key)

	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	blogs = models.Blogs{}
	err = json.Unmarshal(body, &blogs)
	if err != nil {
		panic(err)
	}

	return blogs
}
