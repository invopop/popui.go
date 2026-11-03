package popui

import (
	"bytes"
	"context"
	"io"

	"github.com/a-h/templ"

	"github.com/invopop/popui.go/props"
	"github.com/invopop/popui.go/tailwind"
)

// detailsInlineKey marks the context of inline Details' children.
type detailsInlineKey struct{}

// detailsInline renders its children as the details of inline Details
// when inline is set.
func detailsInline(inline bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		children := templ.GetChildren(ctx)
		ctx = templ.ClearChildren(ctx)
		if inline {
			ctx = context.WithValue(ctx, detailsInlineKey{}, true)
		}
		return children.Render(ctx, w)
	})
}

// isDetailsInline reports whether a detail is in inline Details.
func isDetailsInline(ctx context.Context) bool {
	inline, _ := ctx.Value(detailsInlineKey{}).(bool)
	return inline
}

// Detail is a label and its value, in Details. Use Label and Value for
// text, or give a component, such as a TagStatus, as the children instead
// of a Value. A detail with neither shows a dash.
//
// Details in inline Details, and those that set Inline or Actions, are rows
// on the popui-detail-* styles: a fixed-width label column, a hover-tinted
// value cell filling the rest of the row, and hover-revealed action buttons
// pinned to the row's right edge. Without Inline the label is above the
// value.
func Detail(opts ...props.Detail) templ.Component {
	return renderDetail(props.First(opts), false)
}

// DescriptionList provides a wrapper for description lists (dl element),
// 400px wide with labels above values unless Layout is set.
//
// Deprecated: use Details, which fill their container and are inline by
// default.
func DescriptionList(opts ...props.DescriptionList) templ.Component {
	p := props.First(opts)
	if p.Layout == "" {
		p.Layout = props.DetailsLayoutStacked
		p.Class = tailwind.Merge("w-[400px]", p.Class)
	}
	return Details(p)
}

// DescriptionListItem provides a wrapper for description list items. Its
// children are placed between the label and the value, so they can be its
// own DT and DD, unless it's shown as a row.
//
// Deprecated: use Detail, whose children are its value.
func DescriptionListItem(opts ...props.DescriptionListItem) templ.Component {
	return renderDetail(props.First(opts), true)
}

func renderDetail(p props.Detail, legacy bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		p.Inline = p.Inline || isDetailsInline(ctx)
		// The children are rendered first to tell whether there are any,
		// as a detail without a value shows a dash instead.
		var content bytes.Buffer
		if err := templ.GetChildren(ctx).Render(templ.ClearChildren(ctx), &content); err != nil {
			return err
		}
		has := len(bytes.TrimSpace(content.Bytes())) > 0
		return detail(p, templ.Raw(content.String()), has, legacy).Render(ctx, w)
	})
}
