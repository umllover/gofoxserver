package console

import (
	"bufio"
	"github.com/lovelly/leaf/conf"
	"github.com/lovelly/leaf/network"
	"math"
	"strconv"
	"strings"
	"github.com/lovelly/leaf/log"
	"os"
)

var server *network.TCPServer

func Init() {
	go run()

	if conf.ConsolePort != 0 {
		server = new(network.TCPServer)
		server.Addr = "localhost:" + strconv.Itoa(conf.ConsolePort)
		server.MaxConnNum = int(math.MaxInt32)
		server.PendingWriteNum = 100
		server.NewAgent = newAgent
		server.Start()
	}
}

func run() {
	for {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Error("console ReadString is error: %v", err)
			continue
		}
		line = strings.TrimSuffix(line[:len(line)-1], "\r")

		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}

		name := args[0]
		c := getCommand(name)
		if c == nil {
			log.Error("command not found, try `help` for help\r\n")
			continue
		}
		output := c.run(args[1:])
		if output != "" {
			log.Release("%v cmd run result: %v", name, output)
		}
	}
}

func Destroy() {
	if server != nil {
		server.Close()
	}
}

type Agent struct {
	conn   *network.TCPConn
	reader *bufio.Reader
}

func newAgent(conn *network.TCPConn) network.Agent {
	a := new(Agent)
	a.conn = conn
	a.reader = bufio.NewReader(conn)
	return a
}

func (a *Agent) Run() {
	for {
		if conf.ConsolePrompt != "" {
			a.conn.Write([]byte(conf.ConsolePrompt))
		}

		line, err := a.reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line[:len(line)-1], "\r")

		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}
		if args[0] == "quit" {
			break
		}

		c := getCommand(args[0])
		if c == nil {
			a.conn.Write([]byte("command not found, try `help` for help\r\n"))
			continue
		}
		output := c.run(args[1:])
		if output != "" {
			a.conn.Write([]byte(output + "\r\n"))
		}
	}
}

func (a *Agent) OnClose() {}
