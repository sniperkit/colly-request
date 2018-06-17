package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Note: since database connection/schema is reused, these tests cannot run in parallel
const pgTestConnection = "postgres://rbaskets:pwd@localhost/baskets?sslmode=disable"

func TestPgSQLDatabase_Create(t *testing.T) {
	name := "test1"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	auth, err := db.Create(name, BasketConfig{Capacity: 20})
	defer db.Delete(name)

	if assert.NoError(t, err) {
		assert.NotEmpty(t, auth.Token, "basket token may not be empty")
		assert.False(t, len(auth.Token) < 30, "weak basket token: %v", auth.Token)
	}
}

func TestPgSQLDatabase_Create_NameConflict(t *testing.T) {
	name := "test2"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	db.Create(name, BasketConfig{Capacity: 20})
	defer db.Delete(name)

	auth, err := db.Create(name, BasketConfig{Capacity: 20})

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), ": "+name+" ", "error is not detailed enough")
		assert.Empty(t, auth.Token, "basket token is not expected")
	}
}

func TestPgSQLDatabase_Get(t *testing.T) {
	name := "test3"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	auth, err := db.Create(name, BasketConfig{Capacity: 16})
	defer db.Delete(name)

	assert.NoError(t, err)

	basket := db.Get(name)
	if assert.NotNil(t, basket, "basket with name: %v is expected", name) {
		assert.True(t, basket.Authorize(auth.Token), "basket authorization has failed")
		assert.Equal(t, 16, basket.Config().Capacity, "wrong capacity")
	}
}

func TestPgSQLDatabase_Get_NotFound(t *testing.T) {
	name := "test4"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	basket := db.Get(name)
	assert.Nil(t, basket, "basket with name: %v is not expected", name)
}

func TestPgSQLDatabase_Delete(t *testing.T) {
	name := "test5"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	db.Create(name, BasketConfig{Capacity: 10})
	assert.NotNil(t, db.Get(name), "basket with name: %v is expected", name)

	db.Delete(name)
	assert.Nil(t, db.Get(name), "basket with name: %v is not expected", name)
}

func TestPgSQLDatabase_Delete_Multi(t *testing.T) {
	name := "test6"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	config := BasketConfig{Capacity: 10}
	for i := 0; i < 10; i++ {
		bname := fmt.Sprintf("%s_%v", name, i)
		db.Create(bname, config)
		defer db.Delete(bname)
	}

	dname := name + "_5"

	assert.NotNil(t, db.Get(dname), "basket with name: %v is expected", name)
	assert.Equal(t, 10, db.Size(), "wrong database size")

	db.Delete(dname)

	assert.Nil(t, db.Get(dname), "basket with name: %v is not expected", name)
	assert.Equal(t, 9, db.Size(), "wrong database size")
}

func TestPgSQLDatabase_Size(t *testing.T) {
	name := "test7"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	config := BasketConfig{Capacity: 15}
	for i := 0; i < 25; i++ {
		bname := fmt.Sprintf("%s_%v", name, i)
		db.Create(bname, config)
		defer db.Delete(bname)
	}

	assert.Equal(t, 25, db.Size(), "wrong database size")
}

func TestPgSQLDatabase_GetNames(t *testing.T) {
	name := "test8"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	config := BasketConfig{Capacity: 15}
	for i := 0; i < 45; i++ {
		bname := fmt.Sprintf("%s_%v", name, i)
		db.Create(bname, config)
		defer db.Delete(bname)
	}

	// Get and validate page 1 (test8_0, test8_1, test8_10, test8_11, ... - sorted)
	page1 := db.GetNames(10, 0)
	assert.Equal(t, 45, page1.Count, "wrong baskets count")
	assert.True(t, page1.HasMore, "expected more names")
	assert.Len(t, page1.Names, 10, "wrong page size")
	assert.Equal(t, "test8_10", page1.Names[2], "wrong basket name at index #2")

	// Get and validate page 5 (test8_5, test8_6, test8_7, test8_8, test8_9)
	page5 := db.GetNames(10, 40)
	assert.Equal(t, 45, page5.Count, "wrong baskets count")
	assert.False(t, page5.HasMore, "no more names are expected")
	assert.Len(t, page5.Names, 5, "wrong page size")
	assert.Equal(t, "test8_5", page5.Names[0], "wrong basket name at index #0")

	// Corner cases
	assert.Empty(t, db.GetNames(0, 0).Names, "names are not expected")
	assert.False(t, db.GetNames(5, 40).HasMore, "no more names are expected")
}

func TestPgSQLDatabase_FindNames(t *testing.T) {
	name := "test9"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	config := BasketConfig{Capacity: 5}
	for i := 0; i < 35; i++ {
		bname := fmt.Sprintf("%s_%v", name, i)
		db.Create(bname, config)
		defer db.Delete(bname)
	}

	res1 := db.FindNames("test9_2", 20, 0)
	assert.False(t, res1.HasMore, "no more names are expected")
	assert.Len(t, res1.Names, 11, "wrong number of found names")
	for _, name := range res1.Names {
		assert.Contains(t, name, "test9_2", "invalid name among search results")
	}

	res2 := db.FindNames("test9_1", 5, 0)
	assert.True(t, res2.HasMore, "more names are expected")
	assert.Len(t, res2.Names, 5, "wrong number of found names")

	// Corner cases
	assert.Len(t, db.FindNames("test9_1", 5, 10).Names, 1, "wrong number of returned names")
	assert.Empty(t, db.FindNames("test9_2", 5, 20).Names, "names in this page are not expected")
	assert.False(t, db.FindNames("test9_3", 5, 6).HasMore, "no more names are expected")
	assert.False(t, db.FindNames("abc", 5, 0).HasMore, "no more names are expected")
	assert.Empty(t, db.FindNames("xyz", 5, 0).Names, "names are not expected")
}

func TestPgSQLBasket_Add(t *testing.T) {
	name := "test101"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	db.Create(name, BasketConfig{Capacity: 20})
	defer db.Delete(name)

	basket := db.Get(name)
	if assert.NotNil(t, basket, "basket with name: %v is expected", name) {
		// add 1st HTTP request
		content := "{ \"user\": \"tester\", \"age\": 24 }"
		data := basket.Add(createTestPOSTRequest(
			fmt.Sprintf("http://localhost/%v/demo?name=abc&ver=12", name), content, "application/json"))

		assert.Equal(t, 1, basket.Size(), "wrong basket size")

		// detailed http.Request to RequestData tests should be covered by test of ToRequestData function
		assert.Equal(t, content, data.Body, "wrong body")
		assert.Equal(t, int64(len(content)), data.ContentLength, "wrong content length")

		// add 2nd HTTP request
		basket.Add(createTestPOSTRequest(fmt.Sprintf("http://localhost/%v/demo", name), "Hellow world", "text/plain"))
		assert.Equal(t, 2, basket.Size(), "wrong basket size")
	}
}

func TestPgSQLBasket_Add_ExceedLimit(t *testing.T) {
	name := "test102"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	db.Create(name, BasketConfig{Capacity: 10})
	defer db.Delete(name)

	basket := db.Get(name)
	if assert.NotNil(t, basket, "basket with name: %v is expected", name) {
		// fill basket
		for i := 0; i < 35; i++ {
			basket.Add(createTestPOSTRequest(
				fmt.Sprintf("http://localhost/%v/demo", name), fmt.Sprintf("test%v", i), "text/plain"))
		}
		assert.Equal(t, 10, basket.Size(), "wrong basket size")
	}
}

func TestPgSQLBasket_Clear(t *testing.T) {
	name := "test103"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	db.Create(name, BasketConfig{Capacity: 20})
	defer db.Delete(name)

	basket := db.Get(name)
	if assert.NotNil(t, basket, "basket with name: %v is expected", name) {
		// fill basket
		for i := 0; i < 15; i++ {
			basket.Add(createTestPOSTRequest(
				fmt.Sprintf("http://localhost/%v/demo", name), fmt.Sprintf("test%v", i), "text/plain"))
		}
		assert.Equal(t, 15, basket.Size(), "wrong basket size")

		// clean basket
		basket.Clear()
		assert.Equal(t, 0, basket.Size(), "wrong basket size, empty basket is expected")
	}
}

func TestPgSQLBasket_Update_Shrink(t *testing.T) {
	name := "test104"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	db.Create(name, BasketConfig{Capacity: 30})
	defer db.Delete(name)

	basket := db.Get(name)
	if assert.NotNil(t, basket, "basket with name: %v is expected", name) {
		// fill basket
		for i := 0; i < 25; i++ {
			basket.Add(createTestPOSTRequest(
				fmt.Sprintf("http://localhost/%v/demo", name), fmt.Sprintf("test%v", i), "text/plain"))
		}
		assert.Equal(t, 25, basket.Size(), "wrong basket size")

		// update config with lower capacity
		config := basket.Config()
		config.Capacity = 12
		basket.Update(config)
		assert.Equal(t, config.Capacity, basket.Size(), "wrong basket size")
	}
}

func TestPgSQLBasket_GetRequests(t *testing.T) {
	name := "test105"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	db.Create(name, BasketConfig{Capacity: 25})
	defer db.Delete(name)

	basket := db.Get(name)
	if assert.NotNil(t, basket, "basket with name: %v is expected", name) {
		// fill basket
		for i := 1; i <= 35; i++ {
			basket.Add(createTestPOSTRequest(
				fmt.Sprintf("http://localhost/%v/demo?id=%v", name, i), fmt.Sprintf("req%v", i), "text/plain"))
		}
		assert.Equal(t, 25, basket.Size(), "wrong basket size")

		// Get and validate last 10 requests
		page1 := basket.GetRequests(10, 0)
		assert.True(t, page1.HasMore, "expected more requests")
		assert.Len(t, page1.Requests, 10, "wrong page size")
		assert.Equal(t, 25, page1.Count, "wrong requests count")
		assert.Equal(t, 35, page1.TotalCount, "wrong requests total count")
		assert.Equal(t, "req35", page1.Requests[0].Body, "last request #35 is expected at index #0")

		// Get and validate 10 requests, skip 20
		page3 := basket.GetRequests(10, 20)
		assert.False(t, page3.HasMore, "no more requests are expected")
		assert.Len(t, page3.Requests, 5, "wrong page size")
		assert.Equal(t, 25, page3.Count, "wrong requests count")
		assert.Equal(t, 35, page3.TotalCount, "wrong requests total count")
		assert.Equal(t, "req15", page3.Requests[0].Body, "request #15 is expected at index #0")

		// Get only collected statistics
		page0 := basket.GetRequests(0, 0)
		assert.True(t, page0.HasMore, "expected more requests")
		assert.Empty(t, page0.Requests, "requests are not expected")
		assert.Equal(t, 25, page1.Count, "wrong requests count")
		assert.Equal(t, 35, page1.TotalCount, "wrong requests total count")
	}
}

func TestPgSQLBasket_FindRequests(t *testing.T) {
	name := "test106"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	db.Create(name, BasketConfig{Capacity: 100})
	defer db.Delete(name)

	basket := db.Get(name)
	if assert.NotNil(t, basket, "basket with name: %v is expected", name) {
		// fill basket
		for i := 1; i <= 30; i++ {
			r := createTestPOSTRequest(fmt.Sprintf("http://localhost/%v?id=%v", name, i), fmt.Sprintf("req%v", i), "text/plain")
			r.Header.Add("HeaderId", fmt.Sprintf("header%v", i))
			if i <= 10 {
				r.Header.Add("ChocoPie", "yummy")
			}
			if i <= 20 {
				r.Header.Add("Muffin", "tasty")
			}
			basket.Add(r)
		}
		assert.Equal(t, 30, basket.Size(), "wrong basket size")

		// search everywhere
		s1 := basket.FindRequests("req1", "any", 30, 0)
		assert.False(t, s1.HasMore, "no more results are expected")
		assert.Len(t, s1.Requests, 11, "wrong number of found requests")
		for _, r := range s1.Requests {
			assert.Contains(t, r.Body, "req1", "incorrect request among results")
		}

		// search everywhere (limited output)
		s2 := basket.FindRequests("req2", "any", 5, 5)
		assert.True(t, s2.HasMore, "more results are expected")
		assert.Len(t, s2.Requests, 5, "wrong number of found requests")

		// search everywhere with max = 0
		assert.Empty(t, basket.FindRequests("req2", "any", 0, 0).Requests, "found unexpected requests")

		// search in body (positive)
		assert.Len(t, basket.FindRequests("req3", "body", 100, 0).Requests, 2, "wrong number of found requests")
		// search in body (negative)
		assert.Empty(t, basket.FindRequests("yummy", "body", 100, 0).Requests, "found unexpected requests")

		// search in headers (positive)
		assert.Len(t, basket.FindRequests("yummy", "headers", 100, 0).Requests, 10, "wrong number of found requests")
		assert.Len(t, basket.FindRequests("tasty", "headers", 100, 0).Requests, 20, "wrong number of found requests")
		// search in headers (negative)
		assert.Empty(t, basket.FindRequests("req1", "headers", 100, 0).Requests, "found unexpected requests")

		// search in query (positive)
		assert.Len(t, basket.FindRequests("id=1", "query", 100, 0).Requests, 11, "wrong number of found requests")
		// search in query (negative)
		assert.Empty(t, basket.FindRequests("tasty", "query", 100, 0).Requests, "found unexpected requests")
	}
}

func TestPgSQLBasket_SetResponse(t *testing.T) {
	name := "test107"
	method := "POST"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	db.Create(name, BasketConfig{Capacity: 20})
	defer db.Delete(name)

	basket := db.Get(name)
	if assert.NotNil(t, basket, "basket with name: %v is expected", name) {
		// Ensure no response
		assert.Nil(t, basket.GetResponse(method))

		// Set response
		basket.SetResponse(method, ResponseConfig{Status: 201, Body: "{ 'message' : 'created' }"})
		// Get and validate
		response := basket.GetResponse(method)
		if assert.NotNil(t, response, "response for method: %v is expected", method) {
			assert.Equal(t, 201, response.Status, "wrong HTTP response status")
			assert.Equal(t, "{ 'message' : 'created' }", response.Body, "wrong HTTP response body")
			assert.False(t, response.IsTemplate, "template is not expected")
		}
	}
}

func TestPgSQLBasket_SetResponse_Update(t *testing.T) {
	name := "test108"
	method := "GET"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	db.Create(name, BasketConfig{Capacity: 20})
	defer db.Delete(name)

	basket := db.Get(name)
	if assert.NotNil(t, basket, "basket with name: %v is expected", name) {
		// Set response
		basket.SetResponse(method, ResponseConfig{Status: 200, Body: ""})
		// Update response
		basket.SetResponse(method, ResponseConfig{Status: 200, Body: "welcome", IsTemplate: true})
		// Get and validate
		response := basket.GetResponse(method)
		if assert.NotNil(t, response, "response for method: %v is expected", method) {
			assert.Equal(t, 200, response.Status, "wrong HTTP response status")
			assert.Equal(t, "welcome", response.Body, "wrong HTTP response body")
			assert.True(t, response.IsTemplate, "template is expected")
		}
	}
}

func TestPgSQLBasket_Config_Error(t *testing.T) {
	name := "test120"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	db.Create(name, BasketConfig{Capacity: 30, ForwardURL: "http://localhost:8080"})
	basket := db.Get(name)
	// delete basket
	db.Delete(name)

	// try to get configuration of deleted basket
	config := basket.Config()
	if assert.NotNil(t, config, "configuration is expected") {
		// empty config is expected
		assert.Equal(t, 0, config.Capacity, "Capacity is not expected")
		assert.Empty(t, config.ForwardURL, "ForwardURL is not expected")
	}
}

func TestPgSQLBasket_SetResponse_Error(t *testing.T) {
	name := "test121"
	method := "POSTVERYVERYVERYVERYLONGNAME"
	db := NewSQLDatabase(pgTestConnection)
	defer db.Release()

	db.Create(name, BasketConfig{Capacity: 20})
	defer db.Delete(name)

	basket := db.Get(name)
	if assert.NotNil(t, basket, "basket with name: %v is expected", name) {
		// Ensure no response
		assert.Nil(t, basket.GetResponse(method))

		// Set response
		basket.SetResponse(method, ResponseConfig{Status: 201, Body: "{ 'message' : 'created' }"})

		// Ensure no response
		assert.Nil(t, basket.GetResponse(method), "Response for very long method name is not expected")
	}
}
