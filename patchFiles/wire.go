//go:build wireinject
// +build wireinject

package patchFiles

import (
	"auto-patch-system/patchFiles/controller"
	"auto-patch-system/patchFiles/repository"
	"auto-patch-system/patchFiles/service"
	"github.com/google/wire"
)

func InitPatchFilesController() controller.PatchFileController {
	wire.Build(
		repository.NewPatchFileRepository,
		service.NewPatchFileService,
		controller.NewPatchFileController,
	)
	return controller.PatchFileController{}
}
