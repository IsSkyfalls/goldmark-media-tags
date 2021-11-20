package media

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"testing"
)

func TestRenderVideo(t *testing.T) {
	md := goldmark.New(goldmark.WithExtensions(Extension{}))
	cases := []struct {
		node     Media
		expected string
	}{
		{
			Media{ //bare minimum video
				BaseInline: ast.BaseInline{},
				IsVideo:    true,
				Sources:    Sources{{"example.org", ""}},
			}, "<video><source src=\"example.org\"></video>",
		},
		{
			Media{ //bare minimum audio
				BaseInline: ast.BaseInline{},
				IsVideo:    false,
				Sources:    Sources{{"example.org", ""}},
			}, "<audio><source src=\"example.org\"></audio>",
		}, {
			Media{ //full video
				BaseInline: ast.BaseInline{},
				Controls:   true,
				Autoplay:   true,
				Loop:       true,
				Muted:      true,
				IsVideo:    true,
				Preload:    Auto,
				Sources:    Sources{{"example.org", "mp4"}},
			}, "<video controls autoplay loop muted><source src=\"example.org\" type=\"mp4\"></video>",
		}, {
			Media{ //full audio
				BaseInline: ast.BaseInline{},
				Controls:   true,
				Autoplay:   true,
				Loop:       true,
				Muted:      true,
				IsVideo:    false,
				Preload:    Auto,
				Sources:    Sources{{"example.org", "mp4"}},
			}, "<audio controls autoplay loop muted><source src=\"example.org\" type=\"mp4\"></audio>",
		},
	}

	for _, c := range cases {
		c.node.Dump([]byte{}, 0)
		buf := bytes.NewBuffer([]byte{})
		err := md.Renderer().Render(buf, []byte{}, &c.node)
		assert.NoError(t, err)
		assert.Equal(t, c.expected, buf.String())
		fmt.Println(">>>", buf)
	}
}
