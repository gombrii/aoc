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

	if err = setSession(session); err != nil {
		return fmt.Errorf("setting session: %v", err)
	}

	fmt.Println("Logged in as user", username)

	return nil
}

func LoggedIn() (string, bool) {
	path, ok := cache.Contains(cache.ConfigKey(User), files.Session)
	if !ok {
		return "", false
	}

	data, err := files.Read(path)
	if err != nil {
		return "", false
	}

	return string(data), true
}

func setSession(session string) error {
	cPath, ok := cache.Contains(cache.ConfigKey(User), files.Session)
	if !ok {
		paths, err := files.GenTemp(map[string]string{files.Session: session}, nil)
		if err != nil {
			return fmt.Errorf("creating file: %v", err)
		}

		_, err = cache.Store(cache.ConfigKey(User), files.Session, paths[files.Session])
		if err != nil {
			return fmt.Errorf("caching token: %v", err)
		}

		return nil
	}

	files.Write(cPath, []byte(session))

	return nil
}
