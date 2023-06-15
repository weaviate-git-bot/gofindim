package math

import "math"

func dotProduct(a, b []float32) float32 {
	var sum float32
	sum = 0.0
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

func magnitude(v []float32) float32 {
	var sum float32
	sum = 0.0
	for _, val := range v {
		sum += val * val
	}
	return float32(math.Sqrt(float64(sum)))
}

func CosineSimilarity(a, b []float32) float32 {
	dotProduct := dotProduct(a, b)
	magnitudeA := magnitude(a)
	magnitudeB := magnitude(b)

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0
	}

	return dotProduct / (magnitudeA * magnitudeB)
}
