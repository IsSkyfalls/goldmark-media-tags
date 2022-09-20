package media

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"testing"
)

func TestMedia_Dump(t *testing.T) {
	media := Media{
		BaseInline: ast.BaseInline{},
		MediaType:  TypeVideo,
		Alt:        "video of mount Tai",
		Link:       "https://tai.com/video.mp4",
	}
	init := tagInitsLUT['v']
	opt := Options{
		MediaControls: true,
		MediaAutoplay: false,
		MediaLoop:     true,
		MediaMuted:    false,
		MediaPreload:  "metadata",
	}
	init.initAttributes(&media, opt)
	source := init.makeSourceTag(media, opt)
	media.AppendChild(&media, source)

	source1 := init.makeSourceTag(media, Options{}).(*TagSourceSource) // note the pointer(*)
	source1.Src = "https://transformation.com/640x480/video.mp4"
	media.AppendChild(&media, source1)

	media.Dump([]byte{}, 0)

	md := goldmark.New(goldmark.WithExtensions(Extension{}))
	buf := bytes.NewBuffer([]byte{})
	err := md.Renderer().Render(buf, []byte{}, &media)
	assert.NoError(t, err)
	fmt.Println(buf.String())
}
