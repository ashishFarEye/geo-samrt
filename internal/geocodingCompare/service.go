package geocodingCompare

import (
	"context"

	"geo-smart/api/gen/models"
)

type Service interface {
	GeocodeCompare(context.Context, map[string]models.Address, map[string]models.Success, string) error
}
