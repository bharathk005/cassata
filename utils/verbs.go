package utils

type Verb string

const (
	Create Verb = "Create"
	Read   Verb = "Read"
	Update Verb = "Update"
	Delete Verb = "Delete"
)

var HTTPMethodToVerb = map[string]Verb{
	"GET":    Read,
	"POST":   Create,
	"PUT":    Update,
	"PATCH":  Update,
	"DELETE": Delete,
}
