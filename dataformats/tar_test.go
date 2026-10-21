package dataformats

import (
	"archive/tar"
	"bytes"
	"errors"
	"io"
	"pimtrace"
	"strings"
	"testing"
)

func TestReadTarStream_StrictAndBestEffort(t *testing.T) {
	// Create a tar archive with one valid mock entry and one "invalid" entry
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	// Valid entry
	hdr1 := &tar.Header{
		Name: "valid.txt",
		Mode: 0600,
		Size: int64(len("valid content")),
	}
	tw.WriteHeader(hdr1)
	tw.Write([]byte("valid content"))

	// Invalid entry
	hdr2 := &tar.Header{
		Name: "invalid.txt",
		Mode: 0600,
		Size: int64(len("invalid content")),
	}
	tw.WriteHeader(hdr2)
	tw.Write([]byte("invalid content"))

	// Valid entry 2
	hdr3 := &tar.Header{
		Name: "valid2.txt",
		Mode: 0600,
		Size: int64(len("valid content 2")),
	}
	tw.WriteHeader(hdr3)
	tw.Write([]byte("valid content 2"))
	tw.Close()

	mockNext := func(f io.Reader, fType string, fName string, ops ...any) ([]string, error) {
		strict := false
		for _, op := range ops {
			if s, ok := op.(string); ok && s == "strict" {
				strict = true
			}
		}

		if strings.Contains(fName, "invalid") {
			if strict {
				return nil, errors.New("strict error parsing invalid.txt")
			}
			return nil, &pimtrace.PartialResultError{Skipped: 1, Message: "best-effort skipping invalid.txt"}
		}
		return []string{"parsed " + fName}, nil
	}

	t.Run("best-effort (default)", func(t *testing.T) {
		r := bytes.NewReader(buf.Bytes())
		res, err := ReadTarStream(r, "mock", "test.tar", mockNext, []string{"*"})
		var partial *pimtrace.PartialResultError
		if err != nil && !errors.As(err, &partial) {
			t.Fatalf("expected PartialResultError or no error in best-effort mode, got %v", err)
		}
		if len(res) != 2 || res[0] != "parsed valid.txt" || res[1] != "parsed valid2.txt" {
			t.Errorf("expected 2 valid items, got %v", res)
		}
		if partial == nil || partial.Skipped != 1 {
			t.Errorf("expected 1 skipped item in partial result, got %v", partial)
		}
	})

	t.Run("strict mode", func(t *testing.T) {
		r := bytes.NewReader(buf.Bytes())
		_, err := ReadTarStream(r, "mock", "test.tar", mockNext, []string{"*"}, "strict")
		if err == nil {
			t.Fatalf("expected error in strict mode, got none")
		}
		if !strings.Contains(err.Error(), "strict error parsing invalid.txt") {
			t.Errorf("expected error to mention parsing invalid.txt, got: %v", err)
		}
	})
}
