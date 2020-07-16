package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
)

func handleMetric(m map[string]interface{}){
    fmt.Println()
    fmt.Println("METRIC")
    fmt.Println("  name =", m["name"])
    fmt.Println("  ts =", m["timestamp"])
    fmt.Println("FIELDS")
    for k, v := range m["fields"].(map[string]interface{}){
        fmt.Println(" ", k, "=", v)
    }
    
    fmt.Println("TAGS")
    for k, v := range m["tags"].(map[string]interface{}){
        fmt.Println(" ", k, "=", v)
    }

}

func index(w http.ResponseWriter, r *http.Request) {
    //body, _ := r.GetBody()
    fmt.Println("Request incoming!")
    if r.Body != nil {    
        var b bytes.Buffer
        var dest io.Writer = &b
        _, _ = io.Copy(dest, r.Body)
        //fmt.Println(string(b.Bytes()))
        var result map[string]interface{}
        json.Unmarshal(b.Bytes(), &result)
	_, isBatch := result["metrics"]
        if isBatch {
            for _, m := range result["metrics"].([]interface{}) {
                handleMetric(m.(map[string]interface{}))
            }
        }
    }
}

func main(){
    http.HandleFunc("/", index)
    fmt.Println("Starting server")
    log.Fatal(http.ListenAndServe("vpn:8080", nil))
}
