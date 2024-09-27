package worker

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/svc0a/worker/syncx"
	"testing"
	"time"
)

type User struct {
	Name string
	Age  int
}

func TestNew(t *testing.T) {
	m := syncx.Define[User]()
	list := []User{}
	for i := 0; i < 250; i++ {
		list = append(list, User{
			Name: fmt.Sprintf("name%d", i),
			Age:  i,
		})
	}
	Walk[User](list, func(user User) error {
		m.Store(user.Name, user)
		time.Sleep(1 * time.Second)
		return nil
	}, WithErrHandler(func(err error) {
		logrus.Error(err)
	}), WithWorkerNumber(100), WithChanSize(100))
	logrus.Info(m.Size())
}
