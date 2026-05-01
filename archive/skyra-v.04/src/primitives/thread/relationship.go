package thread

type RelationshipKey struct {
	A, B string
}

func NewRelationshipKey(a, b string) RelationshipKey {
	if a < b {
		return RelationshipKey{A: a, B: b}
	}
	return RelationshipKey{A: b, B: a}
}
