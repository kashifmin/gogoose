package gen

import (
	"context"
	"errors"
	"reflect"

	"github.com/kashifmin/gogoose"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserModel struct {
	dbColl *mongo.Collection
}

type UserDocument struct {
	dbColl *mongo.Collection
	raw    *gogoose.User
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
		fieldName := gogoose.GetBsonName(field)
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
	user := &gogoose.User{}
	err := gogoose.DecodeSingleResult(res, user)
	if err != nil {
		return nil, err
	}
	return NewUserDocument(user, userModel.dbColl), nil
}

func (userModel *UserModel) New(user *gogoose.User) *UserDocument {
	return &UserDocument{raw: user, dbColl: userModel.dbColl}
}

// NewUserDocument ...
func NewUserDocument(user *gogoose.User, coll *mongo.Collection) *UserDocument {
	return &UserDocument{
		dbColl: coll,
		raw:    user,
	}
}

func NewUserModel(collection *mongo.Collection) *UserModel {
	return &UserModel{dbColl: collection}
}
