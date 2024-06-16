package handler

import (
	"github.com/SawitProRecruitment/UserService/repository"
	"math"
	"net/http"
	"sort"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
)

func (s *Server) PostEstate(ctx echo.Context) error {
	createEstateRequest := new(generated.CreateEstateRequest)
	err := ctx.Bind(&createEstateRequest)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}
	// Request Validate
	if err := s.Validator.Struct(createEstateRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	// Create Estate
	output, err := s.Repository.CreateEstate(ctx.Request().Context(), repository.Estate{
		Width:  createEstateRequest.Width,
		Length: createEstateRequest.Length,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, generated.CreateEstateResponse{Id: output.Id})
}

func (s *Server) GetEstateIdStats(ctx echo.Context, id string) error {
	// Request Validate
	if err := s.Validator.Struct(IdPath{ID: id}); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	// Check estate exist
	_, err := s.Repository.GetEstateById(ctx.Request().Context(), repository.GetEstateByIdInput{
		Id: id,
	})
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: "estate is not found"})
		} else {
			return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
		}
	}

	count, min, max, median := 0, 0, 0, 0
	// Get list trees
	trees, err := s.Repository.ListTreesByEstateId(ctx.Request().Context(), repository.ListTreesByEstateIdInput{
		EstateId: id,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}
	count = len(trees)
	if count == 0 {
		return ctx.JSON(http.StatusOK, generated.GetEstateStatsResponse{Count: count, Min: min, Max: max, Median: median})
	}

	// find min, max
	min, max = trees[0].Height, trees[0].Height
	var heights []int
	for _, tree := range trees {
		if tree.Height < min {
			min = tree.Height
		}
		if tree.Height > max {
			max = tree.Height
		}
		heights = append(heights, tree.Height)
	}

	// find median
	sort.Ints(heights)
	if count%2 == 1 {
		// Odd number of elements, return the middle one
		median = int(float64(heights[count/2]))
	}
	// Even number of elements, return the average of the two middle ones
	middle1 := heights[count/2-1]
	middle2 := heights[count/2]
	median = int(float64(middle1+middle2) / 2)
	return ctx.JSON(http.StatusOK, generated.GetEstateStatsResponse{Count: count, Min: min, Max: max, Median: median})
}

func (s *Server) PostEstateIdTree(ctx echo.Context, id string) error {
	createTreeRequest := new(generated.CreateTreeRequest)
	err := ctx.Bind(&createTreeRequest)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}
	// Request Validate
	if err := s.Validator.Struct(IdPath{ID: id}); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}
	if err := s.Validator.Struct(createTreeRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	// Check estate exist
	estate, err := s.Repository.GetEstateById(ctx.Request().Context(), repository.GetEstateByIdInput{
		Id: id,
	})
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: "estate is not found"})
		} else {
			return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
		}
	}
	// Check if plot out of bound
	if estate.Length < createTreeRequest.X || estate.Width < createTreeRequest.Y {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "index out of bound"})
	}

	// Check plot exist
	_, err = s.Repository.GetTreeByPlot(ctx.Request().Context(), repository.GetTreeByPlot{
		EstateId: id,
		X:        createTreeRequest.X,
		Y:        createTreeRequest.Y,
	})
	if err == nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "plot already exist"})
	}

	// Create Tree
	tree, err := s.Repository.CreateTree(ctx.Request().Context(), repository.Tree{
		X:        createTreeRequest.X,
		Y:        createTreeRequest.Y,
		Height:   createTreeRequest.Height,
		EstateId: estate.Id,
	})
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, generated.CreateTreeResponse{Id: tree.Id})
}

func (s *Server) GetEstateIdDronePlan(ctx echo.Context, id string, params generated.GetEstateIdDronePlanParams) error {
	// Request Validate
	if err := s.Validator.Struct(IdPath{ID: id}); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}
	if params.MaxDistance != nil && *params.MaxDistance < 1 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "invalid max distance"})
	}

	// get estate
	estate, err := s.Repository.GetEstateById(ctx.Request().Context(), repository.GetEstateByIdInput{
		Id: id,
	})
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: "estate is not found"})
		} else {
			return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
		}
	}

	// list trees
	trees, err := s.Repository.ListTreesByEstateId(ctx.Request().Context(), repository.ListTreesByEstateIdInput{
		EstateId: id,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	var maxDistance *int
	if params.MaxDistance != nil {
		maxDistance = params.MaxDistance
	}

	// map the trees plot
	mapTrees := make(map[int]map[int]int)
	for _, tree := range trees {
		if mapTrees[tree.X] == nil {
			mapTrees[tree.X] = make(map[int]int)
		}
		mapTrees[tree.X][tree.Y] = tree.Height
	}

	plot := make(map[string]interface{})
	x, y := 1, 1
	turnBack := false // true: move from west to east, false: move from east to west
	distance := 1
	for y <= estate.Width {
		if turnBack {
			x = estate.Length
			for x >= 1 {
				plot["x"], plot["y"] = x, y
				if x == 1 && y == estate.Width {
					break // break if the end of bound
				}
				height := int(math.Abs(float64(mapTrees[x][y] - mapTrees[x-1][y])))
				distance += 10 + height // sum distance by add 10 meter move and up tree
				// checking max move, if reach max break all the loop
				if maxDistance != nil && *maxDistance < distance {
					x, y = estate.Length, estate.Width
					break
				}
				x--
			}
			turnBack = false
		} else {
			x = 1
			for x <= estate.Length {
				plot["x"], plot["y"] = x, y
				if x == estate.Length && y == estate.Width {
					break // break if the end of bound
				}
				height := int(math.Abs(float64(mapTrees[x][y] - mapTrees[x+1][y])))
				distance += 10 + height // sum distance by add 10 meter move and up tree
				// checking max move, if reach max break all the loop
				if maxDistance != nil && *maxDistance < distance {
					x, y = estate.Length, estate.Width
					break
				}
				x++
			}
			turnBack = true
		}
		y++
	}
	distance += 1 // add 1 meter for drone landing

	// put distance as max distance if distance more than max distance
	if maxDistance != nil && *maxDistance < distance {
		distance = *maxDistance
	}

	return ctx.JSON(http.StatusOK, generated.GetEstateDronePlanResponse{Distance: distance, Rest: &plot})
}
