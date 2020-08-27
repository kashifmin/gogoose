# gogoose
Mongoose for Go. Generates mongoose like models and document types based on your type struct.

# Example
```go
func main() {
	db := NewMongoClient()
	userModel := gen.NewUserModel(db.Collection("kuser"))
	oid := primitive.NewObjectID()
	doc := userModel.New(&gogoose.User{Name: "Kashif", Age: 23, ID: &oid})
	err := doc.Save(context.Background())
	if err != nil {
		panic(err)
	}
}
```

# TODO
- [ ] Concept with `.Save()` and `.FindOne()`
- [ ] Implement `Find*`, `Update*`, `Delete*`
- [ ] Implement difftracker for `.Save()`
- [ ] Implement aggregations
- [ ] Create a sample service demonstrating the uses