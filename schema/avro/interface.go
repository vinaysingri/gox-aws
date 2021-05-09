package avro

import (
	"github.com/devlibx/gox-base"
	errors2 "github.com/devlibx/gox-base/errors"
	"github.com/fatih/structs"
	"github.com/linkedin/goavro/v2"
	_ "github.com/linkedin/goavro/v2"
	"io/ioutil"
)

//go:generate mockgen -source=interface.go -destination=mocks/mock_interface.go -package=mock_gox_aws_avro
type SchemaEngine interface {
	ToAvro(data interface{}) ([]byte, error)
	FromAvro(data []byte) (gox.StringObjectMap, error)
}

type schemaEngine struct {
	codec *goavro.Codec
}

func (e *schemaEngine) ToAvro(data interface{}) ([]byte, error) {
	if _, ok := data.(map[string]interface{}); ok {
		return e.codec.BinaryFromNative(nil, data)
	} else {
		return e.codec.BinaryFromNative(nil, structs.Map(data))
	}
}

func (e *schemaEngine) FromAvro(data []byte) (gox.StringObjectMap, error) {
	obj, _, err := e.codec.NativeFromBinary(data)
	if err != nil {
		return nil, errors2.Wrap(err, "failed to parse avro bytes")
	}

	if _obj, ok := obj.(map[string]interface{}); !ok {
		return nil, errors2.Wrap(err, "failed to parse avro bytes to map")
	} else {
		return _obj, nil
	}
}

func NewAvroSchemaEngine(schema string) (SchemaEngine, error) {
	codec, err := goavro.NewCodec(schema)
	if err != nil {
		return nil, errors2.Wrap(err, "input avro schema is not valid")
	}
	return &schemaEngine{codec: codec}, nil
}

func NewAvroSchemaEngineWithFile(schemaFile string) (SchemaEngine, error) {
	data, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		return nil, errors2.Wrap(err, "file not found: %s", schemaFile)
	}
	return NewAvroSchemaEngine(string(data))
}
