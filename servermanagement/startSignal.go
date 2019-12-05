package servermanagement

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/SUMUKHA-PK/Raft-Distributed-Consensus/types"
)

// StartSignal initialises raft server instances and sends out signal to initiate leader election
func StartSignal(config types.Configuration, raftServers map[string]types.RaftServer, delay time.Duration, wg *sync.WaitGroup) error {
	serverDesignations := make(map[string]string)
	for _, server := range config.Servers {
		serverDesignations[server.URI()] = "follower"
	}

	for _, server := range config.Servers {
		raftServers[server.URI()] = types.RaftServer{"follower", server.IP, server.Port, serverDesignations}
	}

	payload, err := json.Marshal(raftServers)
	if err != nil {
		log.Fatal("Can't Marshall Payload to JSON in startSignal.go: %v\n", err)
		return err
	}

	time.Sleep(delay * time.Millisecond)
	sendInitiateSignals(config.Servers, payload)

	wg.Done()
	return nil
}

func sendInitiateSignals(servers []types.Server, payload []byte) {
	for _, server := range servers {
		go sendInitiateSignal(server.URL("http://", "startRaft"), strings.NewReader(string(payload)))
	}
}

func sendInitiateSignal(URL string, payload *strings.Reader) error {
	req, err := http.NewRequest("POST", URL, payload)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Printf("Bad request for %v while sending start signal: %v\n", URL, err)
		return err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("Bad response from %v while sending start signal: %v\n", URL, err)
		return err
	}

	res.Body.Close()
	return nil
}
