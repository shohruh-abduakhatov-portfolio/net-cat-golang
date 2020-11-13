package nc

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"
)

type Client struct {
	name    string
	group   *Group
	message chan string
	conn    net.Conn
}

func NewClient(conn net.Conn, name string) *Client {
	client := &Client{
		message: make(chan string),
		conn:    conn,
	}
	return client
}

func (c *Client) ReadInput(msg string) (string, error) {
	c.conn.Write([]byte(msg))
	text, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		log.Printf("Could not read from user: %v, %v",
			c.conn.RemoteAddr().String(), err)
		return "", err
	}
	text = strings.Trim(text, "\r\n")
	return text, nil
}

func (c *Client) Print(msg interface{}) error {
	var err error
	value := reflect.ValueOf(msg)
	switch value.Kind() {
	case reflect.String:
		_, err = c.conn.Write([]byte(value.String()))
		break
	}
	if err != nil {
		log.Printf("Cannot send to user: %v, %v",
			c.conn.RemoteAddr().String(), err)
		return err
	}
	return nil
}

func (c *Client) Receive() {
	for {
		msg := <-c.message
		if msg == quit {
			if c.group != nil {
				c.group.PrintGroup(fmt.Sprintf(leftRoom, c.name))
			}
			break
		}
		log.Printf("recieve: client(%v) recvd msg: %s ", c.conn.RemoteAddr(), strings.Trim(msg, "\r\n"))
		c.Print(msg)
		if c.message == nil {
			break
		}
	}
}

func (c *Client) Send() {
sender:
	for {
		msg, err := c.ReadInput("")
		if err != nil {
			log.Printf("client %v@%s has quit", c.conn.RemoteAddr(), msg)
			c.leaveGroup()
			break sender
		}
		if msg == quit {
			c.Close()
			log.Printf(leftRoom, c.name)
			break sender
		}
		c.HandleRequest(msg)
	}
}

func (c *Client) AddToGroup(g *Group) error {
	clientID := c.conn.RemoteAddr().String()
	if _, ok := g.members[clientID]; ok {
		// already in this group
		log.Printf("already member of this group; user %v@%v room %v", c.name, c.conn.RemoteAddr(), g.name)
		return New(inThisGroup)
	}
	if len(g.members) == groupLimit {
		// max # of clients exceds in this group
		log.Printf("group limit exceeded '%s' %v", g.name, c.conn.RemoteAddr())
		return New(fmt.Sprintf(overGroupLimit, g.name))
	}
	if c.group != nil {
		// already in a group leave first
		log.Printf("already in a group leaving; user %v@%v room %v", c.name, c.conn.RemoteAddr(), c.group.name)
		c.Print(fmt.Sprintf(leaveFirst, c.group.name))
		c.leaveGroup()
	}
	log.Printf("%v@%v joined '%s'", c.name, c.conn.RemoteAddr(), g.name)
	g.members[clientID] = c
	c.group = g
	return nil
}

func (c *Client) RemoveFromGroup(g *Group) {
	clientID := c.conn.RemoteAddr().String()
	if _, ok := g.members[clientID]; !ok {
		// not in this group
		log.Printf("not in this group to leave; user %v@%v room %v", c.name, c.conn.RemoteAddr(), g.name)
		c.Print(notInThisGroup)
		return
	}
	if c.group == nil {
		// not in a group to leave first
		log.Printf("not in a group to leave; user %v@%v room %v", c.name, c.conn.RemoteAddr(), g.name)
		c.Print(notInGroup)
		return
	}
	log.Printf("left: removing user %v from room %v", c.name, g.name)
	delete(g.members, clientID)
	c.group = nil
}

func (c *Client) Close() {
	c.leaveGroup()
	c.conn.Close()
	c.message <- quit
}

func (c *Client) ValidateClient(name string) error {
	if userExists(name) {
		log.Printf("Username exists! %v@%v", name, c.conn.RemoteAddr().String())
		c.Print(fmt.Sprintf(existingUsername, name))
		return New("Username exists!")
	}
	if !isValid(name) {
		log.Printf("Invalid input! %v@%v", name, c.conn.RemoteAddr().String())
		c.Print("Invalid input!")
		return New("Invalid input!")
	}
	return nil
}
