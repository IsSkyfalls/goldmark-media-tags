package media

// Preload represent the "preload" attribute of the <video> and <audio> elements.
// The default value(if left empty) is Auto
// https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/audio#attr-preload
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/video#attr-preload
type Preload string

const (
	None     = "none"
	Metadata = "metadata"
	Auto     = "auto"
)
