package ws

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"client"`
}

type WebsocketHub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	UnRegister chan *Client
	Broadcast  chan *Message
}

func NewWebsocketHub() *WebsocketHub {
	return &WebsocketHub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *WebsocketHub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if room, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := room.Clients[cl.ID]; !ok {
					room.Clients[cl.ID] = cl
				}
			}
		case cl := <-h.UnRegister:
			if room, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := room.Clients[cl.ID]; ok {
					if len(h.Rooms[cl.RoomID].Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  "user left the chat",
							RoomID:   cl.RoomID,
							Username: cl.Username,
						}
					}
					delete(room.Clients, cl.ID)
					close(cl.Message)
				}
			}
		case m := <-h.Broadcast:
			if room, ok := h.Rooms[m.RoomID]; ok {
				for _, rc := range room.Clients {
					rc.Message <- m
				}
			}
		}
	}
}
