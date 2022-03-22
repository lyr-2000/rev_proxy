package consistenhash

import (
	"errors"
	"hash/crc32"
	"myproxyHttp/utils/strutil"
	"sort"
)

var ErrNodeNotFound = errors.New("node not found")

type Ring struct {
	Nodes Nodes
	//sync.RWMutex
}

func NewRing() *Ring {
	return &Ring{Nodes: Nodes{}}
}

func (r *Ring) Has(id string) bool {
	for _, v := range r.Nodes {
		if id == v.Id {
			return true
		}
	}
	return false
}

func (r *Ring) AddNode(id string, value interface{}) {
	//r.Lock()
	//defer r.Unlock()

	node := NewNode(id, value)
	r.Nodes = append(r.Nodes, node)

	sort.Sort(r.Nodes)
}

func (r *Ring) RemoveNode(id string) error {
	//r.Lock()
	//defer r.Unlock()

	i := r.search(&id)
	if i >= r.Nodes.Len() || r.Nodes[i].Id != id {
		return ErrNodeNotFound
	}

	r.Nodes = append(r.Nodes[:i], r.Nodes[i+1:]...)

	return nil
}

func (r *Ring) Get(id string) interface{} {
	//r.RLock()
	//defer r.RUnlock()
	i := r.search(&id)
	ln := r.Nodes.Len()

	if ln <= 0 {
		return nil
	}
	//ln > 0
	if i >= ln {
		i = 0
	}

	return r.Nodes[i].Value
}

func (r *Ring) search(id *string) int {
	hashid := hashId(id)
	searchfn := func(i int) bool {
		return r.Nodes[i].HashId >= hashid
	}

	return sort.Search(r.Nodes.Len(), searchfn)
}

//----------------------------------------------------------
// Node
//----------------------------------------------------------

type Node struct {
	Id     string
	HashId uint32
	Value  interface{}
}

func NewNode(id string, value interface{}) *Node {
	return &Node{
		Id:     id,
		HashId: hashId(&id),
		Value:  value,
	}
}

type Nodes []*Node

func (n Nodes) Len() int           { return len(n) }
func (n Nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Nodes) Less(i, j int) bool { return n[i].HashId < n[j].HashId }

//----------------------------------------------------------
// Helpers
//----------------------------------------------------------

func hashId(key *string) uint32 {
	return crc32.ChecksumIEEE(strutil.String2Bytes(key))
}
