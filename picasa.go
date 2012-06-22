package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"fmt"
)

type Album struct {
	Title     string
	FileName  string
	NumPhotos int
	Url       string
}

type Item struct {
	FileName string
	Url      string
}

type Picasa struct {
	client *http.Client
	auth   string
}

func (p *Picasa) get(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("GData-Version", "2")
	if p.auth != "" {
		req.Header.Add("Authorization", "GoogleLogin auth="+p.auth)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("request failed (%d) - %q", resp.StatusCode, data)
	}

	return resp.Body, nil
}

func (p *Picasa) getAlbums(userid string) []*Album {
	if userid == "" {
		userid = "default"
	}
	body, err := p.get("https://picasaweb.google.com/data/feed/api/user/" + userid)
	check(err)
	defer body.Close()
	return p.parseAlbums(body)
}

func (p *Picasa) parseAlbums(r io.Reader) []*Album {
	feed, err := AtomParse(r)
	check(err)

	var albums []*Album
	for _, e := range feed.Entry {
		album := &Album{}
		album.NumPhotos = e.NumPhotos
		album.Title = e.Title
		album.FileName = e.Name
		for _, l := range e.Link {
			if l.Rel == "http://schemas.google.com/g/2005#feed" {
				album.Url = l.Href
			}
		}
		albums = append(albums, album)
	}

	return albums
}

func (p *Picasa) getAlbum(url string) []*Item {
	body, err := p.get(url)
	check(err)
	defer body.Close()
	return p.parseAlbum(body)
}

func (p *Picasa) parseAlbum(r io.Reader) []*Item {
	feed, err := AtomParse(r)
	check(err)

	var items []*Item
	for _, e := range feed.Entry {
		item := &Item{}
		item.FileName = e.Title
		for _, c := range e.Media.Content {
			if strings.HasPrefix(c.Type, "video/") {
				item.Url = c.Url
				break
			} else if strings.HasPrefix(c.Type, "image/") {
				item.Url = c.Url
				// Don't break here; keep looking in case there's video.
			}
		}
		items = append(items, item)
	}

	return items
}
