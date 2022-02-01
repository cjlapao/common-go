package mongodb

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var tOID = reflect.TypeOf(primitive.ObjectID{})
var tDateTime = reflect.TypeOf(primitive.DateTime(0))
var tTime = reflect.TypeOf(time.Now())

func objectIDEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tOID {
		return bsoncodec.ValueEncoderError{Name: "ObjectIDEncodeValue", Types: []reflect.Type{tOID}, Received: val}
	}

	s := val.Interface().(primitive.ObjectID).Hex()
	return vw.WriteString(s)
}

func dateTimeEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	const jDateFormat = "2006-01-02T15:04:05.999Z"
	if !val.IsValid() || val.Type() != tDateTime {
		return bsoncodec.ValueEncoderError{Name: "DateTimeEncodeValue", Types: []reflect.Type{tDateTime}, Received: val}
	}

	ints := val.Int()
	t := time.Unix(0, ints*1000000).UTC()

	return vw.WriteString(t.Format(jDateFormat))
}

func timeEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	const jDateFormat = "2006-01-02T15:04:05.999Z"
	if !val.IsValid() || val.Type() != tTime {
		return bsoncodec.ValueEncoderError{Name: "DateTimeEncodeValue", Types: []reflect.Type{tDateTime}, Received: val}
	}

	time := val.Interface().(time.Time)

	return vw.WriteString(time.Format(jDateFormat))
}

func createCustomRegistry() *bsoncodec.RegistryBuilder {
	var primitiveCodecs bson.PrimitiveCodecs
	rb := bsoncodec.NewRegistryBuilder()
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	rb.RegisterTypeEncoder(tDateTime, bsoncodec.ValueEncoderFunc(dateTimeEncodeValue))
	rb.RegisterTypeEncoder(tTime, bsoncodec.ValueEncoderFunc(timeEncodeValue))
	rb.RegisterEncoder(tOID, bsoncodec.ValueEncoderFunc(objectIDEncodeValue))
	primitiveCodecs.RegisterPrimitiveCodecs(rb)
	return rb
}
