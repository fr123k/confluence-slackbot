package utility

import (
	"fmt"
)

//HelloWorld prints Hello World to standard output.
func HelloWorld() {
	message := Hello("World")
	fmt.Printf("%s\n", message)
}

//Hello returns a string `Hello ` + the staring passed in the name string.
func Hello(name string) (string) {
	return fmt.Sprintf("%s %s", "Hello", name)
}
