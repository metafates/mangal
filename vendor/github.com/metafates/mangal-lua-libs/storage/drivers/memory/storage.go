// this storage can be used for projects that do not store much data and do not save memory
package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	lua_json "github.com/metafates/mangal-lua-libs/json"
	interfaces "github.com/metafates/mangal-lua-libs/storage/drivers/interfaces"

	lua "github.com/yuin/gopher-lua"
)

var listOfStorages = &listStorages{list: make(map[string]*Storage)}

type listStorages struct {
	sync.Mutex
	list map[string]*Storage
}

type Storage struct {
	sync.Mutex
	filename     string
	Data         map[string]*storageValue `json:"data"`
	usageCounter int
}

type storageValue struct {
	Value      []byte `json:"value"`        // json value
	MaxValidAt int64  `json:"max_valid_at"` // unix nano
}

func (sv *storageValue) valid() bool {
	return sv.MaxValidAt > time.Now().UnixNano()
}

func (st *Storage) New(filename string) (interfaces.Driver, error) {

	listOfStorages.Lock()
	defer listOfStorages.Unlock()

	if result, ok := listOfStorages.list[filename]; ok {
		result.Lock()
		defer result.Unlock()
		result.usageCounter++
		return result, nil
	}

	s := &Storage{Data: make(map[string]*storageValue, 0)}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// create
		dst, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		dst.Close()
	} else {
		// read && decode
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, s); err != nil {
			return nil, err
		}
	}
	s.filename = filename
	s.usageCounter++
	listOfStorages.list[filename] = s
	go s.loop()
	return s, s.Sync()
}

func (s *Storage) Sync() error {
	s.Lock()
	defer s.Unlock()
	tmpFilename := s.filename + ".tmp"
	// clean
	newData := make(map[string]*storageValue, 0)
	for k, v := range s.Data {
		if v.valid() {
			newData[k] = v
		}
	}
	s.Data = newData
	// clean end
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(tmpFilename, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmpFilename, s.filename)
}

func (s *Storage) Close() error {
	listOfStorages.Lock()
	defer listOfStorages.Unlock()
	if err := s.Sync(); err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	s.usageCounter--
	return nil
}

func (s *Storage) loop() {
	for {
		time.Sleep(time.Minute)
		if err := s.Sync(); err != nil {
			log.Printf("[ERROR] scheduler for memory storage [%p-%s], sync save: %s\n", s, s.filename, err.Error())
		} else {
			if s.usageCounter == 0 {
				listOfStorages.Lock()
				log.Printf("[INFO] close unused memory storage [%p-%s]\n", s, s.filename)
				delete(listOfStorages.list, s.filename)
				listOfStorages.Unlock()
				return
			}
		}
	}
}

func (s *Storage) Keys() ([]string, error) {
	result := []string{}
	s.Lock()
	defer s.Unlock()
	for k, _ := range s.Data {
		result = append(result, k)
	}
	return result, nil
}

func (s *Storage) Dump(L *lua.LState) (map[string]lua.LValue, error) {
	result := make(map[string]lua.LValue, 0)
	s.Lock()
	defer s.Unlock()
	for k, v := range s.Data {
		if v.valid() {
			value, err := lua_json.ValueDecode(L, v.Value)
			if err != nil {
				return nil, err
			}
			result[k] = value
		}
	}
	return result, nil
}
