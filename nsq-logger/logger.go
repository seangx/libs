package nsqlogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

const (
	ENV_NSQD         = "NSQD_HOST"
	DEFAULT_PUB_ADDR = "http://172.17.42.1:4151/pub?topic=LOG"
	MIME             = "application/octet-stream"
)

const (
	FINEST byte = iota
	FINE
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	CRITICAL
)

// log format to sent to nsqd, packed with json
type LogFormat struct {
	Prefix string
	Time   time.Time
	Host   string
	Level  byte
	Msg    string
	Caller string
	LineNo int
}

var (
	_pub_addr string
	_prefix   string
	_ch       chan []byte
)

func init() {
	// get nsqd publish address
	_pub_addr = DEFAULT_PUB_ADDR
	if env := os.Getenv(ENV_NSQD); env != "" {
		_pub_addr = env + "/pub?topic=LOG"
	}
	_ch = make(chan []byte, 4096)
	go publish_task()
}

func publish_task() {
	for {
		// post to nsqd
		bts := <-_ch
		resp, err := http.Post(_pub_addr, MIME, bytes.NewReader(bts))
		if err != nil {
			log.Println(err, string(bts))
			continue
		}

		// read response
		if _, err := ioutil.ReadAll(resp.Body); err != nil {
			log.Println(err, string(bts))
		}

		// close
		resp.Body.Close()
	}
}

// publish to nsqd (localhost nsqd is suggested!)
func publish(msg LogFormat) {
	// fill in the common fields
	hostname, _ := os.Hostname()
	msg.Host = hostname
	msg.Time = time.Now()
	msg.Prefix = _prefix

	// Determine caller func
	if pc, _, lineno, ok := runtime.Caller(2); ok {
		msg.Caller = runtime.FuncForPC(pc).Name()
		msg.LineNo = lineno
	}

	// pack message
	if bts, err := json.Marshal(msg); err == nil {
		_ch <- bts
	} else {
		log.Println(err, msg)
		return
	}
}

// set prefix
func SetPrefix(prefix string) {
	_prefix = prefix
}

// wrappers for diffent loglevels
func Finest(v ...interface{}) {
	msg := LogFormat{Level: FINEST, Msg: fmt.Sprint(v...)}
	publish(msg)
}

func Finestf(format string, v ...interface{}) {
	msg := LogFormat{Level: FINEST, Msg: fmt.Sprintf(format, v...)}
	publish(msg)
}

func Fine(v ...interface{}) {
	msg := LogFormat{Level: FINE, Msg: fmt.Sprint(v...)}
	publish(msg)
}

func Finef(format string, v ...interface{}) {
	msg := LogFormat{Level: FINE, Msg: fmt.Sprintf(format, v...)}
	publish(msg)
}

func Debug(v ...interface{}) {
	msg := LogFormat{Level: DEBUG, Msg: fmt.Sprint(v...)}
	publish(msg)
}

func Debugf(format string, v ...interface{}) {
	msg := LogFormat{Level: DEBUG, Msg: fmt.Sprintf(format, v...)}
	publish(msg)
}

func Trace(v ...interface{}) {
	msg := LogFormat{Level: TRACE, Msg: fmt.Sprint(v...)}
	publish(msg)
}

func Tracef(format string, v ...interface{}) {
	msg := LogFormat{Level: TRACE, Msg: fmt.Sprintf(format, v...)}
	publish(msg)
}

func Info(v ...interface{}) {
	msg := LogFormat{Level: INFO, Msg: fmt.Sprint(v...)}
	publish(msg)
}

func Infof(format string, v ...interface{}) {
	msg := LogFormat{Level: INFO, Msg: fmt.Sprintf(format, v...)}
	publish(msg)
}

func Warning(v ...interface{}) {
	msg := LogFormat{Level: WARNING, Msg: fmt.Sprint(v...)}
	publish(msg)
}

func Warningf(format string, v ...interface{}) {
	msg := LogFormat{Level: WARNING, Msg: fmt.Sprintf(format, v...)}
	publish(msg)
}

func Error(v ...interface{}) {
	msg := LogFormat{Level: ERROR, Msg: fmt.Sprint(v...)}
	publish(msg)
}

func Errorf(format string, v ...interface{}) {
	msg := LogFormat{Level: ERROR, Msg: fmt.Sprintf(format, v...)}
	publish(msg)
}

func Critical(v ...interface{}) {
	msg := LogFormat{Level: CRITICAL, Msg: fmt.Sprint(v...)}
	publish(msg)
}

func Criticalf(format string, v ...interface{}) {
	msg := LogFormat{Level: CRITICAL, Msg: fmt.Sprintf(format, v...)}
	publish(msg)
}
