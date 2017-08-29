package engine

import "fmt"

//ErrResourceNotFound is returned when specified resource is not found
var ErrResourceNotFound error = fmt.Errorf("Resource not found")

//ErrInvalidDeployment is returned when an invalid deployment is specified
var ErrInvalidDeployment error = fmt.Errorf("Invalid deployment")

// TODO replace this error with errors.Wrap(ErrInvalidDeployment, "bla bla")
//ErrInvalidPipeline is returned when an invalid pipeline is specified
var ErrInvalidPipeline error = fmt.Errorf("Invalid pipeline")

var ErrInvalidReleaseNotification error = fmt.Errorf("Invalid release notification")
