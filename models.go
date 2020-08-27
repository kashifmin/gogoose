package gogoose

import (
	"context"
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID   *primitive.ObjectID `bson:"_id"`
	Name string              `bson:"name"`
	Age  int                 `bson:"age"`
}

func decodeSingleResult(res *mongo.SingleResult, dest interface{}) error {
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

type UserModel struct {
	dbColl *mongo.Collection
}

type UserDocument struct {
	dbColl *mongo.Collection
	raw    *User
}

func (userDocument *UserDocument) Save(ctx context.Context) error {
	// TODO: implement difftracker
	if userDocument.raw.ID == nil {
		return errors.New("_id is nil")
	}
	structValueRef := reflect.ValueOf(userDocument.raw).Elem()
	structTypeRef := reflect.TypeOf(userDocument.raw).Elem()

	nFields := structTypeRef.NumField()
	fieldToUpdate := bson.M{}
	for i := 0; i < nFields; i++ {
		field := structTypeRef.Field(i)
		fieldName := GetBsonName(field)
		if fieldName == "_id" {
			continue
		}
		// fmt.Println(structValueRef.Field(i))
		fieldToUpdate[fieldName] = structValueRef.Field(i).Interface()
	}
	_, err := userDocument.dbColl.UpdateOne(ctx, bson.M{"_id": userDocument.raw.ID}, bson.M{"$set": fieldToUpdate}, options.Update().SetUpsert(true))
	return err
}

func (userModel *UserModel) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*UserDocument, error) {
	res := userModel.dbColl.FindOne(ctx, filter)
	user := &User{}
	err := decodeSingleResult(res, user)
	if err != nil {
		return nil, err
	}
	return NewUserDocument(user, userModel.dbColl), nil
}

func (userModel *UserModel) New(user *User) *UserDocument {
	return &UserDocument{raw: user, dbColl: userModel.dbColl}
}

// NewUserDocument ...
func NewUserDocument(user *User, coll *mongo.Collection) *UserDocument {
	return &UserDocument{
		dbColl: coll,
		raw:    user,
	}
}

func NewUserModel(collection *mongo.Collection) *UserModel {
	return &UserModel{dbColl: collection}
}
