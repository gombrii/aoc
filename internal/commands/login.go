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
			fmt.Println("Invalid session token")
			return nil
		}
		return fmt.Errorf("pinging server: %v", err)
	}

	cPath := cache.MakePath(cache.ConfigKey(User), files.Session)
	err = files.Write(cPath, []byte(session))
	if err != nil {
		return fmt.Errorf("caching session token: %v", err)
	}

	fmt.Println("Logged in with user", username)

	return nil
}
