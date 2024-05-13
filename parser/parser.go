package parser

import (
	"io"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/lnk00/prosp/models"
)

func linkedinParse(doc *goquery.Document) []models.Job {
	list := []models.Job{}

	var jobSelector = "body > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr > td > table > tbody > tr > td > table > tbody > tr > td > table > tbody > tr > td:nth-child(2) > a"
	var titleSelector = "table > tbody > tr:nth-child(1) > td > a"
	var locationSelector = "table > tbody > tr:nth-child(2) > td > p"

	jobs := doc.Find(jobSelector)
	jobs.Each(func(_ int, s *goquery.Selection) {
		title := s.Find(titleSelector).Text()
		link := s.Find(titleSelector).AttrOr("href", "link not found")
		link = strings.Split(link, "?")[0]
		location := s.Find(locationSelector).Text()
		list = append(list, models.Job{
			Title:    title,
			Location: location,
			Link:     link,
			Status:   models.TO_APPLY,
		})
	})

	return list
}

func Parse(msg io.Reader) []models.Job {

	doc, err := goquery.NewDocumentFromReader(msg)
	if err != nil {
		log.Fatalf("PARSING failed: %v", err)
	}

	list := linkedinParse(doc)

	return list
}

func ParseAll(messages []io.Reader) []models.Job {
	res := []models.Job{}

	for _, msg := range messages {
		jobs := Parse(msg)
		res = append(res, jobs...)
	}

	return res
}
