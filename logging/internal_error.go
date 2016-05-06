package logging

func Panic(err error) {
	//TODO: log to Sentry or similar service. if we end up here, that's bad!
	panic(err)
}
