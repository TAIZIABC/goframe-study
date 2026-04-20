// render/render.go
// 微信公众号兼容的 HTML 渲染器（内联 CSS）
package render

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// WxRenderer 微信公众号 HTML 渲染器
type WxRenderer struct {
	theme *Theme
}

type Theme struct {
	FontFamily    string
	FontSize      string
	LineHeight    string
	Color         string
	H1Color       string
	H2Color       string
	H3Color       string
	LinkColor     string
	CodeBg        string
	CodeColor     string
	BlockCodeBg   string
	BlockquoteBg  string
	BlockquoteBorder string
	TableBorder   string
	TableHeaderBg string
}

func DefaultTheme() *Theme {
	return &Theme{
		FontFamily:       "system-ui, -apple-system, sans-serif",
		FontSize:         "15px",
		LineHeight:       "1.8",
		Color:            "#333",
		H1Color:          "#1a1a1a",
		H2Color:          "#2b2b2b",
		H3Color:          "#3a3a3a",
		LinkColor:        "#576b95",
		CodeBg:           "#fff5f5",
		CodeColor:        "#c7254e",
		BlockCodeBg:      "#2b2b2b",
		BlockquoteBg:     "#f7f7f7",
		BlockquoteBorder: "#ddd",
		TableBorder:      "#ddd",
		TableHeaderBg:    "#f5f5f5",
	}
}

func NewWxRenderer() renderer.NodeRenderer {
	return &WxRenderer{theme: DefaultTheme()}
}

func (r *WxRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// 块级
	reg.Register(ast.KindDocument, r.renderDocument)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindList, r.renderList)
	reg.Register(ast.KindListItem, r.renderListItem)
	reg.Register(ast.KindThematicBreak, r.renderThematicBreak)
	reg.Register(ast.KindHTMLBlock, r.renderHTMLBlock)
	// 内联
	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindString, r.renderString)
	reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindAutoLink, r.renderAutoLink)
	reg.Register(ast.KindRawHTML, r.renderRawHTML)
}

func (r *WxRenderer) renderDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		t := r.theme
		fmt.Fprintf(w, `<section style="font-family:%s;font-size:%s;line-height:%s;color:%s;padding:10px 0;">`,
			t.FontFamily, t.FontSize, t.LineHeight, t.Color)
	} else {
		w.WriteString("</section>")
	}
	return ast.WalkContinue, nil
}

func (r *WxRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	h := node.(*ast.Heading)
	if entering {
		sizes := map[int]string{1: "24px", 2: "20px", 3: "18px", 4: "16px", 5: "15px", 6: "14px"}
		colors := map[int]string{1: r.theme.H1Color, 2: r.theme.H2Color, 3: r.theme.H3Color}
		size := sizes[h.Level]
		color := colors[h.Level]
		if color == "" {
			color = r.theme.Color
		}
		bottom := "16px"
		if h.Level <= 2 {
			bottom = "20px"
			fmt.Fprintf(w, `<h%d style="font-size:%s;color:%s;font-weight:bold;margin:24px 0 %s 0;padding-bottom:8px;border-bottom:1px solid #eee;">`,
				h.Level, size, color, bottom)
		} else {
			fmt.Fprintf(w, `<h%d style="font-size:%s;color:%s;font-weight:bold;margin:20px 0 %s 0;">`,
				h.Level, size, color, bottom)
		}
	} else {
		fmt.Fprintf(w, "</h%d>", h.Level)
	}
	return ast.WalkContinue, nil
}

func (r *WxRenderer) renderParagraph(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString(`<p style="margin:0 0 16px 0;text-align:justify;">`)
	} else {
		w.WriteString("</p>")
	}
	return ast.WalkContinue, nil
}

func (r *WxRenderer) renderBlockquote(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		fmt.Fprintf(w, `<blockquote style="margin:16px 0;padding:12px 16px;background:%s;border-left:4px solid %s;color:#666;font-size:14px;">`,
			r.theme.BlockquoteBg, r.theme.BlockquoteBorder)
	} else {
		w.WriteString("</blockquote>")
	}
	return ast.WalkContinue, nil
}

func (r *WxRenderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		fmt.Fprintf(w, `<pre style="background:%s;color:#f8f8f2;padding:16px;border-radius:6px;font-size:13px;line-height:1.6;overflow-x:auto;margin:16px 0;font-family:Menlo,Monaco,Consolas,monospace;white-space:pre-wrap;word-wrap:break-word;">`,
			r.theme.BlockCodeBg)
		w.WriteString("<code>")
		lines := node.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			w.Write(escapeHTML(line.Value(source)))
		}
		w.WriteString("</code></pre>")
	}
	return ast.WalkContinue, nil
}

func (r *WxRenderer) renderList(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	list := node.(*ast.List)
	if entering {
		if list.IsOrdered() {
			w.WriteString(`<ol style="margin:0 0 16px 0;padding-left:24px;">`)
		} else {
			w.WriteString(`<ul style="margin:0 0 16px 0;padding-left:24px;">`)
		}
	} else {
		if list.IsOrdered() {
			w.WriteString("</ol>")
		} else {
			w.WriteString("</ul>")
		}
	}
	return ast.WalkContinue, nil
}

func (r *WxRenderer) renderListItem(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString(`<li style="margin:4px 0;">`)
	} else {
		w.WriteString("</li>")
	}
	return ast.WalkContinue, nil
}

func (r *WxRenderer) renderThematicBreak(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString(`<hr style="border:none;border-top:1px solid #eee;margin:24px 0;">`)
	}
	return ast.WalkContinue, nil
}

func (r *WxRenderer) renderHTMLBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		lines := node.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			w.Write(line.Value(source))
		}
	}
	return ast.WalkContinue, nil
}

func (r *WxRenderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.Text)
		w.Write(escapeHTML(n.Segment.Value(source)))
		if n.SoftLineBreak() {
			w.WriteString("\n")
		}
		if n.HardLineBreak() {
			w.WriteString("<br>")
		}
	}
	return ast.WalkContinue, nil
}

func (r *WxRenderer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.String)
		w.Write(n.Value)
	}
	return ast.WalkContinue, nil
}

func (r *WxRenderer) renderCodeSpan(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		fmt.Fprintf(w, `<code style="background:%s;color:%s;padding:2px 6px;border-radius:3px;font-size:90%%;font-family:Menlo,Monaco,monospace;">`,
			r.theme.CodeBg, r.theme.CodeColor)
	} else {
		w.WriteString("</code>")
	}
	return ast.WalkContinue, nil
}

func (r *WxRenderer) renderEmphasis(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	e := node.(*ast.Emphasis)
	if e.Level == 2 {
		if entering {
			w.WriteString(`<strong style="font-weight:bold;color:#1a1a1a;">`)
		} else {
			w.WriteString("</strong>")
		}
	} else {
		if entering {
			w.WriteString(`<em style="font-style:italic;">`)
		} else {
			w.WriteString("</em>")
		}
	}
	return ast.WalkContinue, nil
}

func (r *WxRenderer) renderLink(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.Link)
		fmt.Fprintf(w, `<a href="%s" style="color:%s;text-decoration:none;border-bottom:1px solid %s;">`,
			n.Destination, r.theme.LinkColor, r.theme.LinkColor)
	} else {
		w.WriteString("</a>")
	}
	return ast.WalkContinue, nil
}

func (r *WxRenderer) renderImage(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.Image)
		fmt.Fprintf(w, `<img src="%s" alt="%s" style="max-width:100%%;border-radius:4px;margin:8px 0;">`,
			n.Destination, n.Title)
	}
	return ast.WalkContinue, nil
}

func (r *WxRenderer) renderAutoLink(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.AutoLink)
		url := n.URL(nil)
		fmt.Fprintf(w, `<a href="%s" style="color:%s;">%s</a>`, url, r.theme.LinkColor, url)
	}
	return ast.WalkContinue, nil
}

func (r *WxRenderer) renderRawHTML(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.RawHTML)
		for i := 0; i < n.Segments.Len(); i++ {
			seg := n.Segments.At(i)
			w.Write(seg.Value(source))
		}
	}
	return ast.WalkContinue, nil
}

func escapeHTML(b []byte) []byte {
	s := string(b)
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return []byte(s)
}
