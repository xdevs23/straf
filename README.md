[![Go Report Card](https://goreportcard.com/badge/github.com/SonicRoshan/straf)](https://goreportcard.com/report/github.com/SonicRoshan/straf) [![GoDoc](https://godoc.org/github.com/SonicRoshan/straf?status.svg)](https://godoc.org/github.com/SonicRoshan/straf) [![GoCover](https://gocover.io/_badge/github.com/SonicRoshan/straf)](https://gocover.io/github.com/SonicRoshan/straf)

# Straf
1. Convert Golang Struct To GraphQL Object On The Fly
2. Easily Create GraphQL Schemas

## Example

### Converting struct to GraphQL Object
```go
type UserExtra struct {
    Age int `description:"Age of the user"` // You can use description struct tag to add description
    Gender string `deprecationReason:"Some Reason"` // You can use deprecationReason tag to add a deprecation reason
}

type User struct {
    UserID int
    Username string `unique:"true"` // You can use unique tag to define if a field would be unique
    Extra UserExtra
    Password string `exclude:"true"` // You can use exclude tag to exclude a field
}


func main() {
    // GetGraphQLObject will convert golang struct to a graphQL object
    userType, err := straf.GetGraphQLObject(User{})

    // You can then use userType in your graphQL schema
}
```


### Using The Schema Builder
```go
type User struct {
    UserID int `isArg:"true"` // You can use isArg tag to define a field as a graphql argument
    Username string `isArg:"true"`
}

var database []User = []User{}

func main() {
    // GetGraphQLObject will convert golang struct to a graphQL object
    userType, err := straf.GetGraphQLObject(User{})

    builder := straf.NewSchemaBuilder(userType, User{})
    
    builder.AddArgumentsFromStruct(object2{}) // You can use this function to add more arguments from a struct
    builder.AddFunction("CreateUser", 
                        "Adds a user to database",
                        func(params graphql.ResolveParams) (interface{}, error)) {
                            id := params.Args["UserID"]
                            username := params.Args["Username"]
                            database = append(database, User{UserID: id, Username: Username})
                        })
    schema := builder.Schema
    // You can then use this schema
}
```

### Using Middleware In Schema Builder
```go

func middleware(function func(graphql.ResolveParams) (interface{}, error), params graphql.ResolveParams) (interface{}, error) {

    fmt.Println("This function will run as a middleware")

    return function(params)
}

func main() {
    builder := straf.NewSchemaBuilder(userType, User{}, middleware)

    builder.AddFunction("SomeFunction", 
                        "Does Something",
                        someFunction)

   // Here the middleware function would run everytime before someFunction. middleware function would act as a middleware to all functions added to schema builder.
}

```

## Author
Roshan Jignesh Mehta - sonicroshan122@gmail.com
