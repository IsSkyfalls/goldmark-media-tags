package media

import (
	"errors"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
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

// every tag should be a namedTag
type namedTag interface {
	tagName() string
}

// tagInit is a helper type to create and initialize a Media tag
type tagInit interface {
	namedTag
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
		tag.Src = parent.Link
	}
	return &tag
}

var tagInitsLUT = map[Type]tagInit{
	TypePicture: tagPictureInit{},
	TypeAudio:   tagVideoAndAudioInit{name: "audio"},
	TypeVideo:   tagVideoAndAudioInit{"video"},
}

// TagSourceImg and TagSourceSource should be attributed
type attributed interface {
	updateAttributes()
}

type HasSrc interface {
	GetSrc() string
	SetSrc(src string)
}

type TagSourceImg struct {
	ast.BaseInline
	namedTag
	Src string
	Alt string
}

type TagSourceSource struct {
	ast.BaseInline
	namedTag
	Src string
	// SrcSet is only used in <picture>s
	SrcSet string
}

var KindSource = ast.NewNodeKind("MediaSourceSource")

// Kind implements ast.Node.Kind
func (t TagSourceSource) Kind() ast.NodeKind {
	return KindSource
}

func (t TagSourceSource) tagName() string {
	return "source"
}

func (t *TagSourceSource) updateAttributes() {
	t.SetAttributeString("src", t.Src)
	t.SetAttributeString("srcset", t.SrcSet)
}

func (t TagSourceSource) GetSrc() string {
	return t.Src
}

func (t *TagSourceSource) SetSrc(src string) {
	t.Src = src
}

// Dump implements ast.Node.Dump
func (t TagSourceSource) Dump(source []byte, level int) {
	t.updateAttributes()
	dumpAttributes(&t, source, level)
}

var KindImg = ast.NewNodeKind("MediaSourceImage")

// Kind implements ast.Node.Kind
func (t TagSourceImg) Kind() ast.NodeKind {
	return KindImg
}

func (t TagSourceImg) tagName() string {
	return "img"
}

func (t *TagSourceImg) updateAttributes() {
	t.SetAttributeString("src", t.Src)
	t.SetAttributeString("alt", t.Alt)
}

func (t TagSourceImg) GetSrc() string {
	return t.Src
}

func (t *TagSourceImg) SetSrc(src string) {
	t.Src = src
}

// Dump implements ast.Node.Dump
func (t TagSourceImg) Dump(source []byte, level int) {
	t.updateAttributes()
	dumpAttributes(&t, source, level)
}

// Media represents an inline <video>, <audio> or <picture> node
type Media struct {
	ast.BaseInline
	MediaType Type
	Alt       string
	Link      string
}

var KindMedia = ast.NewNodeKind("MediaParent")

// Kind implements Node.Kind
func (n *Media) Kind() ast.NodeKind {
	return KindMedia
}

// Dump implements Node.Dump
func (n *Media) Dump(source []byte, level int) {
	dumpAttributes(n, source, level)
}

func (n *Media) Text(source []byte) []byte {
	return []byte(n.Alt)
}

// mediaHTMLRenderer implements rendering for Media nodes
type mediaHTMLRenderer struct {
}

func (v mediaHTMLRenderer) RegisterFuncs(registerer renderer.NodeRendererFuncRegisterer) {
	registerer.Register(KindMedia, renderMediaTag)
	registerer.Register(KindSource, renderSourceTag)
	registerer.Register(KindImg, renderSourceTag)
}

func renderSourceTag(writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n.(attributed).updateAttributes()
		writer.WriteString(renderTagWithAttributesSelfClose(n, n.(namedTag).tagName()))
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkSkipChildren, nil
}

func renderMediaTag(writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if media, ok := n.(*Media); ok {
		init, validTag := tagInitsLUT[media.MediaType]
		if !validTag {
			return ast.WalkSkipChildren, errors.New("invalid media type")
		}
		if entering {
			writer.WriteString(renderTagWithAttributesNoClosing(n, init.tagName()))
			writer.WriteString(">")
			return ast.WalkContinue, nil
		}
		writer.WriteString("</" + init.tagName() + ">")
		return ast.WalkContinue, nil
	}
	return ast.WalkSkipChildren, errors.New("invalid tag")
}
