all:
	@go run optparser.go -yn --arg1 value1 #--bool-arg -S short-arg-value --servers 'long arg value'
