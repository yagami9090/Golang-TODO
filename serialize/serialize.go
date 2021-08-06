package serialize

import (
	"encoding/json"
	"io"

	msgpack "github.com/vmihailenco/msgpack/v5"
)

type JSONSerializer struct{}

func NewJSONSerializer() JSONSerializer {
	return JSONSerializer{}
}

func (JSONSerializer) Decode(body io.ReadCloser, v interface{}) error {
	return json.NewDecoder(body).Decode(v)
}

type MsgPkgSerializer struct{}

func NewMsgPkgSerializer() MsgPkgSerializer {
	return MsgPkgSerializer{}
}

func (MsgPkgSerializer) Decode(body io.ReadCloser, v interface{}) error {
	return msgpack.NewDecoder(body).Decode(v)
}
