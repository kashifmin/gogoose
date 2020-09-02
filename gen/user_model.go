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

func (userModel *UserModel) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]*UserDocument, error) {
	cursor, err := userModel.dbColl.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	docs := make([]*UserDocument, 0, 0)
	for cursor.Next(ctx) {
		user := &gogoose.User{}
		err := bson.Unmarshal(cursor.Current, user)
		if err != nil {
			return docs, err
		}
	}
	return docs, cursor.Err()
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

func (userModel *UserModel) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) (*UserDocument, error) {
	res := userModel.dbColl.FindOneAndUpdate(ctx, filter, update, opts...)
	user := &gogoose.User{}
	err := gogoose.DecodeSingleResult(res, user)
	if err != nil {
		return nil, err
	}
	return NewUserDocument(user, userModel.dbColl), nil
}

func (userModel *UserModel) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) (*UserDocument, error) {
	res := userModel.dbColl.FindOneAndDelete(ctx, filter, opts...)
	user := &gogoose.User{}
	err := gogoose.DecodeSingleResult(res, user)
	if err != nil {
		return nil, err
	}
	return NewUserDocument(user, userModel.dbColl), nil
}

func (userModel *UserModel) FindOneAndReplace(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndReplaceOptions) (*UserDocument, error) {
	res := userModel.dbColl.FindOneAndReplace(ctx, filter, update, opts...)
	user := &gogoose.User{}
	err := gogoose.DecodeSingleResult(res, user)
	if err != nil {
		return nil, err
	}
	return NewUserDocument(user, userModel.dbColl), nil
}

func (userModel *UserModel) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return userModel.dbColl.UpdateOne(ctx, filter, update, opts...)
}

func (userModel *UserModel) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return userModel.dbColl.UpdateMany(ctx, filter, update, opts...)
}

func (userModel *UserModel) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return userModel.dbColl.DeleteOne(ctx, filter, opts...)
}

func (userModel *UserModel) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return userModel.dbColl.DeleteMany(ctx, filter, opts...)
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
