/*Tool for finding oplog metrics post-fact

A common issue in mongo (especially with mmap) is that there is a spike of oplog
activity which causes secondary read latency. After the fact all you have from
the oplog metrics is whether it was an insert/update/delete -- which isn't
enough to look into it. This CLI tool connects to the oplog and generates metrics
per namespace to aid in investigation


*/

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Sirupsen/logrus"
	flags "github.com/jessevdk/go-flags"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// time format on CLI
const Layout = "2006-01-02 15:04:05"

var opts struct {
	MongoHost string `long:"mongodb" description:"hostname/ip of the mongo host" default:"127.0.0.1"`

	StartTime string `long:"start" description:"start time" required:"true"`
	EndTime   string `long:"end" description:"end time" required:"true"`

	TimeLocation string `long:"location" description:"time location string" default:"Local"`
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		logrus.Fatalf("Error parsing flags: %v", err)
	}

	session, err := mgo.Dial(opts.MongoHost)
	if err != nil {
		panic(err)
	}

	location, err := time.LoadLocation(opts.TimeLocation)
	if err != nil {
		panic(err)
	}

	// load times
	fromDate, err := time.ParseInLocation(Layout, opts.StartTime, location)
	if err != nil {
		panic(err)
	}

	toDate, err := time.ParseInLocation(Layout, opts.EndTime, location)
	if err != nil {
		panic(err)
	}

	c := session.DB("local").C("oplog.rs")

	iter := c.Find(
		bson.M{
			"ts": bson.M{
				"$gt": TimeToMongoTime(fromDate),
				"$lt": TimeToMongoTime(toDate),
			},
		},
	).Iter()

	// Store metrics
	// we want namespace -> count
	metrics := NewMetrics()

	// this seems messy, but you cna't take address without it being a variable
	tmp := int64(1)
	newValue := &tmp

	entry := OpLog{}
	for iter.Next(&entry) {
		fmt.Println(entry)
		var m *sync.Map
		switch entry.Operation {
		case "i":
			m = metrics.Insert
		case "u":
			m = metrics.Update
		case "d":
			m = metrics.Delete
		default:
			continue
		}
		val, loaded := m.LoadOrStore(entry.Namespace, newValue)
		// if we loaded it, just incrememnt
		if loaded {
			atomic.AddInt64(val.(*int64), 1)
		} else { // if not, then we just take the address
			tmp := int64(1)
			newValue = &tmp
		}
	}
	if err := iter.Close(); err != nil {
		panic(err)
	}

	metrics.Print()
}
