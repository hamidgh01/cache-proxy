package main

import "fmt"

func main() {
	fmt.Println("Developing Cache Proxy Server...")

	// needed components (primary):
	// 1. a Network Interface to handle incoming requests & outgoing responses
	// 2. a Cache Storage system (Redis integrated) + data serialization to store in Redis
	// 3. some rules, filtering and policies for caching
	// 4. a CLI parser for command line args + a Configuration struct
	// 5. a Logger (better monitoring, debugging, recording events, etc.)
	// etc... (new incoming ideas during development process)
}
