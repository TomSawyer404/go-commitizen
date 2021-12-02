package main

import (
	"fmt"
	"git-cz/internal/user"
	//"git-cz/internal/user-old"
)

func main() {
	new_user := user.NewUser()

	for {
		if new_user.Stage >= 0 && new_user.Stage <= 5 {
			new_user.Driver()
		} else if new_user.Stage == 6 {
			fmt.Println("\n\x1b[32m\t====== Done! Good bye! ======\x1b[0m")
			break
		} else {
			fmt.Println("\n\x1b[31mSomething WRONG!!!  UNKNOWN STAGE!\x1b[0m")
			break
		}
	}

}
