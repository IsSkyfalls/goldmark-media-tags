package media

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type Options struct {
	// MediaControls default is true, affects <audio> and <video>
	MediaControls bool
	// MediaAutoplay default is false, affects <audio> and <video>
	MediaAutoplay bool
	// MediaLoop default is false, affects <audio> and <video>
	MediaLoop bool
	// MediaMuted default is false, affects <audio> and <video>
	MediaMuted bool
	// MediaPreload default is empty string(""). Affects <audio> and <video>
	MediaPreload string
}

type Extension struct {
	Options
}

func WithDefaults() Extension {
	return Extension{
		Options: Options{
			MediaControls: true,
			MediaPreload:  "",
		},
	}
}

func WithOptions(options Options) Extension {
	return Extension{
		Options: options,
	}
}

func (e Extension) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(mediaParser{e.Options}, 100)))
	md.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(mediaHTMLRenderer{}, 100),
	))
}
