package main

import (
    "fmt"
    "net/http"
    "os/exec"
    "strings"
    "io/ioutil"
    "encoding/json"
    "flag"
)

type Container struct {
    Name        string
    Ports       string
    Status      string
}

type Config struct {
    Port   string `json:"port"`
    Listen string `json:"listen"`
}

func main() {
    configPath := flag.String("config", "config.json", "path to config file")
    flag.Parse()

    config := Config{}
    configFile, err := ioutil.ReadFile(*configPath)
    if err != nil {
        fmt.Println("Error reading config file:", err)
        config.Port = "8080"
        config.Listen = "0.0.0.0"
    } else {
        json.Unmarshal(configFile, &config)
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        out, err := exec.Command("bash", "-c", "docker ps --format \"{{json .Names}} {{json .Ports}} {{json .Status}}\"").Output()
        if err != nil {
            fmt.Fprintf(w, "Error: %s", err)
            return
        }
        text := string(out)

        lines := strings.Split(text, "\n")
        var containers []Container
        for _, line := range lines {
            parts := strings.Split(line, "\"")
            if len(parts) > 1 {
                container := Container{
                    Name:   parts[1],
                    Ports:  parts[3],
                    Status: parts[5],
                }
                containers = append(containers, container)
            }
        }

        fmt.Fprintf(w, "<html><head><title>Containers</title></head><body>")
        fmt.Fprintf(w, "<table border='1'>")
        fmt.Fprintf(w, "<tr><th>Name</th><th>Ports</th><th>Status</th></tr>")
        for _, container := range containers {
            fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td>%s</td></tr>", container.Name, container.Ports, container.Status)
        }
        fmt.Fprintf(w, "</table>")
        fmt.Fprintf(w, "</body></html>")
    })

    fmt.Println("Listening on", config.Listen+":"+config.Port)
    http.ListenAndServe(config.Listen+":"+config.Port, nil)
}
