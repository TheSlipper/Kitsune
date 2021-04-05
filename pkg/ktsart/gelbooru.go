// LEGAL:
// gelbooru.go mostly made by etw (https://github.com/etw)
// To see who committed which parts of the code check: https://github.com/etw/gobooru/blob/develop/common.go

// Package ktsart provides bindings to the booru APIs.
package ktsart

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/TheSlipper/Kitsune/internal/settings"
)

// GBAuth contains information on gelbooru API authentication
type GBAuth struct {
	User string
	Hash string
}

// GBAPI is a gelbooru connection sub-module that contains connection information and a set of functions.
type GBAPI struct {
	httpClient *http.Client
	prefix     string
	auth       []http.Cookie
}

// GBPosts contains information on a set of gelbooru posts.
type GBPosts struct {
	XMLName xml.Name `xml:"posts"`
	Count   int      `xml:"count,attr"`
	Offset  int      `xml:"offset,attr"`
	List    []GBPost `xml:"post"`
}

// GBPost contains information on a single gelbooru post.
type GBPost struct {
	Height        int    `xml:"height,attr"`
	Width         int    `xml:"width,attr"`
	ParentID      int    `xml:"parent_id,attr"`
	FileURL       string `xml:"file_url,attr"`
	SampleURL     string `xml:"sample_url,attr"`
	SampleHeight  int    `xml:"sample_height,attr"`
	SampleWidth   int    `xml:"sample_width,attr"`
	Score         int    `xml:"score,attr"`
	PreviewURL    string `xml:"preview_url,attr"`
	PreviewHeight int    `xml:"preview_height,attr"`
	PreviewWidth  int    `xml:"preview_width,attr"`
	Rating        string `xml:"rating,attr"`
	ID            int    `xml:"id,attr"`
	Tags          string `xml:"tags,attr"`
	Change        int    `xml:"change,attr"`
	Md5           string `xml:"md5,attr"`
	CreatorID     int    `xml:"creator_id,attr"`
	CreatedAt     string `xml:"created_at,attr"`
	Status        string `xml:"status,attr"`
	Source        string `xml:"source,attr"`
	HasNotes      bool   `xml:"has_notes,attr"`
	HasComments   bool   `xml:"has_comments,attr"`
	HasChildren   bool   `xml:"has_children,attr"`
}

// NewGB creates and returns a new gelbooru api connection module.
func NewGB(c *http.Client, p string, a *GBAuth) *GBAPI {
	api := new(GBAPI)
	api.httpClient = c
	api.prefix = p + "index.php?page=dapi&q=index"
	if a != nil {
		api.auth = append(api.auth, http.Cookie{
			Name: settings.BotSettings.GelbooruUsrID,
			// Name:  "user_id",
			Value: a.User,
		})
		api.auth = append(api.auth, http.Cookie{
			Name:  "pass_hash",
			Value: a.Hash,
		})
	}
	return api
}

// TODO: perhaps change this to a function that's not gallery/struct specific
// metaGet creates and sends an HTTP GET request to gelbooru.
func (api *GBAPI) metaGet(u *string) ([]byte, error) {
	if req, err := http.NewRequest("GET", *u, nil); err != nil {
		return nil, err
	} else {
		for _, c := range api.auth {
			req.AddCookie(&c)
		}
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
func (api *GBAPI) GetByIDRaw(id int) (p *GBPosts, e error) {
	p = new(GBPosts)

	path := fmt.Sprintf("%s&s=post&id=%d", api.prefix, id)

	defer func() {
		if r := recover(); r != nil {
			p, e = nil, errors.New(fmt.Sprintf("Unknown error while getting %s", path))
		}
	}()

	if data, err := api.metaGet(&path); err != nil {
		return nil, err
	} else {
		if err := xml.Unmarshal(data, p); err != nil {
			return nil, err
		} else {
			return p, nil
		}
	}
}

// GetByTagsRaw gets a set of gelbooru posts containing the given tags and unmarshals all of the response data.
func (api *GBAPI) GetByTagsRaw(t []string, n int) (p *GBPosts, e error) {
	p = new(GBPosts)

	path := fmt.Sprintf("%s&s=post&tags=%s&limit="+strconv.Itoa(n), api.prefix, strings.Join(t, " "))

	defer func() {
		if r := recover(); r != nil {
			p, e = nil, errors.New(fmt.Sprintf("Unknown error while getting %s", path))
		}
	}()

	if data, err := api.metaGet(&path); err != nil {
		return nil, err
	} else {
		if err := xml.Unmarshal(data, p); err != nil {
			return nil, err
		} else {
			return p, nil
		}
	}
}

// GetByTagsGeneric calls the raw function and packs the data into a set of generified dataset called BooruPost.
func (api GBAPI) GetByTagsGeneric(t []string, n int) (p *[]BooruPost, e error) {
	pr, e := api.GetByTagsRaw(t, n)
	if e != nil {
		return nil, e
	}

	return pr.PostList(), nil
}

// GetByTagsRandRaw gets a set of randomly selected GBPosts that have the given tags.
func (api GBAPI) GetByTagsRandRaw(t []string, n int) (p *GBPosts, e error) {
	// Since gelbooru has no random post with tags feature unlike danbooru we have to grab 100 latest arts and get 'n' of them with random indices.
	const artAmount = 100
	if n >= artAmount {
		return &GBPosts{}, errors.New("amount of requested posts too big")
	}

	posts, err := api.GetByTagsRaw(t, artAmount) // TODO: Profiling test - how much adding another 100 will affect the performance
	if err != nil {
		return posts, err
	} else if len(posts.List) == 0 {
		return &GBPosts{}, nil
	}

	var randSet GBPosts
	rand.Seed(int64(rand.Uint64()))
	for i := 0; i < n; i++ {
		rngres := rand.Int31n(artAmount)
		randSet.List = append(randSet.List, posts.List[rngres])
	}
	return &randSet, nil
}

// GetByTagsRandGeneric calls the rand raw function and packs the data into a set of generified dataset called BooruPost.
func (api GBAPI) GetByTagsRandGeneric(t []string, n int) (p *[]BooruPost, e error) {
	pr, e := api.GetByTagsRandRaw(t, n)
	if e != nil {
		return nil, e
	}

	return pr.PostList(), nil
}

// PostList returns a pointer to the posts.
func (gbps GBPosts) PostList() *[]BooruPost {
	posts := make([]BooruPost, len(gbps.List))

	for i, p := range gbps.List {
		posts[i] = p
	}

	return &posts
}

// IMGURL returns a pointer to the uncompressed highest resolution url of the image.
func (gbp GBPost) IMGURL() *string {
	return &gbp.FileURL
}

// ComprIMGURL returns a pointer to an URL of a possibly compressed and lower resolution image.
func (gbp GBPost) ComprIMGURL() *string {
	return &gbp.FileURL
}
