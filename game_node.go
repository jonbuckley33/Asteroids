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
	playerAddresses map[string]string
	PlayerId int
}

func NewGameServer(hostAddress string) (*GameNode, error) {
	gs := new(GameNode)
	gs.address = hostAddress
	gs.playerAddresses = make(map[string]string)
	gs.PlayerId = 0

	gs.playerAddresses["0"] = hostAddress

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
	gs.PlayerId = len(gs.playerAddresses)
	gs.playerAddresses[strconv.Itoa(gs.PlayerId)] = myHostAddress

	hostMap := make(map[int]string)
	for k, v := range(gs.playerAddresses) {
		i, _ := strconv.Atoi(k)
		hostMap[i] = v
	}

	// Make a new node as a "replacement" node.
	node, err := paxos.NewPaxosNode(myHostAddress, hostMap,
		gs.PlayerId, gs.PlayerId + 1, 5, true)
	if err != nil {
		// println("Could not connect to server:", err.Error())
		return nil, err
	} else {
		// println("Connected to server")
	}
	gs.node = node

	// Update num players.
	gs.MakeProposal("num_players", strconv.Itoa(gs.PlayerId + 1))

	// Update the player addresses map.
	v, _ := json.Marshal(gs.playerAddresses)
	playerAddressesEncoded := string(v)
	gs.MakeProposal("player_addresses", playerAddressesEncoded)

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
	if propReply.V == nil {
		return "", errors.New("Failed to get non-nil value.")
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


//Send player information to paxos nodes
func (gs *GameNode) SharePlayer(ship *Ship) {
	playerKey := fmt.Sprintf("player_%v", PlayerId)
	playerPos := fmt.Sprintf("(%v,%v,%v,%v,%v,%v,%v,%v)", 
		ship.PosX, ship.PosY, ship.Angle,
		ship.VelocityX, ship.VelocityY,
		ship.TurnRate, ship.AccelerationRate,
		ship.IsAlive())

	_, err := gs.MakeProposal(playerKey, playerPos)	
	if err != nil {
		println("Was not able to share the ship for player", PlayerId)
	}		
}


//Get player info from paxos nodes
func (gs *GameNode) GetPlayers() map[int]*Ship{
	// TODO
	encodedPlayerAddresses, _ := gs.GetValue("player_addresses")
	json.Unmarshal([]byte(encodedPlayerAddresses), &gs.playerAddresses)

	m := make(map[int]*Ship)

	for id := range(gs.playerAddresses) {
		query := fmt.Sprintf("player_%v", id)
		posString, err := gs.GetValue(query)

		// Only worry about players who we have positions for.
		if err == nil {
			var x,y,angle,vX,vY,turnRate,accelerationRate float64
			var isAlive bool

			fmt.Sscanf(posString, "(%v,%v,%v,%v,%v,%v,%v,%v)", &x,&y,&angle,&vX,&vY,&turnRate,&accelerationRate,&isAlive)
			
			newShip:=new(Ship)

			newShip.PosX=x
			newShip.PosY=y
			newShip.Angle=angle
			newShip.VelocityX=vX
			newShip.VelocityY=vY
			newShip.TurnRate=turnRate
			newShip.AccelerationRate=accelerationRate
			newShip.isAlive=isAlive

			i, _ := strconv.Atoi(id)
			m[i] = newShip
		}
	}

	return m
}

func (gs *GameNode) ShareAsteroids(asteroids map[int]*Asteroid) {
	counter := 0
	asteroidIds := make([]int, len(asteroids))
	for i, asteroid := range(asteroids) {
		// Add id to map.
		asteroidIds[counter] = asteroid.Id
		counter += 1

		// Share asteroid data.
		asteroidKey := fmt.Sprintf("asteroid_%v", i)
		asteroidPos := fmt.Sprintf("(%v,%v,%v,%v,%v,%v,%v,%v,%v)", 
			asteroid.PosX, asteroid.PosY, asteroid.Angle,
			asteroid.VelocityX, asteroid.VelocityY,
			asteroid.TurnRate, asteroid.AccelerationRate,
			asteroid.SizeRatio, asteroid.Lives)
		gs.MakeProposal(asteroidKey, asteroidPos)
	}

	asteroidIdEncoded, _ := json.Marshal(asteroidIds)
	gs.MakeProposal("asteroid_ids", string(asteroidIdEncoded))
}

func (gs *GameNode) GetAsteroids() map[int]*Asteroid {
	// Get list of IDs of asteroids.
	var asteroidIds []int
	asteroidIdsEncoded, _ := gs.GetValue("asteroid_ids")

	json.Unmarshal([]byte(asteroidIdsEncoded), &asteroidIds)
	
	// Mapping from ID to asteroid.
	asteroids := make(map[int]*Asteroid)
	for _, id := range(asteroidIds) {
		asteroidKey := fmt.Sprintf("asteroid_%v", id)
		asteroidEncoded, _ := gs.GetValue(asteroidKey)

		var posX, posY, angle, turnRate, vX, vY, acceleration, size float64
		var lives int

		fmt.Sscanf(asteroidEncoded, "(%v,%v,%v,%v,%v,%v,%v,%v,%v)", &posX, &posY,
			&angle, &vX, &vY, &turnRate, &acceleration, &size, &lives)

		asteroid := new(Asteroid)
		asteroid.PosX = posX
		asteroid.PosY = posY
		asteroid.Angle = angle
		asteroid.TurnRate = turnRate
		asteroid.VelocityX = vX
		asteroid.VelocityY = vY
		asteroid.SizeRatio = size
		asteroid.Id = id
		asteroid.Lives = lives

		asteroids[id] = asteroid
	}

	return asteroids
}


// Gets the hostports of all of the players registered with the game
// server located at "server". 
func (gs *GameNode) GetPlayerAddresses(server string) (map[string]string, error) {
	// println("Trying to dial the server", server)
	client, err := rpc.DialHTTP("tcp", server)
	if err != nil {
		// println("Could not open TCP connection to server", server)
		return nil, err
	} else {
		// println("Opened TCP connection with server", server)
	}

	getArgs := &paxosrpc.GetValueArgs{
		Key: "player_addresses",
	}
	getReply := new(paxosrpc.GetValueReply)
	// println("Trying to call Paxos Get Val")
	client.Call("PaxosNode.GetValue", getArgs, getReply)
	// println("Finished call")
	var vals map[string]string
	json.Unmarshal([]byte(getReply.V.(string)), &vals)

	client.Close()

	return vals, nil
}


//Get ship velocities
func (gs *GameNode) GetPlayerVelocities() map[int][]int {
	m := make(map[int][]int)

	for id := range(gs.playerAddresses) {
		query := fmt.Sprintf("player_%v_velocity", id)
		posString, err := gs.GetValue(query)

		// Only worry about players who we have velocities for.
		if err == nil {
			var Vx, Vy int
			fmt.Sscanf(posString, "(%v,%v)", &Vx, &Vy)
			arr := []int{Vx, Vy}

			i, _ := strconv.Atoi(id)
			m[i] = arr
		}
	}

	return m
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

			i, _ := strconv.Atoi(id)
			m[i] = arr
		}
	}

	return m
}


