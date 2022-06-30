package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ID incremented
var ID int = 4

type Student struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Age      int    `json:"age"`
}

var Students = []Student{
	{ID: 1, FullName: "Diego", Age: 19},
	{ID: 2, FullName: "Maria", Age: 19},
	{ID: 3, FullName: "Paulo", Age: 20},
}

func routeHearth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func routePostStudent(c *gin.Context) {
	var student Student

	err := c.Bind(&student)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possivel obter o payload",
		})

		return
	}

	student.ID = ID
	Students = append(Students, student)
	ID++
	c.JSON(http.StatusCreated, student)
}

func routePutStudent(c *gin.Context) {
	var studentPayload Student
	var studentLocal Student
	var newStudents []Student

	err := c.Bind(&studentPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possivel obter o payload",
		})

		return
	}

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possivel obter Param",
		})

		return
	}

	for _, studentElement := range Students {
		if id == studentElement.ID {
			studentLocal = studentElement
		}
	}

	fmt.Println(studentLocal.ID, studentLocal.FullName, studentLocal.Age)

	if studentLocal.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possivel encontrar o estudante",
		})

		return
	}

	studentLocal.FullName = studentPayload.FullName
	studentLocal.Age = studentPayload.Age

	for _, studentElement := range Students {
		if id == studentElement.ID {
			newStudents = append(newStudents, studentLocal)
		} else {
			newStudents = append(newStudents, studentElement)
		}
	}

	Students = newStudents

	c.JSON(http.StatusOK, studentLocal)
}

func routeGetStudents(c *gin.Context) {
	c.JSON(http.StatusOK, Students)
}

func routeDeleteStudent(c *gin.Context) {

	var newStudents []Student
	var studentExits bool

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não possivel obter o Param",
		})
	}

	for _, studentElement := range Students {
		if id == studentElement.ID {
			studentExits = true
		}
	}

	if !studentExits {
		c.JSON(http.StatusNotFound, gin.H{
			"message_error": "Não possivel encontrar o estudante",
		})

		return
	}

	for _, studentElement := range Students {
		if id != studentElement.ID {
			newStudents = append(newStudents, studentElement)
		}
	}

	Students = newStudents

	c.JSON(http.StatusOK, gin.H{
		"message": "Student excluido com sucesso",
	})

}

func main() {
	service := gin.Default()

	getRoutes(service)

	service.Run()
}

func getRoutes(c *gin.Engine) *gin.Engine {
	c.GET("/heart", routeHearth)

	groupStudents := c.Group("/students")
	groupStudents.GET("/", routeGetStudents)
	groupStudents.POST("/", routePostStudent)
	groupStudents.PUT("/:id", routePutStudent)
	groupStudents.DELETE("/:id", routeDeleteStudent)

	return c
}
