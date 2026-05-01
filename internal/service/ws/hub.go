package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	userType "github.com/umesshk/termi-chatt/internal/user"
	"github.com/umesshk/termi-chatt/internal/redisx"
	redis "github.com/redis/go-redis/v9"
)

type Hub struct {
	mu sync.RWMutex

	roomMap map[int][]userType.User
	connMap map[*websocket.Conn]int

	redis *redisx.Client

	// roomId -> pubsub
	subs map[int]*redis.PubSub
}

func NewHub(r *redisx.Client) *Hub {
	return &Hub{
		roomMap: make(map[int][]userType.User),
		connMap: make(map[*websocket.Conn]int),
		redis:   r,
		subs:    make(map[int]*redis.PubSub),
	}
}

func (h *Hub) RoomExists(ctx context.Context, roomID int) (bool, error) {
	// Without Redis, we can only validate rooms that exist in-memory.
	if h.redis == nil {
		h.mu.RLock()
		_, ok := h.roomMap[roomID]
		h.mu.RUnlock()
		return ok, nil
	}

	key := fmt.Sprintf("room:%d", roomID)
	n, err := h.redis.Rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

func (h *Hub) MarkRoomExists(ctx context.Context, roomID int) error {
	if h.redis == nil {
		return nil
	}
	key := fmt.Sprintf("room:%d", roomID)
	return h.redis.Rdb.Set(ctx, key, "1", 7*24*time.Hour).Err()
}

func (h *Hub) AddConn(roomID int, u userType.User) {
	h.mu.Lock()
	h.roomMap[roomID] = append(h.roomMap[roomID], u)
	h.connMap[u.User_conn] = roomID
	h.mu.Unlock()
}

func (h *Hub) RemoveConn(conn *websocket.Conn) {
	h.mu.Lock()
	roomID, ok := h.connMap[conn]
	if !ok {
		h.mu.Unlock()
		return
	}

	users := h.roomMap[roomID]
	idx := -1
	for i, user := range users {
		if user.User_conn == conn {
			idx = i
			break
		}
	}
	if idx != -1 {
		users = append(users[:idx], users[idx+1:]...)
	}
	h.roomMap[roomID] = users
	delete(h.connMap, conn)
	h.mu.Unlock()
}

func (h *Hub) JoinedRoomID(conn *websocket.Conn) (int, bool) {
	h.mu.RLock()
	roomID, ok := h.connMap[conn]
	h.mu.RUnlock()
	return roomID, ok
}

func (h *Hub) RoomUsers(roomID int) []userType.User {
	h.mu.RLock()
	users := h.roomMap[roomID]
	// Copy to avoid holding lock while writing.
	out := make([]userType.User, len(users))
	copy(out, users)
	h.mu.RUnlock()
	return out
}

func (h *Hub) EnsureRoomSub(roomID int) {
	if h.redis == nil {
		return
	}

	h.mu.Lock()
	if _, ok := h.subs[roomID]; ok {
		h.mu.Unlock()
		return
	}
	channel := fmt.Sprintf("room:%d:pubsub", roomID)
	sub := h.redis.Rdb.Subscribe(context.Background(), channel)
	h.subs[roomID] = sub
	h.mu.Unlock()

	go func() {
		for msg := range sub.Channel() {
			var resp userType.ServerResponse
			if err := json.Unmarshal([]byte(msg.Payload), &resp); err != nil {
				log.Println("redis pubsub unmarshal:", err)
				continue
			}
			// Broadcast to local connections in that room.
			for _, u := range h.RoomUsers(roomID) {
				_ = u.User_conn.WriteJSON(resp)
			}
		}
	}()
}

func (h *Hub) Publish(roomID int, resp userType.ServerResponse) {
	// Always broadcast locally.
	for _, u := range h.RoomUsers(roomID) {
		_ = u.User_conn.WriteJSON(resp)
	}

	if h.redis == nil {
		return
	}

	h.EnsureRoomSub(roomID)

	channel := fmt.Sprintf("room:%d:pubsub", roomID)
	b, err := json.Marshal(resp)
	if err != nil {
		return
	}
	if err := h.redis.Rdb.Publish(context.Background(), channel, string(b)).Err(); err != nil {
		// Non-fatal; local broadcast already happened.
		log.Println("redis publish:", err)
	}
}

