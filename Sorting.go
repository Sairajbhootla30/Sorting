package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sort"
    "time"
)

type SortRequest struct {
    ToSort [][]int `json:"to_sort"`
}

type SortResponse struct {
    SortedArrays [][]int `json:"sorted_arrays"`
    TimeNs int64 `json:"time_ns"`
}

func main() {
    http.HandleFunc("/process-single", processSingle)
    http.HandleFunc("/process-concurrent", processConcurrent)
    fmt.Println("Server listening on port 8000...")
    fmt.Println(http.ListenAndServe(":8000", nil))
}

func processSingle(w http.ResponseWriter, r *http.Request) {
    var req SortRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    startTime := time.Now()
    for _, subArray := range req.ToSort {
        sort.Slice(subArray, func(i, j int) bool { return subArray[i] < subArray[j] })
    }
    endTime := time.Now()

    resp := SortResponse{
        SortedArrays: req.ToSort,
        TimeNs: endTime.Sub(startTime).Nanoseconds(),
    }

    json.NewEncoder(w).Encode(resp)
}

func processConcurrent(w http.ResponseWriter, r *http.Request) {
    var req SortRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    startTime := time.Now()

    // Channel to receive sorted sub-arrays
    results := make(chan []int, len(req.ToSort))

    // Launch goroutines for each sub-array
    for _, subArray := range req.ToSort {
        go func(array []int) {
            sort.Slice(array, func(i, j int) bool { return array[i] < array[j] })
            results <- array
        }(subArray)
    }

    // Collect sorted sub-arrays from channel
    var sortedArrays [][]int
    for i := 0; i < len(req.ToSort); i++ {
        sortedArrays = append(sortedArrays, <-results)
    }

    endTime := time.Now()

    resp := SortResponse{
        SortedArrays: sortedArrays,
        TimeNs: endTime.Sub(startTime).Nanoseconds(),
    }

    json.NewEncoder(w).Encode(resp)
}

