package media

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"testing"
)

func TestRender(t *testing.T) {
	md := goldmark.New(goldmark.WithExtensions(Extension{}))
	cases := []struct {
		node     Media
		expected string
	}{
		{
			Media{ //bare minimum video
				BaseInline: ast.BaseInline{},
				MediaType:  Video,
				Sources:    Sources{{"example.org", "", "", "", "", false}},
			}, "<video><source src=\"example.org\"></video>",
		},
		{
			Media{ //bare minimum audio
				BaseInline: ast.BaseInline{},
				MediaType:  Audio,
				Sources:    Sources{{"example.org", "", "", "", "", false}},
			}, "<audio><source src=\"example.org\"></audio>",
		}, {
			node: Media{ //bare minimum picture
				BaseInline: ast.BaseInline{},
				MediaType:  Picture,
				Sources:    Sources{{"example.org", "", "", "", "", true}},
			},
			expected: "<picture><img src=\"example.org\"></picture>",
		}, {
			Media{ //full video
				BaseInline: ast.BaseInline{},
				Controls:   true,
				Autoplay:   true,
				Loop:       true,
				Muted:      true,
				Alt:        "THIS_SHOULD_NOT_APPEAR",
				Preload:    Auto,
				MediaType:  Video,
				Sources:    Sources{{"example.org", "INVALID", "mp4", "INVALID", "INVALID", false}},
			}, "<video controls autoplay loop muted preload=\"auto\"><source src=\"example.org\" type=\"mp4\"></video>",
		}, {
			Media{ //full audio
				BaseInline: ast.BaseInline{},
				Controls:   true,
				Autoplay:   true,
				Loop:       true,
				Muted:      true,
				Alt:        "THIS_SHOULD_NOT_APPEAR",
				MediaType:  Audio,
				Preload:    Auto,
				Sources:    Sources{{"example.org", "INVALID", "mp3", "INVALID", "INVALID", false}},
			}, "<audio controls autoplay loop muted preload=\"auto\"><source src=\"example.org\" type=\"mp3\"></audio>",
		}, {
			Media{ //full picture
				BaseInline: ast.BaseInline{},
				/* INVALID */
				Controls: true,
				Autoplay: true,
				Loop:     true,
				Muted:    true,
				Preload:  Auto,
				// VALID
				Alt:       "IMAGE_ALT",
				MediaType: Picture,
				Sources: Sources{
					{"example1.org", "VALID_SIZE", "image/png", "VALID_SRCSET", "VALID_MEDIA", false},
					{"example2.org", "VALID_SIZE", "image/png", "VALID_SRCSET", "VALID_MEDIA", true},
				},
			}, "<picture><source src=\"example1.org\" type=\"image/png\" media=\"VALID_MEDIA\" sizes=\"VALID_SIZE\" srcset=\"VALID_SRCSET\"><img alt=\"IMAGE_ALT\" src=\"example2.org\" type=\"image/png\" media=\"VALID_MEDIA\" sizes=\"VALID_SIZE\" srcset=\"VALID_SRCSET\"></picture>",
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
