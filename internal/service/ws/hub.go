package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	redis "github.com/redis/go-redis/v9"
	"github.com/umesshk/termi-chatt/internal/redisx"
	userType "github.com/umesshk/termi-chatt/internal/user"
)

type roomMember struct {
	userID   int
	username string
	client   *Client
}

type Hub struct {
	mu sync.RWMutex

	roomMap map[int][]roomMember
	connMap map[*Client]int

	redis *redisx.Client

	subs map[int]*redis.PubSub
}

func NewHub(r *redisx.Client) *Hub {
	return &Hub{
		roomMap: make(map[int][]roomMember),
		connMap: make(map[*Client]int),
		redis:   r,
		subs:    make(map[int]*redis.PubSub),
	}
}

func (h *Hub) RoomExists(ctx context.Context, roomID int) (bool, error) {
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

func (h *Hub) AddClient(roomID int, userID int, username string, client *Client) {
	h.mu.Lock()
	h.roomMap[roomID] = append(h.roomMap[roomID], roomMember{
		userID:   userID,
		username: username,
		client:   client,
	})
	h.connMap[client] = roomID
	h.mu.Unlock()
}

func (h *Hub) RemoveClient(client *Client) {
	h.mu.Lock()
	roomID, ok := h.connMap[client]
	if !ok {
		h.mu.Unlock()
		return
	}

	members := h.roomMap[roomID]
	idx := -1
	for i, m := range members {
		if m.client == client {
			idx = i
			break
		}
	}
	if idx != -1 {
		members = append(members[:idx], members[idx+1:]...)
	}
	h.roomMap[roomID] = members
	delete(h.connMap, client)
	h.mu.Unlock()

	close(client.Send)
}

func (h *Hub) JoinedRoomID(client *Client) (int, bool) {
	h.mu.RLock()
	roomID, ok := h.connMap[client]
	h.mu.RUnlock()
	return roomID, ok
}

func (h *Hub) RoomUsers(roomID int) []userType.User {
	h.mu.RLock()
	members := h.roomMap[roomID]
	out := make([]userType.User, len(members))
	for i, m := range members {
		out[i] = userType.User{UserId: m.userID, Username: m.username}
	}
	h.mu.RUnlock()
	return out
}

func (h *Hub) roomClients(roomID int) []*Client {
	h.mu.RLock()
	members := h.roomMap[roomID]
	out := make([]*Client, len(members))
	for i, m := range members {
		out[i] = m.client
	}
	h.mu.RUnlock()
	return out
}

func (h *Hub) BroadcastToRoom(roomID int, resp userType.ServerResponse) {
	for _, client := range h.roomClients(roomID) {
		client.Enqueue(resp)
	}
}

func (h *Hub) broadcastToRoom(roomID int, resp userType.ServerResponse) {
	h.BroadcastToRoom(roomID, resp)
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
			h.broadcastToRoom(roomID, resp)
		}
	}()
}

func (h *Hub) Publish(roomID int, resp userType.ServerResponse) {
	if h.redis != nil {
		h.EnsureRoomSub(roomID)

		channel := fmt.Sprintf("room:%d:pubsub", roomID)
		b, err := json.Marshal(resp)
		if err != nil {
			return
		}
		if err := h.redis.Rdb.Publish(context.Background(), channel, string(b)).Err(); err != nil {
			log.Println("redis publish:", err)
		}
		return
	}

	h.broadcastToRoom(roomID, resp)
}
