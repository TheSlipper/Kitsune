// Package ktsart provides bindings to the booru APIs.
package ktsart

// BooruAPI is an interface that defines necessary functions of each booru struct.
type BooruAPI interface {
	GetByTagsGeneric(t []string, n int) (p *[]BooruPost, e error)
	GetByTagsRandGeneric(t []string, n int) (p *[]BooruPost, e error)
}

// BooruPosts is an interface that defines necessary functions of each booru posts struct.
type BooruPosts interface {
	PostList() *[]BooruPost
}

// BooruPost is an interface that defines necessary functions of each booru post struct.
type BooruPost interface {
	IMGURL() *string
	ComprIMGURL() *string
}
