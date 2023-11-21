package cache

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type mockWriter response

func newMockWriter() *mockWriter {
	return &mockWriter{
		body: []byte{},
		header: http.Header{},
	}
}

func (mw *mockWriter) Write(b []byte) (int, error) {
	mw.body = make([]byte, len(b))

	for k, v := range b {
		mw.body[k] = v
	}
	return len(b), nil
}

func (mw *mockWriter) WriteHeader(code int) {
	mw.code = code
}

func (mw *mockWriter) Header() http.Header {
	return mw.header
}

func TestWriter(t *testing.T) {
	mw := newMockWriter()

	res := "/test/url?with=params"
	u, err := url.Parse(res)
	if err != nil {
		t.Fatal("Invalid url")
	}

	req := &http.Request{
		URL: u,
	}

	t.Log("test newWriter")
	w := NewWriter(mw, req)
	if w.resource != res {
		t.Errorf("resources are different. expected %s / got %s", res, w.resource)
	}

	if w.writer != mw {
		t.Fatal("writer not assigned")
	}

	t.Log("test Header")

	h := w.Header()
	h.Add("test", "value")
	h2 := w.response.header
	if h2.Get("test") != "value" {
		t.Error("value not stored in the header")
	}

	t.Log("test WriteHeader")
	c := 201
	w.WriteHeader(c)
	if w.response.code != c {
		t.Error("status code not stored")
	}

	if mw.code != c {
		t.Error("status code not written")
	}

	h2 = mw.header
	if h2.Get("test") != "value" {
		t.Error("header not written")
	}

	t.Log("test Write")
	bd := []byte{1,2,3,4,5}
	n, err := w.Write(bd)
	if err != nil {
		t.Fatalf("error while writing: %s", err)
	}

	if n != len(bd) {
		t.Errorf("unexpected number of bytes written. expected %d / got %d", len(bd), n)
	}

	if &w.response.body == &bd {
		t.Error("body assigned, not copied")
	}

	if !reflect.DeepEqual(w.response.body, bd) {
		t.Error("body not copied")
		t.Error(w.response.body)
		t.Error(bd)
	}
	if &mw.body == &bd {
		t.Error("body assigned, not copied")
	}

	if !reflect.DeepEqual(mw.body, bd) {
		t.Error("body not passed through")
		t.Error(mw.body)
		t.Error(bd)
	}
}
