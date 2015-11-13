package httpmock

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	ANY = "__ANY__"
)

// Return "canonical" path. When path is equal "/" it is the same
// as when path is equal "", so to make comparisons easier empty
// string ("") for such case was choosen.
func canonicalPath(path string) string {
	if path == "/" {
		return ""
	}
	return path
}

type HttpMockMatcher struct {
	method    string
	url       *url.URL
	anyUrl    bool
	responder Responder
}

func (m *HttpMockMatcher) MatchMethod(methodToCheck string) bool {
	if m.method == ANY {
		return true
	}

	if strings.ToUpper(methodToCheck) == strings.ToUpper(m.method) {
		return true
	}
	return false
}

func (m *HttpMockMatcher) MatchUrl(urlToCheck string) bool {
	if m.anyUrl {
		return true
	}

	parsedUrl, err := url.ParseRequestURI(urlToCheck)

	if err != nil {
		panic(err)
	}

	if m.url.Scheme != "" && (parsedUrl.Scheme != m.url.Scheme) {
		return false
	}

	if m.url.Host != "" && (parsedUrl.Host != m.url.Host) {
		return false
	}

	if canonicalPath(parsedUrl.Path) != canonicalPath(m.url.Path) {
		return false
	}

	return true
}

func (m *HttpMockMatcher) Match(methodToCheck, urlToCheck string) (bool, error) {
	if m.MatchMethod(methodToCheck) != true {
		return false, nil
	}

	if m.MatchUrl(urlToCheck) != true {
		return false, nil
	}

	return true, nil
}

func (m *HttpMockMatcher) Respond(req *http.Request) (*http.Response, error) {
	return m.responder(req)
}

func NewMatcher(methodToMatch, urlToMatch string, responder Responder) *HttpMockMatcher {
	matcher := &HttpMockMatcher{
		method:    methodToMatch,
		anyUrl:    urlToMatch == ANY,
		responder: responder,
	}

	if matcher.anyUrl != true {
		matcher.url, _ = url.Parse(urlToMatch)
	}

	return matcher
}
