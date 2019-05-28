package republique

func (x SkirmishRating) Decrement() SkirmishRating {
	switch x {
	case SkirmishRating_ADEQUATE:
		return SkirmishRating_POOR
	case SkirmishRating_CRACK_SHOT:
		return SkirmishRating_ADEQUATE
	case SkirmishRating_EXCELLENT:
		return SkirmishRating_CRACK_SHOT
	}
	return x
}
func (x SkirmishRating) Increment() SkirmishRating {
	switch x {
	case SkirmishRating_POOR:
		return SkirmishRating_ADEQUATE
	case SkirmishRating_ADEQUATE:
		return SkirmishRating_CRACK_SHOT
	case SkirmishRating_CRACK_SHOT:
		return SkirmishRating_EXCELLENT
	}
	return x
}
