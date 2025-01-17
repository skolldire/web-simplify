package app_builder

// The App interface defines the methods that the application should expose.
type App interface {
	Run() error
}

// The Builder interface defines how the different parts of the application will be built.
type Builder interface {
	LoadConfig() Builder
	InitRepositories() Builder
	InitUseCases() Builder
	InitHandlers() Builder
	InitRoutes() Builder
	Build() App
}
