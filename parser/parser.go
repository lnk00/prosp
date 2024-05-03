package parser

import (
	"io"
	"log"

	"github.com/PuerkitoBio/goquery"
)

type Job struct {
	title    string
	location string
}

var jobSelector = "body > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr > td > table > tbody > tr > td > table > tbody > tr > td > table > tbody > tr > td:nth-child(2) > a"
var titleSelector = "table > tbody > tr:nth-child(1) > td > a"
var locationSelector = "table > tbody > tr:nth-child(2) > td > p"

func Parse(msg io.Reader) []Job {
	res := []Job{}

	doc, err := goquery.NewDocumentFromReader(msg)
	if err != nil {
		log.Fatalf("PARSING failed: %v", err)
	}

	jobs := doc.Find(jobSelector)
	jobs.Each(func(_ int, s *goquery.Selection) {
		title := s.Find(titleSelector).Text()
		location := s.Find(locationSelector).Text()
		res = append(res, Job{title: title, location: location})
	})

	return res
}

func ParseAll(messages []io.Reader) []Job {
	res := []Job{}

	for _, msg := range messages {
		jobs := Parse(msg)
		res = append(res, jobs...)
	}

	return res
}
