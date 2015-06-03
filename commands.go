package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "io/ioutil"
    "encoding/json"
    "time"
    "strings"
    "strconv"
    "net/http"
    "bytes"

    "github.com/codegangsta/cli"
    "github.com/codeskyblue/go-sh"
)

var Commands = []cli.Command{
    commandPush,
}

var commandPush = cli.Command{
    Name:  "push",
    Usage: "Pushes json data to the mutinymon server",
    Description: `This command pushes a json string to the mutinymon server`,
    Action: doPush,
}

func debug(v ...interface{}) {
    if os.Getenv("DEBUG") != "" {
        log.Println(v...)
    }
}

func assert(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func doPush(c *cli.Context) {
    key, err := ioutil.ReadFile("/etc/mutinymon/mate.key")

    if err != nil {
        log.Fatal(key)
    }

    host, err := ioutil.ReadFile("/etc/mutinymon/mutinymon.host")

    if err != nil {
        log.Fatal(key)
    }

    keyString  := strings.TrimSpace(string(key))
    hostString := strings.TrimSpace(string(host))
    
    type Process struct {
        User string  `json:"user"`
        Comm string  `json:"comm"`
        Pcpu float64 `json:"pcpu"`
        Vsz  int     `json:"vsz"`
    }

    type Snapshot struct {
        UID                    string    `json:"uid"`
        Timestamp              int64     `json:"timestamp"`
        CPUCount               int       `json:"cpuCount"`
        CPULoadMin1            float64   `json:"cpuLoadMin1"`
        CPULoadMin5            float64   `json:"cpuLoadMin5"`
        CPULoadMin15           float64   `json:"cpuLoadMin15"`
        MemoryTotal            int       `json:"memoryTotal"`
        MemoryUsed             int       `json:"memoryUsed"`
        MemoryFree             int       `json:"memoryFree"`
        MemoryBuffersCacheUsed int       `json:"memoryBuffersCacheUsed"`
        MemoryBuffersCacheFree int       `json:"memoryBuffersCacheFree"`
        MemorySwapTotal        int       `json:"memorySwapTotal"`
        MemorySwapUsed         int       `json:"memorySwapUsed"`
        MemorySwapFree         int       `json:"memorySwapFree"`
        DiskTotal              int       `json:"diskTotal"`
        DiskUsed               int       `json:"diskUsed"`
        DiskFree               int       `json:"diskFree"`
        Processes              []Process `json:"processes"`
    }

    timestamp := time.Now().Unix()

    cmd1, err := exec.Command("/bin/grep", "-c", "^processor", "/proc/cpuinfo").Output()
    
    if err != nil {
        log.Fatal(err)
    }

    cpuCount, err := strconv.Atoi(strings.TrimSpace(string(cmd1)))

    if err != nil {
        log.Fatal(err)
    }

    cmd2, err := exec.Command("/bin/cat", "/proc/loadavg").Output()

    if err != nil {
        log.Fatal(err)
    }

    cpu := strings.Split(string(cmd2), " ")

    cpuLoad1Min, err := strconv.ParseFloat(cpu[0], 64)
    
    if err != nil {
        log.Fatal(err)
    }

    cpuLoad5Min, err := strconv.ParseFloat(cpu[1], 64)
    
    if err != nil {
        log.Fatal(err)
    }

    cpuLoad15Min, err := strconv.ParseFloat(cpu[2], 64)
    
    if err != nil {
        log.Fatal(err)
    }

    cmd3, err := sh.Command("/usr/bin/free", "-b").Command("/usr/bin/awk", "NR == 2 {printf $2 \" \" $3 \" \" $4 \" \"} NR == 3 {printf $3 \" \" $4 \" \"} NR ==4 {printf $2 \" \" $3 \" \" $4}").Output()

    if err != nil {
        log.Fatal(err)
    }

    memory := strings.Split(string(cmd3), " ")

    memoryTotal, err := strconv.Atoi(memory[0])
    
    if err != nil {
        log.Fatal(err)
    }

    memoryUsed, err := strconv.Atoi(memory[1])
    
    if err != nil {
        log.Fatal(err)
    }

    memoryFree, err := strconv.Atoi(memory[2])
    
    if err != nil {
        log.Fatal(err)
    }

    memoryBuffersCacheUsed, err := strconv.Atoi(memory[3])
    
    if err != nil {
        log.Fatal(err)
    }

    memoryBuffersCacheFree, err := strconv.Atoi(memory[4])
    
    if err != nil {
        log.Fatal(err)
    }

    memorySwapTotal, err := strconv.Atoi(memory[5])
    
    if err != nil {
        log.Fatal(err)
    }

    memorySwapUsed, err := strconv.Atoi(memory[6])
    
    if err != nil {
        log.Fatal(err)
    }

    memorySwapFree, err := strconv.Atoi(memory[7])
    
    if err != nil {
        log.Fatal(err)
    }

    cmd4, err := sh.Command("/bin/df", "-B1", "--total").Command("/usr/bin/awk", "/total/ {printf $2 \" \" $3 \" \" $4}").Output()

    if err != nil {
        log.Fatal(err)
    }

    disk := strings.Split(string(cmd4), " ")

    diskTotal, err := strconv.Atoi(disk[0])
    
    if err != nil {
        log.Fatal(err)
    }

    diskUsed, err := strconv.Atoi(disk[1])
    
    if err != nil {
        log.Fatal(err)
    }

    diskFree, err := strconv.Atoi(disk[2])
    
    if err != nil {
        log.Fatal(err)
    }

    cmd5, err := sh.Command("/bin/ps", "axo", "user,comm,pcpu,vsz", "--sort", "-pcpu").Command("/usr/bin/awk", "BEGIN{OFS=\":\"} NR>1 {printf $1 \" \" $2 \" \" $3 \" \" $4 \"@@\" }").Output()

    if err != nil {
        log.Fatal(err)
    }

    processLines := strings.Split(string(cmd5), "@@")

    //Pop last element
    processLines = processLines[:len(processLines)-1]

    processes := []Process {}

    for _, processLine := range processLines {

        processArray := strings.Split(string(processLine), " ")

        pcpuFloat, err := strconv.ParseFloat(processArray[2], 64)   

        if err != nil {
            log.Fatal(err)
        }

        vszInt, err := strconv.Atoi(processArray[3])

        if err != nil {
            log.Fatal(err)
        }

        user, comm, pcpu, vsz := processArray[0], processArray[1], pcpuFloat, vszInt

        process := Process{
            user,
            comm,
            pcpu,
            vsz,
        }

        processes = append(processes, process)
    }

    snapshot := Snapshot{
        strings.TrimSpace(string(key)), 
        timestamp,
        cpuCount,
        cpuLoad1Min, 
        cpuLoad5Min, 
        cpuLoad15Min, 
        memoryTotal, 
        memoryUsed, 
        memoryFree, 
        memoryBuffersCacheUsed, 
        memoryBuffersCacheFree, 
        memorySwapTotal, 
        memorySwapUsed, 
        memorySwapFree,
        diskTotal, 
        diskUsed, 
        diskFree, 
        processes,
    }

    output, err := json.Marshal(snapshot)

    url := hostString + "/push/" + keyString

    var jsonStr = []byte(string(output) )
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    
    resp, err := client.Do(req)
    
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()

    if resp.Status != "200 OK" {
        fmt.Println("Your request could not be processesd \n")
        fmt.Println("URL:")
        fmt.Println(url + "\n")
        fmt.Println("Payload:")
        fmt.Println(string(output),"\n")
        fmt.Println("response Status:")
        fmt.Println(resp.Status,"\n")
        fmt.Println("response Headers:")
        fmt.Println(resp.Header,"\n")
        fmt.Println("response Body:")
        body, _ := ioutil.ReadAll(resp.Body)
        fmt.Println(string(body),"\n")
    }
}