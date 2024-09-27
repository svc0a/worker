package syncx

import (
	"github.com/sirupsen/logrus"
	"testing"
)

type User struct {
	Name string
	Age  int
}

func TestDefine(t *testing.T) {
	m := Define[User]()
	m.Store("user1", User{"tom", 18})
	data, err := m.Load("user1")
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info(data)
}
