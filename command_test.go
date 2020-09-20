package command

import (
	"errors"
	"fmt"
	"testing"

	"github.com/sampgo/sampgo"
)

func TestEntryPoint(t *testing.T) {
	cmd := NewCommand(Command{
		Name:   "f",
		Prefix: "/",
	})

	cmd.SetName("foo")

	cmd.SetAlias("bar", "baz")

	err := cmd.Handle(func(ctx Context) (err error) {
		fmt.Println("command fired!")
		fmt.Println(ctx)
		return
	})

	if err != nil {
		t.Error(err)
	}

	SetGeneralCommandBeforeFunc(func(ctx Context) (err error) {
		fmt.Println("before commadn gets fired!")
		fmt.Println(ctx)
		return
	})

	SetGeneralCommandAfterFunc(func(ctx Context) (err error) {
		fmt.Println("after command gets fired!")
		fmt.Println(ctx)
		fmt.Println()
		return
	})

	fail := errors.New("handler call failed")

	success := handler(sampgo.Player{ID: 0}, "/foo first handler call")
	if !success {
		t.Error(fail)
	}

	cmd.Prefix = "$"
	err = cmd.Sync()
	if err != nil {
		t.Error(err)
	}

	success = handler(sampgo.Player{ID: 1}, "$bar second handler call")
	if !success {
		t.Error(fail)
	}

	cmd.SetPrefix("barbazbax")
	err = cmd.Sync()
	if err != nil {
		t.Error(err)
	}

	success = handler(sampgo.Player{ID: 2}, "barbazbaxbaz third handler call")
	if !success {
		t.Error(fail)
	}

	cmd.SetAlias("foobar!")
	err = cmd.Sync()
	if err != nil {
		t.Error(err)
	}

	success = handler(sampgo.Player{ID: 3}, "barbazbaxfoobar! fourth handler call")
	if !success {
		t.Error(fail)
	}
}
