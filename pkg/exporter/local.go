package exporter

// LocalDestination is the data transfer for local
type LocalDestination struct{}

// NewLocalDestination creates a new instance of the local data transfer
func NewLocalDestination() (*LocalDestination, error) {
	return &LocalDestination{}, nil
}

// Transfer is called to transfer the image to the Destination
func (ld *LocalDestination) Transfer() (ProcessingPhase, error) {
	return ProcessingPhaseComplete, nil
}

// Close closes any senders or other open resources.
func (ld *LocalDestination) Close() {
}
