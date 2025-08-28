package main

import (
	"flag"
	"fmt"
	"lazygo/command/genapplication"
	"lazygo/command/gencontroller"
	"lazygo/command/gencrud"
	"lazygo/command/gendomain"
	"lazygo/command/genentity"
	"lazygo/command/genenum"
	"lazygo/command/generror"
	"lazygo/command/gengateway"
	"lazygo/command/genrepository"
	"lazygo/command/genservice"
	"lazygo/command/gentest"
	"lazygo/command/genusecase"
	"lazygo/command/genvalueobject"
	"lazygo/command/genvaluestring"
	"lazygo/command/genweb"
	"lazygo/command/genwebapp"
)

var Version = "v1.0.0"

func main() {

	type C struct {
		Command string
		Func    func(...string) error
	}

	commands := make([]C, 0)

	commands = append(commands,
		C{"domain", gendomain.Run},
		C{"entity", genentity.Run},
		C{"valueobject", genvalueobject.Run},
		C{"valuestring", genvaluestring.Run},
		C{"enum", genenum.Run},
		C{"usecase", genusecase.Run},
		C{"repository", genrepository.Run},
		C{"service", genservice.Run},
		C{"test", gentest.Run},
		C{"gateway", gengateway.Run},
		C{"controller", gencontroller.Run},
		C{"error", generror.Run},
		C{"application", genapplication.Run},
		C{"crud", gencrud.Run},
		C{"webapp", genwebapp.Run},
		C{"web", genweb.Run},
	)

	commandMap := map[string]func(...string) error{}

	for _, c := range commands {
		commandMap[c.Command] = c.Func
	}

	flag.Parse()
	cmd := flag.Arg(0)

	if cmd == "" {
		fmt.Printf("LazyGo %s\n", Version)
		fmt.Printf("Try one of this command to learn how to use it\n")
		for _, k := range commands {
			fmt.Printf("  lazygo %s\n", k.Command)
		}
		return
	}

	var values = make([]string, 0)
	if flag.NArg() > 1 {
		values = flag.Args()[1:]
	}

	f, exists := commandMap[cmd]
	if !exists {
		fmt.Printf("Command %s is not recognized\n", cmd)
		return
	}
	err := f(values...)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

}
