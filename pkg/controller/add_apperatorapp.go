package controller

import (
	"github.com/otaviof/apperator/pkg/controller/apperatorapp"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, apperatorapp.Add)
}
