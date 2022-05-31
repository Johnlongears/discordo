package ui

import (
	"fmt"
	"strings"

	"github.com/ayntgl/astatine"
	"github.com/google/shlex"
	"github.com/rivo/tview"
)

var commandMap = make(map[string]Command)

func HandleCommand(mi *MessageInputField, t string, m *astatine.Message) bool {
	argv, _ := shlex.Split(t)
	argc := len(argv)
	cmd,found := commandMap[argv[0]]
	if found {
		cmd.Execute(mi, argv, argc, m)
	}
    mi.app.SelectedMessage = -1
		if cmd.Terminating {
    mi.SetText("")
		}
    mi.SetTitle("")
		return cmd.Terminating
	} else {
		mi.app.MessagesTextView.Write([]byte("[[#FFFF00]SYSTEM[-]] Unknown command.\n"));
		mi.SetTitle("");
		return true;
	}
}

var commands []Command

func InitCommands() {
	commands = []Command{
		{
			Names:       []string{"/help", "/?"},
			Description: "Display information about available commands",
			Usage:       "%s [command name]",
			Execute: func(mi *MessageInputField, argv []string, argc int, m *astatine.Message) {
				if argc > 1 {
                    var cmd Command;
                        cmd = commandMap[argv[1]]
                    } else {
                        cmd = commandMap["/" + argv[1]]
                    }
                    list := CreateList(mi.app.MessagesTextView)
					list.AddItem("Name:",cmd.Names[0],0,nil)
                    list.AddItem("Description:",cmd.Description,0,nil)
                    list.AddItem("Usage:",fmt.Sprintf(cmd.Usage,cmd.Names[0]),0,nil)
                    list.AddItem("Aliases:",strings.Join(cmd.Names[1:],", "),0,nil);
					mi.app.SetRoot(list, true)
				} else {
					list := CreateList(mi.app.MessagesTextView)
					for _, cmd := range commands {
						dcmd := cmd
						list.AddItem(
							cmd.Names[0]+" -> "+cmd.Description,
							fmt.Sprintf(" ╰ Usage: "+cmd.Usage+
								" | Aliases: "+strings.Join(cmd.Names[1:], ", "), cmd.Names[0]),
							0,
							func(){
								mi.SetText(dcmd.Names[0] +" ")
								mi.app.
									SetRoot(mi.app.MainFlex, true).
									SetFocus(mi)
							})
					}
					mi.app.SetRoot(list, true)
				}
			},
			Terminating: true,
		},
		{
			Names:       []string{"/shrug","/shrugend"},
			Description: "Send or prepend a ¯\\_(ツ)_/¯ to your message. Use /shrugend to append ¯\\_(ツ)_/¯ instead",
			Usage:       "%s [message content]",
			Execute: func(mi *MessageInputField, argv []string, argc int, m *astatine.Message) {
				if argc > 1 {
					if argv[0] == "/shrugend" {
						mi.SetText(strings.Join(argv[1:], " ") + " ¯\\_(ツ)_/¯")
					} else {
						mi.SetText("¯\\_(ツ)_/¯ " + strings.Join(argv[1:], " "))
					}
				} else {
					mi.SetText("¯\\_(ツ)_/¯")
				}
			},
			Terminating: false,
		},
	}
	for _, cmd := range commands {
		for _, name := range cmd.Names {
			commandMap[name] = cmd
		}
	}
}

type Command struct {
	Names       []string
	Description string
	Usage       string
	Execute     func(mi *MessageInputField, argv []string, argc int, m *astatine.Message)
	Terminating bool // If true the command will stop message sending, if false, it will send a message when done.
}

func CreateList(mtv *MessagesTextView) *tview.List {
	list := tview.NewList()
	list.SetDoneFunc(func() {
		mtv.app.
			SetRoot(mtv.app.MainFlex, true).
			SetFocus(mtv.app.MessageInputField)
	})
	list.SetTitle("Press the Escape key to close")
	list.SetTitleAlign(tview.AlignLeft)
	list.SetBorder(true)
	list.SetBorderPadding(0, 0, 1, 1)
	return list
}
