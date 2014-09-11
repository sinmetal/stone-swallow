package stone_swallow

import (
	"appengine/aetest"
	"appengine/datastore"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
}
