package main

import (
	"errors"
	"fmt"
	"net/rpc"
	"encoding/json"
	"strconv"

	"github.com/cmu440-F15/paxosapp/paxos"
	"github.com/cmu440-F15/paxosapp/rpc/paxosrpc"
)

type GameNode struct {
	node paxos.PaxosNode
	address string
	playerAddresses map[int]string
	playerId int
}

func NewGameServer(hostAddress string) (*GameNode, error) {
	gs := new(GameNode)
	gs.address = hostAddress
	gs.playerAddresses = make(map[int]string)
	gs.playerId = 0

	gs.playerAddresses[0] = hostAddress

	hostMap := make(map[int]string)
	hostMap[0] = hostAddress
	node, err := paxos.NewPaxosNode(hostAddress, hostMap, 1, 0, 5, false)
	if err != nil {
		return nil, err
	}
	gs.node = node

	gs.InitializeGame()

	return gs, nil
}

func NewGameClient(myHostAddress, serverHostAddress string) (*GameNode, error) {
	gs := new(GameNode)
	gs.address = myHostAddress

	// Contact server for player info.
	playerAddresses, err := gs.GetPlayerAddresses(serverHostAddress)
	if err != nil {
		return nil, err
	}

	gs.playerAddresses = playerAddresses
	gs.playerId = len(gs.playerAddresses)

	// Make a new node as a "replacement" node.
	node, err := paxos.NewPaxosNode(myHostAddress, gs.playerAddresses,
		gs.playerId, gs.playerId, 5, true)
	if err != nil {
		return nil, err
	}
	gs.node = node

	gs.MakeProposal("num_players", strconv.Itoa(gs.playerId))

	return gs, nil
}

func (gs *GameNode) MakeProposal(key string, value string) (string, error) {
	// Get a proposal number.
	Nargs := &paxosrpc.ProposalNumberArgs{
		Key: key,
	}
	Nreply := new(paxosrpc.ProposalNumberReply)
	err := gs.node.GetNextProposalNumber(Nargs, Nreply)
	if err != nil {
		return "", err
	}
	
	// Propose value with proposal number.
	propArgs := &paxosrpc.ProposeArgs{
		N: Nreply.N,
		Key: key,
		V: value,
	}
	propReply := new(paxosrpc.ProposeReply)
	err = gs.node.Propose(propArgs, propReply)
	if err != nil {
		return "", err
	}

	return propReply.V.(string), nil
}

func (gs *GameNode) GetValue(key string) (string, error) {
	getArgs := &paxosrpc.GetValueArgs{
		Key: key,
	}
	getReply := new(paxosrpc.GetValueReply)
	err := gs.node.GetValue(getArgs, getReply)
	if err != nil {
		return "", err
	}

	if getReply.Status == paxosrpc.KeyNotFound {
		return "", errors.New("Could not find key")
	}

	return getReply.V.(string), nil
}

// Should only be called from master server of a game.
func (gs *GameNode) InitializeGame() {
	_, err := gs.MakeProposal("num_players", "1")
	if err != nil {
		panic("Could not initialize game server")
	}

	v, _ := json.Marshal(gs.playerAddresses)
	playerAddressesEncoded := string(v)
	_, err = gs.MakeProposal("player_addresses", playerAddressesEncoded)
	if err != nil {
		panic("Could not initialize game server")
	}
}

func (gs *GameNode) SharePlayerLocation(x, y float64) {
	encodedCoords := fmt.Sprintf("(%v,%v)", x, y)
	encodedKey := fmt.Sprintf("player_%v_location", gs.playerId)
	gs.MakeProposal(encodedKey, encodedCoords)
}

// Gets the hostports of all of the players registered with the game
// server located at "server". 
func (gs *GameNode) GetPlayerAddresses(server string) (map[int]string, error) {
	client, err := rpc.DialHTTP("tcp", server)
	if err != nil {
		return nil, err
	}

	getArgs := &paxosrpc.GetValueArgs{
		Key: "player_addresses",
	}
	getReply := new(paxosrpc.GetValueReply)
	client.Call("PaxosNode.GetValue", getArgs, getReply)

	var vals map[int]string
	json.Unmarshal([]byte(getReply.V.(string)), &vals)

	client.Close()

	return vals, nil
}

// Returns the positions of every known player with a map of arrays
// where i -> [x, y] for player i at location x, y
func (gs *GameNode) GetPlayerLocations() map[int][]int {
	m := make(map[int][]int)

	for id := range(gs.playerAddresses) {
		query := fmt.Sprintf("player_%v_location", id)
		posString, err := gs.GetValue(query)

		// Only worry about players who we have positions for.
		if err == nil {
			var x, y int
			fmt.Sscanf(posString, "(%v,%v)", &x, &y)
			arr := []int{x, y}

			m[id] = arr
		}
	}

	return m
}

/*func (gs *GameNode) SaveAsteroidLocations(asteroids []byte){
	_, err := gs.MakeProposal("asteroids", asteroids(string))
}*/