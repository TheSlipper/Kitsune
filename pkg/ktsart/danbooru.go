// Package ktsart provides bindings to the booru APIs.
package ktsart

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// DBAuth contains information on danbooru API authentication.
type DBAuth struct {
	User string
	Hash string
}

// DBAPI is a danbooru connection sub-module that contains connection information and a set of functions.
type DBAPI struct {
	httpClient *http.Client
	prefix     string
	auth       []http.Cookie
}

// DBPosts contains information on a set of danbooru posts.
type DBPosts struct {
	List []DBPost
}

// DBPost contains information on a single danbooru post.
type DBPost struct {
	ID                  int         `json:"id"`
	CreatedAt           string      `json:"created_at"`
	UploaderID          int         `json:"uploader_id"`
	Score               int         `json:"score"`
	Source              string      `json:"source"`
	Md5                 string      `json:"md5"`
	LastCommentBumpedAt string      `json:"last_comment_bumped_at"`
	Rating              string      `json:"rating"`
	ImageWidth          int         `json:"image_width"`
	ImageHeight         int         `json:"image_height"`
	TagString           string      `json:"tag_string"`
	IsNoteLocked        bool        `json:"is_note_locked"`
	FavCount            int         `json:"fav_count"`
	FileExt             string      `json:"file_ext"`
	LastNotedAt         interface{} `json:"last_noted_at"`
	IsRatingLocked      bool        `json:"is_rating_locked"`
	ParentID            int         `json:"parent_id"`
	HasChildren         bool        `json:"has_children"`
	ApproverID          interface{} `json:"approver_id"`
	TagCountGeneral     int         `json:"tag_count_general"`
	TagCountArtist      int         `json:"tag_count_artist"`
	TagCountCharacter   int         `json:"tag_count_character"`
	TagCountCopyright   int         `json:"tag_count_copyright"`
	FileSize            int         `json:"file_size"`
	IsStatusLocked      bool        `json:"is_status_locked"`
	FavString           string      `json:"fav_string"`
	PoolString          string      `json:"pool_string"`
	UpScore             int         `json:"up_score"`
	DownScore           int         `json:"down_score"`
	IsPending           bool        `json:"is_pending"`
	IsFlagged           bool        `json:"is_flagged"`
	IsDeleted           bool        `json:"is_deleted"`
	TagCount            int         `json:"tag_count"`
	UpdatedAt           string      `json:"updated_at"`
	IsBanned            bool        `json:"is_banned"`
	PixivID             interface{} `json:"pixiv_id"`
	LastCommentedAt     string      `json:"last_commented_at"`
	HasActiveChildren   bool        `json:"has_active_children"`
	BitFlags            int         `json:"bit_flags"`
	UploaderName        string      `json:"uploader_name"`
	HasLarge            bool        `json:"has_large"`
	TagStringArtist     string      `json:"tag_string_artist"`
	TagStringCharacter  string      `json:"tag_string_character"`
	TagStringCopyright  string      `json:"tag_string_copyright"`
	TagStringGeneral    string      `json:"tag_string_general"`
	HasVisibleChildren  bool        `json:"has_visible_children"`
	ChildrenIds         interface{} `json:"children_ids"`
	FileURL             string      `json:"file_url"`
	LargeFileURL        string      `json:"large_file_url"`
	PreviewFileURL      string      `json:"preview_file_url"`
}

// NewDB creates and returns a new danbooru api connection module.
func NewDB(c *http.Client, p string, a *DBAuth) *DBAPI {
	api := new(DBAPI)
	api.httpClient = c
	api.prefix = p
	if a != nil { // TODO: Check if those names are correct
		api.auth = append(api.auth, http.Cookie{
			Name:  "login",
			Value: a.User,
		})
		api.auth = append(api.auth, http.Cookie{
			Name:  "api_key",
			Value: a.Hash,
		})
	}
	return api
}

// metaGet creates and sends an HTTP GET request to danbooru.
func (api *DBAPI) metaGet(u *string) ([]byte, error) {
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
func (api *DBAPI) GetByIDRaw(id int) (p *DBPosts, e error) {
	p = new(DBPosts)

	path := fmt.Sprintf("%sposts.json?id=%d", api.prefix, id)

	defer func() {
		if r := recover(); r != nil {
			p, e = nil, errors.New(fmt.Sprintf("Unknown error while getting %s", path))
		}
	}()

	if data, err := api.metaGet(&path); err != nil {
		return nil, err
	} else {
		if err := json.Unmarshal(data, p); err != nil {
			return nil, err
		} else {
			return p, nil
		}
	}
}

// GetByTagsRaw gets a set of danbooru posts containing the given tags and unmarshals all of the response data.
func (api *DBAPI) GetByTagsRaw(t []string, n int) (p *DBPosts, e error) {
	p = new(DBPosts)

	path := fmt.Sprintf("%sposts.json?tags=%s&limit="+strconv.Itoa(n), api.prefix, strings.Join(t, "+"))

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
func (api DBAPI) GetByTagsGeneric(t []string, n int) (p *[]BooruPost, e error) {
	pr, e := api.GetByTagsRaw(t, n)
	if e != nil {
		return nil, e
	}

	return pr.PostList(), nil
}

// GetByTagsRandRaw gets a set of randomly selected GBPosts that have the given tags.
func (api *DBAPI) GetByTagsRandRaw(t []string, n int) (p *DBPosts, e error) {
	p = new(DBPosts)

	path := fmt.Sprintf("%sposts.json?tags=%s&random=true&limit="+strconv.Itoa(n), api.prefix, strings.Join(t, "+"))

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

// GetByTagsRandGeneric calls the rand raw function and packs the data into a set of generified dataset called BooruPost.
func (api DBAPI) GetByTagsRandGeneric(t []string, n int) (p *[]BooruPost, e error) {
	pr, e := api.GetByTagsRandRaw(t, n)
	if e != nil {
		return nil, e
	}

	return pr.PostList(), nil
}

// PostList returns a pointer to the posts.
func (dbps DBPosts) PostList() *[]BooruPost {
	posts := make([]BooruPost, len(dbps.List))

	for i, p := range dbps.List {
		posts[i] = p
	}

	return &posts
}

// IMGURL returns a pointer to the uncompressed highest resolution url of the image.
func (dbp DBPost) IMGURL() *string {
	return &dbp.LargeFileURL
}

// ComprIMGURL returns a pointer to an URL of a possibly compressed and lower resolution image.
func (dbp DBPost) ComprIMGURL() *string {
	return &dbp.FileURL
}
