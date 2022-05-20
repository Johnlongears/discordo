package ui

import {
  "strings"
  
  "github.com/ayntgl/astatine"
  "github.com/google/shlex"
    "github.com/rivo/tview"
}

var commandMap := make(map[string]*Command)

func HandleCommand(mi *MessageInputField,t string,m astatine.Message) {
    argv := shlex.Split(t)
    argc := len(argc)
    cmd := commandMap[argv[0]]
    if cmd != nil {
        cmd.Execute(mi,argv,argc,m)   
    }
}

func InitCommands() {
    for _, cmd := range commands {
        for _, name := range cmd.Names {
            commandMap[name] = &cmd   
        }
    }
}

type Command struct{
    Names []string
    Description string
    Usage string
    Execute func(mi *MessageInputField,argv []string, argc int,m astatine.Message)
}

var commands = [...]Command{
    Command{
        Names: []string{"/help","/?"},
        Description: "Display information about available commands",
        Usage: "%s [command name]"
        Execute: func(mi *MessageInputField,argv []string, argc int,m astatine.Message){
            
            if argc > 1 {
                
            } else {
                list := CreateList(mi.app.MessagesTextView)
                for _,cmd := range commands {
                    list.AddItem(
                        cmd.Names[0]+" "+cmd.Description,
                        fmt.Sprintf(cmd.Usage,cmd.Names[0])
                }
                mi.app.SetRoot(dialog,true)
            }
            return
        },
    },  
}

func CreateList(mtv *MessagesTextView) *tview.List {
    list := tview.NewList()
	list.SetDoneFunc(func() {
		mtv.app.
			SetRoot(mtv.app.MainFlex, true).
			SetFocus(mtv.app.MessagesTextView)
	})
	list.SetTitle("Press the Escape key to close")
	list.SetTitleAlign(tview.AlignLeft)
	list.SetBorder(true)
	list.SetBorderPadding(0, 0, 1, 1)
    return list
}
