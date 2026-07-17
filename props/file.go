package props

import (
	"maps"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

// FileDownload Templ component props. FileDownload displays details for
// an existing file with space for actions.
type FileDownload struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Borderless removes the border for embedding the file row directly
	// inside another container.
	Borderless bool
	// Hover highlights the whole row with a background color on hover.
	Hover   bool
	Preview string
}

// PreviewAttributes returns the component attributes with an onclick handler that runs the Preview JavaScript expression, ignoring clicks on nested links and buttons.
func (f FileDownload) PreviewAttributes() templ.Attributes {
	if f.Preview == "" {
		return f.Attributes
	}
	attrs := templ.Attributes{}
	maps.Copy(attrs, f.Attributes)
	attrs["onclick"] = "if (!event.target.closest('a,button')) { " + f.Preview + " }"
	return attrs
}

// FileDownloadInfo Templ component props
type FileDownloadInfo struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Label      string
	// Value renders as smaller, muted text below the Label, e.g. a
	// timestamp or file size.
	Value string
}

// InputFile defines the properties for the InputFile and FileUpload components.
type InputFile struct {
	ID            string
	Class         string
	Attributes    templ.Attributes
	Name          string
	Accept        string
	Capture       string
	Multiple      bool
	Autofocus     bool
	Required      bool
	Disabled      bool
	AvatarURL     string
	AvatarAlt     string
	Text          string
	PreviewSquare bool
}

// GenerateID generates a unique ID if not provided.
func (i InputFile) GenerateID() InputFile {
	if i.ID == "" {
		i.ID = "input-file-" + uuid.NewString()
	}
	return i
}
