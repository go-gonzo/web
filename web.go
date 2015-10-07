package web

import (
	"net/http"

	"github.com/omeid/gonzo"
	"github.com/omeid/gonzo/context"
)

// Gets  the list of urls and passes the results to output channel.
// It reports the progress to the Context using a ReadProgress proxy.
func Get(ctx context.Context, urls ...string) gonzo.Pipe {

	out := make(chan gonzo.File)
	client := &http.Client{}
	go func() {
		defer close(out)

		for _, url := range urls {

			select {
			case <-ctx.Done():
				ctx.Warn(context.Canceled)
				return
			default:
				ctx.Infof("Downloading %s", url)

				file, err := get(ctx, client, url)
				if err != nil {
					ctx.Error(err)
					break
				}

				//TODO: Add progress meter.
				//s, _ := file.Stat()
				//file.Reader = c.ReadProgress(file.Reader, "Downloading "+file.Path, s.Size())
				out <- file
			}
		}
	}()

	return gonzo.NewPipe(ctx, out)
}
