package httpmock

import (
	"io/ioutil"
	"net/http"
	"testing"
)

var dummyResponder = NewStringResponder(200, "OK")

// Checks if "request" is matched without error against "target"
func assertMatch(t *testing.T, targetMethod, targetUrl, requestMethod, requestUrl string) {
	matcher := NewMatcher(targetMethod, targetUrl, dummyResponder)

	matched, err := matcher.Match(requestMethod, requestUrl)

	if err != nil {
		t.Errorf(err.Error())
	}

	if matched != true {
		t.Errorf("Matcher %s %s failed to match %s %s", targetMethod, targetUrl, requestMethod, requestUrl)
	}
}

// Checks if both pairs of url and method are matching against each other without error
func assertMatchBoth(t *testing.T, firstMethod, firstUrl, secondMethod, secondUrl string) {
	assertMatch(t, firstMethod, firstUrl, secondMethod, secondUrl)
	assertMatch(t, secondMethod, secondUrl, firstMethod, firstUrl)
}

// Checks if "request" will not be matched (without error) against "target"
func assertNoMatch(t *testing.T, targetMethod, targetUrl, requestMethod, requestUrl string) {
	matcher := NewMatcher(targetMethod, targetUrl, dummyResponder)

	matched, err := matcher.Match(requestMethod, requestUrl)

	if err != nil {
		t.Errorf(err.Error())
	}

	if matched {
		t.Errorf("Matcher %s %s unexpectedly matched %s %s", targetMethod, targetUrl, requestMethod, requestUrl)
	}
}

// Check if both pairs of url and method are not matched against each other without error
func assertNoMatchBoth(t *testing.T, firstMethod, firstUrl, secondMethod, secondUrl string) {
	assertNoMatch(t, firstMethod, firstUrl, secondMethod, secondUrl)
	assertNoMatch(t, secondMethod, secondUrl, firstMethod, firstUrl)
}

func TestUrlMatching(t *testing.T) {
	assertMatchBoth(t,
		"GET", "http://www.test.com",
		"GET", "http://www.test.com")
	assertMatchBoth(t,
		"GET", "http://www.test.com",
		"GET", "http://www.test.com/")
	assertMatchBoth(t,
		"GET", "http://www.test.com/abc",
		"GET", "http://www.test.com/abc")
	assertMatchBoth(t,
		"GET", "http://www.test.com:5000/abc",
		"GET", "http://www.test.com:5000/abc")

	assertNoMatchBoth(t,
		"GET", "https://www.test.com",
		"GET", "http://www.test.com")
	assertNoMatchBoth(t,
		"GET", "http://www.test.com/abc",
		"GET", "http://www.test.com")
	assertNoMatchBoth(t,
		"GET", "http://test.com",
		"GET", "http://www.test.com")
	assertNoMatchBoth(t,
		"GET", "http://test.com",
		"GET", "http://www.test.com")
	assertNoMatchBoth(t,
		"GET", "http://test.com/abc",
		"GET", "http://www.test.com/abc/")
	assertNoMatchBoth(t,
		"GET", "http://test.com/abc/",
		"GET", "http://www.test.com/abc")
	assertNoMatchBoth(t,
		"GET", "http://test.com:5000/abc/",
		"GET", "http://www.test.com/abc")
	assertNoMatchBoth(t,
		"GET", "http://test.com/abc/",
		"GET", "http://www.test.com:5000/abc")
}

func TestSubsetMatch(t *testing.T) {
	assertMatch(t,
		"GET", "/path",
		"GET", "http://www.test.com/path")
	assertMatch(t,
		"GET", "/path",
		"GET", "http://www.test.com/path")
	assertMatch(t,
		"GET", "//www.test.com/path",
		"GET", "http://www.test.com/path")
	assertMatch(t,
		"GET", "//www.test.com/path",
		"GET", "https://www.test.com/path")
}

func TestMethodMatch(t *testing.T) {
	url := "http://www.test.com/path"

	assertNoMatch(t,
		"GET", url,
		"POST", url)

	methods := []string{
		"PUT",
		"put",
		"PuT",
		"puT",
	}

	for _, method := range methods {
		assertMatch(t,
			"PUT", url,
			method, url)
	}
}

func TestMatchAnyUrl(t *testing.T) {
	url := "http://www.test.com/path"

	assertNoMatch(t,
		"POST", url,
		"GET", url)

	urls := []string{
		"http://google.com",
		"http://www.test.com",
		"http://some-url.dot.com/with/path",
		"/just/path",
		"ftp://127.0.0.1:8000",
	}

	for _, url := range urls {
		assertMatch(t,
			"GET", ANY,
			"GET", url)
	}
}

func TestMatchAnyMethod(t *testing.T) {
	url := "http://www.test.com/path"

	assertNoMatch(t,
		ANY, url,
		"GET", url+"/more/path")

	methods := []string{
		"OPTIONS",
		"PUT",
		"GET",
		"POST",
		"PATCH",
		"HEAD",
		"CREAZY",
	}

	for _, method := range methods {
		assertMatch(t,
			ANY, url,
			method, url)
	}
}

func TestCaconicalPath(t *testing.T) {
	if canonicalPath("/") != "" {
		t.Errorf("canonicalPath '/' != ''")
	}

	if canonicalPath("") != "" {
		t.Errorf("canonicalPath '' != ''")
	}

	if canonicalPath("/path") != "/path" {
		t.Errorf("canonicalPath '/path' != '/path'")
	}
}

func TestMatcherReturnsResponder(t *testing.T) {
	body := "It's OK"
	testResponder := NewStringResponder(200, body)
	matcher := NewMatcher(ANY, ANY, testResponder)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := matcher.Respond(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Response status is %s instead of 200", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != string(body) {
		t.Errorf("Response body %s is different from expected body %s", data, body)
	}

}
