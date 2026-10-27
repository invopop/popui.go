// Package playground provides the experimental jsonui playground: a page
// that renders JSON documents to PopUI HTML in the browser with the jsonui
// WebAssembly build. It is served by `popui serve` and published by
// `popui build` under Path.
package playground

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/invopop/popui.go/jsonui"
	"github.com/invopop/popui.go/jsonui/schema"
)

// Path is where the playground is mounted on the docs site.
const Path = "/jsonui/"

// actions are the actions documents may use in the playground: the built-in
// ones plus "save", which the playground page registers itself.
var actions = append(append([]string{}, jsonui.BuiltinActions...), "save")

//go:embed index.html samples.js editor.js
var static embed.FS

// files returns the published files that don't need compiling.
func files() (map[string][]byte, error) {
	out := map[string][]byte{"schema.json": schema.Document}
	for _, name := range []string{"index.html", "samples.js", "editor.js"} {
		data, err := static.ReadFile(name)
		if err != nil {
			return nil, err
		}
		out[name] = data
	}
	runtime, err := fs.ReadFile(jsonui.Assets, "assets/jsonui.js")
	if err != nil {
		return nil, err
	}
	out["jsonui.js"] = runtime
	goroot, err := goEnv("GOROOT")
	if err != nil {
		return nil, err
	}
	if out["wasm_exec.js"], err = os.ReadFile(filepath.Join(goroot, "lib", "wasm", "wasm_exec.js")); err != nil {
		return nil, err
	}
	return out, nil
}

// Publish writes the playground, including the compiled WebAssembly module,
// to dir. It needs the Go toolchain and must run from the module root.
func Publish(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	fsFiles, err := files()
	if err != nil {
		return err
	}
	for name, data := range fsFiles {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
			return err
		}
	}
	return BuildWASM(filepath.Join(dir, "popui.wasm"))
}

// BuildWASM compiles the jsonui renderer to WebAssembly at out.
func BuildWASM(out string) error {
	cmd := exec.Command("go", "build", "-trimpath", "-ldflags=-s -w", "-o", out, "./jsonui/wasm")
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("building jsonui wasm: %w: %s", err, stderr.String())
	}
	return nil
}

func goEnv(key string) (string, error) {
	out, err := exec.Command("go", "env", key).Output()
	if err != nil {
		return "", fmt.Errorf("go env %s: %w", key, err)
	}
	return strings.TrimSpace(string(out)), nil
}

// Handler serves the playground for local development. The WebAssembly
// module is compiled from the current source on first request. It also
// offers POST render, a server-side render of the same document for
// comparison, which the static published site doesn't have.
func Handler() http.Handler {
	var (
		once    sync.Once
		wasm    []byte
		wasmErr error
	)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /popui.wasm", func(w http.ResponseWriter, _ *http.Request) {
		once.Do(func() {
			dir, err := os.MkdirTemp("", "popui-jsonui-")
			if err != nil {
				wasmErr = err
				return
			}
			defer os.RemoveAll(dir) //nolint:errcheck
			out := filepath.Join(dir, "popui.wasm")
			if wasmErr = BuildWASM(out); wasmErr == nil {
				wasm, wasmErr = os.ReadFile(out)
			}
		})
		if wasmErr != nil {
			http.Error(w, wasmErr.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/wasm")
		w.Write(wasm) //nolint:errcheck
	})
	mux.HandleFunc("POST /render", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		opts := jsonui.Options{Actions: actions, Trusted: r.URL.Query().Get("trusted") == "1"}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := opts.Render(r.Context(), body, w); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
	})
	mux.HandleFunc("GET /{file}", func(w http.ResponseWriter, r *http.Request) {
		serveFile(w, r, r.PathValue("file"))
	})
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		serveFile(w, r, "index.html")
	})
	return http.StripPrefix(strings.TrimSuffix(Path, "/"), mux)
}

func serveFile(w http.ResponseWriter, r *http.Request, name string) {
	fsFiles, err := files()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, ok := fsFiles[name]
	if !ok {
		http.NotFound(w, r)
		return
	}
	types := map[string]string{
		".html": "text/html; charset=utf-8",
		".js":   "text/javascript; charset=utf-8",
		".json": "application/schema+json",
	}
	w.Header().Set("Content-Type", types[filepath.Ext(name)])
	w.Write(data) //nolint:errcheck
}
