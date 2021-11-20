package media

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		str   string
		_type ast.NodeKind
	}{
		{"!v[video](https://example.org/test.mp4)", kindMedia},
		{"!a[audio](https://example.org/test.mp3)", kindMedia},
		{"!v[](https://example.org/test.mp4)", kindMedia}, //no alt text, still valid
		{"!v[test]()", ast.KindText},                      //no url
		{"!x[test](1)", ast.KindText},                     //wrong type
	}

	md := goldmark.New(goldmark.WithExtensions(Extension{}))
	for _, c := range cases {
		src := []byte(c.str)
		doc := md.Parser().Parse(text.NewReader(src))
		n := doc.FirstChild().FirstChild()
		assert.Equal(t, c._type, n.Kind())
		doc.Dump(src, 0)
	}
}

func TestParsePreserve(t *testing.T) {
	cases := []struct {
		str   string
		pre   string
		after string
	}{
		{"pre!v[video](h)after", "pre", "after"},
		{"pre!v[](h)after", "pre", "after"},    // no alt text, valid
		{"pre!v[]()after", "pre!v", "after"},   // no url, should be handled as Link
		{"pre!x[1](1)after", "pre!x", "after"}, // invalid, should be handled as Link
	}
	md := goldmark.New(goldmark.WithExtensions(Extension{}))
	for _, c := range cases {
		src := []byte(c.str)
		doc := md.Parser().Parse(text.NewReader(src))
		p := doc.FirstChild()
		c1 := p.FirstChild()
		c2 := p.LastChild()
		t1 := string(c1.Text(src))
		t2 := string(c2.Text(src))
		assert.Equal(t, ast.KindText, c1.Kind())
		assert.Equal(t, ast.KindText, c2.Kind())
		assert.Equal(t, c.pre, t1)
		assert.Equal(t, c.after, t2)
	}
}
