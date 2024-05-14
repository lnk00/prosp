package models

import (
	"io"
	"time"
)

type Message struct {
	Date   time.Time
	Reader io.Reader
}
