package main

import "gitlab.itoodev.com/wrkgrp/binlog/system"

func main() {
	// You can think of a booter somewhat as a dependency injection mechanism.
	// It is responsible for starting new things, and providing those things
	// with everything they need to start.
	boot := system.NewSystemBooter()

	// You'll see this a lot. Whenever there is one return value (in go you can have more)
	// you can wrap allocation and evaluation of the return together like below.
	// So below reads like: If boot.Kick() is false then execute what is in the braces.
	if ok := boot.Kick(); !ok {
		// The only place the system may and should panic.
		// If you are here, something is wrong with your design or implementation.
		// Fault tolerance is not about error handling, it's about recovering clean state.
		panic(boot.Inspect())
	}
}
