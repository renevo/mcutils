package minecraft

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/millkhan/mcstatusgo/v2"
)

// Status will return the servers status from the status query
// This is done over UDP and is configured under the query.port and enable-query server properties
func (s *Server) Status() (mcstatusgo.FullQueryResponse, error) {
	port := uint16(25565)

	if portCfg := s.Properties["query.port"]; portCfg != "" {
		portParsed, _ := strconv.ParseUint(portCfg, 10, 16)
		if portParsed > 0 {
			port = uint16(portParsed)
		}
	}

	if value, ok := s.Properties["enable-query"]; !ok || !strings.EqualFold(value, "true") {
		return mcstatusgo.FullQueryResponse{}, errors.New("Query Not Enabled")
	}

	// hard coded to localhost if not configured
	address := s.Properties["server-ip"]
	if address == "" {
		address = "localhost"
	}

	// using liberal timeouts, its against a server running locally
	return mcstatusgo.FullQuery(address, port, time.Second*5, time.Second*10)
}
