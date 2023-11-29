package helpers

type CreatePlayerRequestBody struct {
	Name     string `json:"name,omitempty"`
	Grade    string `json:"grade,omitempty"`
	Position string `json:"position,omitempty"`
}

func (player *CreatePlayerRequestBody) IsCreatePlayerRequestBodyValid() bool {
	if player.Name != "" && IsGradeValid(player.Grade) && IsPositionValid(player.Position) {
		return true
	}
	return false
}

func IsGradeValid(grade string) bool {
	if grade == "" {
		return false
	}
	asciiValueOfGrade := []rune(grade)
	if len(asciiValueOfGrade) > 1 {
		return false
	}
	if asciiValueOfGrade[0] < 65 || asciiValueOfGrade[0] > 67 {
		return false
	}
	return true
}

func IsPositionValid(position string) bool {
	if position == "" {
		return false
	}
	validPositions := []string{
		"ST", "CF", "LW", "RW", "RS", "LS",
		"CM", "RCM", "LCM", "CAM", "LM", "RM",
		"CDM", "CB", "RCB", "LCB", "RB", "LB", "GK",
	}
	isValidFlag := false
	for _, validPosition := range validPositions {
		if position == validPosition {
			isValidFlag = true
			break
		}
	}
	return isValidFlag
}
