#### Schema

This example shows how to use a schema and convert a golang struct to Avro byte data. Then convert
it back from Avro byte data to original object. 

```go
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
```