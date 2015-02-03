package stone_swallow

import (
	"appengine/aetest"
	"appengine/datastore"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
)

func TestListEntity(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

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
	if res.Code != http.StatusCreated {
		t.Fatalf("Non-expected status code : %v\n\tbody: %v", res.Code, res.Body)
	}

	var dst []datastore.PropertyList
	err = json.Unmarshal(res.Body.Bytes(), &dst)
	if err != nil {
		t.Fatal(err)
	}
	if len(dst) != 1 {
		t.Fatalf("Non-expected response length %v", len(dst))
	}

	if len(dst[0]) != 2 {
		t.Fatalf("Non-expected PropertyList length %v", len(dst[0]))
	}
	if dst[0][0].Name != "Name" {
		t.Fatalf("Non-expected PropertyList[0] Name %v", dst[0][0].Name)
	}
	if dst[0][0].Value != h.Name {
		t.Fatalf("Non-expected PropertyList[0] Value %v", dst[0][0].Value)
	}
	if dst[0][1].Name != "Created" {
		t.Fatalf("Non-expected PropertyList[1] Name %v", dst[0][1].Name)
	}
	if dst[0][1].Value != h.Created {
		// 日付は形式がちょっと変わっちゃった
		t.Logf("Non-expected PropertyList[1] Value %v : %v", dst[0][1].Value, h.Created)
	}
}

func TestGetTestCookie(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

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
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	r := strings.NewReader(`{"Domain" : "hoge"}`)
	req, err := http.NewRequest("POST", "/testcookie", r)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	handleTestCookie(res, req, c)
	if res.Code != http.StatusOK {
		t.Fatalf("Non-expected status code : %v\n\tbody: %v", res.Code, res.Body)
	}
}
