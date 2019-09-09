package main

import (
    "flag"
    "sync"
)

var (
    cmd = flag.String("cmd", "list", "Start | List | Remove  on docker container")
    container_id = flag.String("id", "", "container_id")
    wait_group sync.WaitGroup
)

func main(){
    flag.Parse()
    defer wait_group.Wait()
    wait_group.Add(1)
    go func (cmd string, id string) {
        defer wait_group.Done()
        Run(cmd, id)
    }(*cmd, *container_id)
}