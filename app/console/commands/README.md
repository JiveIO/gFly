# Command `app/console/commands/`

Command pool only accept instance which implemented interface `ICommand`

```go
// ICommand The interface command.
type ICommand interface {
    Validate(parameters CommandParameter) error // Validate input parameter from console
    Handle() // Method to handle the command
}
```

Let example you want to create a command name `EchoCommand`. A struct `EchoCommand` will be created. And the command will receive parameter `msg` from console. 

You should embed `console.Command` inside `EchoCommand` to get some default methods were implemented.


```go
// EchoCommand struct for echo command.
type EchoCommand struct{
    console.Command
    Msg     string
}
```

Some case you only create command without parameter to be validated. So, embed parent command `console.Command` will help you reduce creating method `Validate(parameters CommandParameter) error` for your command `EchoCommand`. 

Back to this example. So, we need validate input parameters.

```go
// Validate Check required parameters.
func (c *EchoCommand) Validate(parameters console.CommandParameter) error {
    // Get parameter
    data, ok := parameters["msg"]
    if !ok {
        return errors.New("invalid msg parameter")
    }

    // Get parameters
    c.Msg = data.(string)

    return nil
}
```
Final method for `EchoCommand` to handle data

```go
// Handle Process command.
func (c *EchoCommand) Handle() {
    log.Infof("Hello %s", c.Msg)

    log.Infof("EchoCommand :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}
```

The latest important thing is register command

```go
// Auto-register command.
func init() {
    console.RegisterCommand(&EchoCommand{}, "echo")
}
```

## Job file `echo_command.go`

```go
package commands

import (
    "app/core/console"
    "app/core/log"
    "errors"
    "time"
)

// ---------------------------------------------------------------
//                      Register command.
// ./build/artisan cmd:run echo --msg="Gopher"
// ---------------------------------------------------------------

// Auto-register command.
func init() {
    console.RegisterCommand(&EchoCommand{}, "echo")
}

// ---------------------------------------------------------------
//                      EchoCommand struct.
// ---------------------------------------------------------------

// EchoCommand struct for echo command.
type EchoCommand struct {
    console.Command
    Msg string
}

// Validate Check required parameters.
func (c *EchoCommand) Validate(parameters console.CommandParameter) error {
    // Get parameter
    data, ok := parameters["msg"]
    if !ok {
        return errors.New("invalid msg parameter")
    }

    // Get parameters
    c.Msg = data.(string)

    return nil
}

// Handle Process command.
func (c *EchoCommand) Handle() {
    log.Infof("Hello %s", c.Msg)

    log.Infof("EchoCommand :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}
```

**Note:** Put the job `echo_command.go` correct folder `app/console/commands/`

## Check `EchoCommand` command

Run command

```bash
./build/artisan cmd:run echo --msg="Gopher"
```

## Auto command name

When you change register command without second parameter. The system will auto use the name exactly the name of command struct. In this case will be `EchoCommand`

```go
// Auto-register command.
func init() {
    console.RegisterCommand(&EchoCommand{})
}
```

Run command
```bash
./build/artisan cmd:run EchoCommand --msg="Gopher"
```