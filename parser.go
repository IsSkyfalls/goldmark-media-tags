package media

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"regexp"
)

var regex = regexp.MustCompile(`^!(.)\[(.*?)]\((.+?)\)`) // only from the trigger position by prepending ^

// mediaParser implements parser.InlineParser interface
type mediaParser struct {
	Options
}

func (p mediaParser) Trigger() []byte {
	return []byte{'!'}
}

func (p mediaParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	match := MatchAndAdvance(block, regex)
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
		return &media
	}
	block.Advance(1)
	return nil
}

func MatchAndAdvance(reader text.Reader, reg *regexp.Regexp) []string {
	// workaround: block.FindSubMatch spits out completely wrong values when provided with unicode characters
	oldline, oldseg := reader.Position()
	matches := reg.FindReaderSubmatchIndex(reader)
	if matches == nil {
		return nil
	}
	reader.SetPosition(oldline, oldseg)
	segments := make([]string, 0, len(matches)/2)
	line, _ := reader.PeekLine()

	for i := 0; i < len(matches)/2; i++ {
		s := line[matches[i*2]:matches[i*2+1]]
		segments = append(segments, string(s))
	}
	reader.Advance(matches[1])
	return segments
}
