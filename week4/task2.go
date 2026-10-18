package main

import "fmt"

type Person struct {
	Name  string
	Email string
}

func (p Person) Contact() string {
	return fmt.Sprintf("%s <%s>", p.Name, p.Email)
}

type Student struct {
	Person
	GPA float64
}

type Teacher struct {
	Person
	Department string
}

func (t Teacher) Contact() string {
	return fmt.Sprintf(
		"%s <%s> - Department: %s",
		t.Name,
		t.Email,
		t.Department,
	)
}


type StudentInfo struct {
	ID int
}

type LibraryInfo struct {
	ID int
}

type User struct {
	StudentInfo
	LibraryInfo
}

func main() {
	student := Student{
		Person: Person{
			Name:  "Tyler",
			Email: "tyler@kbtu.kz",
		},
		GPA: 3.7,
	}

	teacher := Teacher{
		Person: Person{
			Name:  "Walter",
			Email: "heisenberg@kbtu.kz",
		},
		Department: "Information Technology",
	}

	person := Person{
		Name:  "Dexter",
		Email: "dexter@kbtu.kz",
	}

	fmt.Println(student.Contact())
	fmt.Println(teacher.Contact())
	fmt.Println(teacher.Person.Contact())
	fmt.Println(person.Contact())
	fmt.Println(student.Person.Contact())
	fmt.Println(student.GPA)
	fmt.Println(student)



	u := User{
		StudentInfo: StudentInfo{ID: 10},
		LibraryInfo: LibraryInfo{ID: 20},
	}

// fmt.Println(u.ID)
// compile error: ambiguous selector u.ID

	fmt.Println(u.StudentInfo.ID)
	fmt.Println(u.LibraryInfo.ID)
}