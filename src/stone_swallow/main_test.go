package stone_swallow

import (
	"github.com/favclip/testerator"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListEntity(t *testing.T) {
	_, c, err := testerator.SpinUp()
	defer testerator.SpinDown()

	var h Hoge
	h.Id = "gufu"
	h.Name = "グフ"
	_, err = h.save(c)
	if err != nil {
		t.Fatal(err)
	}
	_, err = getHoge(c, createHogeKey(c, h.Id))
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/entity?kind=Hoge", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	listEntity(res, req, c)
	if res.Code != http.StatusOK {
		t.Fatalf("Non-expected status code : %v\n\tbody: %v", res.Code, res.Body)
	}
}

func TestGetTestCookie(t *testing.T) {
	_, c, err := testerator.SpinUp()
	defer testerator.SpinDown()

	req, err := http.NewRequest("GET", "/testcookie", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	handleTestCookie(res, req, c)
	if res.Code != http.StatusNotFound {
		t.Fatalf("Non-expected status code : %v\n\tbody: %v", res.Code, res.Body)
	}
}

func TestPostTestCookie(t *testing.T) {
	_, c, err := testerator.SpinUp()
	defer testerator.SpinDown()

	reader := strings.NewReader(`{"Domain" : "hoge"}`)
	req, err := http.NewRequest("POST", "/testcookie", reader)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	handleTestCookie(res, req, c)
	if res.Code != http.StatusOK {
		t.Fatalf("Non-expected status code : %v\n\tbody: %v", res.Code, res.Body)
	}
}
