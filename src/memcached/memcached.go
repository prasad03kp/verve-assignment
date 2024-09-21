package memcached

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/prasad03kp/verve-assignment/kafka"
)

var (
	client    *memcache.Client
	logFile   *os.File
	dataMutex sync.Mutex
	keyList   = "memcache_keys"
)

func init() {
	var err error
	client = getMemCacheClient()
	logFile, err = os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	go periodicTask()
}

func getMemCacheClient() *memcache.Client {
	if client == nil {
		client = memcache.New("127.0.0.1:11211")
	}
	return client
}

func WriteToMemCache(id int) {
	timestamp := time.Now().Format(time.RFC3339)
	value := fmt.Sprintf("ID: %d, Timestamp: %s", id, timestamp)
	client = getMemCacheClient()
	item := &memcache.Item{Key: strconv.Itoa(id), Value: []byte(value)}

	err := client.Set(item)
	if err != nil {
		log.Println(err.Error())
		log.Fatalf("Failed to write to Memcached for id %d", id)
	}

	err = addKeyToList(item.Key)
	if err != nil {
		log.Printf("Failed to add key to key list: %v", err)
	}
}

func addKeyToList(key string) error {
	existingKeysItem, err := client.Get(keyList)
	if err != nil && err != memcache.ErrCacheMiss {
		return fmt.Errorf("failed to retrieve key list: %w", err)
	}

	var keys []string
	if existingKeysItem != nil {
		keys = strings.Split(string(existingKeysItem.Value), ",")
	}

	if !containsKey(keys, key) {
		keys = append(keys, key)
		newKeyList := strings.Join(keys, ",")
		err = client.Set(&memcache.Item{Key: keyList, Value: []byte(newKeyList)})
		if err != nil {
			return fmt.Errorf("failed to update key list: %w", err)
		}
	}
	return nil
}

func containsKey(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

func getAllKeys() []string {
	client = getMemCacheClient()
	item, err := client.Get(keyList)
	if err != nil {
		log.Printf("Failed to retrieve key list: %v", err)
		return []string{}
	}

	if item != nil {
		return strings.Split(string(item.Value), ",")
	}

	return []string{}
}

func periodicTask() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		countUniqueAndLog()
		clearLastMinuteData()
	}
}

func countUniqueAndLog() {
	dataMutex.Lock()
	defer dataMutex.Unlock()

	client = getMemCacheClient()
	keys := getAllKeys()
	if len(keys) == 0 || (len(keys) == 1 && keys[0] == "") {
		log.Printf("No keys found")
		return
	}

	items, err := client.GetMulti(keys)
	if err != nil {
		log.Printf("Failed to retrieve items: %v", err)
		return
	}

	now := time.Now()
	lastMinute := now.Add(-1 * time.Minute)

	uniqueIDs := make(map[string]struct{})

	for _, item := range items {
		var id int
		var timestampStr string
		_, err := fmt.Sscanf(string(item.Value), "ID: %d, Timestamp: %s", &id, &timestampStr)
		if err != nil {
			log.Printf("Failed to parse item value: %v", err)
			continue
		}

		timestamp, err := time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			log.Printf("Failed to parse timestamp: %v", err)
			continue
		}

		if timestamp.After(lastMinute) && timestamp.Before(now) {
			uniqueIDs[item.Key] = struct{}{}
		}
	}

	count := len(uniqueIDs)
	logEntry := fmt.Sprintf("%s: Unique ID count is %d\n", now.Format(time.RFC3339), count)
	kafka.WriteToKafka(logEntry)
}


func clearLastMinuteData() {
	dataMutex.Lock()
	defer dataMutex.Unlock()

	client = getMemCacheClient()
	keys := getAllKeys()
	if len(keys) == 0 || (len(keys) == 1 && keys[0] == "") {
		log.Printf("No keys found")
		return
	}

	items, err := client.GetMulti(keys)
	if err != nil {
		log.Printf("Failed to retrieve items: %v", err)
		return
	}

	lastMinute := time.Now().Add(-1 * time.Minute).Format(time.RFC3339)

	for _, item := range items {
		var id int
		var timestamp string
		fmt.Sscanf(string(item.Value), "ID: %d, Timestamp: %s", &id, &timestamp)
		if timestamp < lastMinute {
			client.Delete(item.Key)
			removeKeyFromList(item.Key)
		}
	}
}

func removeKeyFromList(key string) {
	existingKeysItem, err := client.Get(keyList)
	if err != nil {
		log.Printf("Failed to retrieve key list: %v", err)
		return
	}

	if existingKeysItem != nil {
		keys := strings.Split(string(existingKeysItem.Value), ",")
		updatedKeys := make([]string, 0)
		for _, k := range keys {
			if k != key {
				updatedKeys = append(updatedKeys, k)
			}
		}
		newKeyList := strings.Join(updatedKeys, ",")
		err := client.Set(&memcache.Item{Key: keyList, Value: []byte(newKeyList)})
		if err != nil {
			log.Printf("Failed to update key list: %v", err)
		}
	}
}

func CountUniqueIDsInCurrentMinute() int {
	dataMutex.Lock()
	defer dataMutex.Unlock()

	client = getMemCacheClient()
	keys := getAllKeys()
	if len(keys) == 0 || (len(keys) == 1 && keys[0] == "") {
		log.Printf("No keys found")
		return 0
	}

	items, err := client.GetMulti(keys)
	if err != nil {
		log.Printf("Failed to retrieve items: %v", err)
		return 0
	}

	now := time.Now()
	currentMinuteStart := now.Truncate(time.Minute)

	uniqueIDs := make(map[string]struct{})

	for _, item := range items {
		var id int
		var timestampStr string
		_, err := fmt.Sscanf(string(item.Value), "ID: %d, Timestamp: %s", &id, &timestampStr)
		if err != nil {
			log.Printf("Failed to parse item value: %v", err)
			continue
		}

		timestamp, err := time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			log.Printf("Failed to parse timestamp: %v", err)
			continue
		}

		if timestamp.After(currentMinuteStart) && timestamp.Before(now) {
			uniqueIDs[item.Key] = struct{}{}
		}
	}

	return len(uniqueIDs)
}
