package reality

import (
	"math"
	"sort"
)

type VecIndex struct {
	Vectors map[string][]float64 `json:"vectors"`
}

func NewVecIndex() *VecIndex {
	return &VecIndex{Vectors: make(map[string][]float64)}
}

func (v *VecIndex) Add(id string, vec []float64) {
	v.Vectors[id] = vec
}

func (v *VecIndex) Remove(id string) {
	delete(v.Vectors, id)
}

type VecResult struct {
	ID    string
	Score float64
}

func (v *VecIndex) Search(query []float64, k int) []VecResult {
	var results []VecResult
	for id, vec := range v.Vectors {
		score := cosineSimilarity(query, vec)
		results = append(results, VecResult{ID: id, Score: score})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	if len(results) > k {
		results = results[:k]
	}
	return results
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
