package main

import (
	"io"
	"io/ioutil"
	"encoding/xml"
)

type AtomFeed struct {
	Entry []AtomEntry `xml:"entry"`
}

type AtomEntry struct {
	ETag      string `xml:"etag,attr"`
	Title     string `xml:"title"`
	Name      string `xml:"name"`
	Id        string `xml:"id"`
	Link      []AtomLink `xml:"link"`
	NumPhotos int    `xml:"numphotos"`
	Media     AtomMedia  `xml:"group"`
}

type AtomLink struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

type AtomMedia struct {
	Content []AtomContent `xml:"content"`
}

type AtomContent struct {
	Type string `xml:"type,attr"`
	Url  string `xml:"url,attr"`
}

func AtomParse(r io.Reader) (*AtomFeed, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var feed AtomFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}
