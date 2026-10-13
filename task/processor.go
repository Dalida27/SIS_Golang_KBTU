package task

type Processor interface {
	Process(scores []int) Result
}

type ScoreProcessor struct {
	PassScore int
}

func (p ScoreProcessor) Process(scores []int) Result {
	var result Result

	if len(scores) == 0 {
		return result
	}

	result.Count = len(scores)
	result.Min = scores[0]
	result.Max = scores[0]

	sum := 0
	for _, score := range scores {
		sum += score

		if score < result.Min {
			result.Min = score
		}
		if score > result.Max {
			result.Max = score
		}
		if score >= p.PassScore {
			result.Passed++
		}
	}

	result.Average = float64(sum) / float64(result.Count)

	return result
}
