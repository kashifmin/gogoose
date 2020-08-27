package gogoose

import (
	"errors"
	"reflect"
	"strings"
	"unicode"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
)

func DecodeSingleResult(res *mongo.SingleResult, dest interface{}) error {
	if res == nil {
		return errors.New("decodeSingleResult: nil SingleResult passed")
	}
	raw, err := res.DecodeBytes()
	if err != nil {
		return err
	}
	decodeContext := bsoncodec.DecodeContext{
		Registry: bson.DefaultRegistry,
		Truncate: true,
	}
	return bson.UnmarshalWithContext(decodeContext, raw, dest)
}

//Lower cases first char of string
func lowerInitial(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func GetBsonName(field reflect.StructField) string {
	tag := field.Tag.Get("bson")
	tags := strings.Split(tag, ",")

	if len(tags[0]) > 0 {
		return tags[0]
	} else {
		return lowerInitial(field.Name)
	}

}
