package stone_swallow

import (
	"testing"

	"appengine/aetest"
	"net/http"
	"net/http/httptest"
)

func TestFirst(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	testSave(c, t)
	testPutHoge(c, t)
	testHoge(c, t)
	testGetAllHoge(c, t)
}

func TextEmpty(t *testing.T) {

}

func testSave(c aetest.Context, t *testing.T) {
	var h Hoge
	h.Id = "gufu"
	s, err := h.save(c)
	if err != nil {
		t.Fatal(err)
	}

	if s.Id != "gufu" {
		t.Errorf("Non-expected key:%v", s.Id)
	}
}

func testPutHoge(c aetest.Context, t *testing.T) {
	req, err := http.NewRequest("POST", "/sample", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	putHoge(res, req, c)

	if res.Code != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %v", "200", res.Code)
	}
}

func testHoge(c aetest.Context, t *testing.T) {
	var h Hoge
	h.Id = "gufu"
	_, err := h.save(c)
	if err != nil {
		t.Fatal(err)
	}
	_, err = getHoge(c, createHogeKey(c, h.Id))
	if err != nil {
		t.Fatal(err)
	}
}

func testGetAllHoge(c aetest.Context, t *testing.T) {
	var h Hoge
	h.Id = "gufu"
	_, err := h.save(c)
	if err != nil {
		t.Fatal(err)
	}
	_, err = getHoge(c, createHogeKey(c, h.Id))
	if err != nil {
		t.Fatal(err)
	}

	hoges, err := getAllHoge(c)
	if err != nil {
		t.Fatal(err)
	}

	if len(hoges) < 1 {
		t.Errorf("hoges size = %d", len(hoges))
	}
}
