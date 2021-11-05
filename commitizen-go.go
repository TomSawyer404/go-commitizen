// commitizen-go.go - version 0.1
// @Note
//  - Main moduel of the project
//  - First we construct header of commit message
//  - Second we construct body of commit message
//  - Third we construct footer of commit message
//  - Last we execve `git commit -m "balabla"`
//
// @Author:  MrBanana
// @Date:    2021-8-13
// @Licence: The MIT Licence

package main

import "fmt"

func main() {
	new_user := NewUser()

	for {
		if new_user.stage >= 0 && new_user.stage <= 5 {
			new_user.driver()
		} else if new_user.stage == 6 {
			fmt.Println("\n\x1b[32m\t======== Done! Good bye! ========")
			break
		} else {
			fmt.Println("\nSomething WRONG!!!  UNKNOWN STAGE!")
			break
		}
	}

}
