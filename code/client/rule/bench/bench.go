package bench

import (
	"fmt"
	"net"
	"net/http"

	"github.com/jkstack/natpass/code/client/global"
	"github.com/jkstack/natpass/code/client/pool"
	"github.com/jkstack/natpass/code/client/rule"
	"github.com/lwch/logging"
	"github.com/lwch/runtime"
)

// Bench benchmark handler
type Bench struct {
	Name string
	cfg  global.Rule
}

// Link bench link
type Link struct {
	id string
}

// GetID get link id
func (link *Link) GetID() string {
	return link.id
}

// GetBytes get send and recv bytes
func (link *Link) GetBytes() (uint64, uint64) {
	return 0, 0
}

// GetPackets get send and recv packets
func (link *Link) GetPackets() (uint64, uint64) {
	return 0, 0
}

// New new benchmark handler
func New(cfg global.Rule) *Bench {
	return &Bench{
		Name: cfg.Name,
		cfg:  cfg,
	}
}

// NewLink new link
func (bench *Bench) NewLink(id, remote string, remoteIdx uint32, localConn net.Conn, remoteConn *pool.Conn) rule.Link {
	return &Link{id: id}
}

// GetName get bench rule name
func (bench *Bench) GetName() string {
	return bench.Name
}

// GetTypeName get bench rule type name
func (bench *Bench) GetTypeName() string {
	return "bench"
}

// GetTarget get target of this rule
func (bench *Bench) GetTarget() string {
	return bench.cfg.Target
}

// GetLinks get rule links
func (bench *Bench) GetLinks() []rule.Link {
	return nil
}

// GetRemote get remote target name
func (bench *Bench) GetRemote() string {
	return bench.cfg.Target
}

// GetPort get listen port
func (bench *Bench) GetPort() uint16 {
	return bench.cfg.LocalPort
}

// Handle handle shell
func (bench *Bench) Handle(pl *pool.Pool) {
	defer func() {
		if err := recover(); err != nil {
			logging.Error("close shell: %s, err=%v", bench.Name, err)
		}
	}()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bench.http(pl, w, r)
	})
	svr := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", bench.cfg.LocalAddr, bench.cfg.LocalPort),
		Handler: mux,
	}
	runtime.Assert(svr.ListenAndServe())
}
