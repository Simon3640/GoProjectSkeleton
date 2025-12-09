package azure

import (
	"github.com/simon3640/goprojectskeleton/src/adapters"
)

func NewHTTPAdapter() adapters.Adapter {
	return adapters.NewHTTPAdapter()
}

func GetDefaultAdapter() adapters.Adapter {
	return adapters.GetDefaultAdapter()
}

func CreateAdapter(adapterType adapters.AdapterType) adapters.Adapter {
	factory := adapters.NewFactory()
	return factory.CreateAdapter(adapterType)
}
