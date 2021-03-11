通过嵌入匿名结构体的方式来实现动态类型


期望

```go
type view struct {
    FieldName inject.PathParam(IsPhoneNumber, IsNULL)
}
```

不能接受的方案

```go
type view struct {
    FieldName inject.PathParam `valid:"IsPhoneNumber, IsNULL"`
}
```
    
目前的方案

```go
type view struct {
    FieldName struct { 
        inject.PathParam
        schema.IsPhoneNumber
        schema.IsNULL
    }
}
```
    
其中 schema.IsPhoneNumber 是空结构体

其中 Field 可以缩写成 

`FieldName struct { inject.PathParam; schema.IsPhoneNumber; schema.IsNULL }`