package nc

import (
	"fmt"
	"log"
	"net"
	"time"
)

var groupsList = map[string]*Group{}

func Start(ip, port string) error {
	log.Printf(startupMessage, port)
	fmt.Printf(startupMessage, port)
	listener, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client %s: ", err)
			return err
		}
		go createNewClient(conn)
	}
}

func catchPanic() {
	if a := recover(); a != nil {
		log.Printf("%v", a)
	}
}

func createNewClient(conn net.Conn) {
	log.Printf("New connection from %v", conn.RemoteAddr().String())
	client := &Client{
		message: make(chan string),
		conn:    conn,
	}
	name, err := client.ReadInput(fmt.Sprint(linuxIcon, yourName))
	if err != nil {
		log.Printf("Error creating client %v@%v", name, conn.RemoteAddr().String())
		return
	}
	if err := client.ValidateClient(name); err != nil {
		return
	}
	client.name = name
	log.Printf("New client: '%v' @%v", client.name, conn.RemoteAddr().String())

	go client.Send()
	go client.Receive()
}

func (c *Client) HandleRequest(command string) {
	log.Printf("%v command from %v@%v", command, c.name, c.conn.RemoteAddr().String())
	switch command {
	case groups:
		c.allGroups()
		break
	case join:
		c.joinGroup()
		break
	case create:
		c.createNewGroup()
		break
	case leave:
		c.leaveGroup()
		break
	case changeName:
		c.rename()
		break
	default:
		if c.group != nil && command != "" {
			c.group.exclude = c
			send := fmt.Sprintf("[%v][%v]:%v\n", time.Now().Format(timeFormat), c.name, command)
			c.group.PrintGroup(send)
		}
	}
}

func (c *Client) createNewGroup() {
	name, err := c.ReadInput(groupName)
	if err != nil {
		log.Printf("error reading group input %v - %v", c.conn.RemoteAddr(), err)
		return
	}
	if !isValid(name) {
		log.Printf("Empty name")
		c.Print(fmt.Sprintf(invalidGroupName, name))
		return
	}
	if _, ok := groupsList[name]; ok {
		log.Printf("Existing group" + name)
		c.Print(fmt.Sprintf(groupExists, name))
		return
	}
	group, err := NewGroup(name)
	if err != nil {
		log.Printf("error reading group input %v - %v", c.conn.RemoteAddr(), err)
		c.Print(internalError)
		return
	}
	c.AddToGroup(group)
	groupsList[group.name] = group
	c.group.PrintGroup(fmt.Sprintf(joinedGroup, c.name))
}

func (c *Client) rename() {
	name, err := c.ReadInput(yourName)
	if err != nil {
		log.Printf("Error renaming client %v@%v", c.name, c.conn.RemoteAddr().String())
		return
	}
	if err := c.ValidateClient(name); err != nil {
		return
	}
	oldName := c.name
	c.name = name
	if c.group != nil {
		c.group.PrintGroup(fmt.Sprintf(changedName, oldName, c.name))

	}
}

func (c *Client) joinGroup() {
	name, err := c.ReadInput(joinGroup)
	if err != nil {
		log.Printf("error reading group input %v - %v", c.conn.RemoteAddr(), err)
		return
	}
	group, ok := groupsList[name]
	if !ok {
		// no such group
		log.Printf("no such group '%s' %v", name, c.conn.RemoteAddr())
		c.Print(fmt.Sprintf(noSuchGroup, name))
		return
	}
	if err := c.AddToGroup(group); err != nil {
		//
		c.Print(err)
		return
	}
	c.Print(c.group.history)
	c.group.PrintGroup(fmt.Sprintf(joinedGroup, c.name))
}

func (c *Client) leaveGroup() {
	if c.group == nil {
		return
	}
	group, ok := groupsList[c.group.name]
	if ok {
		c.RemoveFromGroup(group)
		group.PrintGroup(fmt.Sprintf(leftRoom, c.name))
		c.Print(fmt.Sprintf(leavingGroup, group.name))

	}
}

func (c *Client) allGroups() {
	for name, _ := range groupsList {
		c.conn.Write([]byte(name + "\n"))
	}

}

func (g *Group) DeleteGroup() {
	delete(groupsList, g.name)
}

func isValid(name string) bool {
	if name == "" {
		return false
	}
	return true
}

func userExists(name string) bool {
	for _, g := range groupsList {
		for _, m := range g.members {
			if m.name == name {
				return true
			}
		}
	}
	return false
}
