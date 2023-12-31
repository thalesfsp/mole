package tunnel_test

import (
	"log"

	"github.com/thalesfsp/mole/tunnel"
)

// This example shows the basic usage of the package: define both the source and
// destination endpoints, the ssh server and then start the tunnel that will
// exchange data from the local address to the remote address through the
// established ssh channel.
func Example() {
	sourceEndpoints := []string{"127.0.0.1:8080"}
	destinationEndpoints := []string{"user@example.com:80"}

	// Initialize the SSH Server configuration providing all values so
	// tunnel.NewServer will not try to lookup any value using $HOME/.ssh/config
	server, err := tunnel.NewServer("user", "172.17.0.20:2222", "/home/user/.ssh/key", "", "", "/home/user/.ssh/config")
	if err != nil {
		log.Fatalf("error processing server options: %v\n", err)
	}

	t, err := tunnel.New("local", server, sourceEndpoints, destinationEndpoints, "/home/user/.ssh/key")
	if err != nil {
		log.Fatalf("error creating tunnel: %v\n", err)
	}

	// Start the tunnel
	err = t.Start()
	if err != nil {
		log.Fatalf("error starting tunnel: %v\n", err)
	}
}
