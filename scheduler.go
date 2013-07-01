package rhube

// Todo: Check computer clock infrequently, expire old things if clock changes

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

type Scheduler struct {
	minutes  map[int64]*expiresCollection
	refresh  chan (bool)
	nextTick time.Time
	running  bool
	db       *DB
	sync.Mutex
}

func NewScheduler() *Scheduler {
	s := &Scheduler{
		minutes:  make(map[int64]*expiresCollection),
		nextTick: time.Now().Add(time.Minute),
		refresh:  make(chan (bool)),
	}

	go s.Run()
	return s
}

func (s *Scheduler) Add(e *Expire) {
	expireMinute := e.expireMinute()

	if e.TTL() < 0 {
		fmt.Println("Expiring", ("'" + (*e.key) + "',"), "late:", time.Now().Sub(e.at))
		s.db.Del(*e.key)
		fmt.Println("Expired", ("'" + (*e.key) + "',"), "late:", time.Now().Sub(e.at))
		return
	}

	if _, exist := s.minutes[expireMinute]; !exist {
		s.minutes[expireMinute] = &expiresCollection{make([]*Expire, 0, 1)}
	}

	s.Lock()
	s.minutes[expireMinute].expires = append(s.minutes[expireMinute].expires, e)
	sort.Sort(s.minutes[expireMinute])
	s.Unlock()

	if s.running && e.at.UnixNano() <= s.nextTick.UnixNano() {
		s.refresh <- true
	}
}

func (s *Scheduler) Remove(e *Expire) bool {
	found := -1
	min := e.expireMinute()
	s.Lock()
	for index, exp := range s.minutes[min].expires {
		if *exp.key == *e.key {
			found = index
			break
		}
	}

	if found >= 0 {
		if _, exists := s.minutes[min]; exists {
			s.minutes[min].expires = append(s.minutes[min].expires[:found], s.minutes[min].expires[found+1:]...)
			s.Unlock()
			return true
		}
		return false
	}
	s.Unlock()
	return false
}

func (s *Scheduler) Run() {
	s.running = true
	for {
		min := time.Now().Unix() / 60

		s.expireExpires(min - 1)
		s.expireExpires(min)

		expireIn := s.nextExpire(min)
		if expireIn < 0 {
			continue
		}
		// fmt.Println("Next expire check in:", expireIn)
		s.nextTick = time.Now().Add(expireIn)
		timeout := time.After(expireIn)
		select {
		case <-timeout:
			continue
		case <-s.refresh:
			continue
		}
	}
}

func (s *Scheduler) expireExpires(min int64) {
	s.Lock()
	if _, exists := s.minutes[min]; exists {
		for _, exp := range s.minutes[min].expires {
			if exp.at.UnixNano() <= time.Now().UnixNano() {
				s.Unlock()
				s.db.Del(*exp.key)
				fmt.Println("Expired", ("'" + (*exp.key) + "',"), "late:", time.Now().Sub(exp.at))
				s.Lock()
			} else {
				break
			}
			if len(s.minutes[min].expires) == 0 {
				delete(s.minutes, min)
			}
		}
	}
	s.Unlock()
}

func (s *Scheduler) nextExpire(min int64) time.Duration {
	s.Lock()
	for minutes := 0; minutes < 1440; minutes++ {
		if _, exists := s.minutes[min+int64(minutes)]; exists && len(s.minutes[min+int64(minutes)].expires) > 0 {
			s.Unlock()
			return s.minutes[min+int64(minutes)].expires[0].pttlDur() - time.Millisecond
		}
	}
	s.Unlock()

	return time.Hour
}

// Used for sorting
type expiresCollection struct {
	expires []*Expire
}

func (e *expiresCollection) Len() int {
	return len(e.expires)
}

func (e *expiresCollection) Swap(i, j int) {
	e.expires[i], e.expires[j] = e.expires[j], e.expires[i]
}

func (e *expiresCollection) Less(i, j int) bool {
	return e.expires[i].at.UnixNano() < e.expires[j].at.UnixNano()
}

type Expire struct {
	at  time.Time
	key *string
}

func (e *Expire) TTL() int {
	return int(e.at.Sub(time.Now()).Seconds())
}

func (e *Expire) PTTL() int {
	return int(e.at.Sub(time.Now()).Seconds() / 1000)
}

func (e *Expire) pttlDur() time.Duration {
	return e.at.Sub(time.Now())
}

func (e *Expire) expireMinute() int64 {
	t := e.at.Unix()
	rem := math.Abs(math.Remainder(float64(t), float64(60.0)))
	return int64(float64(t)-rem) / 60
}

func (db *DB) expireKeyIn(key string, milliseconds int) bool {
	newExpireTime := time.Now().Add(time.Duration(milliseconds) * time.Millisecond)
	return db.expireKeyAt(key, newExpireTime)
}

func (db *DB) expireKeyAt(key string, newExpireTime time.Time) bool {
	if newExpireTime.UnixNano() < time.Now().UnixNano() {
		return false
	}

	if expire, ok := db.ExpiresMap[key]; ok {
		db.Scheduler.Remove(&expire)
		expire.at = newExpireTime
		db.ExpiresMap[key] = expire
		db.Scheduler.Add(&expire)
	} else {
		if !db.Exists(key) {
			return false
		}
		exp := Expire{key: &key, at: newExpireTime}
		db.ExpiresMap[key] = exp
		db.Scheduler.Add(&exp)
	}

	return true
}

func (db *DB) cancelExpireKey(key string) bool {
	if exp, ok := db.ExpiresMap[key]; ok {
		db.Scheduler.Remove(&exp)
		delete(db.ExpiresMap, key)
		return true
	}

	return false
}

func (db *DB) renameExpire(key, newKey string) bool {
	if _, ok := db.ExpiresMap[key]; ok {
		exp := db.ExpiresMap[key]
		db.Scheduler.Remove(&exp)
		db.ExpiresMap[newKey] = db.ExpiresMap[key]
		*db.ExpiresMap[newKey].key = newKey
		exp = db.ExpiresMap[newKey]
		db.Scheduler.Add(&exp)
		delete(db.ExpiresMap, key)
		return true
	}
	return false
}
