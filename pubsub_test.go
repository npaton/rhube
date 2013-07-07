package rhube

import (
	"testing"
)

func TestSubscribePublish(t *testing.T) {
	db := NewDB()
	client := NewClient()
	channel := NewChannel("toto")

	if db.Subscribe(client, channel) != 1 {
		t.Fail()
	}
	var result []byte
	go func() {
		result = client.Receive()
	}()

	db.Publish(channel, []byte("salut!"))
	if string(result) != "salut!" {
		t.Fail()
	}

	channel1 := NewChannel("tata")
	db.Publish(channel1, []byte("salut!"))

	channel2 := NewChannel("toto")
	go func() {
		result = client.Receive()
	}()

	db.Publish(channel2, []byte("coucou!"))
	if string(result) != "coucou!" {
		t.Fail()
	}

}

func TestUnsubscribePublish(t *testing.T) {
	db := NewDB()
	client := NewClient()
	channel := NewChannel("toto")

	if db.Subscribe(client, channel) != 1 {
		t.Fail()
	}
	var result []byte
	go func() {
		result = client.Receive()
	}()

	db.Publish(channel, []byte("salut!"))
	if string(result) != "salut!" {
		t.Fail()
	}

	if db.Unsubscribe(client, channel) != 1 {
		t.Fail()
	}

	db.Publish(channel, []byte("salut!"))
}

func TestSubscribePublishPattern(t *testing.T) {
	db := NewDB()
	client := NewClient()
	pattern := NewPatternChannel("tot*")
	channel := NewChannel("toto")

	if db.Psubscribe(client, pattern) != 1 {
		t.Fail()
	}
	var result []byte
	go func() {
		result = client.Receive()
	}()

	db.Publish(channel, []byte("salut!"))
	if string(result) != "salut!" {
		t.Fail()
	}

	channel1 := NewChannel("tata")
	db.Publish(channel1, []byte("salut!"))

	channel2 := NewChannel("tota")
	go func() {
		result = client.Receive()
	}()

	db.Publish(channel2, []byte("coucou!"))
	if string(result) != "coucou!" {
		t.Fail()
	}

}

func TestUnsubscribePublishPattern(t *testing.T) {
	db := NewDB()
	client := NewClient()
	pattern := NewPatternChannel("tot*")
	channel := NewChannel("toto")

	if db.Psubscribe(client, pattern) != 1 {
		t.Fail()
	}
	var result []byte
	go func() {
		result = client.Receive()
	}()

	db.Publish(channel, []byte("salut!"))
	if string(result) != "salut!" {
		t.Fail()
	}

	if db.Punsubscribe(client, pattern) != 1 {
		t.Fail()
	}

	db.Publish(channel, []byte("salut!"))
}
