package redis

import (
	"fmt"
	"strings"
	"sync"
)

type Store struct {
	mu    sync.RWMutex
	items map[string]string
}

func NewStore() *Store {
	return &Store{
		items: make(map[string]string),
	}
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = value
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exists := s.items[key]
	return value, exists
}

func (s *Store) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.items[key]; exists {
		delete(s.items, key)
		return true
	}
	return false
}

type CommandHandler struct {
	store *Store
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		store: NewStore(),
	}
}

func (h *CommandHandler) HandleCommand(args []string) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no command provided")
	}

	cmd := strings.ToUpper(args[0])
	switch cmd {
	case "SET":
		if len(args) != 3 {
			return nil, fmt.Errorf("wrong number of arguments for SET")
		}
		h.store.Set(args[1], args[2])
		return "OK", nil

	case "GET":
		if len(args) != 2 {
			return nil, fmt.Errorf("wrong number of arguments for GET")
		}
		value, exists := h.store.Get(args[1])
		if !exists {
			return nil, nil
		}
		return value, nil

	case "DEL":
		if len(args) != 2 {
			return nil, fmt.Errorf("wrong number of arguments for DEL")
		}
		success := h.store.Del(args[1])
		if success {
			return 1, nil
		}
		return 0, nil

	case "PING":
		return "PONG", nil

	default:
		return nil, fmt.Errorf("unknown command: %s", cmd)
	}
}
