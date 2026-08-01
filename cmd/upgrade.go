package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade odu to the latest version",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		current := rootCmd.Version
		fmt.Println("Checking for latest version...")

		latest, err := fetchLatestVersion()
		if err != nil {
			return fmt.Errorf("could not check latest version: %w", err)
		}

		if current != "dev" && current == latest {
			fmt.Printf("✓ Already on the latest version (%s)\n", current)
			return nil
		}

		fmt.Printf("Upgrading %s → %s...\n", current, latest)

		goos := runtime.GOOS
		arch := runtime.GOARCH
		if arch == "x86_64" {
			arch = "amd64"
		}

		binaryName := fmt.Sprintf("odu-%s-%s", goos, arch)
		url := fmt.Sprintf("https://github.com/niranjangowdaj/odu/releases/download/%s/%s", latest, binaryName)

		tmp, err := downloadBinary(url)
		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}
		defer os.Remove(tmp)

		self, err := os.Executable()
		if err != nil {
			return fmt.Errorf("could not locate current binary: %w", err)
		}

		if err := replaceBinary(tmp, self); err != nil {
			return fmt.Errorf("could not replace binary: %w\nTry: sudo mv %s %s", err, tmp, self)
		}

		fmt.Printf("✓ odu upgraded to %s\n", latest)
		return nil
	},
}

func fetchLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/niranjangowdaj/odu/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned HTTP %d — check your internet connection", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("could not parse GitHub API response: %w", err)
	}
	if release.TagName == "" {
		return "", fmt.Errorf("no releases found")
	}
	return release.TagName, nil
}

func downloadBinary(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d — check that the release exists for your platform", resp.StatusCode)
	}

	tmp, err := os.CreateTemp("", "odu-upgrade-*")
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		return "", err
	}
	if err := os.Chmod(tmp.Name(), 0755); err != nil {
		return "", err
	}
	return tmp.Name(), nil
}

func replaceBinary(src, dst string) error {
	dstNew := dst + ".new"
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dstNew, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		os.Remove(dstNew)
		return err
	}
	out.Close()

	if err := os.Rename(dstNew, dst); err != nil {
		os.Remove(dstNew)
		return err
	}
	return nil
}
