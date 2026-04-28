package enum

type ServerEnv int

const (
	ServerEnvLocal ServerEnv = iota
	ServerEnvDev
	ServerEnvProd
	serverEnvCount
)

const (
	serverEnvLocalStr = "local"
	serverEnvDevStr   = "dev"
	serverEnvProdStr  = "prod"
)

var serverEnvStrList = [serverEnvCount]string{
	serverEnvLocalStr,
	serverEnvDevStr,
	serverEnvProdStr,
}

func (s ServerEnv) String() string {
	return serverEnvStrList[s]
}
