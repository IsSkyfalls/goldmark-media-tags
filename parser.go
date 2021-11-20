package media

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"regexp"
)

const (
	video = 'v'
	audio = 'a'
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
		flag := match[1][0] // one character only
		//alt := string(match[2])
		url := string(match[3])
		if flag == video || flag == audio {
			return &Media{
				BaseInline: ast.BaseInline{},
				Controls:   p.MediaControls,
				Autoplay:   p.MediaAutoplay,
				Loop:       p.MediaLoop,
				Preload:    p.MediaPreload,
				Muted:      p.MediaMuted,
				IsVideo:    flag == video,
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
