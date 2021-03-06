package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
)

func copyWithProgress(text string, r io.Reader, w io.Writer, total int64) error {
	defer fmt.Printf("\n")

	var cur int64 = 0
	buf := make([]byte, 64 << 10)
	for {
		nr, err := r.Read(buf)
		if err != nil {
			return err
		}
		nw, err := w.Write(buf[0:nr])
		if err != nil {
			return err
		}
		if nw < nr {
			return fmt.Errorf("short write")
		}
		cur += int64(nw)

		fmt.Printf("\r\x1b[K")
		fmt.Printf("%4dk ", cur / 1000)
		if total >= 0 && cur < total {
			fmt.Printf("[")
			frac := float32(cur + 1) / float32(total + 1)
			width := 10
			for i := 0; i < width; i++ {
				if i < int(float32(width) * frac) {
					fmt.Printf("#")
				} else {
					fmt.Printf(" ")
				}
			}
			fmt.Printf("] ")
		}
		fmt.Printf("%s", text)
	}
	return nil
}

func download(client *http.Client, url string, path string, etag string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if etag != "" {
		req.Header.Add("If-None-Match", etag)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 304 {
		return nil
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("fetch returned status %d", resp.StatusCode)
	}

	dlpath := path + ".download"
	f, err := os.Create(dlpath)
	if err != nil {
		return err
	}

	err = copyWithProgress(path, resp.Body, f, resp.ContentLength)

	f.Close()
	if err != nil {
		os.Remove(dlpath)
		return err
	}
	err = os.Rename(dlpath, path)
	if err != nil {
		return err
	}

	return nil
}
