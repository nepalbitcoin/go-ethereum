// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package discover implements the Node Discovery Protocol.
//
// The Node Discovery protocol provides a way to find RLPx nodes that
// can be connected to. It uses a Kademlia-like protocol to maintain a
// distributed database of the IDs and endpoints of all listening
// nodes.
package discover

import (
	"container/heap"
	"time"
)

const MaxEntries = 10000
const MaxEntriesPerTopic = 50

type Topic uint64

type topicEntry struct {
	topic   Topic
	fifoIdx uint64
	node    *Node
	timer   *time.Timer
}

type topicInfo struct {
	entries            map[uint64]*topicEntry
	fifoHead, fifoTail uint64
	rqItem             *topicRequestQueueItem
}

// removes tail element from the fifo
func (t *topicInfo) getFifoTail() *topicEntry {
	for t.entries[t.fifoTail] == nil {
		t.fifoTail++
	}
	tail := t.entries[t.fifoTail]
	t.fifoTail++
	return tail
}

type nodeInfo struct {
	entries map[Topic]*topicEntry
}

type TopicTable struct {
	nodes   map[*Node]*nodeInfo
	topics  map[Topic]*topicInfo
	globalEntries uint64
	requested     topicRequestQueue
	requestCnt    uint64
	closed        chan struct{} // closed when loop is done
	expired       chan *topicEntry
}

func NewTopicTable() *TopicTable {
	return &TopicTable{
		nodes:  make(map[*Node]*nodeInfo),
		topics: make(map[Topic]*topicInfo),
		closed:       make(chan struct{}),
		expired:      make(chan *topicEntry),
	}
}

func (t *TopicTable) Stop() {
	close(t.closed)
}

func (t *TopicTable) GetEntries(topic Topic) []*Node {
	te := t.topics[topic]
	if te == nil {
		return nil
	}
	nodes := make([]*Node, len(te.entries))
	i := 0
	for _, e := range te.entries {
		nodes[i] = e.node
		i++
	}
	t.requestCnt++
	t.requested.update(te.rqItem, t.requestCnt)
	return nodes
}

func (t *TopicTable) AddEntries(topics []Topic, node *Node, expiry time.Duration) {
	// clear previous entries by the same node
	if n, ok := t.nodes[node]; ok {
		for _, entry := range n.entries {
			delete(t.topics[entry.topic].entries, entry.fifoIdx)
			t.globalEntries--
		}
	} else {
		t.nodes[node] = &nodeInfo{}
	}
	entries := make(map[Topic]*topicEntry)
	t.nodes[node].entries = entries
	for _, topic := range topics {
		te := t.topics[topic]
		if te == nil {
			rqItem := &topicRequestQueueItem{
				topic:    topic,
				priority: t.requestCnt,
			}
			te = &topicInfo{
				entries: make(map[uint64]*topicEntry),
				rqItem:  rqItem,
			}
			t.topics[topic] = te
			heap.Push(&t.requested, rqItem)
		}

		if len(te.entries) == MaxEntriesPerTopic {
			t.deleteEntry(te.getFifoTail())
		}

		if t.globalEntries == MaxEntries {
			t.deleteEntry(t.leastRequested()) // not empty, no need to check for nil
		}

		fifoIdx := te.fifoHead
		te.fifoHead++
		entry := &topicEntry{
			topic:   topic,
			fifoIdx: fifoIdx,
			node:    node,
		}
		te.entries[fifoIdx] = entry
		entries[topic] = entry
		t.globalEntries++
		t.expireEntry(expiry, entry)
	}
}

// removes least requested element from the fifo
func (t *TopicTable) leastRequested() *topicEntry {
	for t.requested.Len() > 0 && t.topics[t.requested[0].topic] == nil {
		heap.Pop(&t.requested)
	}
	if t.requested.Len() == 0 {
		return nil
	}
	return t.topics[t.requested[0].topic].getFifoTail()
}

func (t *TopicTable) deleteEntry(e *topicEntry) {
	ne := t.nodes[e.node].entries
	delete(ne, e.topic)
	if len(ne) == 0 {
		delete(t.nodes, e.node)
	}
	te := t.topics[e.topic]
	delete(te.entries, e.fifoIdx)
	if len(te.entries) == 0 {
		delete(t.topics, e.topic)
		heap.Remove(&t.requested, te.rqItem.index)
	}
	t.globalEntries--
	e.timer.Stop()
}

func (t *TopicTable) DeleteExpiredEntry(d time.Duration, e *topicEntry) {
	if n := t.nodes[e.node]; n != nil {
		if ee := n.entries[e.topic]; ee == e {
			t.deleteEntry(e)
		}
	}
}

func (t *TopicTable) expireEntry(d time.Duration, e *topicEntry) {
	e.timer = time.AfterFunc(d, func() {
		select {
		case t.expired <- e:
		case <-t.closed:
		}
	})
}

type topicRequestQueueItem struct {
	topic    Topic
	priority uint64
	index    int
}

// A topicRequestQueue implements heap.Interface and holds topicRequestQueueItems.
type topicRequestQueue []*topicRequestQueueItem

func (tq topicRequestQueue) Len() int { return len(tq) }

func (tq topicRequestQueue) Less(i, j int) bool {
	return tq[i].priority < tq[j].priority
}

func (tq topicRequestQueue) Swap(i, j int) {
	tq[i], tq[j] = tq[j], tq[i]
	tq[i].index = i
	tq[j].index = j
}

func (tq *topicRequestQueue) Push(x interface{}) {
	n := len(*tq)
	item := x.(*topicRequestQueueItem)
	item.index = n
	*tq = append(*tq, item)
}

func (tq *topicRequestQueue) Pop() interface{} {
	old := *tq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*tq = old[0 : n-1]
	return item
}

func (tq *topicRequestQueue) update(item *topicRequestQueueItem, priority uint64) {
	item.priority = priority
	heap.Fix(tq, item.index)
}