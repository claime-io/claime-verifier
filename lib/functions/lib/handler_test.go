package lib

import (
	"os"
	"testing"
)

func Test_allowedOrigin(t *testing.T) {
	t.Run("at dev env, allowed origin should be the origin which is specified at request header", func(t *testing.T) {
		os.Setenv("EnvironmentId", "dev")
		defer os.Setenv("EnvironmentId", "")
		want := "https://claime.io"
		got := allowedOrigin(want)
		if want != got {
			t.Errorf("allowedOrigin() = %v, want %v", got, want)
		}
	})
	t.Run("at prod env, allowed origin should be the origin which is specified at request header, too", func(t *testing.T) {
		os.Setenv("EnvironmentId", "prod")
		want := "https://auroradao.org"
		os.Setenv("AllowedOrigin", want)
		defer os.Setenv("EnvironmentId", "")
		defer os.Setenv("AllowedOrigin", "")
		got := allowedOrigin("https://auroradao.org")
		if want != got {
			t.Errorf("allowedOrigin() = %v, want %v", got, want)
		}
	})
}
