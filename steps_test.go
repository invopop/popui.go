package popui

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/invopop/popui.go/props"
)

// TestSteps checks that steps say their status, in words too for screen
// readers, except the current one, which aria-current announces, that a
// completed step has a tick in place of its number, and that the list of sub-steps has
// nothing in it, whitespace included, for a step without any, as it is
// hidden while empty.
func TestSteps(t *testing.T) {
	var buf bytes.Buffer
	children := templ.Join(
		Step(props.Step{Label: "Details", Status: props.StepStatusCompleted}),
		Step(props.Step{Label: "Sign", Status: props.StepStatusCurrent}),
		Step(props.Step{Label: "Done", StatusLabel: "Pendiente"}),
	)
	ctx := templ.WithChildren(context.Background(), children)
	if err := Steps(props.Steps{Label: "Progress", Orientation: props.StepsOrientationResponsive}).Render(ctx, &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{
		`<ol aria-label="Progress" data-orientation="responsive" class="popui-steps">`,
		`<li data-status="completed" class="popui-step">`,
		`<li data-status="current" aria-current="step" class="popui-step">`,
		`<li class="popui-step">`,
		`<span class="popui-step-marker" aria-hidden="true"></span>`,
		`<span class="popui-step-label">Sign</span> </div>`,
		`<span class="popui-step-label">Details</span> <span class="popui-step-status">Completed</span>`,
		`<span class="popui-step-label">Done</span> <span class="popui-step-status">Pendiente</span>`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in:\n%s", want, out)
		}
	}
	if n := strings.Count(out, "<svg"); n != 1 {
		t.Errorf("want 1 tick, got %d in:\n%s", n, out)
	}
	if n := strings.Count(out, `<ol class="popui-substeps"></ol>`); n != 3 {
		t.Errorf("want 3 empty lists of sub-steps, got %d in:\n%s", n, out)
	}
}

// TestSubSteps checks that sub-steps are listed inside their step and
// that the current one is marked.
func TestSubSteps(t *testing.T) {
	var buf bytes.Buffer
	children := templ.Join(
		SubStep(props.SubStep{Label: "Representation", Current: true}),
		SubStep(props.SubStep{Label: "Identity"}),
	)
	ctx := templ.WithChildren(context.Background(), children)
	if err := Step(props.Step{Label: "Sign", Status: props.StepStatusCurrent}).Render(ctx, &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	want := `<ol class="popui-substeps"><li data-current aria-current="step" class="popui-substep">Representation</li><li class="popui-substep">Identity</li></ol></li>`
	if !strings.Contains(out, want) {
		t.Errorf("missing %s in:\n%s", want, out)
	}
}
