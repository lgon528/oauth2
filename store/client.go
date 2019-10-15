package store

import (
	"errors"
	"fmt"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/generates"
	"sync"
)

// NewClientStore create client store
func NewClientStore() *ClientStore {
	return &ClientStore{
		index: 10000,
		data:  make(map[string]oauth2.ClientInfo),
	}
}

// ClientStore client information store
type ClientStore struct {
	sync.RWMutex
	index uint64
	data  map[string]oauth2.ClientInfo
}

func (cs *ClientStore) isExistedByName(name string) bool {
	return false
}

func (cs *ClientStore) generateSecret(info oauth2.ClientInfo) string {
	return generates.Secret(32)
	//return info.GetID()
}

func (cs *ClientStore) AddClient(info oauth2.ClientInfo) error {
	cs.Lock()
	defer cs.Unlock()

	if cs.isExistedByName(info.GetName()) {
		return fmt.Errorf("client name already used")
	}
	info.SetID(fmt.Sprintf("%d", cs.index))
	info.SetSecret(cs.generateSecret(info))
	cs.data[info.GetID()] = info

	cs.index++

	return nil
}

func (cs *ClientStore) GetClients() []oauth2.ClientInfo {
	infos := make([]oauth2.ClientInfo, 0)

	cs.RLock()
	defer cs.RUnlock()

	for _, info := range cs.data {
		infos = append(infos, info)
	}

	return infos
}

// GetByID according to the ID for the client information
func (cs *ClientStore) GetByID(id string) (cli oauth2.ClientInfo, err error) {
	cs.RLock()
	defer cs.RUnlock()
	if c, ok := cs.data[id]; ok {
		cli = c
		return
	}
	err = errors.New("not found")
	return
}

// Set set client information
func (cs *ClientStore) Set(id string, cli oauth2.ClientInfo) (err error) {
	cs.Lock()
	defer cs.Unlock()
	cs.data[id] = cli
	return
}
