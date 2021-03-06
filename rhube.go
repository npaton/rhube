package rhube

type DB struct {
	StringsMap     map[string][]byte
	HashesMap      map[string]map[string]string
	SetsMap        map[string]map[string]bool
	ListsMap       map[string][][]byte
	ExpiresMap     map[string]Expire
	ZsetsMap       map[string]map[string]float64
	Scheduler      *Scheduler
	Subscriptions  Subscriptions
	Psubscriptions Subscriptions
}

func NewDB() *DB {
	s := NewScheduler()
	db := &DB{
		StringsMap:     make(map[string][]byte),
		HashesMap:      make(map[string]map[string]string),
		SetsMap:        make(map[string]map[string]bool),
		ZsetsMap:       make(map[string]map[string]float64),
		ListsMap:       make(map[string][][]byte),
		ExpiresMap:     make(map[string]Expire),
		Scheduler:      s,
		Subscriptions:  Subscriptions(make(map[Channel][]*Client)),
		Psubscriptions: Subscriptions(make(map[Channel][]*Client)),
	}
	s.db = db
	return db
}
