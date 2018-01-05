package goddbench

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Data struct {
	Bool   bool      `json:"bool" msg:"bool"`
	Int    int       `json:"int" msg:"int"`
	Float  float64   `json:"float" msg:"float"`
	String string    `json:"string" msg:"string"`
	Time   time.Time `json:"time" msg:"time"`
}

type DataList []Data

func GenN(n int) []Data {
	list := make([]Data, n)
	for i := 0; i < n; i++ {
		b := make([]byte, 8+rand.Intn(64))
		rand.Read(b)
		list[i] = Data{
			rand.Int()&1 == 0,
			rand.Int(),
			rand.NormFloat64(),
			fmt.Sprintf("%X", b),
			time.Unix(rand.Int63n(1<<31), rand.Int63n(1000000000)).UTC(),
		}
	}
	return list
}

func Hash(obj interface{}) string {
	hash := sha1.New()
	enc := json.NewEncoder(hash)
	enc.Encode(obj)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func TimeParse(s string) time.Time {
	t, _ := time.Parse(time.RFC3339Nano, s)
	return t
}
