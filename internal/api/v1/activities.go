package v1

import (
	"invoice-backend/internal/repositories/activities"
)

type ActivitiesHandler struct {
	activitiesRepo activities.Repository
}

func NewActivitiesHandler(activitiesRepo activities.Repository) *ActivitiesHandler {
	return &ActivitiesHandler{
		activitiesRepo: activitiesRepo,
	}
}
