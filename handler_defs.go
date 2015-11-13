package isumm

const (
	PostAction   = "ae"
	DeleteAction = "d"
	ActionParam  = "action"
)

type handlingError struct {
	Msg  string
	Code int
}
