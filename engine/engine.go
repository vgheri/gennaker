package engine

type engine struct {
	db        DeploymentRepository
	chartsDir string
}

func New(repository DeploymentRepository, savedChartsDir string) DeploymentEngine {
	return &engine{
		db:        repository,
		chartsDir: savedChartsDir,
	}
}
