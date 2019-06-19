package url

import (
	"errors"
	"flag"
	"golang.org/x/net/html"
	"io"
	"net/url"
)

func appendHrefToLinks(n *html.Node, links *[]string) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				*links = append(*links, a.Val)
			}
		}
	}
}

// FetchLinks holt alle Links aus einem Html reader
func FetchLinks(httpBody io.Reader) ([]string, error) {
	links := make([]string, 0)

	doc, err := html.Parse(httpBody)
	if err != nil {
		return nil, err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		appendHrefToLinks(n, &links)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links, nil
}

// ParseURL parsed eine URL, die aus einem Link von einer Seite gewonnen wurde.
// Wenn raw eine absolute URL ist, wird diese genommen. Sonst wird
// versucht, aus der base und dem relativen raw Teil eine absolute URL zusammen
// zu setzen.
func ParseURL(baseRaw, raw string) (*url.URL, error) {
	uri, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	if uri.IsAbs() {
		return uri, nil
	}

	base, err := url.Parse(baseRaw)
	if err != nil {
		return nil, err
	} else if !base.IsAbs() {
		return nil, errors.New("Weder base, noch raw URL ist absolut")
	}

	return base.ResolveReference(uri), nil
}

// ReadURLFromCommandLine Liest eine URL aus den Command line Arguments
func ReadURLFromCommandLine() (string, error) {
	url := flag.String("url", "", "Die Start URL")
	flag.Parse()

	if url == nil || *url == "" {
		return "", errors.New("Die URL ist leer")
	}

	return *url, nil
}
