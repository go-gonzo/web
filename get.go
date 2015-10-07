package web

import (
	"fmt"
	"mime"
	"net/http"
	"path"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"

	"github.com/omeid/gonzo"
)

func get(ctx context.Context, client *http.Client, url string) (gonzo.File, error) {

	resp, err := ctxhttp.Get(ctx, client, url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, fmt.Errorf("%s (%s)", resp.Status, url)
	}

	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	name, ok := params["filename"]
	if !ok || err != nil {
		name = path.Base(url)
	}

	file := gonzo.NewFile(resp.Body, gonzo.NewFileInfo())
	file.FileInfo().SetName(name)
	file.FileInfo().SetSize(resp.ContentLength)

	return file, nil
}
