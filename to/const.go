package to

const (
	FaceRkg    = "Face Recognition"
	GestureRkg = "Gesture Recognition"
)

// vacant, idle, prepared
var ConstDelay = [3]float64{5.0, 50.0, 500.0}

var TypeDeadlineDict = map[string]float64{
	FaceRkg:    2500,
	GestureRkg: 2000,
}

var TypeResourceDict = map[string][3]float64{
	FaceRkg:    {0.3, 1.5, 180},
	GestureRkg: {0.15, 0.8, 75},
}

var TypeDataSizeDict = map[string]float64{
	FaceRkg:    1,
	GestureRkg: 0.5,
}
