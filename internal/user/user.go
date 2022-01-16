// user.go - version 0.4
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
//  version 0.3: Some details updated
//  version 0.4: Add Windows System support
//
// @Author:  MrBanana
// @Date:    2022-1-16
// @Licence: The MIT Licence

package user

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/c-bata/go-prompt"
)

var (
	headType = map[string]struct{}{
		"feat":   {},
		"fix":    {},
		"docs":   {},
		"style":  {},
		"factor": {},
		"test":   {},
		"chore":  {},
	}

	headScope = map[string]struct{}{
		"repo":    {},
		"model":   {},
		"logic":   {},
		"handler": {},
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
	Stage  int
	header string
	body   string
	footer string
}

// Constructor of User struct
func NewUser() *User {
	new_user := &User{
		Stage:  0,
		header: "",
		body:   "\n",
		footer: "\n",
	}
	return new_user
}

func (me *User) Driver() {
	switch me.Stage {
	case 0:
		t := input_prompt("Header Type", completer_header_type)

		if _, ok := headType[t]; ok {
			me.header += t
			me.Stage = 1
		} else {
			fmt.Println("Wrong header-type, please try again!")
			me.Stage = 0
		}

	case 1: // header-scope
		fmt.Println()
		t := input_prompt("Header Scope", completer_header_scope)

		if t == "None" {
			me.header += `: `
			me.Stage = 2
		} else if _, ok := headScope[t]; ok {
			me.header += `(` + t + `): `
			me.Stage = 2
		} else {
			fmt.Println("Wrong header-scope, please try again!")
			me.Stage = 1
		}

	case 2: // header-subject
		fmt.Println()
		t := input_prompt("Header Subject", completer_header_subject)

		if t == "" {
			fmt.Println("header-subject is a must! please try again!")
			me.Stage = 2
		} else {
			me.header += t + "\n"
			me.Stage = 3
		}

	case 3: // Body
		fmt.Println()
		for {
			t := input_prompt("Message Body", completer_multiline)
			if t == "" {
				break
			}
			me.body += t + "\n"
		}
		me.Stage = 4

	case 4: // Footer
		fmt.Println()
		for {
			t := input_prompt("Message Footer", completer_multiline)
			if t == "" {
				break
			}
			me.footer += t + "\n"
		}
		me.Stage = 5

	case 5: // Write message to a file and <git commit -F> it
		// Open a file in `../git-message.txt`
		git_msg_path := "../git-message.txt"
		fd, err := os.OpenFile(git_msg_path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalln("ERROR in os.OpenFile ->", err)
		}
		defer func() {
			// Delete the tmp file
			if err = fd.Close(); err != nil {
				log.Fatalln("ERROR in os.Close() ->", err)
			}
			if err = os.Remove(git_msg_path); err != nil {
				log.Fatalln("ERROR in os.Remove ->", err)
			}
		}()

		// Write message to this file
		written_num, err2 := fd.WriteString(me.header + me.body + me.footer)
		if err2 != nil || written_num == 0 {
			log.Fatalln("ERROR in fd.WriteString ->", err2)
		}

		// Execve("/bin/bash", "git commit -F ~/message")
		cmd := exec.Command("git", "commit", "-F", git_msg_path)
		if err = cmd.Run(); err != nil {
			log.Fatalln("ERROR in exec.Command ->", err)
		}

		me.Stage = 6

	default:
		// Something wrong...
		log.Printf("\x1b[31m\t-------- What's that stage? %d\n", me.Stage)
	}
}

func completer_header_type(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "feat", Description: "A new feature"},
		{Text: "fix", Description: "A bug fix"},
		{Text: "docs", Description: "Documentation only changes"},
		{Text: "style", Description: "Changes that do not affect the meaning of the code"},
		{Text: "factor", Description: "A code change that neither fixes a bug or adds a feature"},
		{Text: "perf", Description: "A code change that improves performance"},
		{Text: "test", Description: "Add some test cases"},
		{Text: "chore", Description: "Changes to the build process or auxiliary tools"},
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
		{Text: "Anything", Description: "Type anything to descript your commit"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func completer_multiline(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return s
}

func input_prompt(str string, completer prompt.Completer) string {
	return prompt.Input(">> "+str+": ", completer,
		prompt.OptionTitle("git-cz"),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray))

}
