//Package http provides helpful http functions.
package http

import (
	"fmt"
	"mime"
	"net/http"
	"path"

	"github.com/omeid/gonzo"
)

func name(url string, response *http.Response) string {

	_, params, err := mime.ParseMediaType(response.Header.Get("Content-Disposition"))

	name, ok := params["filename"]
	if !ok || err != nil {
		name = path.Base(url)
	}

	return name
}

func Get(url string) (gonzo.File, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, fmt.Errorf("%s (%s)", resp.Status, url)
	}

	name := name(url, resp)

	file := gonzo.NewFile(resp.Body, gonzo.NewFileInfo())
	file.FileInfo().SetName(name)
	file.FileInfo().SetSize(resp.ContentLength)

	return file, nil
}
