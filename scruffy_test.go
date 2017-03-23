package main

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satori/go.uuid"
)

func TestScruffy(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dir := path.Join(cwd, "cmd", "test")
	config := path.Join(dir, "config.yml")

	t.Run("Help", func(t *testing.T) {
		tmp := build(t)
		defer os.RemoveAll(tmp)

		cmd := exec.Command("scruffy", "--help")
		b, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal("Expected exit code 0 but got 1")
		}

		s := string(b)
		if !strings.Contains(s, "Scruffy - build your blueprints from mess to order!") {
			t.Fatalf("Expected %v, but %v", "Usage: ", s)
		}
	})

	t.Run("Bulding", func(t *testing.T) {

		t.Run("Help", func(t *testing.T) {
			tmp := build(t)
			defer os.RemoveAll(tmp)

			cmd := exec.Command("scruffy", "build", "--help")
			b, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err)
			}

			watchers := []string{"--config scruffy.yml", "--env value", "--watch false"}
			s := string(b)
			for _, w := range watchers {
				if !strings.Contains(s, w) {
					t.Fatalf("Expected %v, but %v", w, s)
				}
			}
		})

		t.Run("Empty environment", func(t *testing.T) {
			tmp := build(t)
			defer os.RemoveAll(tmp)

			os.Chdir(dir)
			cmd := exec.Command("scruffy", "build", "--config", config)
			b, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatal("Expected error code 1 got 0")
			}

			s := string(b)
			if !strings.Contains(s, "Please pass non-empty --env") {
				t.Fatalf("Expected %v, but %v", "Please pass non-empty --env", s)
			}
		})

		t.Run("Unexisting config", func(t *testing.T) {
			tmp := build(t)
			defer os.RemoveAll(tmp)

			os.Chdir(dir)
			cmd := exec.Command("scruffy", "build", "--config", "/path/to/unknown/config", "--env", "public")
			b, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatal("Expected error code 1 got 0")
			}

			s := string(b)
			if !strings.Contains(s, "config parsing error:") {
				t.Fatalf("Expected %v, but %v", "Config parsing error:", s)
			}
		})

		t.Run("Unexisting environment", func(t *testing.T) {
			tmp := build(t)
			defer os.RemoveAll(tmp)

			os.Chdir(dir)
			cmd := exec.Command("scruffy", "build", "--config", config, "--env", "unknown")
			b, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatal("Expected error code 1 got 0")
			}

			s := string(b)
			if !strings.Contains(s, "Unknown environment: unknown") {
				t.Fatalf("Expected %v, but %v", "Unknown environment: unknown", s)
			}
		})

		t.Run("Success", func(t *testing.T) {
			tmp := build(t)
			defer os.RemoveAll(tmp)

			os.Chdir(dir)
			cmd := exec.Command("scruffy", "build", "--config", config, "--env", "public")
			b, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err)
			}

			s := string(b)
			if !strings.Contains(s, "Build available at:") {
				t.Fatalf("Expected %v, but %v", "Build available at:", s)
			}
		})

		t.Run("Watching", func(t *testing.T) {
			t.Skip("Implement someday!")
		})

	})

	t.Run("Publishing", func(t *testing.T) {

		t.Run("Help", func(t *testing.T) {
			tmp := build(t)
			defer os.RemoveAll(tmp)

			cmd := exec.Command("scruffy", "publish", "--help")
			b, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err)
			}

			watchers := []string{"--config scruffy.yml", "--env value", "--watch false", "--release false"}
			s := string(b)
			for _, w := range watchers {
				if !strings.Contains(s, w) {
					t.Fatalf("Expected %v, but %v", w, s)
				}
			}
		})

		t.Run("Empty environment", func(t *testing.T) {
			tmp := build(t)
			defer os.RemoveAll(tmp)

			os.Chdir(dir)
			cmd := exec.Command("scruffy", "publish", "--config", config)
			b, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatal("Expected error code 1 got 0")
			}

			s := string(b)
			if !strings.Contains(s, "Please pass non-empty --env") {
				t.Fatalf("Expected %v, but %v", "Please pass non-empty --env", s)
			}
		})

		t.Run("Unexisting config", func(t *testing.T) {
			tmp := build(t)
			defer os.RemoveAll(tmp)

			os.Chdir(dir)
			cmd := exec.Command("scruffy", "publish", "--config", "/path/to/unknown/config", "--env", "public")
			b, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatal("Expected error code 1 got 0")
			}

			s := string(b)
			if !strings.Contains(s, "config parsing error:") {
				t.Fatalf("Expected %v, but %v", "Config parsing error:", s)
			}
		})

		t.Run("Unexisting environment", func(t *testing.T) {
			tmp := build(t)
			defer os.RemoveAll(tmp)

			os.Chdir(dir)
			cmd := exec.Command("scruffy", "publish", "--config", config, "--env", "unknown")
			b, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatal("Expected error code 1 got 0")
			}

			s := string(b)
			if !strings.Contains(s, "Unknown environment: unknown") {
				t.Fatalf("Expected %v, but %v", "Unknown environment: unknown", s)
			}
		})

		t.Run("Success", func(t *testing.T) {
			tmp := build(t)
			defer os.RemoveAll(tmp)

			os.Chdir(dir)
			cmd := exec.Command("scruffy", "publish", "--config", config, "--env", "public")
			b, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err)
			}

			s := string(b)
			if !strings.Contains(s, "Public docs changed:") {
				t.Fatalf("Expected %v, but %v", "Public docs changed:", s)
			}
		})

		t.Run("Success release", func(t *testing.T) {
			tmp := build(t)
			defer os.RemoveAll(tmp)

			os.Chdir(dir)
			cmd := exec.Command("scruffy", "publish", "--config", config, "--env", "public", "--release")
			b, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err)
			}

			s := string(b)
			if !strings.Contains(s, "Public docs changed:") {
				t.Fatalf("Expected %v, but %v", "Public docs changed:", s)
			}
		})

		t.Run("Watching", func(t *testing.T) {
			t.Skip("Implement someday!")
		})

	})
}

func build(t *testing.T) string {
	tmp := os.TempDir()
	tmp = filepath.Join(tmp, uuid.NewV4().String())
	run(t, "go", "build", "-o", filepath.Join(tmp, "bin", "scruffy"), "github.com/m1ome/scruffy")
	os.Setenv("PATH", filepath.Join(tmp, "bin")+string(filepath.ListSeparator)+os.Getenv("PATH"))
	os.MkdirAll(filepath.Join(tmp, "src"), 0755)
	return tmp
}

func run(t *testing.T, cmd string, args ...string) []byte {
	b, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("Expected %v, but %v: %v", nil, err, string(b))
	}
	return b
}
