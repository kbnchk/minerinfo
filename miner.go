package minerinfo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/netip"
	"strings"
	"time"
)

// Miner represents general miner machine interface
type Miner interface {
	Type() string
	Model() string
	Version() string
	Hashrate() float64
	Pools() []Pools
}

func New(ip string, port uint16, timeout time.Duration) (Miner, error) {
	ipaddr, err := netip.ParseAddr(ip)
	if err != nil {
		return nil, fmt.Errorf("error parsing ip address: %w", err)
	}
	addr := netip.AddrPortFrom(ipaddr, port)

	info, err := getInfo(addr, timeout)
	if err != nil {
		return nil, err
	}
	return detectMiner(info), nil
}

// Detects miner type
// Determining the type of miner using the cgminer api is difficult, since different manufacturers and models
// have their own api implementation and the amount of information returned can vary significantly.
// Therefore, this implementation is based only on known types of miners and the information returned by them.
func detectMiner(info minerInfo) Miner {
	var typestring string
	if len(info.Stats) > 0 && info.Stats[0].Type != "" {
		typestring = info.Stats[0].Type
	} else {
		if len(info.Devdetails) > 0 && info.Devdetails[0].Name != "" {
			typestring = info.Devdetails[0].Name
		}
	}
	if strings.Split(typestring, " ")[0] == "Antminer" {
		return antminer{genericMiner{info}}
	}

	if typestring == "SM" {
		return whatsminer{genericMiner{info}}

	}
	return genericMiner{info}
}

func getInfo(addr netip.AddrPort, timeout time.Duration) (minerInfo, error) {
	var info minerInfo
	var reqErr error
	commands := []string{"stats", "summary", "devdetails", "pools"}

	for _, command := range commands {
		req := minerRequest{
			Command: command,
		}

		conn, err := net.DialTimeout("tcp", addr.String(), timeout)
		if err != nil {
			return info, fmt.Errorf("error connecting miner: %w", err)
		}
		defer conn.Close()

		conn.SetDeadline(time.Now().Add(timeout))
		reqErr = sendCommand(conn, req, &info)
		if reqErr != nil {
			break
		}
	}
	return info, reqErr
}

func sendCommand(conn io.ReadWriter, req minerRequest, resp *minerInfo) error {
	reqBytes, _ := json.Marshal(req)

	_, err := conn.Write(reqBytes)
	if err != nil {
		return fmt.Errorf("error sending command %s: %w", req.Command, err)
	}

	respBytes, err := io.ReadAll(conn)
	if err != nil {
		return fmt.Errorf("error reading response for command %s: %w", req.Command, err)
	}

	respBytes = bytes.Replace(respBytes, []byte("}{"), []byte("},{"), 1) // L3 buggy json output fix
	respBytes = bytes.TrimRight(respBytes, "\x00")

	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		return fmt.Errorf("error umarshalling response for command %s: %w", req.Command, err)
	}

	return nil
}
