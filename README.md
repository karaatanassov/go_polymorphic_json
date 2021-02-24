# Polymorphic JSON with Go

In this article and example project we look at the problem of handling
polymorphic JSON structures in the Golang programming language.

How to deal with polymorphic serialized JSON data is a common problem discussed
in the Go community. This can be found in many existing APIs and is sometimes
needed for new APIs. Unfortunately the default Go JSON code does not handle
polymorphism out of the box and hence the motivation for this work.

First we need to understand what is polymorphism and how it maps to Go language.

As per Wikipedia [Polymorphism](https://en.wikipedia.org/wiki/Polymorphism_(computer_science))
is

> In programming languages and type theory, polymorphism is the provision of a
> single interface to entities of different types

In the JSON serialization format Polymorphism is used to select among one of
multiple data structure types. Those data types may be a flat enumeration or
represent a  hierarchical system such as the classification of living things or
exception hierarchy.

In this work we would detail the more complex case of deep hierarchical
structure. The same approach is applicable to a flat enumeration of data
structure types.

One example of polymorphic JSON API is [RFC 7946 Geo JSON](https://tools.ietf.org/html/rfc7946).
```json
{
    "type": "Point",
    "coordinates": [102.0, 0.5]
}
// OR
{
    "type": "Polygon",
    "coordinates": [
        [
            [100.0, 0.0],
            [101.0, 0.0],
            [101.0, 1.0],
            [100.0, 1.0],
            [100.0, 0.0]
        ]
    ]
}
```
In this specification we would look into a simple error model that would be a
common chore for API authors and consumers.
```json
{
    "errors": [
        {
            "Kind" : "Fault",
            "Message": "Something went wrong.",
            "Cause": {  "Kind": "Fault", "Message": "Missing file" }
        },{
            "Kind" : "RuntimeFault",
            "Message": "Unexpected error"
        },{
            "Kind" : "Not Found",
            "Message": "The cat Lucie is missing.",
            "Obj": "Lucie",
            "ObjKind": "Cat"
        }
    ]
}
```
## Data Model in Go
The first task is to define the data objects in Go. This is rather
straightforward job:
```go
type Fault struct {
	Message string
	Cause   interfaces.Fault
}
type RuntimeFault struct {
	Fault
}
type NotFound struct {
	RuntimeFault
	ObjKind string
	Obj     string
}
```
Note how embedding can be used to construct more complex types without repeating
the definitions. This is similar to the `allOf` construct in JSON schema.

We have now defined a schema of three objects that extend each other. However we
do not get Polymorphic behavior. In other words [the following results in error](https://play.golang.org/p/q55Z6zdA63Y)
```go
func printFault(f *Fault) {
	fmt.Println("The fault message is:", f.Message);
}

func main() {
	printFault(&RuntimeFault{ Fault { Message: "test" } })
}
```
We cannot use `RuntimeFault` or `NotFound` in place where `Fault` is needed.
Let's see how to solve this.

## Polymorphism in Go

The Go language provides us with the concept of [interfaces](https://golang.org/doc/effective_go#interfaces})
implemented through methods on various types. It is with interfaces that
different data structures can exhibit polymorphic behavior.

Applying interfaces to our example above produces [working code](https://play.golang.org/p/UckCCDz4wYT)
```go
type FaultInterface interface {
	GetMessage() string
}

func (f *Fault) GetMessage() string {
	return f.Message
}
func printFault(f FaultInterface) {
	fmt.Println("The fault message is:", f.GetMessage())
}
func main() {
	printFault(&RuntimeFault{Fault{Message: "test"}})
}
```
In similar way an interface can be used as return type, member of another
structure or element type of a slice.

In addition to data model we would need to define interfaces for each of
out type to have polymorphic behavior.

I prefer to keep things tidy and kept the models and their implementations in
the `models` package. While the interface definitions reside in the `interface`
package. Thus we get the following definitions

```go
type Fault interface {
	GetMessage() string
	SetMessage(string)
	GetCause() Fault
	SetCause(Fault)
}
type RuntimeFault interface {
	Fault
}
type NotFound interface {
	RuntimeFault
	GetObjKind() string
	SetObjKind(string)
	GetObj() string
	SetObj(string)
}
```
Again note that our code is [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself).
There is no need to include members of generic abstractions into
specialized ones. Instead Go provides us with embedding to reuse existing code.

Of course we will need to implement the methods in the `models` package at
least once on the most generic type they appear on. We could also implement some
methods more than once on more specific types if there is specific business
logic to include at that level of abstraction.

## Rendering JSON

Rendering JSON from Go structure even referred through an interface is a
straightforward task. In our case we want to add the class name or discriminator
in the JSON payload for the recipients to know what type of fault they are
receiving. The following simple pattern works in Go:
```go
func (nfo *NotFound) MarshalJSON() ([]byte, error) {
	type marshalable NotFound
	return json.Marshal(struct {
		Kind string
		marshalable
	}{
		Kind:        "NotFound",
		marshalable: marshalable(*nfo),
	})
}
```
Here is what this method does:

1. It declares a method for marshalling `NotFound` objects. Go discovers it
dynamically by checking if the object implements `encoding/json.Unmarshaler`.
1. A new type is declared that has no marshaling logic - `marshalable` this is
needed to evade recursion
1. An anonymous struct type is declared that provides  additional field for the
discriminator.
1. An object of the new type is created that embeds the current `NotFound`
instance data and adds the `NotFound` value for the `Kind` field.

It is good idea to have the field `Kind` on top. This way Go will render it
first on the wire and clients will consume the output more efficiently.

After we add these methods to our error types the serialization starts to work.

```go
func main() {
	var fi FaultInterface = &RuntimeFault{Fault{Message: "test"}}
	data, err := json.Marshal(fi)
	if err != nil {
		fmt.Println("Cannot write JSON", err)
		return
	}
	fmt.Println("JSON:", string(data))
}
```
prints:
```json
{"Kind":"RuntimeFault","Message":"test","Cause":null}
```
Here is the [evolving example](https://play.golang.org/p/1A0uMnBfcYV).


## Reading JSON

Reading back the JSON values is slightly more convoluted. The biggest obstacle
is that Go has no easy way to associate functionality to an interface type. Thus
it is not possible to implement unmarshal method on our interface types. Instead
we need to take care of unmarshaling interface one layer above them. There are
several scenarios where we can see the need to unmarshal interfaces:

1. Read an top level JSON object into an interface
2. Read member field of an object into an interface
3. Read array member into an interface

Reading a JSON document like the one we have above into na interface is simplest.
We need a function that receives a `[]byte` and returns `interface.Fault`. Let
us call this `UnmarshalFault`:

```go
func UnmarshalFault(in []byte) (interfaces.Fault, error) {
	d := &struct {
		Kind string
	}{}
	// Double pointer detects null values
	err := json.Unmarshal(in, &d)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, nil
	}
	kind := d.Kind

	var res interfaces.Fault
	switch kind {
	case "NotFound":
		res = &NotFound{}
	case "RuntimeFault":
		res = &RuntimeFault{}
	default: // Error on default or try to use base type?
		res = &Fault{}
	}
	json.Unmarshal(in, res)

	return res, nil
}
```
The basic operation of the function is as follows:

1. Scan the incoming JSON bytes for the discriminator field `Kind`
2. Depending on the type the proper struct type is initialized and read from the
wire

This pretty much does the trick for basic interface unmarshaling.

Here is the [sample getting bigger](https://play.golang.org/p/6k2nFYd-dAN)

### Reading JSON fields

The last challenge is reading interfaces when they are members of an object or
an array.

To achieve this each object that has field of polymorphic type will need custom unmarshaling logic that reads into a type whose fields are all struct types. Later it will convert into the proper type.

For [example](models/field_test.go)

```go
type Container struct {
	FaultField interfaces.Fault
}

var _ json.Unmarshaler = &Container{}

func (c *Container) UnmarshalJSON(in []byte) error {
	// Deserialize into temp object of utility class
	temp := struct {
		FaultField FaultField
	}{}
	err := json.Unmarshal(in, &temp)
	if err != nil {
		return err
	}
	// Re-assign all fields to the container unwrapping the util classes
	c.FaultField = temp.FaultField.Fault
	return nil
}
```
Our goal in the code above is to unmarshal `Container` structure. To achieve it we first unmarshal into  spacial structure that has `FaultField` field instead of the `Fault` interface type. That `FaultField` in its unmarshal method will use the `UnmarshalFault` we introduced earlier:

```go
type FaultField struct {
	interfaces.Fault
}

func (ff *FaultField) UnmarshalJSON(in []byte) error {
	var err error
	ff.Fault, err = UnmarshalFault(in)
	return err
}
```

In this way we can now easily read our `Container` object form a JSON file and get the correct `Fault` instance.

```go
	c := Container{}
	err = json.Unmarshal(b, &c)
```

### Reading JSON arrays

Working with arrays is similar to object fields. We will need to first fead into special slice of `FaultField` elements and then build the proper slice of `Fault` interface implementations.

You ca see the code in [array_test.go](models/array_test.go)

```go
type ArrayContainer struct {
	Faults []interfaces.Fault
}

var _ json.Unmarshaler = &ArrayContainer{}

func (c *ArrayContainer) UnmarshalJSON(in []byte) error {
	// Deserialize into temp object of utility class
	temp := struct {
		Faults []FaultField
	}{}
	err := json.Unmarshal(in, &temp)
	if err != nil {
		return err
	}
	// Re-assign all fields to the container unwrapping the util classes
	c.Faults = ToFaultsArray(temp.Faults)
	return nil
}
```
To make this work we will need to add one more utility to our `Fault` implementation:

```go
func ToFaultsArray(faults []FaultField) []interfaces.Fault {
	var items []interfaces.Fault
	for _, tmp := range faults {
		items = append(items, tmp.Fault)
	}
	return items
}
```

### What about the Fault cause

The `Fault` object contains a `Cause` field that links to a related `Fault` object. We need to change the field from pointer to struct type to interface as to allow polymorphic behavior. Further we need to add `UnmarshalJSON` operations to all types in the hierarchy as to correctly deserialize. The methods on every type need to care about all fields and cannot delegate to the embedded types.

```go
func (nfo *NotFound) UnmarshalJSON(in []byte) error {
	pxy := &struct {
		Message string
		Cause   FaultField
		ObjKind string
		Obj     string
	}{}
	err := json.Unmarshal(in, pxy)
	if err != nil {
		return err
	}
	nfo.Message = pxy.Message
	nfo.Cause = pxy.Cause.Fault
	nfo.Obj = pxy.Obj
	nfo.ObjKind = pxy.ObjKind
	return nil
}
```

### How about type safety

Using deep hierarchies in Go may be problematic as Go interface conversions are based on the presence methods. In our example any `Fault` object works well as `RuntimeFault`. This is different from the behavior of C++ and Java where objects with no members are used to provide type safety and classification. This functionality can be emulated by adding synthetic member functions for each interface type in a hierarchy. For example the following prevents a `Fault` to be converted to `RuntimeFault`

```go
type RuntimeFault interface {
	Fault
    ZzRuntimeFault()
}

// ZzRuntimeFault is a marker
func (rf *RuntimeFault) ZzRuntimeFault() {
}
```

### How is unmarshaling working in the hierarchy?

 `RuntimeFault` fields are polymorphic and are not root of a hierarchy. One option is to have `UnmarshalRuntimeFault` bank on the root objects like `Fault` and invoke `UnmarshalFault` function.  This will save code size as in a hug hierarchy the unmashal methods may become too big.

```go
 func UnmarshalRuntimeFault(in []byte) (interfaces.RuntimeFault, error) {
	fault, err := UnmarshalFault(in)
	if err != nil {
		return nil, err
	}
	if runtimeFault, ok := fault.(interfaces.RuntimeFault); ok {
		return runtimeFault, nil
	}
	return nil, fmt.Errorf("Cannot unmarshal RuntimeFault %v", fault)
}
```

## Conclusion and next steps

This article and sample code illustrate the basic handling of polymorphic JSON in Go. We see that out of the box support is lacking but a little bit of creativity helps us get near native experience with polymorphic unmarshaling in Go.

The article and sample code leave out some details.

One area to discuss is how this work can be mapped onto polymorphic OpenAPI schema. The combination of `allOf` and `discriminator` constructs used in OpenAPI generator with Java provides good base.

The performance of the switch statement in `UnmarshalFault` on a string when thousands of classes exist in a hierarchy may require optimized implementation. For example use of state machine that iterates the characters to discern different options and assert valid sequences.

## References

This article builds on a number of previous implementations. Some of those are
listed below:

* [StackOverflow: Polymorphic JSON unmarshalling of embedded structs](https://stackoverflow.com/questions/44380095/polymorphic-json-unmarshalling-of-embedded-structs)
* [unmarshal_interface.go](https://gist.github.com/tkrajina/aec8d1b15b088c20f0df4afcd5f0c511)
* [go-swagger](https://github.com/go-swagger/go-swagger)
* [JSON polymorphism in Go](https://alexkappa.medium.com/json-polymorphism-in-go-4cade1e58ed1)
* [StackOverflow: Can I use MarshalJSON to add arbitrary fields to a json encoding in golang?](https://stackoverflow.com/questions/23045884/can-i-use-marshaljson-to-add-arbitrary-fields-to-a-json-encoding-in-golang)