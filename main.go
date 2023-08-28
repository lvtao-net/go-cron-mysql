package main

import (
    "database/sql"
    "fmt"
    "reflect"
    "time"

    "github.com/robfig/cron"
)

type Task struct {
    ID        int
    Name      string
    Schedule  string
    Status    string
    Module    string
    Method    string
    Arguments string
}

func main() {
    // 连接数据库
    db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/database")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer db.Close()

    // 初始化 cron 实例
    c := cron.New()

    // 读取任务列表
    tasks, err := getTasksFromDB(db)
    if err != nil {
        fmt.Println(err)
        return
    }

    // 注册任务
    for _, task := range tasks {
        if task.Status == "enabled" {
            c.AddFunc(task.Schedule, func() {
                // 动态调用任务方法
                callModuleMethod(task.Module, task.Method, task.Arguments)
            })
        }
    }

    // 启动 cron
    c.Start()

    // 等待任务执行完毕
    time.Sleep(5 * time.Minute)

    // 停止 cron
    c.Stop()
}

func getTasksFromDB(db *sql.DB) ([]Task, error) {
    // 从数据库中读取任务列表
    tasks := []Task{}
    rows, err := db.Query("SELECT * FROM tasks")
    if err != nil {
        return tasks, err
    }
    defer rows.Close()

    for rows.Next() {
        var task Task
        err := rows.Scan(&task.ID, &task.Name, &task.Schedule, &task.Status, &task.Module, &task.Method, &task.Arguments)
        if err != nil {
            return tasks, err
        }
        tasks = append(tasks, task)
    }

    return tasks, nil
}

func callModuleMethod(module string, method string, arguments string) {
    // 动态调用模块任务方法
    moduleValue := reflect.ValueOf(module)
    methodValue := moduleValue.MethodByName(method)
    if methodValue.IsValid() {
        argumentValues := make([]reflect.Value, 1)
        argumentValues[0] = reflect.ValueOf(arguments)
        methodValue.Call(argumentValues)
    } else {
        fmt.Printf("Method %s not found in module %s\n", method, module)
    }
}

// 示例模块任务方法
type ExampleModule struct{}

func (m ExampleModule) ExampleMethod(arguments string) {
    fmt.Printf("ExampleMethod called with arguments %s\n", arguments)
}
