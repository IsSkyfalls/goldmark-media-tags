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
				MediaType:  Video,
				Sources:    Sources{{"example.org", "", "", "", ""}},
			}, "<video><source src=\"example.org\"></video>",
		},
		{
			Media{ //bare minimum Audio
				BaseInline: ast.BaseInline{},
				MediaType:  Audio,
				Sources:    Sources{{"example.org", "", "", "", ""}},
			}, "<audio><source src=\"example.org\"></audio>",
		}, {
			Media{ //full Video
				BaseInline: ast.BaseInline{},
				Controls:   true,
				Autoplay:   true,
				Loop:       true,
				Muted:      true,
				MediaType:  Video,
				Preload:    Auto,
				Sources:    Sources{{"example.org", "(max-width:480px)", "", "mp4", ""}},
			}, "<video controls autoplay loop muted preload=\"auto\"><source media=\"(max-width:480px)\" src=\"example.org\" type=\"mp4\"></video>",
		}, {
			Media{ //full Audio
				BaseInline: ast.BaseInline{},
				Controls:   true,
				Autoplay:   true,
				Loop:       true,
				Muted:      true,
				MediaType:  Audio,
				Preload:    Auto,
				Sources:    Sources{{"example.org", "", "", "mp3", ""}},
			}, "<audio controls autoplay loop muted preload=\"auto\"><source src=\"example.org\" type=\"mp3\"></audio>",
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
