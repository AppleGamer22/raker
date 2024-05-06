package shared

import (
	"io"
	"os/exec"
)

func Stream2MP4(response io.ReadCloser, path string) error {
	process := exec.Command("ffmpeg", "-i", "pipe:", path)
	process.Stdin = response
	return process.Run()
}
