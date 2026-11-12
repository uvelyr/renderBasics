package main

import (
	"math"
	"fmt"
//	"github.com/gdamore/tcell/v2"
)

//type Color int

type Vector3 struct {
	X, Y, Z float64
}

type Coord3D struct {
	X, Y, Z float64
}

type Coord2D struct {
	X, Y int
}

type Ray struct {
	origin Coord3D
	direction Vector3
}

type Object interface {
//	getColor() Color
	intersectsRay(Ray) (bool, float64)
	getTexture() rune
}

func dotProductVectors(a, b Vector3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func crossProductVectors(a, b Vector3) Vector3 {
	return Vector3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func dotProductVectorAndPoint(a Vector3, b Coord3D) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z	
}

func (a Vector3) scale(b float64) Vector3 {
	return Vector3{
		X: a.X*b,
		Y: a.Y*b,
		Z: a.Z*b,
	}
}


func addVectors(a, b Vector3) Vector3 {
	return Vector3{
		X: a.X+b.X,
		Y: a.Y+b.Y,
		Z: a.Z+b.Z,
	}
}

func (a Vector3) magnitude() float64{
	return math.Sqrt(dotProductVectors(a, a))
}

func (a Vector3) standartise() Vector3 {
	magnitude := a.magnitude()
	if magnitude == 0 {
		return Vector3{}
	}
	return Vector3{
		X: a.X / magnitude,
		Y: a.Y / magnitude,
		Z: a.Z / magnitude,
	}
}		

func subtractPoints(a, b Coord3D) Vector3 {
	return Vector3{
		X: a.X-b.X,
		Y: a.Y-b.Y,
		Z: a.Z-b.Z,
	}
}

//func subtractVectorFromCoord(a Coord3D, b Vector3) Vector3 {
//	return Vector3{
//		X: a.X-b.X,
//		Y: a.Y-b.Y,
//		Z: a.Z-b.Z,
//	}
//}

func addVectorToPoint(a Coord3D, b Vector3) Coord3D {
	return Coord3D{
		X: a.X+b.X,
		Y: a.Y+b.Y,
		Z: a.Z+b.Z,
	}
}

//type Axis int
//
//const (
//	X Axis = iota
//	Y
//	Z
//}

//Could use later for walls, perhaps?
type Plane struct {
	//normal is directly linked to rotation, both attributes are currently
	//present for convenience
	normal Vector3
	rotation Rotation
}

type AngleDegree float64

type Rotation struct {
	Pitch AngleDegree
	Yaw AngleDegree
	Roll AngleDegree
}

type Camera struct {
	focal_point Coord3D
	rotation Rotation
	FOV AngleDegree
	bottomleft_corner Coord3D
	right, up Vector3
}

func degreeToRadian(degree AngleDegree) float64 {
	return float64(degree*math.Pi/180)
}

func rotationToVector(rotation Rotation) Vector3 {
	sin_pitch, cos_pitch := math.Sincos(degreeToRadian(-rotation.Pitch))
	sin_yaw, cos_yaw := math.Sincos(degreeToRadian(-rotation.Yaw))
	x := cos_pitch * cos_yaw
	y := cos_pitch * sin_yaw
	z := sin_pitch
	return Vector3{x, y, z}
}

//func getVisibleCubes(scene) []Cube
const (
	ASCIIramp string = " .:=+*#@"
	characterAspectRatio float64 = 0.5
	imageAspectRatio float64 = 2.0
	defaultScreenFactor float64 = 8.0
)

func (camera *Camera) updateScreenPosition() (error) {
	plane_normal := rotationToVector(camera.rotation).scale(defaultScreenFactor)
	plane_center := addVectorToPoint(camera.focal_point, plane_normal)
	plane_constant := dotProductVectorAndPoint(plane_normal, plane_center)
	global_up := Vector3{X:0, Y:0, Z: 1}
	projection_point := addVectorToPoint(plane_center, global_up)
	projection_direction := subtractPoints(projection_point, camera.focal_point)
	t := -(dotProductVectorAndPoint(plane_normal, camera.focal_point)-plane_constant)/dotProductVectors(plane_normal, projection_direction)	
	plane_vertical_point := addVectorToPoint(camera.focal_point, projection_direction.scale(t))
	vertical_screen_vector := subtractPoints(plane_vertical_point, plane_center)
	//perpendicular to normal (plane_normal) and to the vertical_screen_vector
	helper_vector := crossProductVectors(plane_normal, vertical_screen_vector).scale(1/crossProductVectors(plane_normal, vertical_screen_vector).magnitude())
	rollRad := degreeToRadian(camera.rotation.Roll)
	cos_roll, sin_roll := math.Sincos(rollRad)
	vertical_screen_vector = addVectors(vertical_screen_vector.scale(sin_roll), helper_vector.scale(cos_roll))
	horisontal_screen_vector := crossProductVectors(plane_normal, vertical_screen_vector)
	//points up
	vertical_screen_vector = vertical_screen_vector.standartise()
	//points right
	horisontal_screen_vector = horisontal_screen_vector.standartise()
	height := 2 * plane_normal.magnitude() * math.Tan(0.5*degreeToRadian(camera.FOV))
	width := imageAspectRatio * height
	bottomleft_corner := addVectorToPoint(plane_center, vertical_screen_vector.scale(-0.5*height))
	bottomleft_corner = addVectorToPoint(bottomleft_corner, horisontal_screen_vector.scale(-0.5*width))
	camera.right = horisontal_screen_vector.scale(width)
	camera.up = vertical_screen_vector.scale(height)
	camera.bottomleft_corner = bottomleft_corner
	fmt.Println("Focal point:              ", camera.focal_point)
	fmt.Println("Camera direction:         ", plane_normal)
	fmt.Println("Plane_center:             ", plane_center)
	fmt.Println("Plane constant:           ", plane_constant)
	fmt.Println("Projection point:         ", projection_point)
	fmt.Println("Projection direction:     ", projection_direction)
	fmt.Println("Factor:                   ", t)
	fmt.Println("Plane vertical point:     ", plane_vertical_point)
	fmt.Println("Vertical_screen_vector:   ", vertical_screen_vector)
	fmt.Println("Horisontal_screen_vector: ", horisontal_screen_vector)
	fmt.Println("bottomleft_corner         ", bottomleft_corner)
	fmt.Println("width, height:            ", width, height)
	fmt.Println("up                        ", camera.up)
	fmt.Println("right                     ", camera.right)
	return nil
}
func (camera Camera) getRayOriginsAndPixels(resolutionX, resolutionY int) ([]Coord3D, []Coord2D) {
//	height, width := camera.up.magnitude(), camera.right.magnitude()
	vertical_step := camera.up.scale(1/float64(resolutionY))
	horisontal_step := camera.right.scale(1/float64(resolutionX))
	//let i indicate number of vertical steps, while j number of horisontal
	//then a pair of i and j denotes a point on the plane
	//a and b are factors of i and j respectively
	var origins []Coord3D
	var pixelPositions []Coord2D
	fmt.Println(camera.bottomleft_corner)
	for x := 0; x < resolutionX; x++ {
		for y := 0; y < resolutionY; y++ {
			horisontal_step_number, vertical_step_number := 0.5+float64(x), 0.5+float64(y)
			if (x == 0 && y == 0) {
				fmt.Println("hor_st_nu: ", horisontal_step_number)
				fmt.Println("hor step: ", horisontal_step)
			}
			origin := addVectorToPoint(camera.bottomleft_corner, horisontal_step.scale(horisontal_step_number))
			origin = addVectorToPoint(origin, vertical_step.scale(vertical_step_number))
			origins = append(origins, origin)
			pixelPositions = append(pixelPositions, Coord2D{X: x, Y: y})
		}
	}
	return origins, pixelPositions
} 

func move(dx, dy, dz float64) {
	return
}

//WIP
func (camera Camera) rotate(rotationX, rotationY, rotationZ AngleDegree) error {
	camera.updateScreenPosition()
	return nil
}

type Sphere struct {
	radius float64
	center Coord3D
	texture rune
}

func (sphere Sphere) intersectsRay(ray Ray) (bool, float64) {
	//quadratic equation coefficients
	oc := subtractPoints(ray.origin, sphere.center)
	a := dotProductVectors(ray.direction, ray.direction)
	b := dotProductVectors(ray.direction, oc) * 2.0
	c := dotProductVectors(oc, oc) - sphere.radius*sphere.radius

	discriminant := b*b - 4 * a * c
	//does not intersect
	if discriminant < 0 {
		return false, -1.0
	}
	//we have to check if t is positive, so the sphere is ahead, not behind the ray
	discriminant_root := math.Sqrt(discriminant)

	t1 := (-b + discriminant_root) / (2 * a)
	t2 := (-b - discriminant_root) / (2 * a)	
	if t1 < 0 || t2 < 0 {
		return false, -1.0
	}
	return true, min(t1, t2)
}


func (sphere Sphere) getTexture() rune {
	return sphere.texture
}

func main() {
//	screen, err := tcell.NewScreen()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if err := screen.Init(); err != nil {
//		log.Fatal(err)
//	}
//
//	defer screen.Fini()
//	defer screen.Clear()

	var objects []Object
	focal_point := Coord3D{X: 30, Y: -30, Z: 10}
	
	const resolutionX int = 308
	var resolutionY int = int(math.Round(float64(resolutionX)/imageAspectRatio*characterAspectRatio))

	image := make([][]rune, resolutionX)
	for x := range image {
		image[x] = make([]rune, resolutionY)
	}

	closestIntersectFactors := make([][]float64, resolutionX)
	for x := range closestIntersectFactors {
		closestIntersectFactors[x] = make([]float64, resolutionY)
	}

	//initiliasing arrays
	for i := 0; i < resolutionX; i++ {
		for y := 0; y < resolutionY; y++ {
			image[i][y] = ' '
			closestIntersectFactors[i][y] = -1.0
		}
	}
	rotation := Rotation{
		Pitch: AngleDegree(0.0),
		Yaw: AngleDegree(0.0),
		Roll: AngleDegree(0.0),
	}
	FOV := AngleDegree(40)
	camera := Camera{focal_point, rotation, FOV, Coord3D{0, 0, 0}, Vector3{0, 0, 0}, Vector3{0, 0, 0}}
	camera.updateScreenPosition()

	sphere := Sphere{center: Coord3D{X: 10, Y: 20, Z: 10}, radius: 15, texture: '#'}
	objects = append(objects, sphere)
	for true {
		for i := 0; i < resolutionX; i++ {
			for y := 0; y < resolutionY; y++ {
				image[i][y] = ' '
				closestIntersectFactors[i][y] = -1.0
			}
		}
		rayOrigins, pixelPositions := camera.getRayOriginsAndPixels(resolutionX, resolutionY)
		minimum_t, maximum_t := math.Inf(1), -1.0
		for i, rayOrigin := range rayOrigins {
			var ray Ray
			if i == 1 {
				fmt.Println("rayOrigin: ", rayOrigin)
				fmt.Println("pos:       ", pixelPositions[i])
			}
			ray.origin = rayOrigin
			ray.direction = subtractPoints(rayOrigin, camera.focal_point)
			ray.direction = ray.direction.standartise()
			for _, object := range objects {
				hit, intersectFactor := object.intersectsRay(ray)
				if !hit {
					continue
				}
				minimum_t, maximum_t = min(float64(minimum_t), intersectFactor), max(float64(maximum_t), intersectFactor)
				pos := pixelPositions[i]
				smallestFactor := closestIntersectFactors[pos.X][pos.Y]
				if smallestFactor == -1.0 || intersectFactor < smallestFactor {
					closestIntersectFactors[pos.X][pos.Y] = intersectFactor
					image[pos.X][pos.Y] = object.getTexture()
				}
			}
		}
		for y := 0; y < resolutionY; y++ {
			for x := 0; x < resolutionX; x++ {
				factor := closestIntersectFactors[x][y] - minimum_t
				step := (maximum_t - minimum_t)/6.0
				switch {
					case factor < 0:
						fmt.Print(string(ASCIIramp[0]))
					case factor >= 0 && factor < step:
						fmt.Print(string(ASCIIramp[7]))
					case factor >= step && factor < step*2.0:
						fmt.Print(string(ASCIIramp[6]))
					case factor >= step*2.0 && factor < step*3.0:
						fmt.Print(string(ASCIIramp[5]))
					case factor >= step*3.0 && factor < step*4.0:
						fmt.Print(string(ASCIIramp[4]))
					case factor >= step*4.0 && factor < step*5.0:
						fmt.Print(string(ASCIIramp[3]))
					case factor >= step*5.0 && factor < step*6.0:
						fmt.Print(string(ASCIIramp[2]))
					case factor >= step*6.0:
						fmt.Print(string(ASCIIramp[1]))
				}				
			}
			fmt.Println()
		}
		var input string
		fmt.Scan(&input)
		switch input {
			case "w": camera.rotation.Pitch += 5.0
			case "s": camera.rotation.Pitch -= 5.0
			case "d": camera.rotation.Yaw += 5.0
			case "a": camera.rotation.Yaw -= 5.0
			case "f": camera.FOV -= 5.0
			case "g": camera.FOV += 5.0
			case "v": camera.rotation.Roll += 5.0
			case "b": camera.rotation.Roll -= 5.0
		}
		camera.updateScreenPosition()
	}
}

