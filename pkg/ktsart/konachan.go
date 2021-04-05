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

// KAPI is a konachan connection sub-module that contains connection information and a set of functions.
type KAPI struct {
	httpClient *http.Client
	prefix     string
}

// KPosts contains information on a set of konachan posts.
type KPosts struct {
	List []KPost
}

// KPost contains information on a single gelbooru post.
type KPost struct {
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
	Rating              rune   `json:"rating"`
	HasChildren         bool   `json:"has_children"`
	ParentID            int    `json:"parent_id"`
	Status              string `json:"status"`
	Width               int    `json:"width"`
	Height              int    `json:"height"`
	IsHeld              bool   `json:"is_held"`
	FramesPendingString string `json:"frames_pending_string"`
	// 	"frames_pending":[

	// 	],
	// 	"frames_string":"",
	// 	"frames":[

	// 	]
	//  },
}

// NewKB creates and returns a new konachan api connection module.
func NewKB(c *http.Client, p string) *KAPI {
	api := new(KAPI)
	api.httpClient = c
	api.prefix = p
	return api
}

// metaGet creates and sends an HTTP GET request to konachan.
func (api *KAPI) metaGet(u *string) ([]byte, error) {
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

// GetByIDRaw gets a set of gelbooru posts with a given id.
// func (api *DBAPI) GetByIDRaw(id int) (p *DBPosts, e error) {
// 	p = new(DBPosts)

// 	path := fmt.Sprintf("%sposts.json?id=%d", api.prefix, id)

// 	defer func() {
// 		if r := recover(); r != nil {
// 			p, e = nil, errors.New(fmt.Sprintf("Unknown error while getting %s", path))
// 		}
// 	}()

// 	if data, err := api.metaGet(&path); err != nil {
// 		return nil, err
// 	} else {
// 		if err := json.Unmarshal(data, p); err != nil {
// 			return nil, err
// 		} else {
// 			return p, nil
// 		}
// 	}
// }

// GetByTagsRaw gets a set of danbooru posts containing the given tags and unmarshals all of the response data.
func (api *KAPI) GetByTagsRaw(t []string, n int) (p *KPosts, e error) {
	p = new(KPosts)

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
			return nil, err
		} else {
			return p, nil
		}
	}
}

// GetByTagsGeneric calls the raw function and packs the data into a set of generified dataset called BooruPost.
func (api KAPI) GetByTagsGeneric(t []string, n int) (p *[]BooruPost, e error) {
	pr, e := api.GetByTagsRaw(t, n)
	if e != nil {
		return nil, e
	}

	return pr.PostList(), nil
}

// GetByTagsRandRaw gets a set of randomly selected GBPosts that have the given tags.
func (api *KAPI) GetByTagsRandRaw(t []string, n int) (p *KPosts, e error) {
	// Since gelbooru has no random post with tags feature unlike danbooru we have to grab 100 latest arts and get 'n' of them with random indices.
	const artAmount = 100
	if n >= artAmount {
		return &KPosts{}, errors.New("amount of requested posts too big")
	}

	posts, err := api.GetByTagsRaw(t, artAmount) // TODO: Profiling test - how much adding another 100 will affect the performance
	if err != nil {
		return posts, err
	} else if len(posts.List) == 0 {
		return &KPosts{}, nil
	}

	var randSet KPosts
	rand.Seed(int64(rand.Uint64()))
	for i := 0; i < n; i++ {
		rngres := rand.Int31n(artAmount)
		randSet.List = append(randSet.List, posts.List[rngres])
	}
	return &randSet, nil
}

// GetByTagsRandGeneric calls the rand raw function and packs the data into a set of generified dataset called BooruPost.
func (api KAPI) GetByTagsRandGeneric(t []string, n int) (p *[]BooruPost, e error) {
	pr, e := api.GetByTagsRandRaw(t, n)
	if e != nil {
		return nil, e
	}

	return pr.PostList(), nil
}

// PostList returns a pointer to the posts.
func (kbps KPosts) PostList() *[]BooruPost {
	posts := make([]BooruPost, len(kbps.List))

	for i, p := range kbps.List {
		posts[i] = p
	}

	return &posts
}

// IMGURL returns a pointer to the uncompressed highest resolution url of the image.
func (kbp KPost) IMGURL() *string {
	return &kbp.FileURL
}

// ComprIMGURL returns a pointer to an URL of a possibly compressed and lower resolution image.
func (kbp KPost) ComprIMGURL() *string {
	return &kbp.JpegURL
}
