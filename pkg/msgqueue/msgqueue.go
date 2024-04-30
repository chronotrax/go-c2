package msgqueue

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// Message contains a command to send to an agent.
type Message struct {
	MsgID   uuid.UUID `json:"msgID"`
	Command string    `json:"command"`
	Args    []string  `json:"args"`
}

// NewMessage is a Message constructor.
func NewMessage(command string, args ...string) Message {
	id := uuid.New()
	return Message{MsgID: id, Command: command, Args: args}
}

// MessageQueue is a message queue abstraction for sending messages to agents.
type MessageQueue interface {
	// Register creates a message queue for the agentID.
	Register(agentID uuid.UUID)

	// AddMsg adds a [Message] to the agentID's queue.
	// agentID must be registered before adding a [Message].
	AddMsg(agentID uuid.UUID, msg Message) error

	// AddMsgAll adds a [Message] to all agentIDs' queues.
	AddMsgAll(msg Message)

	// GetMsg gets the next message in queue for the agentID.
	GetMsg(agentID uuid.UUID) (Message, error)

	// DeleteMsg deletes the most recent message from agentID's queue.
	DeleteMsg(agentID uuid.UUID) error

	// DeleteMsgAll deletes the most recent message from all agentIDs' queues.
	DeleteMsgAll()
}

// messageQueue is an internal implementation of MessageQueue
// It contains a [sync.Map] of all agentQueue.
type messageQueue struct {
	qMap sync.Map
}

// NewMessageQueue is a messageQueue constructor.
func NewMessageQueue() MessageQueue {
	return &messageQueue{qMap: sync.Map{}}
}

// agentQueue is a queue of agent commands.
type agentQueue struct {
	mutex sync.Mutex
	msgs  []Message
}

// newAgentQueue is a agentQueue constructor.
func newAgentQueue() *agentQueue {
	return &agentQueue{
		mutex: sync.Mutex{},
		msgs:  []Message{},
	}
}

// Error message: "message queue has not been registered yet for id: <id>"
func newNotRegisteredErr(id uuid.UUID) error {
	return fmt.Errorf("message queue has not been registered yet for id: %s", id)
}

// Register creates a message queue for the agentID.
func (q *messageQueue) Register(id uuid.UUID) {
	q.qMap.Store(id, newAgentQueue())
}

// AddMsg adds a [Message] to the agentID's queue.
// agentID must be registered before adding a [Message].
func (q *messageQueue) AddMsg(agentID uuid.UUID, msg Message) error {
	value, _ := q.qMap.Load(agentID)
	agentQ, ok := value.(*agentQueue)
	if !ok {
		return newNotRegisteredErr(agentID)
	}

	agentQ.mutex.Lock()
	defer agentQ.mutex.Unlock()

	agentQ.msgs = append(agentQ.msgs, msg)
	q.qMap.Store(agentID, agentQ)
	return nil
}

// AddMsgAll adds a [Message] to all agentIDs' queues.
func (q *messageQueue) AddMsgAll(msg Message) {
	q.qMap.Range(func(key, value interface{}) bool {
		agentQ, ok := value.(*agentQueue)
		if !ok {
			return false
		}

		agentQ.mutex.Lock()
		defer agentQ.mutex.Unlock()

		agentQ.msgs = append(agentQ.msgs, msg)
		q.qMap.Store(key, agentQ)
		return true
	})
}

// GetMsg gets the next message in queue for the agentID.
func (q *messageQueue) GetMsg(agentID uuid.UUID) (Message, error) {
	value, _ := q.qMap.Load(agentID)
	agentQ, ok := value.(*agentQueue)
	if !ok {
		return Message{}, newNotRegisteredErr(agentID)
	}

	agentQ.mutex.Lock()
	defer agentQ.mutex.Unlock()

	// Agent exists but has no messages
	if len(agentQ.msgs) == 0 {
		return Message{}, nil
	}

	r := agentQ.msgs[0]
	agentQ.msgs = agentQ.msgs[1:]
	q.qMap.Store(agentID, agentQ)

	return r, nil
}

// DeleteMsg deletes the most recent message from agentID's queue.
func (q *messageQueue) DeleteMsg(agentID uuid.UUID) error {
	value, _ := q.qMap.Load(agentID)
	agentQ, ok := value.(*agentQueue)
	if !ok {
		return newNotRegisteredErr(agentID)
	}

	agentQ.mutex.Lock()
	defer agentQ.mutex.Unlock()

	if len(agentQ.msgs) == 0 {
		return fmt.Errorf("message queue empty for id: %s", agentID)
	}

	agentQ.msgs = agentQ.msgs[:len(agentQ.msgs)-1]
	q.qMap.Store(agentID, agentQ)

	return nil
}

// DeleteMsgAll deletes the most recent message from all agentIDs' queues.
func (q *messageQueue) DeleteMsgAll() {
	q.qMap.Range(func(key, value interface{}) bool {
		agentQ, ok := value.(*agentQueue)
		if !ok {
			return false
		}

		agentQ.mutex.Lock()
		defer agentQ.mutex.Unlock()

		agentQ.msgs = agentQ.msgs[:len(agentQ.msgs)-1]
		q.qMap.Store(key, agentQ)
		return true
	})
}
