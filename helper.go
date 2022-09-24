package media

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
	"sort"
)

func dumpAttributes(n ast.Node, source []byte, level int) {
	attrs := n.Attributes()
	list := map[string]string{}
	for _, attr := range attrs {
		name := util.BytesToReadOnlyString(attr.Name)
		// fallback for values that are not []byte
		s := fmt.Sprintf("%v", attr.Value)
		value := util.BytesToReadOnlyString(util.EscapeHTML([]byte(s)))
		list[name] = value
	}
	ast.DumpHelper(n, source, level, list, nil)
}

// renderTagWithAttributesNoClosing will render a self-closing tag ending in />
func renderTagWithAttributesSelfClose(n ast.Node, tagName string) string {
	return renderTagWithAttributesNoClosing(n, tagName) + " />"
}

// renderTagWithAttributesNoClosing will render a tag without closing it with >
func renderTagWithAttributesNoClosing(n ast.Node, tagName string) string {
	attrs := n.Attributes()
	sort.Slice(attrs, func(i, j int) bool {
		return bytes.Compare(attrs[i].Name, attrs[j].Name) > 0
	})
	html := "<" + tagName
	for _, e := range attrs {
		var value string
		if s, isStr := e.Value.(string); isStr {
			if s == "" {
				continue
			}
		} else if b, isBool := e.Value.(bool); isBool {
			if !b {
				continue
			}
			// The values “true” and “false” are not allowed on boolean attributes.
			// To represent a false value, the attribute has to be omitted altogether.
			html += " " + string(e.Name)
			continue
		}
		value = fmt.Sprintf("%s", e.Value)
		escaped := util.BytesToReadOnlyString(util.EscapeHTML([]byte(value)))
		html += " " + string(e.Name)
		html += "=\"" + string(util.EscapeHTML([]byte(escaped))) + "\""
	}
	return html
}
