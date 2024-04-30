package msgqueue

import (
	"reflect"
	"sync"
	"testing"

	"github.com/google/uuid"
)

func Test_msgQueue(t *testing.T) {
	mq := NewMessageQueue()
	id := uuid.New()
	msg1 := NewMessage("test1")
	msg2 := NewMessage("test2")
	msg3 := NewMessage("test3")

	// getting message when id not registered should error
	_, err := mq.GetMsg(id)
	if err == nil {
		t.Error("GetMsg should have returned an error")
	}

	// add some messages
	mq.Register(id)
	err = mq.AddMsg(id, msg1)
	if err != nil {
		t.Error(err)
	}
	err = mq.AddMsg(id, msg2)
	if err != nil {
		t.Error(err)
	}

	// should get msg1
	newMsg, err := mq.GetMsg(id)
	if err != nil {
		t.Errorf("GetMsg() error = %v", err)
	}
	if !reflect.DeepEqual(newMsg, msg1) {
		t.Errorf("GetMsg() got = %v, want %v", newMsg, msg1)
	}

	// add msg3
	err = mq.AddMsg(id, msg3)
	if err != nil {
		t.Error(err)
	}

	// delete most recent (msg3)
	err = mq.DeleteMsg(id)
	if err != nil {
		t.Error(err)
	}

	// should get msg2 then msg3
	newMsg, err = mq.GetMsg(id)
	if err != nil {
		t.Errorf("GetMsg() error = %v", err)
	}
	if !reflect.DeepEqual(newMsg, msg2) {
		t.Errorf("GetMsg() got = %v, want %v", newMsg, msg2)
	}

	// should get an empty message
	newMsg, err = mq.GetMsg(id)
	if err != nil || !reflect.DeepEqual(newMsg, Message{}) {
		t.Error("GetMsg should have returned an error")
	}
}

func Test_msgQueue_DiffAgents(t *testing.T) {
	mq := NewMessageQueue()
	id1 := uuid.New()
	id2 := uuid.New()
	msg1 := NewMessage("test1")
	msg2 := NewMessage("test2")

	// add some messages
	mq.Register(id1)
	mq.Register(id2)
	err := mq.AddMsg(id1, msg1)
	if err != nil {
		t.Error(err)
	}
	err = mq.AddMsg(id2, msg2)
	if err != nil {
		t.Error(err)
	}

	// id1 should get msg1
	newMsg, err := mq.GetMsg(id1)
	if err != nil {
		t.Errorf("GetMsg() error = %v", err)
	}
	if !reflect.DeepEqual(newMsg, msg1) {
		t.Errorf("GetMsg() got = %v, want %v", newMsg, msg1)
	}

	// id2 should get msg2
	newMsg, err = mq.GetMsg(id2)
	if err != nil {
		t.Errorf("GetMsg() error = %v", err)
	}
	if !reflect.DeepEqual(newMsg, msg2) {
		t.Errorf("GetMsg() got = %v, want %v", newMsg, msg2)
	}
}

func Test_msgQueue_AddMsg(t *testing.T) {
	mq := NewMessageQueue()
	wg := &sync.WaitGroup{}
	id := uuid.New()
	msg := NewMessage("test")

	mq.Register(id)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := mq.AddMsg(id, msg)
			if err != nil {
				t.Error(err)
			}

			// remove half the time
			if i%2 == 0 {
				newMsg, err := mq.GetMsg(id)
				if err != nil {
					t.Errorf("GetMsg() error = %v", err)
				}
				if !reflect.DeepEqual(newMsg, msg) {
					t.Errorf("GetMsg() got = %v, want %v", newMsg, msg)
				}
			}
		}()
	}
	wg.Wait()
}

func Test_msgQueue_AddMsgAll(t *testing.T) {
	mq := NewMessageQueue()
	wg := &sync.WaitGroup{}
	id1 := uuid.New()
	id2 := uuid.New()
	msg := NewMessage("test")

	// create 2 agents
	mq.Register(id1)
	mq.Register(id2)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mq.AddMsgAll(msg)

			// remove half the time
			if i%2 == 0 {
				newMsg, err := mq.GetMsg(id1)
				if err != nil {
					t.Errorf("GetMsg() error = %v", err)
				}
				if !reflect.DeepEqual(newMsg, msg) {
					t.Errorf("GetMsg() got = %v, want %v", newMsg, msg)
				}
			}
		}()
	}
	wg.Wait()
}

func Test_msgQueue_DeleteMsg(t *testing.T) {
	mq := NewMessageQueue()
	id := uuid.New()
	msg1 := NewMessage("test1")
	msg2 := NewMessage("test2")

	mq.Register(id)

	// should get error deleting from empty queue
	err := mq.DeleteMsg(id)
	if err == nil {
		t.Error("DeleteMsg should have returned an error")
	}

	// add msg1
	err = mq.AddMsg(id, msg1)
	if err != nil {
		t.Error(err)
	}

	// add msg2
	err = mq.AddMsg(id, msg2)
	if err != nil {
		t.Error(err)
	}

	// delete most recent message (msg2)
	err = mq.DeleteMsg(id)
	if err != nil {
		t.Error(err)
	}

	// should get msg1
	newMsg, err := mq.GetMsg(id)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(newMsg, msg1) {
		t.Errorf("GetMsg() got = %v, want %v", newMsg, msg1)
	}

	// should get error deleting from empty queue
	err = mq.DeleteMsg(id)
	if err == nil {
		t.Error("DeleteMsg should have returned an error")
	}
}

func Test_msgQueue_DeleteMsgAll(t *testing.T) {
	mq := NewMessageQueue()
	id1 := uuid.New()
	id2 := uuid.New()
	msg1 := NewMessage("test1")
	msg2 := NewMessage("test2")

	mq.Register(id1)
	mq.Register(id2)

	// should get error deleting from empty queue
	err := mq.DeleteMsg(id1)
	if err == nil {
		t.Error("DeleteMsg should have returned an error")
	}

	// add msg1
	err = mq.AddMsg(id1, msg1)
	if err != nil {
		t.Error(err)
	}
	err = mq.AddMsg(id2, msg1)
	if err != nil {
		t.Error(err)
	}

	// add msg2
	err = mq.AddMsg(id1, msg2)
	if err != nil {
		t.Error(err)
	}
	err = mq.AddMsg(id2, msg2)
	if err != nil {
		t.Error(err)
	}

	// delete most recent message (msg2)
	err = mq.DeleteMsg(id1)
	if err != nil {
		t.Error(err)
	}
	err = mq.DeleteMsg(id2)
	if err != nil {
		t.Error(err)
	}

	// should get msg1
	newMsg, err := mq.GetMsg(id1)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(newMsg, msg1) {
		t.Errorf("GetMsg() got = %v, want %v", newMsg, msg1)
	}
	newMsg, err = mq.GetMsg(id2)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(newMsg, msg1) {
		t.Errorf("GetMsg() got = %v, want %v", newMsg, msg1)
	}
}
