package main

import (
	"net/http"
	"path"
	"strings"
	"fmt"
	"os"
	"flag"
	// "log"
)

type syncItem struct {
	album  string
	file string
	url  string
}

func addImgmax(url string) string {
	if strings.Index(url, "?") > 0 {
		url += "&"
	} else {
		url += "?"
	}
	return url + "imgmax=d"
}

func syncAlbum(picasa *Picasa, album *Album, ch chan *syncItem) {
	url := addImgmax(album.Url)
	items := picasa.getAlbum(url)
	for _, item := range items {
		ch <- &syncItem{
			album: album.FileName,
			file: item.FileName,
			url: item.Url,
		}
	}
}

func sync(picasa *Picasa, userid string, requestedAlbum string, ch chan *syncItem) {
	albums := picasa.getAlbums(userid)
	for _, album := range albums {
		if requestedAlbum == "" {
			fmt.Printf("- %s (%s)\n", album.Title, album.FileName)
		} else if requestedAlbum == "all" || requestedAlbum == album.FileName {
			syncAlbum(picasa, album, ch)
		}
	}
	ch <- nil
	close(ch)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [userid] [album-name]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "userid 'default' means the auth user\n")
		fmt.Fprintf(os.Stderr, "leaving out album name lists albums\n")
		fmt.Fprintf(os.Stderr, "album name 'all' fetches all albums\n")
		flag.PrintDefaults()
		return
	}
	user := flag.String("user", "", "auth username (email)")
	pass := flag.String("pass", "", "auth password")
	flag.Parse()
	userid := flag.Arg(0)
	album := flag.Arg(1)

	client := &http.Client{}
	auth, err := ClientLogin(client, *user, *pass, "lh2")
	check(err)

	picasa := &Picasa{client: client, auth: auth}

	ch := make(chan *syncItem, 5)
	go sync(picasa, userid, album, ch)

	for {
		si := <-ch
		if si == nil {
			break
		}

		p := path.Join("dl", si.album)
		check(os.MkdirAll(p, 0777))
		p = path.Join(p, si.file)
		download(client, si.url, p, "")
	}
}
