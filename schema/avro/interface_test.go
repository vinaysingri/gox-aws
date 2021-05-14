package avro

import (
	"flag"
	mockAwsAvro "github.com/devlibx/gox-aws/schema/avro/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Employee struct {
	Name string
	Age  int
}

func init() {
	ignore := ""
	flag.StringVar(&ignore, "real.sqs.queue", "false", "run all database tests for dynamo")
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

func TestAvroEngine_Versioning(t *testing.T) {

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

	// Test 1 - Schema with V1
	engine, err := NewAvroSchemaEngine(schema)
	assert.NoError(t, err)

	// Make sure we are able to convert to avro and back to object
	input := map[string]interface{}{"Name": "user_1", "Age": 10}
	data, err := engine.ToAvro(input)
	assert.NoError(t, err)
	backObject, err := engine.FromAvro(data)
	assert.NoError(t, err)
	assert.Equal(t, "user_1", backObject.StringOrEmpty("Name"))
	assert.Equal(t, 10, backObject.IntOrDefault("Age", 0))

	// Test 2 - Schema with V2 with a BAD schema change
	schemaWithBadChanges := `
		{
		   "type" : "record",
		   "namespace" : "gox.aws",
		   "name" : "Employee",
		   "fields" : [
			  { "name" : "Name" , "type" : "string" },
			  { "name" : "Age" , "type" : "string" }
		   ]
		}`
	engineWithBadSchema, err := NewAvroSchemaEngine(schemaWithBadChanges)
	assert.NoError(t, err)
	backObject, err = engineWithBadSchema.FromAvro(data)
	assert.Error(t, err)

	// Test 3 - Schema with V3 with new attribute added
	schemaWithNewAttribute := `
		{
		   "type" : "record",
		   "namespace" : "gox.aws",
		   "name" : "Employee",
		   "fields" : [
			  { "name" : "Name" , "type" : "string" },
			  { "name" : "Age" , "type" : "int" },
              { "name" : "Dob" , "type" : "string" }
		   ]
		}`
	input = map[string]interface{}{"Name": "user_1", "Age": 10, "Dob": "1700"}
	enginesWithNewAttribute, err := NewAvroSchemaEngine(schemaWithNewAttribute)
	assert.NoError(t, err)
	dataWithDob, err := enginesWithNewAttribute.ToAvro(input)
	assert.NoError(t, err)
	backObject, err = enginesWithNewAttribute.FromAvro(dataWithDob)
	assert.NoError(t, err)
	assert.Equal(t, "user_1", backObject.StringOrEmpty("Name"))
	assert.Equal(t, 10, backObject.IntOrDefault("Age", 0))
	assert.Equal(t, "1700", backObject.StringOrEmpty("Dob"))

	// Test 4 - See if the client with Old schema can work with the new data or now
	backObject, err = engine.FromAvro(dataWithDob)
	assert.NoError(t, err)
	assert.Equal(t, "user_1", backObject.StringOrEmpty("Name"))
	assert.Equal(t, 10, backObject.IntOrDefault("Age", 0))

	// Test 4.1 - since we are using "engine" from old schema it should work with new data
	// but it should not see the new values
	assert.Equal(t, "", backObject.StringOrEmpty("Dob"))
}

func TestAvroEngineWithFile(t *testing.T) {
	engine, err := NewAvroSchemaEngineWithFile("test_avro_schema.json")
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
