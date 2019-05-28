package republique

import (
	rp "github.com/steveoc64/republique5/proto"
)

func (x rp.SkirmishRating) Decrement() rp.SkirmishRating {
	switch x {
	case rp.SkirmishRating_ADEQUATE:
		return rp.SkirmishRating_POOR
	case rp.SkirmishRating_CRACK_SHOT:
		return rp.SkirmishRating_ADEQUATE
	case rp.SkirmishRating_EXCELLENT:
		return rp.SkirmishRating_CRACK_SHOT
	}
	return x
}
func (x rp.SkirmishRating) Increment() rp.SkirmishRating {
	switch x {
	case rp.SkirmishRating_POOR:
		return rp.SkirmishRating_ADEQUATE
	case rp.SkirmishRating_ADEQUATE:
		return rp.SkirmishRating_CRACK_SHOT
	case rp.SkirmishRating_CRACK_SHOT:
		return rp.SkirmishRating_EXCELLENT
	}
	return x
}
