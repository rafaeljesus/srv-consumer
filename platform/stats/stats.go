package stats

import (
	"log"
	"time"
)

type Client struct{}

func (c *Client) Start() time.Time {
	return time.Now()
}

func (c *Client) Track(t time.Time, err bool) {
	log.Print("sending stats timing metric")
}
