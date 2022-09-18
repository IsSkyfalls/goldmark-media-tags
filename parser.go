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
		alt := string(match[2])
		url := string(match[3])

		init, isValid := tagInitsLUT[flag]
		if !isValid {
			block.Advance(1)
			return nil
		}
		media := Media{
			BaseInline: ast.BaseInline{},
			MediaType:  flag,
			Alt:        alt,
			Link:       url,
		}
		init.initAttributes(&media, p.Options)
		source := init.makeSourceTag(media, p.Options)
		media.AppendChild(&media, source)
		parent.AppendChild(parent, &media)
	}
	return nil
}
