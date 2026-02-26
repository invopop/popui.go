package components

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/invopop/icons"
)

// Icons renders a grid of all available icons by iterating the icons registry.
func Icons() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<div style="display: grid; grid-template-columns: repeat(auto-fill, minmax(80px, 1fr)); gap: 12px;">`); err != nil {
			return err
		}
		for _, def := range icons.Defs {
			if err := iconCard(def.Name, def.Component).Render(ctx, w); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(w, `</div>`); err != nil {
			return err
		}
		return nil
	})
}
