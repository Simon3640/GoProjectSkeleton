package status

type ApplicationStatusEnum string

type ApplicationStatusDictionary map[ApplicationStatusEnum]string

const (
	Success                   ApplicationStatusEnum = "OK"
	Updated                   ApplicationStatusEnum = "UD"
	Created                   ApplicationStatusEnum = "CD"
	PartialContent            ApplicationStatusEnum = "PA_CO"
	InvalidInput              ApplicationStatusEnum = "BA_RE"
	Unauthorized              ApplicationStatusEnum = "UNAU"
	NotFound                  ApplicationStatusEnum = "NO_FO"
	Conflict                  ApplicationStatusEnum = "CONF"
	InternalError             ApplicationStatusEnum = "IN_ER"
	NotImplemented            ApplicationStatusEnum = "NO_IM"
	ProviderError             ApplicationStatusEnum = "PR_ER"
	ChatProviderError         ApplicationStatusEnum = "CH_PR_ER"
	ProviderEmptyResponse     ApplicationStatusEnum = "PR_EM_RE"
	ProviderEmptyCacheContext ApplicationStatusEnum = "PR_EM_CA_CO"
)
