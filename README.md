# worker


```
	Walk[User, error](list, func(user User) error {
		logrus.Info(user)
		return nil
	})
```