git-cz: cmd/commitizen-go.go 
	go build -o $@ $^

clean:
	rm git-cz
