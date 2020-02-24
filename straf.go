package straf

import (
	"errors"
	"reflect"

	"github.com/graphql-go/graphql"
)

// GetGraphQLObject Converts struct into graphql object
func GetGraphQLObject(object interface{}) (*graphql.Object, error) {
	objectType := reflect.TypeOf(object)
	fields := ConvertStruct(objectType)

	output := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   objectType.Name(),
			Fields: fields,
		},
	)

	return output, nil
}

// ConvertStructToObject converts simple struct to graphql object
func ConvertStructToObject(
	objectType reflect.Type) *graphql.Object {

	fields := ConvertStruct(objectType)

	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   objectType.Name(),
			Fields: fields,
		},
	)
}

// ConvertStruct converts struct to graphql fields
func ConvertStruct(objectType reflect.Type) graphql.Fields {
	fields := graphql.Fields{}

	for i := 0; i < objectType.NumField(); i++ {
		currentField := objectType.Field(i)
		if GetTagValue(currentField, "exclude") != "true" {
			fieldType := GetFieldType(currentField)
			fields[currentField.Name] = &graphql.Field{
				Name:              currentField.Name,
				Type:              fieldType,
				DeprecationReason: GetTagValue(currentField, "deprecationReason"),
				Description:       GetTagValue(currentField, "description"),
			}
		}
	}

	return fields
}

var fieldTypeCache = map[string]graphql.Output{}

// GetFieldType Converts object to a graphQL field type
func GetFieldType(object reflect.StructField) graphql.Output {
	cachedOutput := fieldTypeCache[object.Name]
	if cachedOutput != nil {
		return cachedOutput
	}
	
	isID, ok := object.Tag.Lookup("unique")
	if isID == "true" && ok {
		return graphql.ID
	}

	objectType := object.Type
	if objectType.Kind() == reflect.Struct {
		return ConvertStructToObject(objectType)

	} else if objectType.Kind() == reflect.Slice &&
		objectType.Elem().Kind() == reflect.Struct {

		elemType := ConvertStructToObject(objectType.Elem())
		return graphql.NewList(elemType)

	} else if objectType.Kind() == reflect.Slice {
		elemType, _ := ConvertSimpleType(objectType.Elem())
		return graphql.NewList(elemType)
	}

	output, _ := ConvertSimpleType(objectType)
	
	cachedOutput[object.Name] = output
	return output
}

// ConvertSimpleType converts simple type  to graphql field
func ConvertSimpleType(objectType reflect.Type) (*graphql.Scalar, error) {

	typeMap := map[reflect.Kind]*graphql.Scalar{
		reflect.String:  graphql.String,
		reflect.Bool:    graphql.Boolean,
		reflect.Int:     graphql.Int,
		reflect.Int8:    graphql.Int,
		reflect.Int16:   graphql.Int,
		reflect.Int32:   graphql.Int,
		reflect.Int64:   graphql.Int,
		reflect.Float32: graphql.Float,
		reflect.Float64: graphql.Float,
	}

	graphqlType, ok := typeMap[objectType.Kind()]

	if !ok {
		return &graphql.Scalar{}, errors.New("Invalid Type")
	}

	return graphqlType, nil
}

// getTagValue returns tag value of a struct
func GetTagValue(objectType reflect.StructField, tagName string) string {
	tag := objectType.Tag
	value, ok := tag.Lookup(tagName)
	if !ok {
		return ""
	}
	return value
}
