package contractsproviders

import application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"

type IRendererProvider[D any] interface {
	Render(template string, data D) (string, *application_errors.ApplicationError)
}
