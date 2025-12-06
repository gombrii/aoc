package commands

import (
	"errors"
	"fmt"

	"github.com/gombrii/aoc/internal/cache"
	"github.com/gombrii/aoc/internal/com"
	"github.com/gombrii/aoc/internal/files"
)

const User = "user"

func (c Commands) Login(session string) error {
	username, err := com.Ping(com.NewClient(session))
	if err != nil {
		if errors.Is(err, com.ErrUnauthorized) {
			return errors.New("invalid session token")
		}
		return fmt.Errorf("pinging server: %v", err)
	}

	paths, err := files.GenTemp(map[string]string{files.Session: session}, nil)
	if err != nil {
		return fmt.Errorf("creating session file: %v", err)
	}

	_, err = cache.Store(cache.ConfigKey(User), files.Session, paths[files.Session])
	if err != nil {
		return fmt.Errorf("caching session token: %v", err)
	}

	fmt.Println("Logged in with user", username)

	return nil
}
