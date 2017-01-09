package web

import (
	"fmt"
	"mime"
	"net/http"
	"path"

	"context"

	"github.com/omeid/gonzo"
)

func get(ctx context.Context, client *http.Client, url string) (gonzo.File, error) {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := client.Do(req)
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
