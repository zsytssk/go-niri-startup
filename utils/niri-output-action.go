package utils

type OutputActionPositionSpecific struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type OutputActionPositionDetail struct {
	Automatic *Empty                        `json:"Automatic,omitempty"`
	Specific  *OutputActionPositionSpecific `json:"Specific,omitempty"`
}
type OutputActionPosition struct {
	Position *OutputActionPositionDetail `json:"position,omitempty"`
}

type OutputActionCon struct {
	Off      *Empty                `json:"Off,omitempty"`
	On       *Empty                `json:"On,omitempty"`
	Position *OutputActionPosition `json:"Position,omitempty"`
}
type OutputAction struct {
	Output string           `json:"output"`
	Action *OutputActionCon `json:"action"`
}
