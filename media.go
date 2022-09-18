package media

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
)

type Type byte

const (
	TypePicture  Type = 'p'
	TypeVideo    Type = 'v'
	TypeAudio    Type = 'a'
	AttrControls      = "controls"
	AttrAutoplay      = "autoplay"
	AttrLoop          = "loop"
	AttrMuted         = "muted"
	AttrPreload       = "preload"
)

// tagInit is a helper type to create and initialize a Media tag
type tagInit interface {
	tagName() string
	// makeSourceTag initializes attributes on the parent tag
	initAttributes(parent *Media, options Options)
	// makeSourceTag creates a <source> tag for the parent element, for the first children of <picture>, an <img> tag will be created
	// it should not append the created tag as children
	makeSourceTag(parent Media, options Options) ast.Node
}

// tagPictureInit is the initializer for <picture>
type tagPictureInit struct{}

func (t tagPictureInit) tagName() string {
	return "picture"
}

func (t tagPictureInit) initAttributes(parent *Media, options Options) {
}

func (t tagPictureInit) makeSourceTag(parent Media, options Options) ast.Node {
	if parent.ChildCount() == 0 {
		tag := &TagSourceImg{
			BaseInline: ast.BaseInline{},
			Src:        parent.Link,
			Alt:        parent.Alt,
		}
		return tag
	}
	return &TagSourceSource{}
}

// tagVideoAndAudioInit is the initializer for <video> and <audio>
type tagVideoAndAudioInit struct {
	name string
}

func (t tagVideoAndAudioInit) tagName() string {
	return t.name
}

func (t tagVideoAndAudioInit) initAttributes(parent *Media, options Options) {
	parent.SetAttributeString(AttrControls, options.MediaControls)
	parent.SetAttributeString(AttrAutoplay, options.MediaAutoplay)
	parent.SetAttributeString(AttrLoop, options.MediaLoop)
	parent.SetAttributeString(AttrMuted, options.MediaMuted)
	parent.SetAttributeString(AttrPreload, options.MediaPreload)
}

func (t tagVideoAndAudioInit) makeSourceTag(parent Media, options Options) ast.Node {
	tag := TagSourceSource{}
	if parent.ChildCount() == 0 {
		tag.SetAttributeString("src", parent.Link)
	}
	return &tag
}

var tagInitsLUT = map[Type]tagInit{
	TypePicture: tagPictureInit{},
	TypeAudio:   tagVideoAndAudioInit{name: "audio"},
	TypeVideo:   tagVideoAndAudioInit{"video"},
}

type TagSourceImg struct {
	ast.BaseInline
	Src string
	Alt string
}

type TagSourceSource struct {
	ast.BaseInline
	Src string
	// SrcSet is only used in <picture>s
	SrcSet string
}

var kindSource = ast.NewNodeKind("MediaSourceSource")

// Kind implements ast.Node.Kind
func (t TagSourceSource) Kind() ast.NodeKind {
	return kindSource
}

func (t *TagSourceSource) updateAttributes() {
	t.SetAttributeString("src", t.Src)
	t.SetAttributeString("srcset", t.SrcSet)
}

// Dump implements ast.Node.Dump
func (t TagSourceSource) Dump(source []byte, level int) {
	t.updateAttributes()
	DumpAttributes(&t, source, level)
}

var kindImg = ast.NewNodeKind("MediaSourceImage")

// Kind implements ast.Node.Kind
func (t TagSourceImg) Kind() ast.NodeKind {
	return kindImg
}

func (t *TagSourceImg) updateAttributes() {
	t.SetAttributeString("src", t.Src)
	t.SetAttributeString("alt", t.Alt)
}

// Dump implements ast.Node.Dump
func (t TagSourceImg) Dump(source []byte, level int) {
	t.updateAttributes()
	DumpAttributes(&t, source, level)
}

// Media represents an inline <video>, <audio> or <picture> node
type Media struct {
	ast.BaseInline
	MediaType Type
	Alt       string
	Link      string
}

var kindMedia = ast.NewNodeKind("MediaParent")

// Kind implements Node.Kind
func (n *Media) Kind() ast.NodeKind {
	return kindMedia
}

// Dump implements Node.Dump
func (n *Media) Dump(source []byte, level int) {
	DumpAttributes(n, source, level)
}

func (n *Media) Text(source []byte) []byte {
	return []byte(n.Alt)
}

// mediaHTMLRenderer implements rendering for Media nodes
type mediaHTMLRenderer struct {
}

func (v mediaHTMLRenderer) RegisterFuncs(registerer renderer.NodeRendererFuncRegisterer) {
	registerer.Register(kindMedia, nil)
}
