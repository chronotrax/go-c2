package msgqueue

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// Msg contains a command to send to an agent.
type Msg struct {
	Command string `json:"command"`
}

// NewMsg is a Msg constructor.
func NewMsg(command string) Msg {
	return Msg{Command: command}
}

// MsgQueue is a message queue abstraction for sending messages to agents.
type MsgQueue interface {
	// Register creates a message queue for the id.
	Register(id uuid.UUID)

	// AddMsg adds a [Msg] to the id's queue.
	// ID must be registered before adding a [Msg].
	AddMsg(id uuid.UUID, msg Msg) error

	// AddMsgAll adds a [Msg] to all id's queues.
	AddMsgAll(msg Msg)

	// GetMsg gets the next message in queue for the id.
	GetMsg(id uuid.UUID) (Msg, error)

	// DeleteMsg deletes the most recent message from id's queue.
	DeleteMsg(id uuid.UUID) error

	// DeleteMsgAll deletes the most recent message from all id's queues.
	DeleteMsgAll()
}

// msgQueue is an internal implementation of MsgQueue
// It contains a [sync.Map] of all agentQueue.
type msgQueue struct {
	qMap sync.Map
}

// NewMsgQueue is a msgQueue constructor.
func NewMsgQueue() MsgQueue {
	return &msgQueue{qMap: sync.Map{}}
}

// agentQueue is a queue of agent commands.
type agentQueue struct {
	mutex sync.Mutex
	msgs  []Msg
}

// newAgentQueue is a agentQueue constructor.
func newAgentQueue() *agentQueue {
	return &agentQueue{
		mutex: sync.Mutex{},
		msgs:  []Msg{},
	}
}

// Error message: "message queue has not been registered yet for id: <id>"
func newNotRegisteredErr(id uuid.UUID) error {
	return fmt.Errorf("message queue has not been registered yet for id: %s", id)
}

// Register creates a message queue for the id.
func (q *msgQueue) Register(id uuid.UUID) {
	q.qMap.Store(id, newAgentQueue())
}

// AddMsg adds a [Msg] to the id's queue.
// ID must be registered before adding a [Msg].
func (q *msgQueue) AddMsg(id uuid.UUID, msg Msg) error {
	value, _ := q.qMap.Load(id)
	agentQ, ok := value.(*agentQueue)
	if !ok {
		return newNotRegisteredErr(id)
	}

	agentQ.mutex.Lock()
	defer agentQ.mutex.Unlock()

	agentQ.msgs = append(agentQ.msgs, msg)
	q.qMap.Store(id, agentQ)
	return nil
}

// AddMsgAll adds a [Msg] to all id's queues.
func (q *msgQueue) AddMsgAll(msg Msg) {
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

// GetMsg gets the next message in queue for the id.
func (q *msgQueue) GetMsg(id uuid.UUID) (Msg, error) {
	value, _ := q.qMap.Load(id)
	agentQ, ok := value.(*agentQueue)
	if !ok {
		return Msg{}, newNotRegisteredErr(id)
	}

	agentQ.mutex.Lock()
	defer agentQ.mutex.Unlock()

	// Agent exists but has no messages
	if len(agentQ.msgs) == 0 {
		return Msg{}, nil
	}

	r := agentQ.msgs[0]
	agentQ.msgs = agentQ.msgs[1:]
	q.qMap.Store(id, agentQ)

	return r, nil
}

// DeleteMsg deletes the most recent message from id's queue.
func (q *msgQueue) DeleteMsg(id uuid.UUID) error {
	value, _ := q.qMap.Load(id)
	agentQ, ok := value.(*agentQueue)
	if !ok {
		return newNotRegisteredErr(id)
	}

	agentQ.mutex.Lock()
	defer agentQ.mutex.Unlock()

	if len(agentQ.msgs) == 0 {
		return fmt.Errorf("message queue empty for id: %s", id)
	}

	agentQ.msgs = agentQ.msgs[:len(agentQ.msgs)-1]
	q.qMap.Store(id, agentQ)

	return nil
}

// DeleteMsgAll deletes the most recent message from all id's queues.
func (q *msgQueue) DeleteMsgAll() {
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
