package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		fmt.Println("origin")
		fmt.Println(origin)
		fmt.Println("origin")
		return true
		//return func() bool {
		//	origin := r.Header.Get("Origin")
		//	return origin == "http://localhost:8080"
		//}()
	},
}

type Handler struct {
	wsHub *WebsocketHub
}

type CreateRoomRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoomResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ClientResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func NewHandler(wsHub *WebsocketHub) *Handler {
	return &Handler{
		wsHub: wsHub,
	}
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var createRoomRequest CreateRoomRequest
	err := c.ShouldBindJSON(&createRoomRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	h.wsHub.Rooms[createRoomRequest.ID] = &Room{
		ID:      createRoomRequest.ID,
		Name:    createRoomRequest.Name,
		Clients: make(map[string]*Client),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    createRoomRequest,
	})
}

func (h *Handler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	roomID := c.Query("roomId")
	clientID := c.Query("userId")
	userName := c.Query("username")

	client := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: userName,
	}

	m := &Message{
		Content:  "A new User has joined the room",
		RoomID:   roomID,
		Username: userName,
	}

	// Register a new client through the register channel
	h.wsHub.Register <- client
	// Broadcast that message
	h.wsHub.Broadcast <- m
	// writeMessage
	go client.writeMessage()
	// readMessage
	client.readMessage(h.wsHub)

}

func (h *Handler) GetRooms(c *gin.Context) {
	rooms := make([]RoomResponse, 0)
	for _, room := range h.wsHub.Rooms {
		rooms = append(rooms, RoomResponse{
			ID:   room.ID,
			Name: room.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    rooms,
	})
}

func (h *Handler) GetClients(c *gin.Context) {
	var clients []ClientResponse
	roomID := c.Param("roomId")

	if _, ok := h.wsHub.Rooms[roomID]; !ok {
		clients = make([]ClientResponse, 0)
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"data":    clients,
		})
	}

	for _, client := range h.wsHub.Rooms[roomID].Clients {
		clients = append(clients, ClientResponse{
			ID:       client.ID,
			Username: client.Username,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    clients,
	})

}

func (h *Handler) LeaveRoom(c *gin.Context) {}
