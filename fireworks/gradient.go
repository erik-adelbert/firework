package fireworks

func PhaseGradient2(x float64) float64 {
	if x < 0.087 {
		return 150.0 * x * x
	}

	return -0.8*x + 1.2
}

func PhaseScaledGradient2(x float64) float64 {
	if x < 0.087 {
		return 150 * x * x * 0.6
	}

	return (-0.8*x + 1.2) * 0.6
}

func PhaseGradient5(x float64) float64 {
	switch {
	case x < 0.067:
		return 5*x + 0.1
	case x < 0.2:
		return 2*x + 0.3
	case x < 0.5:
		return x + 0.5
	case x < 0.684:
		return x*0.5 + 0.75
	}

	return -7*(x-0.65)*(x-0.65) + 1.1
}

func LinearGradient(x float64) float64 {
	return -0.7*x + 1
}
