package jsonapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type echo struct {
	Method string `json:"method"`
	Token  string `json:"token"`
	Auth   string `json:"auth"`
	Body   string `json:"body"`
}

func echoServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		json.NewEncoder(w).Encode(echo{Method: r.Method, Token: r.Header.Get("Token"), Auth: r.Header.Get("Authorization"), Body: string(b)})
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestFunctions(t *testing.T) {
	srv := echoServer(t)

	code, res, err := PostStructWithToken[echo]("Token", "abc", map[string]string{"a": "b"}, srv.URL)
	if err != nil || code != 200 || res.Method != "POST" || res.Token != "abc" || res.Body != `{"a":"b"}` {
		t.Errorf("PostStructWithToken: code=%d res=%+v err=%v", code, res, err)
	}
	_, res, err = Get[echo](srv.URL)
	if err != nil || res.Method != "GET" {
		t.Errorf("Get: res=%+v err=%v", res, err)
	}
	_, res, err = GetWithBearer[echo]("xyz", srv.URL)
	if err != nil || res.Auth != "Bearer xyz" {
		t.Errorf("GetWithBearer: res=%+v err=%v", res, err)
	}
	_, res, err = GetStructWithToken[echo]("Token", "def", srv.URL)
	if err != nil || res.Token != "def" {
		t.Errorf("GetStructWithToken: res=%+v err=%v", res, err)
	}
}

func TestMarshalErrorReturned(t *testing.T) {
	srv := echoServer(t)
	if _, _, err := PostStructWithToken[echo]("Token", "x", make(chan int), srv.URL); err == nil {
		t.Error("struct yang tidak bisa di-marshal harus menghasilkan error, bukan body kosong")
	}
}

func TestInvalidJSONErrorIsTruncatedAndRedacted(t *testing.T) {
	big := strings.Repeat("<html>", 1000)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(big)) }))
	defer srv.Close()

	u := strings.Replace(srv.URL, "http://", "http://user:secretpw@", 1)
	_, _, err := Get[echo](u)
	if err == nil {
		t.Fatal("expected error")
	}
	if len(err.Error()) > 1024 || !strings.Contains(err.Error(), "truncated") {
		t.Errorf("error harus dipotong, panjang=%d", len(err.Error()))
	}
	if strings.Contains(err.Error(), "secretpw") {
		t.Error("password di URL tidak boleh muncul di pesan error")
	}
	_, _, err = GetWithBearer[echo]("t", srv.URL)
	if err == nil || !strings.Contains(err.Error(), " | Not A Valid JSON") {
		t.Errorf("GetWithBearer harus menyertakan penyebab: %v", err)
	}
}

func TestResponseSizeLimit(t *testing.T) {
	old := MaxResponseBytes
	MaxResponseBytes = 64
	t.Cleanup(func() { MaxResponseBytes = old })
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"body":"` + strings.Repeat("a", 200) + `"}`))
	}))
	defer srv.Close()
	if _, _, err := Get[echo](srv.URL); err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Errorf("body melebihi batas harus ditolak: %v", err)
	}
}

func TestTimeout(t *testing.T) {
	old := Client
	Client = &http.Client{Timeout: 50 * time.Millisecond}
	t.Cleanup(func() { Client = old })
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond)
	}))
	defer srv.Close()
	start := time.Now()
	if _, _, err := Get[echo](srv.URL); err == nil {
		t.Error("request lambat harus timeout")
	}
	if time.Since(start) > 250*time.Millisecond {
		t.Error("timeout tidak dihormati")
	}
}
