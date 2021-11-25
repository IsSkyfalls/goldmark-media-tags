package media

import (
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

// Media represents an inline <video> or <audio> node in the ast
type Media struct {
	ast.BaseInline
	Controls bool
	Autoplay bool
	Loop     bool
	Muted    bool
	Preload  string
	IsVideo  bool
	Sources  Sources
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

		//opening
		if v.IsVideo {
			_, _ = writer.WriteString("<video")
		} else {
			_, _ = writer.WriteString("<audio")
		}

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
		_, _ = writer.WriteString(">")

		//<source> tags
		for _, s := range v.Sources {
			s.writeHTMLTag(writer)
		}
		//closing
		if v.IsVideo {
			_, _ = writer.WriteString("</video>")
		} else {
			_, _ = writer.WriteString("</audio>")
		}
	}
	//should not have any children
	return ast.WalkSkipChildren, nil
}

// Source represents the <source> element
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/source#attributes
type Source struct {
	Src    string
	Media  string
	Sizes  string
	Type   string
	SrcSet string
}

type Sources []Source

func (s Sources) String() string {
	sources := make([]string, len(s))
	for i, e := range s {
		sources[i] = e.Type + "->" + e.Src
	}
	return strings.Join(sources, ",")
}

func (s Source) writeHTMLTag(writer util.BufWriter) {
	_, _ = writer.WriteString("<source")
	if s.Media != "" {
		_, _ = writer.WriteString(" media=\"" + s.Media + "\"")
	}
	if s.Sizes != "" {
		_, _ = writer.WriteString(" sizes=\"" + s.Sizes + "\"")
	}
	if s.Src != "" {
		_, _ = writer.WriteString(" src=\"" + s.Src + "\"")
	}
	if s.SrcSet != "" {
		_, _ = writer.WriteString(" srcset=\"" + s.SrcSet + "\"")
	}
	if s.Type != "" {
		_, _ = writer.WriteString(" type=\"" + s.Type + "\"")
	}
	_, _ = writer.WriteString(">")
}
