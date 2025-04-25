//go:build wireinject
// +build wireinject

package reservations

import (
	"auto-patch-system/notification"
	"auto-patch-system/reservations/controller"
	"auto-patch-system/reservations/repository"
	reservations "auto-patch-system/reservations/service"
	"github.com/google/wire"
)

func InitReservationController() controller.ReservationController {
	wire.Build(
		repository.NewReservationRepository,
		reservations.NewReservationService,
		notification.NewSlackNotificationService,
		controller.NewReservationController,
	)
	return controller.ReservationController{}
}
