package rhube

import (
	"testing"
	"time"
)

func TestScheduleStuff(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("foo"))
	db.expireKeyAt("ho", time.Now().Add(time.Duration(1)*time.Second))
	if _, ok := db.ExpiresMap["ho"]; !ok {
		t.Fatalf("db.ExpiresMap[\"ho\"] should NOT be nil", db.ExpiresMap)
	}
	
	db.cancelExpireKey("ho")
	if _, ok := db.ExpiresMap["ho"]; ok {
		t.Fatalf("db.ExpiresMap[\"ho\"] should BE nil", db.ExpiresMap)
	}
}

func TestScheduleStuffRename(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("foo"))
	db.expireKeyAt("ho", time.Now().Add(time.Duration(1)*time.Second))
	if _, ok := db.ExpiresMap["ho"]; !ok {
		t.Fatalf("db.ExpiresMap[\"ho\"] should NOT be nil", db.ExpiresMap)
	}
	
	db.Rename("ho", "ha")
	if _, ok := db.ExpiresMap["ho"]; ok {
		t.Fatalf("db.ExpiresMap[\"ho\"] should BE nil", db.ExpiresMap)
	}
	if _, ok := db.ExpiresMap["ha"]; !ok {
		t.Fatalf("db.ExpiresMap[\"ha\"] should NOT be nil", db.ExpiresMap)
	}
	
	db.Renamenx("ha", "ho")
	if _, ok := db.ExpiresMap["ha"]; ok {
		t.Fatalf("db.ExpiresMap[\"ho\"] should BE nil", db.ExpiresMap)
	}
	if _, ok := db.ExpiresMap["ho"]; !ok {
		t.Fatalf("db.ExpiresMap[\"ha\"] should NOT be nil", db.ExpiresMap)
	}
	
	db.Set("hi", []byte("foo"))
	db.Renamenx("ho", "hi")
	if _, ok := db.ExpiresMap["ha"]; ok {
		t.Fatalf("db.ExpiresMap[\"ho\"] should BE nil", db.ExpiresMap)
	}
	if _, ok := db.ExpiresMap["ho"]; !ok {
		t.Fatalf("db.ExpiresMap[\"ha\"] should NOT be nil", db.ExpiresMap)
	}
}
