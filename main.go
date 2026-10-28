package main

import (
	"math"
	"fmt"
	"github.com/gdamore/tcell/v2"
)

//type Color int

type Vector3 struct {
	X, Y, Z float64
}

type Coord struct {
	X, Y, Z float64
}

type Ray struct {
	origin Coord
	direction Vector3
}

type Object interface {
//	getColor() Color
	intersectsRay(Ray) bool
}

func dotProduct(a, b Vector3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func scaleVector(a Vector3, b float64) Vector3 {
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

func standartiseVector(a Vector3) Vector3 {
	magnitude := dotProduct(a, a)
	return Vector3{
		X: a.X / magnitude,
		Z: a.Y / magnitude,
		Y: a.Z / magnitude,
	}
}		

func subtractCoordinates(a, b Coord) Vector3 {
	return Vector3{
		X: a.X-b.X,
		Y: a.Y-b.Y,
		Z: a.Z-b.Z,
	}
}

func subtractVectorFromCoord(a Coord, b Vector3) Vector3 {
	return Vector3{
		X: a.X-b.X,
		Y: a.Y-b.Y,
		Z: a.Z-b.Z,
	}
}

type Sphere struct {
	radius float64
	center Coord
}

func (sphere Sphere) intersectsRay(ray Ray) bool {
	//quadratic equation coefficients
	oc := subtractCoordinates(ray.origin, sphere.center)
	a := dotProduct(ray.direction, ray.direction)
	b := dotProduct(ray.direction, oc) * 2.0
	c := dotProduct(oc, oc) - sphere.radius*sphere.radius

	discriminant := b*b - 4 * a * c
	//does not intersect
	if discriminant < 0 {
		return false
	}
	//we have to check if t is positive, so the sphere is ahead, not behind the ray
	discriminant_root := math.Sqrt(discriminant)

	t1 := (-b + discriminant_root) / (2 * a)
	if t1 >= 0 {
		return true
	}
	t2 := (-b - discriminant_root) / (2 * a)	
	if t2 >= 0 {
		return true
	}
	return false
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}

	if err := screen.Init(); err != nil {
		log.Fatal(err)
	}

	defer screen.Fini()
	defer screen.Clear()

	var scene []Object
//	var screen [20][20]char
	sphere := Sphere{center: Coord{X: 10, Y: 20, Z: 10}, radius: 15}
	scene = append(scene, sphere)
	y := 0
	focal_point := Coord{X: 10, Y: -10, Z: 10}
	for z := 19; z >= 0; z-- {
		for x := 0; x <= 20; x++ {
			var ray Ray
			ray.origin = Coord{X: float64(x), Y: float64(y), Z: float64(z)}
			ray.direction = subtractCoordinates(ray.origin, focal_point)
			//ray.direction = standartise(ray.direction)
			if sphere.intersectsRay(ray) {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	//fmt.Println("sup")
}
