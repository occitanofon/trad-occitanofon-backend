package activity

import (
	"sync"
	"time"
)

const DELAY time.Duration = 10 * time.Minute

type Service struct {
	translatorLastSeen map[string]time.Time
	mu                 *sync.RWMutex
}

func NewService() Service {
	service := Service{
		translatorLastSeen: make(map[string]time.Time),
		mu:                 &sync.RWMutex{},
	}

	go service.removeInactiveTranslator()

	return service
}

// removeInactiveTranslator removes translator who is inactive every 10 minutes
func (s *Service) removeInactiveTranslator() {
	for {
		<-time.After(DELAY)

		for translatorID, lastSeen := range s.translatorLastSeen {
			now := time.Now()
			if now.After(lastSeen.Add(DELAY)) {
				s.Delete(translatorID)
			}
		}
	}
}

// AddOrKeepActive adds new translator to the list of active translators or keep him/she active
func (s *Service) AddOrKeepActive(translatorID string) {
	s.mu.Lock()
	s.translatorLastSeen[translatorID] = time.Now()
	s.mu.Unlock()
}

// Delete deletes a given translator from the list of active translators
func (s *Service) Delete(translatorID string) {
	s.mu.Lock()
	delete(s.translatorLastSeen, translatorID)
	s.mu.Unlock()
}

// Total returns how many translators are active
func (s *Service) Total() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.translatorLastSeen)
}
