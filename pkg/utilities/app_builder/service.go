package app_builder

// Builder ´director´
// The "Director" is the one who decides in what order the builder methods should be called to build the object.
// In this case, the director could simply be a function that accepts a Builder and uses it to build the application.
func Apply(builder Builder) App {
	return builder.
		LoadConfig().
		InitRepositories().
		InitUseCases().
		InitHandlers().
		InitRoutes().
		Build()
}
