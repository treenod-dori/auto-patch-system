//go:build wireinject
// +build wireinject

package reservations

import (
	"auto-patch-system/reservations/controller"
	"auto-patch-system/reservations/repository"
	"auto-patch-system/reservations/service"
	"github.com/google/wire"
)

func InitReservationController() controller.ReservationController {
	wire.Build(
		repository.NewReservationRepository,
		service.NewReservationService,
		controller.NewReservationController,
	)
	return controller.ReservationController{}
}
