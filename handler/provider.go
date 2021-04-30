package handler

import "github.com/DoNewsCode/core/di"

func Provides() di.Deps {
	return di.Deps{
		newDefaultHandler,
		newFooHandler,
	}
}
