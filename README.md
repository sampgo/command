# command

Command handler library for [sampgo](https://github.com/sampgo/sampgo)

## API

### `func NewCommand(command Command) (cmd *Command)`

Creates a new command with passed `Command` struct, and returns a pointer to it

Example:

``` go
cmd := command.NewCommand(Command{Name: "foo", Prefix: "!"})
```

### `func New() (cmd *Command)`

Creates a new command with empty `Command` struct, and returns a pointer to it

Example:

``` go
cmd := command.New()
```

### `func (cmd *Command) SetName(name string) *Command`

Sets the name for `cmd`, and returns a pointer to it

Example:

``` go
cmd := command.New()
cmd.SetName("foo")
```

### `func (cmd *Command) SetAlias(aliases ...string) *Command`

Adds aliases for `cmd`, and returns a pointer to it

Example:

``` go
cmd := command.New()
cmd.SetName("foo")
    .SetAlias("bar", "baz")
```

### `func (cmd *Command) SetPrefix(prefix string) *Command`

Sets the prefix for `cmd`, and returns a pointer to it

Example:

``` go
cmd := command.New()
cmd.SetName("foo")
    .SetAlias("bar", "baz")
    .SetPrefix("!")
```

### `func (cmd *Command) Handle(fn Func) (err error)`

Registers callback for `cmd`, and returns an error if any

Example:

``` go
cmd := command.New()
cmd.SetName("foo")
    .SetAlias("bar", "baz")
    .SetPrefix("!")
err := cmd.Handle(func(ctx command.Context) (err error) {
    fmt.Println(ctx.Player.ID)
    return
})
// handle err
```

### `func (cmd *Command) Sync() (err error)`

Synchronises command data to internal command handler, and returns an error if any

Example:

``` go
cmd := command.New()
cmd.SetName("foo")
    .SetAlias("bar", "baz")
    .SetPrefix("!")
err := cmd.Handle(func(ctx command.Context) (err error) {
    fmt.Println(ctx.Player.ID)
    return
})
// handle err

cmd.Name = "fofo"
err = cmd.Sync()
// handle err
```

### `func SetGeneralCommandErrorFunc(fn ErrorFunc)`

Registers `fn` as the callback that's fired when function passed to `command.Handle` returns an error

Example:

``` go
command.SetGeneralCommandErrorFunc(func(ctx command.ErrorContext) bool {
    // handle ctx.Error
    return true
})

cmd := command.New()
cmd.SetName("foo")
    .SetAlias("bar", "baz")
    .SetPrefix("!")
err := cmd.Handle(func(ctx command.Context) (err error) {
    fmt.Println(ctx.Player.ID)
    return
})
// handle err

cmd.Name = "fofo"
err = cmd.Sync()
// handle err
```

### `func SetGeneralCommandBeforeFunc(fn BeforeFunc)`

Registers `fn` as the callback that's fired before execution of *any* command callback

Example:

``` go
command.SetGeneralCommandBeforeFunc(func(ctx command.Context) bool {
    fmt.Println("before we execute command below:")
    fmt.Println(ctx.Command.Name)
    return true
})

cmd := command.New()
cmd.SetName("foo")
    .SetAlias("bar", "baz")
    .SetPrefix("!")
err := cmd.Handle(func(ctx command.Context) (err error) {
    fmt.Println(ctx.Player.ID)
    return
})
// handle err

cmd.Name = "fofo"
err = cmd.Sync()
// handle err
```

### `func SetGeneralCommandAfterFunc(fn AfterFunc)`

Registers `fn` as the callback that's fired after execution of *any* command callback

Example:

``` go
command.SetGeneralCommandAfterFunc(func(ctx command.Context) bool {
    fmt.Println("after we execute command below:")
    fmt.Println(ctx.Command.Name)
    return true
})

cmd := command.New()
cmd.SetName("foo")
    .SetAlias("bar", "baz")
    .SetPrefix("!")
err := cmd.Handle(func(ctx command.Context) (err error) {
    fmt.Println(ctx.Player.ID)
    return
})
// handle err

cmd.Name = "fofo"
err = cmd.Sync()
// handle err
```
