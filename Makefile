v2: commitizen-go.go user-v2.go
	go build -o git-cz $^

v1: commitizen-go.go user.go
	go build -o git-cz $^

clean:
	rm git-cz
