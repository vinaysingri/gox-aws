package avro

import (
	mockAwsAvro "github.com/devlibx/gox-aws/schema/avro/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Employee struct {
	Name string
	Age  int
}

func TestAvroEngine(t *testing.T) {
	schema := `
		{
		   "type" : "record",
		   "namespace" : "gox.aws",
		   "name" : "Employee",
		   "fields" : [
			  { "name" : "Name" , "type" : "string" },
			  { "name" : "Age" , "type" : "int" }
		   ]
		}`
	engine, err := NewAvroSchemaEngine(schema)
	assert.NoError(t, err)

	// Create Avro binary data - This data can be sent over wire e.g. over kafka
	obj := Employee{Name: "user", Age: 10}
	data, err := engine.ToAvro(obj)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// Get back data from Avro to original - We have our data back as a map
	backObject, err := engine.FromAvro(data)
	assert.NoError(t, err)
	assert.Equal(t, obj.Name, backObject.StringOrEmpty("Name"))
	assert.Equal(t, obj.Age, backObject.IntOrDefault("Age", 0))
}

func TestMockAvroEngine(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSchemaEngine := mockAwsAvro.NewMockSchemaEngine(ctrl)

	input := map[string]interface{}{"data": 10}
	avroByteData := []byte{}
	mockSchemaEngine.EXPECT().ToAvro(gomock.Eq(input)).Return(avroByteData, nil)
	mockSchemaEngine.EXPECT().FromAvro(gomock.Eq(avroByteData)).Return(input, nil)

	data, err := mockSchemaEngine.ToAvro(input)
	assert.NoError(t, err)
	output, err := mockSchemaEngine.FromAvro(data)
	assert.NoError(t, err)
	assert.Equal(t, input["data"], output["data"])
}
