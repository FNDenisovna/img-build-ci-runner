package api

import (
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	t.Run("simple get without errors", func(t *testing.T) {
		testapi := New("https://rdb.altlinux.org/api/license")
		resp, err := testapi.Get()
		if err != nil {
			t.Errorf("api.Get() returns unexpected error = %v", err)
			return
		}
		if len(resp) > 0 && strings.Contains(string(resp), "Version") {
			return
		}

		t.Errorf("api.Get() returns something wrong:\n %x", resp)
	})
}
