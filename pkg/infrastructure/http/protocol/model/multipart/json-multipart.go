package multipart

import (
	"bytes"
	"encoding/json"

	"github.com/ldez/mimetype"
	http "github.com/nullexp/finman-api-gateway/pkg/infrastructure/http/protocol"
)

type jsonMultipart struct {
	object   interface{}
	r        *bytes.Reader
	partName string
}

func (j *jsonMultipart) Read(p []byte) (n int, err error) {
	if j.r == nil {
		data, _ := json.Marshal(j.object)
		j.r = bytes.NewReader(data)
	}
	return j.r.Read(p)
}

func (j *jsonMultipart) GetPartName() string {
	return j.partName
}

func (j *jsonMultipart) GetMimeType() string {
	return mimetype.ApplicationJSON
}
func (jsonMultipart) Close() error { return nil }

func NewJsonPart(obj interface{}, partName string) http.Multipart {
	return &jsonMultipart{object: obj, partName: partName}
}
