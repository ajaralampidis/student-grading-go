package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Grade string

const (
	A Grade = "A"
	B Grade = "B"
	C Grade = "C"
	F Grade = "F"
)

type student struct {
	firstName, lastName, university                string
	test1Score, test2Score, test3Score, test4Score int
}

func (s student) String() string {
	return fmt.Sprintf("First Name: %s\nLast Name: %s\nUniversity: %s\nTest 1 Score: %v\nTest 2 Score: %v\nTest 3 Score: %v\nTest 4 Score: %v\n",
		s.firstName, s.lastName, s.university,
		s.test1Score, s.test2Score, s.test3Score, s.test4Score)
}

type studentStat struct {
	student
	finalScore float32
	grade      Grade
}

func (s studentStat) String() string {
	return fmt.Sprintf("%vFinal Score: %v\nGrade: %v\n", s.student, s.finalScore, s.grade)
}

func parseCSV(filePath string) []student {
	file, ferr := os.Open(filePath)
	if ferr != nil {
		panic(ferr)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	count := 0
	var students []student
	for scanner.Scan() {
		if count == 0 {
			count++
			continue
		}
		count++
		line := scanner.Text()
		items := strings.Split(line, ",")
		firstName := items[0]
		lastName := items[1]
		university := items[2]
		test1Score, errScore1 := strconv.Atoi(items[3])
		test2Score, errScore2 := strconv.Atoi(items[4])
		test3Score, errScore3 := strconv.Atoi(items[5])
		test4Score, errScore4 := strconv.Atoi(items[6])

		if errScore1 != nil || errScore2 != nil || errScore3 != nil || errScore4 != nil {
			panic("errScore for " + firstName)
		}

		students = append(students, student{
			firstName:  firstName,
			lastName:   lastName,
			university: university,
			test1Score: test1Score,
			test2Score: test2Score,
			test3Score: test3Score,
			test4Score: test4Score,
		})
	}

	return students
}

func grade(score float32) (Grade, error) {
	if score < 35 {
		return "F", nil
	}

	if score >= 35 && score < 50 {
		return "C", nil
	}

	if score >= 50 && score < 70 {
		return "B", nil
	}

	if score >= 70 {
		return "A", nil
	}

	return "", errors.New("invalid score")
}

func calculateGrade(students []student) []studentStat {
	var sStats []studentStat
	for _, s := range students {
		finalScore := float32(s.test1Score+s.test2Score+s.test3Score+s.test4Score) / 4
		grade, err := grade(finalScore)

		if err != nil {
			panic(err)
		}

		sStats = append(sStats, studentStat{
			student:    s,
			finalScore: finalScore,
			grade:      grade,
		})

	}
	return sStats
}

func findOverallTopper(gradedStudents []studentStat) studentStat {
	var t studentStat
	t.finalScore = -1
	for _, s := range gradedStudents {
		if s.finalScore > t.finalScore {
			t = s
		}
	}
	if t.finalScore == -1 {
		panic("No topper found")
	}
	return t
}

func findTopperPerUniversity(gs []studentStat) map[string]studentStat {
	r := make(map[string]studentStat)
	for _, s := range gs {
		if currentTopper, ok := r[s.university]; ok {
			if currentTopper.finalScore < s.finalScore {
				r[s.university] = s
			}
			continue
		}

		r[s.university] = s
	}
	return r
}

func main() {
	students := parseCSV("grades.csv")
	studentSts := calculateGrade(students)
	fmt.Println(studentSts[0])
}
