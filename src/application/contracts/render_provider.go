package contracts

import application_errors "gormgoskeleton/src/application/shared/errors"

type IRenderProvider[D any] interface {
	Render(template string, data D) (string, *application_errors.ApplicationError)
}
