package jsonapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client dipakai semua fungsi di package ini. Satu client bersama memakai ulang koneksi
// (keep-alive), dan Timeout mencegah request menggantung. Boleh diganti sesuai kebutuhan.
var Client = &http.Client{Timeout: 30 * time.Second}

// MaxResponseBytes membatasi ukuran body respons yang dibaca (10 MiB).
var MaxResponseBytes int64 = 10 << 20

func PostStructWithToken[T any](tokenkey string, tokenvalue string, structname interface{}, urltarget string) (statusCode int, result T, err error) {
	mJson, err := json.Marshal(structname)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, urltarget, bytes.NewReader(mJson))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(tokenkey, tokenvalue)
	return do[T](req, false)
}

func Get[T any](urltarget string) (statusCode int, result T, err error) {
	req, err := http.NewRequest(http.MethodGet, urltarget, nil)
	if err != nil {
		return
	}
	return do[T](req, false)
}

func GetWithBearer[T any](tokenbearer string, urltarget string) (statusCode int, result T, err error) {
	req, err := http.NewRequest(http.MethodGet, urltarget, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenbearer)
	return do[T](req, true)
}

func GetStructWithToken[T any](tokenkey string, tokenvalue string, urltarget string) (statusCode int, result T, err error) {
	req, err := http.NewRequest(http.MethodGet, urltarget, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(tokenkey, tokenvalue)
	return do[T](req, false)
}

// do mengirim request, membaca body dengan batas ukuran, lalu mendekode JSON ke T.
// withCause menyertakan pesan error json di depan (perilaku lama GetWithBearer).
func do[T any](req *http.Request, withCause bool) (statusCode int, result T, err error) {
	resp, err := Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxResponseBytes+1))
	if err != nil {
		return
	}
	if int64(len(body)) > MaxResponseBytes {
		err = fmt.Errorf("response from %s exceeds %d bytes", req.URL.Redacted(), MaxResponseBytes)
		return
	}
	if jerr := json.Unmarshal(body, &result); jerr != nil {
		msg := "Not A Valid JSON Response from " + req.URL.Redacted() + " . CONTENT: " + snippet(body)
		if withCause {
			msg = jerr.Error() + " | " + msg
		}
		err = errors.New(msg)
		return
	}
	return
}

// snippet memotong body untuk pesan error supaya log tidak membengkak.
func snippet(b []byte) string {
	const max = 512
	if len(b) > max {
		return string(b[:max]) + "...(truncated)"
	}
	return string(b)
}
