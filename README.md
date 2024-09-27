# worker

```
	Walk[User, error](list, func(user User) error {
		logrus.Info(user)
		return nil
	})
```

# map

```
	m := Define[User]()
	m.Store("user1", User{"tom", 18})
	data, err := m.Load("user1")
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info(data)
```