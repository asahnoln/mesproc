package store

// StepStore is a store for saving messages
type Step interface {
	Save(string) error
}
