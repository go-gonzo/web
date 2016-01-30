package web

import (
	"net/http"

	"github.com/omeid/gonzo"
	"github.com/omeid/gonzo/context"
)

// Gets  the list of urls and passes the results to output channel.
// It reports the progress to the Context using a ReadProgress proxy.
func Get(ctx context.Context, urls ...string) gonzo.Pipe {

	ctx, cancel := context.WithCancel(ctx)
	out := make(chan gonzo.File)
	client := &http.Client{}
	go func() {
		defer close(out)

		for _, url := range urls {

			if url == "" {
				ctx.Error("Empty URL.")
				cancel()
				return
			}
			select {
			case <-ctx.Done():
				ctx.Warn(context.Canceled)
				return
			default:
				ctx.Infof("Downloading %s", url)

				file, err := get(ctx, client, url)
				if err != nil {
					ctx.Error(err)
					cancel()
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
