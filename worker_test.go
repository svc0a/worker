package worker

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
)

type User struct {
	Name string
	Age  int
}

func TestNew(t *testing.T) {
	list := []User{}
	for i := 0; i < 10; i++ {
		list = append(list, User{
			Name: fmt.Sprintf("name%d", i),
			Age:  i,
		})
	}
	Walk[User, error](list, func(user User) error {
		logrus.Info(user)
		return nil
	})
}
