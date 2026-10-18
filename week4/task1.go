package main
import "fmt"

type Student struct {
	Name string
	ID int
	GPA float64
	Courses []string
}

type Team struct {
    Name    string
    Scores  [3]int
    Members []string
}

func addCourseValue(s Student, course string) {
	s.Courses = append(s.Courses, course)
}

func addCoursePointer(s *Student, course string) {
	s.Courses = append(s.Courses, course)
}

func main() {
	var s1 Student

	s2 := Student{
		Name: "Tyler",
		ID: 1,
		GPA: 3,
		Courses: []string{"Chemistry"},
	}

	s3 := Student{
		"Walter",
		2,
		3.5,
		[]string{"Math", "Physics"},
	}

	p := &Student{Name: "Dexter"}
	p.ID = 3
	(*p).GPA = 4.0

	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s3)
	fmt.Println(p)
	fmt.Printf("%v\n", s1)
	fmt.Printf("%+v\n", s2)
	fmt.Printf("%#v\n", s3)

	a := Student{Name: "Ali"}
	b := a
	fmt.Println(a)
	fmt.Println(b)
	b.Name = "Ahmed"
	fmt.Println(a)
	fmt.Println(b)


	t1 := Team{
    Name:    "Gophers",
    Scores:  [3]int{1, 2, 3},
    Members: []string{"Ali", "Bota"},
	}

	t2 := t1

	t2.Name = "Rustaceans"
	t2.Scores[0] = 100
	t2.Members[0] = "Nurlan"

	fmt.Println(t1) //здесь нурлан тоже, потому что оба объекта ссылаются на один и тот же андерлайн массив, а не копия
	fmt.Println(t2) 
	addCourseValue(s2, "Biology")
	addCoursePointer(&s3, "Chemistry")
	fmt.Println(s2)
	fmt.Println(s3)
}