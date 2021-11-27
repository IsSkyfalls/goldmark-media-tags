package media

import (
	"fmt"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
	"strconv"
	"strings"
)

var kindMedia = ast.NewNodeKind("media")

// MediaKind returns the ast.Nodekind of the Media node
func MediaKind() ast.NodeKind {
	return kindMedia
}

type Type byte

const (
	Video   Type = 'v'
	Audio   Type = 'a'
	Picture Type = 'p'
)

// Media represents an inline <video> or <audio> node in the ast
type Media struct {
	ast.BaseInline
	//Controls for media only
	Controls bool
	//Autoplay for media only
	Autoplay bool
	//Loop for media only
	Loop bool
	//Muted for media only
	Muted bool
	//Preload for media only
	Preload string
	//Alt for <img> inside <picture>
	Alt       string
	MediaType Type
	Sources   Sources
}

// Kind implements Node.Kind
func (n *Media) Kind() ast.NodeKind {
	return kindMedia
}

// Dump implements Node.Dump
func (n *Media) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{
		"Controls": strconv.FormatBool(n.Controls),
		"Autoplay": strconv.FormatBool(n.Autoplay),
		"Loop":     strconv.FormatBool(n.Loop),
		"Muted":    strconv.FormatBool(n.Muted),
		"Preload":  n.Preload,
		"Sources":  n.Sources.String(),
	}, nil)
}

func (n Media) Playable() bool {
	return n.MediaType == Video || n.MediaType == Audio
}

func (n Media) IsPicture() bool {
	return n.MediaType == Picture
}

// mediaHTMLRenderer implements rendering for Media nodes
type mediaHTMLRenderer struct {
}

func (v mediaHTMLRenderer) RegisterFuncs(registerer renderer.NodeRendererFuncRegisterer) {
	registerer.Register(kindMedia, renderMedia)
}

// renderMedia is the actual rendering code, implementing renderer.NodeRendererFunc
func renderMedia(writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		v := n.(*Media)

		switch v.MediaType {
		case Video:
			_, _ = writer.WriteString("<video")
		case Audio:
			_, _ = writer.WriteString("<audio")
		case Picture:
			_, _ = writer.WriteString("<picture")
		default:
			return ast.WalkContinue, fmt.Errorf("invalid media type %b", v.MediaType)
		}
		if v.Playable() {
			if v.Controls {
				_, _ = writer.WriteString(" controls")
			}
			if v.Autoplay {
				_, _ = writer.WriteString(" autoplay")
			}
			if v.Loop {
				_, _ = writer.WriteString(" loop")
			}
			if v.Muted {
				_, _ = writer.WriteString(" muted")
			}
			if v.Preload != "" {
				_, _ = writer.WriteString(" preload=\"" + v.Preload + "\"")
			}
		}

		_, _ = writer.WriteString(">")

		//<source> tags
		for _, s := range v.Sources {
			s.writeHTMLTag(writer, v)
		}
		//closing
		switch v.MediaType {
		case Video:
			_, _ = writer.WriteString("</video>")
		case Audio:
			_, _ = writer.WriteString("</audio>")
		case Picture:
			_, _ = writer.WriteString("</picture>")
		}
	}
	//should not have any children
	return ast.WalkSkipChildren, nil
}

// Source represents the <source> element, or the fallback <img> inside <picture>
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/source#attributes
type Source struct {
	Src   string
	Sizes string
	Type  string
	// SrcSet for <picture>, experimental
	SrcSet string
	// Media for <picture>, experimental
	Media string
	// IsDefault only used for the default <img> inside <source>
	IsDefault bool
}

type Sources []Source

func (s Sources) String() string {
	sources := make([]string, len(s))
	for i, e := range s {
		sources[i] = e.Type + "->" + e.Src
	}
	return strings.Join(sources, ",")
}

func (s Source) writeHTMLTag(writer util.BufWriter, parent *Media) {
	if s.IsDefault && parent.IsPicture() {
		_, _ = writer.WriteString("<img")
		if parent.Alt != "" {
			_, _ = writer.WriteString(" alt=\"" + parent.Alt + "\"")
		}
	} else {
		_, _ = writer.WriteString("<source")
	}

	if s.Src != "" {
		_, _ = writer.WriteString(" src=\"" + s.Src + "\"")
	}
	if s.Type != "" {
		_, _ = writer.WriteString(" type=\"" + s.Type + "\"")
	}

	if parent.IsPicture() {
		// these only works on when <source> is inside <picture>
		if s.Media != "" && parent.IsPicture() {
			_, _ = writer.WriteString(" media=\"" + s.Media + "\"")
		}
		if s.Sizes != "" && parent.IsPicture() {
			_, _ = writer.WriteString(" sizes=\"" + s.Sizes + "\"")
		}
		if s.SrcSet != "" && parent.IsPicture() {
			_, _ = writer.WriteString(" srcset=\"" + s.SrcSet + "\"")
		}
	}
	_, _ = writer.WriteString(">")
}
