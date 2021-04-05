// Package ktsart provides bindings to the booru APIs.
package ktsart

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

// YAAPI is a yande.re connection sub-module that contains connection information and a set of functions.
type YAAPI struct {
	httpClient *http.Client
	prefix     string
}

// YAPosts contains information on a set of yande.re posts.
type YAPosts struct {
	List []YAPost
}

// YAPost contains information on a single yande.re post.
type YAPost struct {
	ID                  int    `json:"id"`
	Tags                string `json:"tags"`
	CreatedAt           int64  `json:"created_at"`
	CreatorID           int    `json:"creator_id"`
	Author              string `json:"author"`
	Change              int    `json:"change"`
	Source              string `json:"source"`
	Score               int    `json:"score"`
	MD5                 string `json:"md5"`
	FileSize            int    `json:"file_size"`
	FileURL             string `json:"file_url"`
	IsShownInIndex      bool   `json:"is_shown_in_index"`
	PreviewURL          string `json:"preview_url"`
	PreviewWidth        int    `json:"preview_width"`
	PreviewHeight       int    `json:"preview_height"`
	ActualPreviewWidth  int    `json:"actual_preview_width"`
	ActualPreviewHeight int    `json:"actual_preview_height"`
	SampleURL           string `json:"sample_url"`
	SampleWidth         int    `json:"sample_width"`
	SampleHeight        int    `json:"sample_height"`
	SampleFileSize      int    `json:"sample_file_size"`
	JpegURL             string `json:"jpeg_url"`
	JpegWidth           int    `json:"jpeg_width"`
	JpegHeight          int    `json:"jpeg_height"`
	JpegFileSize        int    `json:"jpeg_file_size"`
	Rating              string `json:"rating"`
	HasChildren         bool   `json:"has_children"`
	ParentID            int    `json:"parent_id"`
	Status              string `json:"status"`
	Width               int    `json:"width"`
	Height              int    `json:"height"`
	IsHeld              bool   `json:"is_held"`
	FramesPendingString string `json:"frames_pending_string"`
}

// NewYA creates and returns a new konachan api connection module.
func NewYA(c *http.Client, p string) *YAAPI {
	api := new(YAAPI)
	api.httpClient = c
	api.prefix = p
	return api
}

// metaGet creates and sends an HTTP GET request to yande.re.
func (api *YAAPI) metaGet(u *string) ([]byte, error) {
	if req, err := http.NewRequest("GET", *u, nil); err != nil {
		return nil, err
	} else {
		if resp, err := api.httpClient.Do(req); err != nil {
			return nil, err
		} else {
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return nil, errors.New(resp.Status)
			}

			if body, err := ioutil.ReadAll(resp.Body); err != nil {
				return nil, err
			} else {
				return body, nil
			}
		}
	}
}

// GetByTagsRaw gets a set of danbooru posts containing the given tags and unmarshals all of the response data.
func (api *YAAPI) GetByTagsRaw(t []string, n int) (p *YAPosts, e error) {
	p = new(YAPosts)

	path := fmt.Sprintf("%spost.json?tags=%s&limit="+strconv.Itoa(n), api.prefix, strings.Join(t, "+"))

	defer func() {
		if r := recover(); r != nil {
			p, e = nil, errors.New(fmt.Sprintf("Unknown error while getting %s", path))
		}
	}()

	if data, err := api.metaGet(&path); err != nil {
		return nil, err
	} else {
		if err := json.Unmarshal(data, &p.List); err != nil {
			return nil, err // TODO: Change this to "err"
		} else {
			return p, nil
		}
	}
}

// GetByTagsGeneric calls the raw function and packs the data into a set of generified dataset called BooruPost.
func (api YAAPI) GetByTagsGeneric(t []string, n int) (p *[]BooruPost, e error) {
	pr, e := api.GetByTagsRaw(t, n)
	if e != nil {
		return nil, e
	}

	return pr.PostList(), nil
}

// GetByTagsRandRaw gets a set of randomly selected GBPosts that have the given tags.
func (api *YAAPI) GetByTagsRandRaw(t []string, n int) (p *YAPosts, e error) {
	// Since gelbooru has no random post with tags feature unlike danbooru we have to grab 100 latest arts and get 'n' of them with random indices.
	const artAmount = 100
	if n >= artAmount {
		return &YAPosts{}, errors.New("amount of requested posts too big")
	}

	posts, err := api.GetByTagsRaw(t, artAmount) // TODO: Profiling test - how much adding another 100 will affect the performance
	if err != nil {
		return posts, err
	} else if len(posts.List) == 0 {
		return &YAPosts{}, nil
	}

	var randSet YAPosts
	rand.Seed(int64(rand.Uint64()))
	for i := 0; i < n; i++ {
		rngres := rand.Int31n(artAmount)
		randSet.List = append(randSet.List, posts.List[rngres])
	}
	return &randSet, nil
}

// GetByTagsRandGeneric calls the rand raw function and packs the data into a set of generified dataset called BooruPost.
func (api YAAPI) GetByTagsRandGeneric(t []string, n int) (p *[]BooruPost, e error) {
	pr, e := api.GetByTagsRandRaw(t, n)
	if e != nil {
		return nil, e
	}

	return pr.PostList(), nil
}

// PostList returns a pointer to the posts.
func (yabps YAPosts) PostList() *[]BooruPost {
	posts := make([]BooruPost, len(yabps.List))

	for i, p := range yabps.List {
		posts[i] = p
	}

	return &posts
}

// IMGURL returns a pointer to the uncompressed highest resolution url of the image.
func (yabp YAPost) IMGURL() *string {
	return &yabp.FileURL
}

// TODO: Check in all of the galleries if u can use preview for ComprIMGURL

// ComprIMGURL returns a pointer to an URL of a possibly compressed and lower resolution image.
func (yabp YAPost) ComprIMGURL() *string {
	return &yabp.FileURL
}
