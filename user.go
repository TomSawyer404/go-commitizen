// user.go - version 0.1
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
// @Author:  MrBanana
// @Date:    2021-8-13
// @Licence: The MIT Licence

package main


import(
    "fmt"
    "os"
    "os/exec"
    "bufio"
    "unicode"
    "strconv"
)

// stage -1: ERROR stage
// stage 0: Ready to choose options of <header-type>;
// stage 1: finish <header-type>, ready to write <header-scope>;
// stage 2: finish <header-scope>, ready to write <header-subject>
// stage 3: finish <header-struct>, ready to write <Body>;
// stage 4: finish <Body>, ready to write <Footer>
// stage 5: finish all ... exit stage
type User struct {
    stage   int
    header  string
    body    string
    footer  string
}

// Constructor of User struct
func NewUser() *User {
    new_user := &User{
        stage: 0,
        header: "",
        body: "\n",
        footer: "\n",
    }
    return new_user
}

func (this *User) driver() {
    switch this.stage {
    case 0:
        show_header_type()

        // Get user input in digit
        var user_input string
        user_input = this.get_digit_input()

        if user_choice, err := strconv.ParseInt(user_input, 10, 32); err == nil {
            switch user_choice {
            case 1:
                this.header += "feat"
                this.stage = 1
                break
            case 2:
                this.header += "fix"
                this.stage = 1
                break
            case 3:
                this.header += "docs"
                this.stage = 1
                break
            case 4:
                this.header += "style"
                this.stage = 1
                break
            case 5:
                this.header += "factor"
                this.stage = 1
                break
            case 6:
                this.header += "test"
                this.stage = 1
                break
            case 7:
                this.header += "chore"
                this.stage = 1
                break
            default:
                fmt.Println(">>>>>>>>>>>>>>>> Out of range baby!\n")
                this.stage = 0
            }
        } else {
            fmt.Println("ERROR strconv.ParseInt():", err)
            os.Exit(1)
        }
        break

    case 1:  // header-scope
        show_header_scope()

        var user_input string
        user_input = this.get_digit_input()

        if user_choice, err := strconv.ParseInt(user_input, 10, 32); err == nil {
            switch user_choice {
            case 0:
                this.header += ": "
                this.stage = 2
                break
            case 1:
                this.header += "(Data-Layer): "
                this.stage = 2
                break
            case 2:
                this.header += "(Control-Layer): "
                this.stage = 2
                break
            case 3:
                this.header += "(View-Layer): "
                this.stage = 2
                break
            default:
                fmt.Println(">>>>>>>>>>>>>>>> Out of range baby!\n")
                this.stage = 1
            }
        } else {
            fmt.Println("ERROR strconv.ParseInt(): ", err)
            os.Exit(1)
        }
        break

    case 2: // header-subject
        fmt.Println("\nPlease input your header-subject: ")

        reader := bufio.NewReader(os.Stdin)
        if line, err := reader.ReadString('\n'); err == nil {

            if len(line) == 1 && line == "\n" {
                this.stage = 2
                fmt.Println(">>>>>>>> header-subject is a must!")
                break
            }

            this.stage = 3
            this.header += line
        } else{
            fmt.Println(err.Error())
            os.Exit(1)
        }

        break

    case 3:  // Body
        reader := bufio.NewScanner(os.Stdin)
        fmt.Println("\nPlease input your message-body<quit to stop>: ")

        for reader.Scan() {
            line := reader.Text()

            if line == "quit" {
                break
            }

            this.body = this.body + line + "\n"
        }

        this.stage = 4
        // debug only
        //fmt.Println("\n>>>>>>>>>>>>>>> what I read:")
        //fmt.Print(this.body)

        break

        case 4:  // Footer
        reader := bufio.NewScanner(os.Stdin)
        fmt.Println("\nPlease input your message-footer<quit to stop>: ")

        for reader.Scan() {
            line := reader.Text()

            if line == "quit" {
                break
            }

            this.footer = this.footer + line + "\n"
        }

        this.stage = 5

        // debug only
        //fmt.Println("\n>>>>>>>>>>>>>>> what I read:")
        //fmt.Print(this.footer)

        break

        case 5:  // Write message to a file and <git commit -F> it
        // Open a file in /tmp/message
        fd, err := os.OpenFile("/tmp/git-message.txt", os.O_RDWR|os.O_CREATE, 0666)
        if err != nil {
            fmt.Println("ERROR in os.OpenFile:", err)
            os.Exit(1)
        }

        // Write message to this file
        written_num, err2 := fd.WriteString(this.header + this.body + this.footer)
        if err2 != nil || 0 == written_num {
            fmt.Println("ERROR in fd.WriteString")
            os.Exit(1)
        }

        // Execve("/bin/bash", "git commit -F ~/message")
        cmd := exec.Command("/usr/bin/git", "commit", "-F", "/tmp/git-message.txt")
        if err = cmd.Run(); err != nil {
            fmt.Printf("Command finished with error: %v", err)
            os.Exit(1)
        }

        // Delete the tmp file
        if err = os.Remove("/tmp/git-message.txt"); err != nil {
            fmt.Println("ERROR in os.Remove:", err)
            os.Exit(1)
        }

        this.stage = 6
        break

    default:
        // Something wrong...
        fmt.Printf("What's that stage? %d\n", this.stage)
        os.Exit(1)
    }
}


func is_num(user_input string) bool {
    for _, ch_val := range user_input {
        if !( unicode.IsDigit(ch_val) ) {
            return false
        }
    }
    return true
}

func (this *User) get_digit_input() string {
    var user_input string
    if _, err := fmt.Scanln(&user_input); err != nil {
        fmt.Println("ERROR fmt.Scanln():", err)
        os.Exit(1)
    }

    if !( is_num(user_input) ) {
        fmt.Println("ERROR is_num(): please input a number!")
        os.Exit(1)
    }
    return user_input
}


func show_header_type() {
    fmt.Println("\nType. Please choose your header-type: ")
    fmt.Println("  1. feat")
    fmt.Println("  2. fix")
    fmt.Println("  3. docs")
    fmt.Println("  4. style")
    fmt.Println("  5. factor")
    fmt.Println("  6. test")
    fmt.Println("  7. chore")
    fmt.Print("Your option: ")
}


func show_header_scope() {
    fmt.Println("\nType. Please choose your header-scope: ")
    fmt.Println("  0. None, do not add a scope in header")
    fmt.Println("  1. Data-Layer")
    fmt.Println("  2. Control-Layer")
    fmt.Println("  3. View-Layer")
    fmt.Print("Your option: ")
}

