package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	city := "moscow"
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
        count int 
        want  int 
    }{
       {count: 0, want: 0},
	   {count: 1, want: 1},
	   {count: 2, want: 2},
	   {count: 100, want: len(cafeList[city])},
    }

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=%s&count=%d", city, v.count), nil)

		handler.ServeHTTP(response, req)

		var responseCity []string

		cityRes := response.Body.String()

		if cityRes == "" {
			responseCity = []string{}
		} else {
			responseCity = strings.Split(cityRes, ",")
		}
 
		fmt.Println(responseCity, len(responseCity), v.want)
		assert.Equal(t, len(responseCity), v.want)
	}
}