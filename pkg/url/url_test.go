package url

import (
	"fmt"
	"strings"
	"testing"
)

func TestFetchLinks(t *testing.T) {
	html := "<html><body><a href=\"http://test.com\"/><a href=\"http://provinzial.com\"/></body></html>"
	links, err := FetchLinks(strings.NewReader(html))
	if err != nil {
		t.Error("Konnte keine Links fetchen")
		return
	}

	expected := []string{"http://test.com", "http://provinzial.com"}
	if len(links) != len(expected) {
		t.Errorf("falsche Link Anzahl: Erwartet: %d, aber war: %d", len(expected), len(links))
		return
	}

	for i := 0; i < len(expected); i++ {
		if links[i] != expected[i] {
			t.Errorf("Link sollte %s sein, aber war %s", expected[i], links[i])
		}
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
