package interfaces

// RuntimeFault represents all runtime faults that can be thrown
// To be generated
type RuntimeFault interface {
	Fault
	// ZzRuntimeFault disallows converting Fault struct to RuntimeFault interface
	ZzRuntimeFault()
}
