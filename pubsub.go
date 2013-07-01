package rhube

import (
	// "net"
	"regexp"
	"strings"
)

type Pubsuber struct {
	patterns   map[string][]subscription
	channels   map[string][]subscription
	publishers []string
}

type subscription struct {
	pipe chan ([]byte)
	id int
}

var pubsubUniqId int

func init() {
	pubsubUniqId = 0
}

func pubsubUniqID() int {
	pubsubUniqId++
	return pubsubUniqId
}

func (p *Pubsuber) Psubscribe(pipe chan ([]byte), patterns ...string) int {
	count := 0
	for _, pattern := range patterns {
		if _, exist := p.patterns[pattern]; !exist {
			p.patterns[pattern] = make([]subscription, 0)
		}
		s := subscription{pipe, pubsubUniqID()}
		p.patterns[pattern] = append(p.patterns[pattern], s)
		n := p.subscribePattern(pattern, s)
		count += n
	}
	return count
}

func (p *Pubsuber) Publish(pipe chan ([]byte), channel string, message []byte) int {
	count := 0
	if subs, exists := p.channels[channel]; exists {
		for _, sub := range subs {
			sub.pipe <- message
			count++
		}
	}
	return count
}

func (p *Pubsuber) Subscribe(pipe chan ([]byte), channels ...string) int {
	count := 0
	for _, channel := range channels {
		if _, exist := p.channels[channel]; !exist {
			p.channels[channel] = make([]subscription, 0)
		}
		s := subscription{pipe, pubsubUniqID()}
		p.channels[channel] = append(p.channels[channel], s)
		count++
	}
	return count
}

func (p *Pubsuber) Unsubscribe(pipe chan ([]byte), channels ...string) int {
	count := 0
	for _, channel := range channels {
		if _, exist := p.channels[channel]; !exist {
			p.channels[channel] = make([]subscription, 0)
		}
		s := subscription{pipe, pubsubUniqID()}
		p.channels[channel] = append(p.channels[channel], s)
		count++
	}
	return count
}

func (p *Pubsuber) subscribePattern(pattern string, sub subscription) int {
	count := 0
	pattern = strings.Replace(pattern, "?", ".", -1)
	pattern = strings.Replace(pattern, "*", ".*", -1)
	pattern = "^" + pattern + "$"

	for _, publisher := range p.publishers {
		match, err := regexp.Match(pattern, []byte(publisher))
		if match && err == nil {
			if _, exist := p.channels[publisher]; !exist {
				p.channels[publisher] = make([]subscription, 0)
			}
			p.channels[publisher] = append(p.channels[publisher], sub)
			count++
		}
	}

	return count
}

func (p *Pubsuber) unsubscribePattern(pattern string, sub subscription) int {
	count := 0
	pattern = strings.Replace(pattern, "?", ".", -1)
	pattern = strings.Replace(pattern, "*", ".*", -1)
	pattern = "^" + pattern + "$"

	for _, publisher := range p.publishers {
		match, err := regexp.Match(pattern, []byte(publisher))
		if match && err == nil {
			if _, exist := p.channels[publisher]; exist {
				for i, s := range p.channels[publisher] {
					if sub.id == s.id {
						p.channels[publisher] = append(p.channels[publisher][:i], p.channels[publisher][i+1:]...)
						break
					}
				}
				count++
			}
		}
	}

	return count
}
