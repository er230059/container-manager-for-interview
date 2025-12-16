package domain

// Greeting represents the core domain model.
type Greeting struct {
	Message string
}

// GreetingRepository is an interface for fetching greetings.
// It's defined in the domain layer, and implemented in the infrastructure layer.
type GreetingRepository interface {
	GetGreeting() (*Greeting, error)
}
