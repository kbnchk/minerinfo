package minerinfo

import (
	"bufio"
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
	h := hostNew(addr, timeout)
	err = h.getInfo()
	if err != nil {
		return nil, err
	}
	return h.detectMiner(), nil
}

type host struct {
	addr    netip.AddrPort
	info    minerInfo
	timeout time.Duration
}

func hostNew(addr netip.AddrPort, timeout time.Duration) host {
	return host{
		addr:    addr,
		info:    minerInfo{},
		timeout: timeout,
	}
}

// Detects miner type
// Determining the type of miner using the cgminer api is difficult, since different manufacturers and models
// have their own api implementation and the amount of information returned can vary significantly.
// Therefore, this implementation is based only on known types of miners and the information returned by them.
func (h *host) detectMiner() Miner {
	var typestring string
	if len(h.info.Stats) > 0 && h.info.Stats[0].Type != "" {
		typestring = h.info.Stats[0].Type
	} else {
		if len(h.info.Devdetails) > 0 && h.info.Devdetails[0].Name != "" {
			typestring = h.info.Devdetails[0].Name

		}
	}

	if strings.Split(typestring, " ")[0] == "Antminer" {
		return antminer{genericMiner{h.info}}
	}

	if typestring == "SM" {
		return whatsminer{genericMiner{h.info}}
	}
	return genericMiner{h.info}
}

// collects miner information
func (h *host) getInfo() error {
	commands := []string{"stats", "summary", "devdetails", "pools"}
	for _, cmd := range commands {
		err := h.exec(cmd, &h.info)
		if err != nil {
			return err
		}
	}
	return nil
}

// Establishes connection sends command and recieves responce
func (h *host) exec(cmd string, v any) error {
	conn, err := net.DialTimeout("tcp", h.addr.String(), h.timeout)
	if err != nil {
		return fmt.Errorf("error connecting to host: %w", err)
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(h.timeout))
	req := minerRequest{
		Command: cmd,
	}
	err = json.NewEncoder(conn).Encode(req)
	if err != nil {
		return err
	}
	r, err := connReader(conn)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(r)
	return decoder.Decode(v)
}

func connReader(c net.Conn) (io.Reader, error) {
	result, err := bufio.NewReader(c).ReadBytes(0x00)
	if err != nil && err != io.EOF {
		return nil, err
	}
	result = bytes.Replace(result, []byte("}{"), []byte("},{"), 1) //Antminer L3 buggy json output fix
	return bytes.NewReader(bytes.TrimRight(result, "\x00")), nil
}
