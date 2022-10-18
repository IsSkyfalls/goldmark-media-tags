package media

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark/ast"
	"testing"
)

func TestAttributeEscape(t *testing.T) {
	tag := &TagSourceSource{
		BaseInline: ast.BaseInline{},
		Src:        "\"'>,<",
		SrcSet:     "",
	}
	tag.SetAttributeString("enabled", true)
	tag.SetAttributeString("autoplay", false)
	tag.SetAttributeString("bytes", []byte("bytes 統一碼<"))
	tag.updateAttributes()
	s := renderTagWithAttributesNoClosing(tag, "source")
	assert.Equal(t, "<source src=\"&amp;quot;'&amp;gt;,&amp;lt;\" enabled bytes=\"bytes 統一碼&amp;lt;\"", s)
}
