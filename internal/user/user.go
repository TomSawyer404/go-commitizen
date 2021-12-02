// user.go - version 0.3
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
//
// @Author:  MrBanana
// @Date:    2021-12-2
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

func (this *User) Driver() {
	switch this.Stage {
	case 0:
		t := input_prompt("Header Type", completer_header_type)

		if _, ok := headType[t]; ok {
			this.header += t
			this.Stage = 1
		} else {
			fmt.Println("Wrong header-type, please try again!")
			this.Stage = 0
		}
		break

	case 1: // header-scope
		fmt.Println()
		t := input_prompt("Header Scope", completer_header_scope)

		if "None" == t {
			this.header += `: `
			this.Stage = 2
		} else if _, ok := headScope[t]; ok {
			this.header += `(` + t + `): `
			this.Stage = 2
		} else {
			fmt.Println("Wrong header-scope, please try again!")
			this.Stage = 1
		}
		break

	case 2: // header-subject
		fmt.Println()
		t := input_prompt("Header Subject", completer_header_subject)

		if "" == t {
			fmt.Println("header-subject is a must! please try again!")
			this.Stage = 2
		} else {
			this.header += t + "\n"
			this.Stage = 3
		}

		break

	case 3: // Body
		fmt.Println()
		for {
			t := input_prompt("Message Body", completer_multiline)
			if "" == t {
				break
			}
			this.body += t + "\n"
		}
		this.Stage = 4

		break

	case 4: // Footer
		fmt.Println()
		for {
			t := input_prompt("Message Footer", completer_multiline)
			if "" == t {
				break
			}
			this.footer += t + "\n"
		}
		this.Stage = 5

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

		this.Stage = 6
		break

	default:
		// Something wrong...
		log.Printf("\x1b[31m\t-------- What's that stage? %d\n", this.Stage)
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
