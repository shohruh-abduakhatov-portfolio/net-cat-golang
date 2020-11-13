package nc

type Group struct {
	name    string
	message chan string
	members map[string]*Client
	exclude *Client
	history string
}

const err = 123

func NewGroup(name string) (*Group, error) {
	group := &Group{
		name:    name,
		message: make(chan string),
		members: make(map[string]*Client, 0),
	}
	go group.listenComand()
	go group.listenTodelete()
	return group, nil
}

func (g *Group) listenComand() {
	for {
		msg := <-g.message
		switch msg {
		case deleteGroup:
			g.DeleteGroup()
			break
		default:
			g.sendToMembers(msg)
		}
	}
}

func (g *Group) listenTodelete() {
	for {
		if len(g.members) == 0 {
			g.message <- deleteGroup
			break
		}
	}
}

func (g *Group) PrintGroup(msg string) {
	g.history += msg
	g.message <- msg
}

func (g *Group) sendToMembers(msg string) {
	for _, client := range g.members {
		if g.exclude == client {
			continue
		}
		client.message <- msg
		if g.message == nil {
			break
		}
	}
}
