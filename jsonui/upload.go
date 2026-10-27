package jsonui

import (
	"fmt"
	"reflect"

	"github.com/a-h/templ"
	"github.com/invopop/popui.go/alpine"
	"github.com/invopop/popui.go/props"
)

const uploadFileType = "upload_file"

// UploadFile are the props of the "upload_file" node, a file picker row with
// a preview image and Upload/Remove buttons. The node's "model" is the data
// path holding the file URL, and its "on" handlers take "select" and
// "remove" actions, e.g.
//
//	{
//	  "type": "upload_file",
//	  "model": "config.logo_url",
//	  "props": {"loading": "ui.uploading"},
//	  "on": {
//	    "select": {"action": "asset.choose", "accept": "image/*", "into": "config.logo_url", "busy": "ui.uploading"},
//	    "remove": {"action": "set", "path": "config.logo_url", "value": ""}
//	  }
//	}
type UploadFile struct {
	// Loading is the data path that is true while a file is being chosen or
	// uploaded, e.g. the "busy" path of the select action.
	Loading string `json:"loading,omitempty"`
}

// uploadEvents are the "on" handlers an upload_file node accepts.
var uploadEvents = map[string]bool{"select": true, "remove": true}

func init() {
	// Built by docBuilder.uploadFile; the entry provides the props type for
	// Components and the schema.
	registry[uploadFileType] = entry{props: reflect.TypeFor[UploadFile]()}
	modelProps[uploadFileType] = ""
}

// uploadFile builds alpine.UploadFile, whose props are Alpine expressions,
// from the node's data paths and actions so documents never supply script.
func (b *docBuilder) uploadFile(n *Node) (templ.Component, error) {
	var p UploadFile
	if len(n.Props) > 0 {
		if err := decodeStrict(n.Props, &p); err != nil {
			return nil, fmt.Errorf("%s props: %w", n.Type, err)
		}
	}
	if len(n.Bind) > 0 || len(n.Children) > 0 {
		return nil, fmt.Errorf("%s: bind and children are not supported", n.Type)
	}
	if n.Model == "" {
		return nil, fmt.Errorf("%s: model is required", n.Type)
	}
	if b.data == nil {
		return nil, fmt.Errorf("%s model: document data must be set", n.Type)
	}
	if !dataPath.MatchString(n.Model) {
		return nil, fmt.Errorf("%s model: %q is not a data path", n.Type, n.Model)
	}
	if _, err := lookup(b.data, n.Model); err != nil {
		return nil, fmt.Errorf("%s model: %w", n.Type, err)
	}
	loading := "false"
	if p.Loading != "" {
		if !dataPath.MatchString(p.Loading) {
			return nil, fmt.Errorf("%s loading: %q is not a data path", n.Type, p.Loading)
		}
		loading = p.Loading
	}
	handlers := map[string]string{}
	for event, raw := range n.On {
		if !uploadEvents[event] {
			return nil, fmt.Errorf("%s on: unknown event %q, expected select or remove", n.Type, event)
		}
		h, err := b.actions(event, raw)
		if err != nil {
			return nil, fmt.Errorf("%s on: %w", n.Type, err)
		}
		handlers[event] = h
	}
	return alpine.UploadFile(props.UploadFile{
		XModel:   n.Model,
		Loading:  loading,
		OnSelect: handlers["select"],
		OnRemove: handlers["remove"],
	}), nil
}
