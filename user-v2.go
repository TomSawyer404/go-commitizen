// user.go - version 0.2
// @Note
//  - Definition of User Struct
//  - Format: <Header> ... <Body> ... <Footer>
//          <Header>: <Type>(<scope>): <subject>
//          <Body>:     ... anything ... could have many lines
//          <Footer>:   Close #12345
//
//  - First we show menu to user;
//  - Then we take user's input from stdin;
//  - Make a choice, keep going until we construct all the message;
//  - Rules come from [here](https://www.ruanyifeng.com/blog/2016/01/commit_message_change_log.html);
//
//  version 0.2: Refactor with plugin `go-prompt`
//
// @Author:  MrBanana
// @Date:    2021-8-18
// @Licence: The MIT Licence

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/c-bata/go-prompt"
)

var (
	headType = map[string]struct{}{
		"feat":   struct{}{},
		"fix":    struct{}{},
		"docs":   struct{}{},
		"style":  struct{}{},
		"factor": struct{}{},
		"test":   struct{}{},
		"chore":  struct{}{},
	}

	headScope = map[string]struct{}{
		"repo":    struct{}{},
		"model":   struct{}{},
		"logic":   struct{}{},
		"handler": struct{}{},
	}
)

// stage -1: ERROR stage
// stage 0: Ready to choose options of <header-type>;
// stage 1: finish <header-type>, ready to write <header-scope>;
// stage 2: finish <header-scope>, ready to write <header-subject>
// stage 3: finish <header-struct>, ready to write <Body>;
// stage 4: finish <Body>, ready to write <Footer>
// stage 5: finish all ... exit stage
type User struct {
	stage  int
	header string
	body   string
	footer string
}

// Constructor of User struct
func NewUser() *User {
	new_user := &User{
		stage:  0,
		header: "",
		body:   "\n",
		footer: "\n",
	}
	return new_user
}

func (this *User) driver() {
	switch this.stage {
	case 0:
		t := prompt.Input("  Please choose your header-type: ", completer_header_type)
		if _, ok := headType[t]; ok {
			this.header += t
			this.stage = 1
		} else {
			fmt.Println("Wrong header-type, please try again!")
			this.stage = 0
		}
		break

	case 1: // header-scope
		fmt.Println()
		t := prompt.Input("  Please choose your header-scope: ", completer_header_scope)
		if "None" == t {
			this.header += `: `
			this.stage = 2
		} else if _, ok := headScope[t]; ok {
			this.header += `(` + t + `): `
			this.stage = 2
		} else {
			fmt.Println("Wrong header-type, please try again!")
			this.stage = 1
		}
		break

	case 2: // header-subject
		fmt.Println()
		fmt.Println("  Please enter your header-subject: ")
		t := prompt.Input(">>> ", completer_header_subject)
		if "" == t {
			fmt.Println("header-subject is a must! please try again!")
			this.stage = 2
		} else {
			this.header += t + "\n"
			this.stage = 3
		}

		break

	case 3: // Body
		fmt.Println("\nPlease input your message body: ")
		for {
			t := prompt.Input("body> ", completer_multiline)
			if "" == t {
				break
			}
			this.body += t + "\n"
		}
		this.stage = 4

		break

	case 4: // Footer
		fmt.Println("\nPlease input your message footer: ")
		for {
			t := prompt.Input("footer> ", completer_multiline)
			if "" == t {
				break
			}
			this.footer += t + "\n"
		}
		this.stage = 5

		break

	case 5: // Write message to a file and <git commit -F> it
		// Open a file in /tmp/message
		fd, err := os.OpenFile("/tmp/git-message.txt", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalln("ERROR in os.OpenFile ->", err)
		}

		// Write message to this file
		written_num, err2 := fd.WriteString(this.header + this.body + this.footer)
		if err2 != nil || 0 == written_num {
			log.Fatalln("ERROR in fd.WriteString ->", err2)
		}

		// Execve("/bin/bash", "git commit -F ~/message")
		cmd := exec.Command("/usr/bin/git", "commit", "-F", "/tmp/git-message.txt")
		if err = cmd.Run(); err != nil {
			log.Fatalln("ERROR in exec.Command ->", err)
		}

		// Delete the tmp file
		if err = os.Remove("/tmp/git-message.txt"); err != nil {
			log.Fatalln("ERROR in os.Remove ->", err)
		}

		this.stage = 6
		break

	default:
		// Something wrong...
		log.Printf("\x1b[31m\t-------- What's that stage? %d\n", this.stage)
	}
}

func completer_header_type(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "feat", Description: "I am trying to add a new feature"},
		{Text: "fix", Description: "Some bugs have fixed"},
		{Text: "docs", Description: "The documents has been modified"},
		{Text: "style", Description: "Changes that do not affect code operation"},
		{Text: "factor", Description: "It is neither a new function nor a code change to modify a bug"},
		{Text: "test", Description: "Add some test cases"},
		{Text: "chore", Description: "Changes in the construction process or auxiliary tools"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func completer_header_scope(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "None", Description: "No effect on each layer"},
		{Text: "repo", Description: "Operations on persistent data storage"},
		{Text: "model", Description: "Assemble and operate data"},
		{Text: "logic", Description: "Realize specific business logic on demand"},
		{Text: "handler", Description: "Control business process"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func completer_header_subject(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "type", Description: "Type something to descript your commit"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func completer_multiline(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return s
}
