package server

import (
	"fmt"
	"net"

	"github.com/Olatokumbo/blueis/internal/redis"
	"github.com/Olatokumbo/blueis/internal/resp"
)

func handleConnections(c net.Conn) error {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	defer c.Close()

	parser := resp.NewParser(c)
	handler := redis.NewCommandHandler()

	for {
		value, err := parser.Parse()
		if err != nil {
			return fmt.Errorf("error parsing command: %v", err)
		}

		if value.Type != resp.Array {
			return fmt.Errorf("expected array type, got %c", value.Type)
		}

		args := value.Data.([]*resp.Value)
		cmdArgs := make([]string, len(args))
		for i, arg := range args {
			if arg.Type != resp.BulkString {
				return fmt.Errorf("expected bulk string type for argument, got %c", arg.Type)
			}
			if arg.Data == nil {
				return fmt.Errorf("nil argument")
			}
			cmdArgs[i] = arg.Data.(string)
		}

		result, err := handler.HandleCommand(cmdArgs)
		if err != nil {
			errorValue := &resp.Value{
				Type: resp.Error,
				Data: err.Error(),
			}
			if err := resp.WriteValue(c, errorValue); err != nil {
				return fmt.Errorf("error writing error response: %v", err)
			}
			continue
		}

		var response *resp.Value
		switch v := result.(type) {
		case string:
			response = &resp.Value{
				Type: resp.BulkString,
				Data: v,
			}
		case int:
			response = &resp.Value{
				Type: resp.Integer,
				Data: int64(v),
			}
		case nil:
			response = &resp.Value{
				Type: resp.BulkString,
				Data: nil,
			}
		default:
			return fmt.Errorf("unsupported response type: %T", v)
		}

		if err := resp.WriteValue(c, response); err != nil {
			return fmt.Errorf("error writing response: %v", err)
		}
	}
}

func StartServer(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("error starting server: %s", err)
	}

	defer listener.Close()

	fmt.Printf("Listening at PORT: %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting: %s", err)
		}

		fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr().String())
		go handleConnections(conn)
	}
}
