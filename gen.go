package goddbench

//go:generate msgp -file data.go -marshal=false -tests=false
//go:generate go run embed_gen.go -seed 100 -len 100000
//go:generate go run json_gen.go -seed 100 -len 100000
//go:generate go run msgpack_gen.go -seed 100 -len 100000
//go:generate go run gob_gen.go -seed 100 -len 100000
//go:generate go run msgp_gen.go -seed 100 -len 100000
