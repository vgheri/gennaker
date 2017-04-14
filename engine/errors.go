package engine

import "fmt"

//ErrResourceNotFound is returned when specified resource is not found
var ErrResourceNotFound error = fmt.Errorf("Resource not found")
