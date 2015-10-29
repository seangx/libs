package nsqredo

import (
	"bytes"
	log "github.com/gonet2/libs/nsq-logger"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/vmihailenco/msgpack.v2"
	"net/http"
	"os"
)

const (
	ENV_NSQD         = "NSQD_HOST"
	DEFAULT_PUB_ADDR = "http://172.17.42.1:4151/pub?topic=REDOLOG"
	MIME             = "application/octet-stream"
)

var (
	_junk []byte
)

// a data change
type Change struct {
	Collection string // collection
	Field      string // field "a.b.c.1.d"
	Doc        []byte // msgpack serialized data
}

// a redo record represents complete transaction
type RedoRecord struct {
	API     string   // the api name
	UID     int32    // userid
	TS      uint64   // timestamp should get from snowflake
	Changes []Change // changes
}

var (
	_pub_addr string
	_prefix   string
)

func init() {
	// get nsqd publish address
	_pub_addr = DEFAULT_PUB_ADDR
	if env := os.Getenv(ENV_NSQD); env != "" {
		_pub_addr = env + "/pub?topic=REDOLOG"
	}
	_junk = make([]byte, 1024)
}

// add a change with o(old value) and n(new value)
func (r *RedoRecord) AddChange(collection, field string, doc interface{}) {
	doc_bin, err := bson.Marshal(doc)
	if err != nil {
		log.Error(err)
		return
	}
	r.Changes = append(r.Changes, Change{Collection: collection, Field: field, Doc: doc_bin})
}

func NewRedoRecord(uid int32, api string, ts uint64) *RedoRecord {
	return &RedoRecord{UID: uid, API: api, TS: ts}
}

// publish to nsqd (localhost nsqd is suggested!)
func Publish(r *RedoRecord) {
	// pack message
	pack, err := msgpack.Marshal(r)
	if err != nil {
		log.Critical(err)
		return
	}

	// post to nsqd
	resp, err := http.Post(_pub_addr, MIME, bytes.NewReader(pack))
	if err != nil {
		log.Critical(err)
		return
	}

	// read & discard
	for {
		if _, err := resp.Body.Read(_junk); err != nil {
			break
		}
	}
	resp.Body.Close()
}
