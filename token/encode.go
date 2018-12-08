package token

const (
	tokenTag = "'s token"
)

// Encode :
func Encode(id string) (string, error) {

	// TODO?
	return id + tokenTag, nil
}
