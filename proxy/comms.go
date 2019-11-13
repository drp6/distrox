package proxy

import (
    "log"
    "time"
    "net"
    "encoding/json"
    "hash/fnv"
)

const (
    UNICAST_MESSAGE = 0 // Message to a single node
    MULTICAST_MESSAGE = 1 // Message to every node
    JOIN_REQUEST_MESSAGE = 2 // Message used for joining the cluster (request)
    JOIN_NOTIFY_MESSAGE = 3 // Message used for noifiying the cluster that a new node has joined
    LEAVE_NOTIFY_MESSAGE = 4 // Message used for notifying the cluster that a node has died
)

type Message struct {
    Timestamp time.Time
    Data []byte
    SenderUrl string
    MessageType int
}

type TCPMessenger struct {
    Listener net.Listener
    RecentMessageHashes map[uint32]time.Time
}

type HTTPRequest struct {

}

type HTTPResponse struct {

}

func InitTCPMessenger(url string) *TCPMessenger {
    rv := &TCPMessenger{}
    rv.RecentMessageHashes = make(map[uint32]time.Time)
    l, err := net.Listen("tcp", url)
    if err != nil {
	log.Fatal(err)
    }
    rv.Listener = l
    return rv
}

func CreateMessage(message []byte, sender_url string, message_type int) Message {
    rv := Message {
	Timestamp: time.Now(),
	SenderUrl: sender_url,
	Data: message,
	MessageType: message_type,
    }
    return rv
}

func MessageToBytes(message Message) []byte {
    b, err := json.Marshal(message)
    if err != nil {
	log.Fatal(err)
    }
    return b
}

func BytesToMessage(bytes []byte) Message {
    rv := Message{}
    json.Unmarshal(bytes, &rv)
    return rv
}

func HashBytes(b []byte) uint32 {
    h := fnv.New32a()
    h.Write(b)
    return h.Sum32()
}

func (m TCPMessenger) PruneStoredMessages() {
    now := time.Now()
    for key := range m.RecentMessageHashes {
	if now.After(m.RecentMessageHashes[key].Add(time.Duration(1.0 * time.Second))) {
	    delete(m.RecentMessageHashes, key)
	}
    }
}

func (m TCPMessenger) HasMessageStored(hash uint32) bool {
    _, ok := m.RecentMessageHashes[hash]
    return ok
}