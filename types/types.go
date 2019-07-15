package types

// Status holds the results of the build
type Status uint

const (
	status Status = iota
	Pending
	Success
	Failure
	Error
	Target
	Log
)

// Event holds an event that has occured
type Event struct {
	Sha    string
	Status Status
	Data   []byte
}

func StatusOnly(e interface{}) bool {
	switch t := e.(type) {
	case Event:
		return t.Status < Target
	case *Event:
		return t.Status < Target
	}
	return false
}

func LogOnly(e interface{}) bool {
	switch t := e.(type) {
	case Event:
		return t.Status == Log
	case *Event:
		return t.Status == Log
	}
	return false
}

func (s Status) Description() string {
	switch s {
	case Pending:
		return "the build is pending..."
	case Success:
		return "the build was sucessful!"
	case Failure:
		return "something went wrong."
	}
	return "the build failed."
}

func (s Status) String() string {
	switch s {
	case Pending:
		return "pending"
	case Success:
		return "success"
	case Failure:
		return "failure"
	case Log:
		return "log"
	case Target:
		return "target"
	}
	return "error"
}
