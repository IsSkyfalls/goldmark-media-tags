package media

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"regexp"
)

var regex = regexp.MustCompile(`!(.)\[(.*?)]\((.+?)\)`)

// mediaParser implements parser.InlineParser interface
type mediaParser struct {
	Options
}

func (p mediaParser) Trigger() []byte {
	return []byte{'!'}
}

func (p mediaParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	match := block.FindSubMatch(regex)
	if len(match) > 0 {
		flag := (Type)(match[1][0]) // one character only
		//alt := string(match[2])
		url := string(match[3])
		if flag == Video || flag == Audio || flag == Picture {
			return &Media{
				BaseInline: ast.BaseInline{},
				Controls:   p.MediaControls,
				Autoplay:   p.MediaAutoplay,
				Loop:       p.MediaLoop,
				Preload:    p.MediaPreload,
				Muted:      p.MediaMuted,
				MediaType:  flag,
				Sources: []Source{{
					Src:  url,
					Type: "",
				}},
			}
		} else {
			block.Advance(1)
		}
	}
	return nil
}
