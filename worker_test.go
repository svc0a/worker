package worker

import (
	"github.com/sirupsen/logrus"
	"testing"
)

type User struct {
	Name string
	Age  int
}

func TestNew(t *testing.T) {
	wp := New[User, error](1000, WithErrHandler[User, error](func(err error) {
		if err != nil {
			logrus.Error(err)
		}
	}))
	for i := 0; i < 100; i++ {
		wp.Submit(User{
			Name: "test",
			Age:  1,
		})
	}
	wp.Start(func(data User) error {
		return nil
	})
	wp.Stop()
}
