package main

import (
	"testing"

	"github.com/go-ini/ini"
)

func TestGetCredentials(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testIni := []byte(`[default]
aws_access_key_id=key
aws_secret_access_key=secret`)
		cfg, err := ini.Load(testIni)
		if err != nil {
			t.Fatal(err)
		}

		sec := cfg.Section("default")

		id, secret, err := getCredentials(sec)
		if err != nil {
			t.Error(err)
		}
		if id != "key" {
			t.Errorf("expected key id 'key', got %s", id)
		}

		if secret != "secret" {
			t.Errorf("expected secret key 'secret', got %s", secret)
		}
	})
}
