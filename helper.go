package media

import (
	"fmt"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

func DumpAttributes(n ast.Node, source []byte, level int) {
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
