// render/table.go
// 表格渲染扩展（goldmark 默认不启用表格，需要扩展）
package render

import (
	"fmt"

	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type WxTableRenderer struct {
	theme *Theme
}

func NewWxTableRenderer() renderer.NodeRenderer {
	return &WxTableRenderer{theme: DefaultTheme()}
}

func (r *WxTableRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(east.KindTable, r.renderTable)
	reg.Register(east.KindTableHeader, r.renderTableHeader)
	reg.Register(east.KindTableRow, r.renderTableRow)
	reg.Register(east.KindTableCell, r.renderTableCell)
}

func (r *WxTableRenderer) renderTable(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		fmt.Fprintf(w, `<table style="width:100%%;border-collapse:collapse;margin:16px 0;font-size:14px;border:1px solid %s;">`,
			r.theme.TableBorder)
	} else {
		w.WriteString("</table>")
	}
	return ast.WalkContinue, nil
}

func (r *WxTableRenderer) renderTableHeader(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString("<thead>")
	} else {
		w.WriteString("</thead>")
	}
	return ast.WalkContinue, nil
}

func (r *WxTableRenderer) renderTableRow(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString("<tr>")
	} else {
		w.WriteString("</tr>")
	}
	return ast.WalkContinue, nil
}

func (r *WxTableRenderer) renderTableCell(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	cell := node.(*east.TableCell)
	tag := "td"
	bg := ""
	if cell.Parent() != nil && cell.Parent().Parent() != nil {
		if _, ok := cell.Parent().Parent().(*east.TableHeader); ok {
			tag = "th"
			bg = fmt.Sprintf("background:%s;", r.theme.TableHeaderBg)
		}
	}

	align := ""
	switch cell.Alignment {
	case east.AlignLeft:
		align = "text-align:left;"
	case east.AlignCenter:
		align = "text-align:center;"
	case east.AlignRight:
		align = "text-align:right;"
	}

	if entering {
		fmt.Fprintf(w, `<%s style="padding:8px 12px;border:1px solid %s;%s%s">`,
			tag, r.theme.TableBorder, bg, align)
	} else {
		fmt.Fprintf(w, "</%s>", tag)
	}
	return ast.WalkContinue, nil
}
