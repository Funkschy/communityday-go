package url

import (
	"fmt"
	"strings"
	"testing"
)

func TestFetchLinks(t *testing.T) {
	html := "<html><body><a href=\"test.com\"/><a href=\"test.com\"/></body></html>"
	links, err := FetchLinks(strings.NewReader(html))
	if err != nil {
		t.Error("Konnte keine Links fetchen")
		return
	}
	if len(links) != 2 {
		t.Errorf("falsche Link Anzahl: Erwartet: 2, aber war: %d", len(links))
		return
	}
	if links[0] != "test.com" {
		t.Errorf("Link sollte test.com sein, aber war %s", links[0])
	}
}

func TestParseURL(t *testing.T) {
	url, err := ParseURL("", "")
	if url != nil || err == nil {
		t.Error("Leerer string sollte error returnen")
		return
	}

	url, err = ParseURL("test", "test")
	if url != nil || err == nil {
		t.Error("nicht absoulte URL sollte error returnen")
		return
	}

	hostname := "www.google.com"
	url, err = ParseURL("", "http://"+hostname)

	if err != nil {
		t.Error("Sollte keinen Error returnen")
		return
	} else if url.Hostname() != hostname {
		t.Errorf("Hostname sollte %s sein, ist aber %s", hostname, url.Hostname())
		return
	}

	hostname = "www.google.com"
	url, err = ParseURL("http://"+hostname, "test")
	expected := fmt.Sprintf("http://%s/test", hostname)

	if err != nil {
		t.Error("Sollte keinen Error returnen")
		return
	} else if url.String() != expected {
		t.Errorf("sollte %s sein, ist aber %s", expected, url.String())
	}
}
