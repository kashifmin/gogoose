package main

import (
	"flag"
	"io"
	"os"
	"os/exec"
	"text/template"

	"github.com/joncalhoun/pipe"
)

type data struct {
	StructType   string
	StructImport string
	Name         string
}

func main() {
	var d data
	flag.StringVar(&d.StructType, "type", "gogoose.User", "The struct type for model being generated")
	flag.StringVar(&d.StructImport, "import", "github.com/kashifmin/gogoose", "The struct import path for model being generated")
	flag.StringVar(&d.Name, "name", "User", "The prefix name for model structs")
	flag.Parse()

	t := template.Must(template.New("gogoose").Parse(queueTemplate))
	rc, wc, _ := pipe.Commands(
		exec.Command("gofmt"),
		exec.Command("goimports"),
	)
	t.Execute(wc, d)
	wc.Close()
	io.Copy(os.Stdout, rc)
}

var queueTemplate = `
package gen

import (
	"context"
	"errors"
	"reflect"

	"{{.StructImport}}"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type {{.Name}}Model struct {
	dbColl *mongo.Collection
}

type {{.Name}}Document struct {
	dbColl *mongo.Collection
	raw    *{{.StructType}}
}

func (userDocument *{{.Name}}Document) Save(ctx context.Context) error {
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

func (userModel *{{.Name}}Model) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*{{.Name}}Document, error) {
	res := userModel.dbColl.FindOne(ctx, filter)
	user := &{{.StructType}}{}
	err := gogoose.DecodeSingleResult(res, user)
	if err != nil {
		return nil, err
	}
	return New{{.Name}}Document(user, userModel.dbColl), nil
}

func (userModel *{{.Name}}Model) New(user *{{.StructType}}) *{{.Name}}Document {
	return &{{.Name}}Document{raw: user, dbColl: userModel.dbColl}
}

// New{{.Name}}Document ...
func New{{.Name}}Document(user *{{.StructType}}, coll *mongo.Collection) *{{.Name}}Document {
	return &{{.Name}}Document{
		dbColl: coll,
		raw:    user,
	}
}

func New{{.Name}}Model(collection *mongo.Collection) *{{.Name}}Model {
	return &{{.Name}}Model{dbColl: collection}
}

`
