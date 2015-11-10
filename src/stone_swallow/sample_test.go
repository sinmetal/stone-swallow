package stone_swallow

import (
	"github.com/favclip/testerator"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TextEmpty(t *testing.T) {

}

func TestSave(t *testing.T) {
	_, c, err := testerator.SpinUp()
	defer testerator.SpinDown()

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

func TestPutHoge(t *testing.T) {
	_, c, err := testerator.SpinUp()
	defer testerator.SpinDown()

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

func TestHoge(t *testing.T) {
	_, c, err := testerator.SpinUp()
	defer testerator.SpinDown()

	var h Hoge
	h.Id = "gufu"
	_, err = h.save(c)
	if err != nil {
		t.Fatal(err)
	}
	_, err = getHoge(c, createHogeKey(c, h.Id))
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAllHoge(t *testing.T) {
	_, c, err := testerator.SpinUp()
	defer testerator.SpinDown()

	var h Hoge
	h.Id = "gufu"
	_, err = h.save(c)
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
