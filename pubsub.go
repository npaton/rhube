package rhube

import (
	"regexp"
	"strings"
)

type Subscriptions map[Channel][]*Client

type Client chan ([]byte)

func NewClient() *Client {
	c := Client(make(chan ([]byte)))
	return &c
}

func (c *Client) Send(msg []byte) {
	(chan ([]byte))(*c) <- msg
}

func (c *Client) Receive() []byte {
	return <-(chan ([]byte))(*c)
}

type Channel string

func NewChannel(name string) Channel {
	return Channel(name)
}

func NewPatternChannel(pattern string) Channel {
	pattern = strings.Replace(pattern, "?", ".", -1)
	pattern = strings.Replace(pattern, "*", ".*", -1)
	return Channel("^" + pattern + "$")
}

func (c Channel) MatchPattern(pattern Channel) bool {
	match, err := regexp.Match(string(pattern), []byte(string(c)))
	return match && err == nil
}

func (db *DB) Subscribe(client *Client, channels ...Channel) int {
	count := 0
	for _, channel := range channels {
		if _, exist := db.Subscriptions[channel]; !exist {
			db.Subscriptions[channel] = make([]*Client, 0, 1)
		}
		db.Subscriptions[channel] = append(db.Subscriptions[channel], client)
		count++
	}
	return count
}

func (db *DB) Psubscribe(client *Client, patterns ...Channel) int {
	count := 0
	for _, pattern := range patterns {
		if _, exist := db.Psubscriptions[pattern]; !exist {
			db.Psubscriptions[pattern] = make([]*Client, 0, 1)
		}
		db.Psubscriptions[pattern] = append(db.Psubscriptions[pattern], client)
		count++
	}
	return count
}

func (db *DB) Unsubscribe(client *Client, channels ...Channel) int {
	count := 0
	for _, channel := range channels {
		if _, exist := db.Subscriptions[channel]; !exist {
			continue
		}

		found := -1
		for index, subscriber := range db.Subscriptions[channel] {
			if subscriber == client {
				found = index
				break
			}
		}

		if found >= 0 {
			db.Subscriptions[channel] = append(db.Subscriptions[channel][:found], db.Subscriptions[channel][found+1:len(db.Subscriptions[channel])]...)
			count++
		}
	}
	return count
}

func (db *DB) Punsubscribe(client *Client, patterns ...Channel) int {
	count := 0
	for _, pattern := range patterns {
		if _, exist := db.Psubscriptions[pattern]; !exist {
			continue
		}

		found := -1
		for index, subscriber := range db.Psubscriptions[pattern] {
			if subscriber == client {
				found = index
				break
			}
		}

		if found >= 0 {
			db.Psubscriptions[pattern] = append(db.Psubscriptions[pattern][:found], db.Psubscriptions[pattern][found+1:len(db.Psubscriptions[pattern])]...)
			count++
		}
	}
	return count
}

func (db *DB) Publish(channel Channel, bytes []byte) int {
	count := 0
	if _, exist := db.Subscriptions[channel]; exist {
		for _, client := range db.Subscriptions[channel] {
			client.Send(bytes)
			count++
		}
	}

	for pattern, subscriptions := range db.Psubscriptions {
		if channel.MatchPattern(pattern) {
			for _, client := range subscriptions {
				client.Send(bytes)
				count++
			}
		}
	}
	return count
}
