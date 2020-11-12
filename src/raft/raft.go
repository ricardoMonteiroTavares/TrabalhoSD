package raft

import (
	"labrpc"
	"math/rand"
	"sync"
	"time"
)

type ApplyMsg struct {
	Index int
	Command interface{}
	UseSnapshot bool
	Snapshot []byte
}

type Raft struct {
	mu sync.Mutex
	peers []*labrpc.ClientEnd
	persister *Persister
	me int
	currentTerm int
	votedFor int
	log []LogEntry
	commitIndex int
	lastApplied int
	nextIndex []int
	matchIndex []int
	state int
	chanGrantVote chan bool
	voteCount int
	chanLeader chan bool
	chanHeartbeat chan bool
}

type LogEntry struct {
	LogIndex int
	LogTerm int
}

const (
	STATE_LEADER = iota
	STATE_CANDIDATE
	STATE_FOLLOWER
	HEARTBEATS = 100 * time.Millisecond
)

func (rf *Raft) getLastTerm() int {
	return rf.log[len(rf.log)-1].LogTerm
}

func (rf *Raft) GetState() (int, bool) {
	var term int
	var isleader bool
	term = rf.currentTerm
	isleader = (rf.state == STATE_LEADER)
	return term, isleader
}

func (rf *Raft) persist() {
}

func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 {
		return
	}
}

type RequestVoteArgs struct {
	Term int
	CandidateId int
	LastLogIndex int
	LastLogTerm int
}

type RequestVoteReply struct {
	Term int
	VoteGranted bool
}

type AppendEntriesArgs struct {
	Term int
	LeaderId int
	PrevLogIndex int
	PrevLogTerm int
	Entries []LogEntry
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term int
	Success bool
}

func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	reply.VoteGranted = false
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		return
	}
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.state = STATE_FOLLOWER
		rf.votedFor = -1
	}
	reply.Term = rf.currentTerm
	upToDate := false
	if args.LastLogTerm > rf.getLastTerm() || args.LastLogTerm == rf.getLastTerm() {
		upToDate = true
	}

	if upToDate && (rf.votedFor == -1 || rf.votedFor == args.CandidateId) {
		rf.chanGrantVote <- true
		rf.state = STATE_FOLLOWER
		reply.VoteGranted = true
		rf.votedFor = args.CandidateId
	}
}

func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if ok {
		if rf.state != STATE_CANDIDATE || args.Term != rf.currentTerm {
			return ok
		}

		if reply.Term > rf.currentTerm {
			rf.currentTerm = reply.Term
			rf.state = STATE_FOLLOWER
			rf.votedFor = -1
			rf.persist()
		}
		if reply.VoteGranted {
			rf.voteCount++
			if rf.state == STATE_CANDIDATE && rf.voteCount > len(rf.peers)/2 {
				rf.state = STATE_FOLLOWER
				rf.chanLeader <- true
			}
		}
	}
	return ok
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	defer rf.persist()
	reply.Success = false
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		return
	}
	rf.chanHeartbeat <- true
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.state = STATE_FOLLOWER
		rf.votedFor = -1
	}
	reply.Term = args.Term
	return
}

func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if ok {
		if rf.state != STATE_LEADER || args.Term != rf.currentTerm {
			return ok
		}
	
		if reply.Term > rf.currentTerm {
			rf.currentTerm = reply.Term
			rf.state = STATE_FOLLOWER
			rf.votedFor = -1
			rf.persist()
			return ok
		}
	}
	return ok
}

func (rf *Raft) bcRequestVote() {
	var args RequestVoteArgs
	rf.mu.Lock()
	args.Term = rf.currentTerm
	args.CandidateId = rf.me
	args.LastLogTerm = rf.getLastTerm()
	rf.mu.Unlock()
	for i := range rf.peers {
		if i != rf.me && rf.state == STATE_CANDIDATE {
			go func(i int) {
				var reply RequestVoteReply
				rf.sendRequestVote(i, &args, &reply)
			}(i)
		}
	}
}

func (rf *Raft) bcAppendEntries() {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	for i := range rf.peers {
		if i != rf.me && rf.state == STATE_LEADER {
			var args AppendEntriesArgs
			args.Term = rf.currentTerm
			args.LeaderId = rf.me
			go func(i int, args AppendEntriesArgs) {
				var reply AppendEntriesReply
				rf.sendAppendEntries(i, &args, &reply)
			}(i, args)
		}
	}
}



func (rf *Raft) setup() {
	for {
		switch rf.state {
		case STATE_FOLLOWER:
			select {
			case <-rf.chanHeartbeat:
			case <-rf.chanGrantVote:
			case <-time.After(time.Duration(rand.Int63() % 200 + 500) * time.Millisecond):
				rf.state = STATE_CANDIDATE
			}
		case STATE_LEADER:
			rf.bcAppendEntries()
			time.Sleep(HEARTBEATS)
		case STATE_CANDIDATE:
			rf.mu.Lock()
			rf.currentTerm++
			rf.votedFor = rf.me
			rf.voteCount = 1
			rf.persist()
			rf.mu.Unlock()
			go rf.bcRequestVote()
			select {
			case <-time.After(time.Duration(rand.Int63()% 200 + 500) * time.Millisecond):
			case <-rf.chanHeartbeat:
				rf.state = STATE_FOLLOWER
			case <-rf.chanLeader:
				rf.mu.Lock()
				rf.state = STATE_LEADER
				rf.mu.Unlock()
			}
		}
	}
}

func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true
	return index, term, isLeader
}

func (rf *Raft) Kill() {
}

func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	rf.state = STATE_FOLLOWER
	rf.votedFor = -1
	rf.log = append(rf.log, LogEntry{LogTerm: 0})
	rf.currentTerm = 0
	rf.chanGrantVote = make(chan bool, 100)
	rf.chanLeader = make(chan bool, 100)
	rf.chanHeartbeat = make(chan bool, 100)
	rf.readPersist(persister.ReadRaftState())
	go rf.setup()
	return rf
}