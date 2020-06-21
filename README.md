# gophertunnel
A Minecraft library containing packages to create clients, servers, proxies and other tools, and a proxy implementation using them.

[Module Documentation](https://pkg.go.dev/mod/github.com/sandertv/gophertunnel)

![telescope gopher](https://github.com/Sandertv/gophertunnel/blob/master/gophertunnel_telescope_coloured.png)

## Overview
gophertunnel is composed of several packages that may be of use for creating Minecraft related tools.

package [query](https://pkg.go.dev/github.com/sandertv/gophertunnel/query?tab=doc): A package implementing the sending of queries
to servers that implement the UT3/GameSpy Query Protocol.

package [minecraft](https://pkg.go.dev/github.com/sandertv/gophertunnel/minecraft?tab=doc): A package implementing connecting
to Minecraft Bedrock Edition servers and listening for Minecraft Bedrock Edition clients using a TCP style interface.

* package [minecraft/auth](https://pkg.go.dev/github.com/sandertv/gophertunnel/minecraft/auth?tab=doc): A package implementing
Microsoft, XBOX Live and Minecraft account authentication.

* package [minecraft/nbt](https://pkg.go.dev/github.com/sandertv/gophertunnel/minecraft/nbt?tab=doc): A package implementing the
Minecraft NBT format. Three variants of the format are implemented: The Java Edition variant (Big Endian) and
the Bedrock Edition variants (Little Endian, both with and without varints)

* package [minecraft/protocol](https://pkg.go.dev/github.com/sandertv/gophertunnel/minecraft/protocol?tab=doc): A package
implementing the reading, writing and handling of packets found in the Minecraft Bedrock Edition protocol.

* package [minecraft/resource](https://pkg.go.dev/github.com/sandertv/gophertunnel/minecraft/resource?tab=doc): A package handling
the reading and compiling of Minecraft resource packs.

* package [minecraft/text](https://pkg.go.dev/github.com/sandertv/gophertunnel/minecraft/text?tab=doc): A package containing utility
functions related to Minecraft text formatting.

## Examples
Creating a Minecraft client that authenticates using an XBOX Live account and connects to a server:
```go
package main

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func main() {
    conn, err := minecraft.Dialer{
        Email: "some@email.address",
        Password: "password",
    }.Dial("raknet", "mco.mineplex.com:19132")
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    for {
    	pk, err := conn.ReadPacket()
    	if err != nil {
    		break
    	}
    	// Handle the incoming packet.
    	_ = pk
    	
    	// Send a packet to the server.
    	if err := conn.WritePacket(&packet.RequestChunkRadius{ChunkRadius: 32}); err != nil {
    		break
    	}
    }
}
```

Creating a Minecraft listener that can accept incoming clients and adapts the MOTD from another server:
```go
package main

import (
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft"
)

func main() {
	listener, err := minecraft.Listen("raknet", "0.0.0.0:19132")
	if err != nil {
		panic(err)
	}
	_ = listener.HijackPong("mco.mineplex.com:19132")

	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		go func() {
			// Process the connection on another goroutine as you would with TCP connections.
			defer conn.Close()
			for {
				// Read a packet from the client.
				if _, err := conn.(*minecraft.Conn).ReadPacket(); err != nil {
					return
				}
			}
		}()
	}
}
```

## Versions
Gophertunnel supports only one version at a time. Generally, a new minor version is tagged when gophertunnel
supports a new Minecraft version that was not previously supported. A list of the recommended gophertunnel
versions for past Minecraft versions is listed below.

| Version | Tag     |
|---------|---------|
| 1.16.0  | Latest  |
| 1.14.60 | v1.6.5  |
| 1.14.0  | v1.3.20 |
| 1.13.0  | v1.3.5  |
| 1.12.0  | v1.2.11 |

## Proxy
A MITM proxy program is implemented in the main.go file. It uses the gophertunnel libraries to create a proxy
that provides user authentication and proxying a connection to another server.