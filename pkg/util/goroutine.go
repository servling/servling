package util

func RunInBackground(task func() error, errorHandler func(error)) {
	go func() {
		if err := task(); err != nil {
			errorHandler(err)
		}
	}()
}
