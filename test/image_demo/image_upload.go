package main

import (
    "fmt"
    "net/http"
    "os"
    "os/exec"
    "log"
)

func main() {

    logFile, err := os.OpenFile("upload_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Error opening log file:", err)
        return
    }
    defer logFile.Close()

    // 设置日志输出到文件
    log.SetOutput(logFile)
    log.Println("Log file opened successfully.")

    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
    http.HandleFunc("/", serveIndex)
    http.HandleFunc("/upload", handleUpload)

    fmt.Println("Server is running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }
    

    // 解析表单数据
    err := r.ParseForm()
    if err != nil {
        log.Println("Error parsing form: %v", err)
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    imageName := r.FormValue("imageName")
    imageTag := r.FormValue("imageTag")
   
    if imageName == "" {
        log.Println("Image Name is empty")
    }
    if imageTag == "" {
        log.Println("Image Tag is empty")
    }

 
    log.Printf("Image Name: %s", imageName)
    log.Printf("Image Tag: %s", imageTag)

    username := r.FormValue("username")
    password := r.FormValue("password")
    harborURL := "47.121.126.55:8000" // 替换为您的 Harbor 地址
    project := "production"              // 替换为您的项目名

    // 构建 Docker 镜像
/*    cmd := exec.Command("docker", "build", "-t", fmt.Sprintf("%s/%s:%s", harborURL+"/"+project, imageName, imageTag), ".")
    if err := cmd.Run(); err != nil {
        http.Error(w, "Failed to build Docker image", http.StatusInternalServerError)
        return
    }*/

    // 登录到 Harbor
    cmd := exec.Command("docker", "login", harborURL, "-u", username, "-p", password)
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Printf("Failed to login to Harbor: %v, Output: %s", err, output)
        http.Error(w, "Failed to login to Harbor", http.StatusInternalServerError)
        return
    }

    // 推送 Docker 镜像
    cmd = exec.Command("docker", "push", fmt.Sprintf("%s/%s/%s:%s", harborURL, project, imageName, imageTag))
    imageReference := fmt.Sprintf("%s/%s/%s:%s", harborURL, project, imageName, imageTag)
    log.Printf("Pushing Docker image: %s", imageReference)
    output, err = cmd.CombinedOutput()
    if err != nil {
        log.Printf("Failed to push Docker image: %v, Output: %s", err, output)
        http.Error(w, "Failed to push Docker image", http.StatusInternalServerError)
        return
    }
    log.Printf("Docker image %s/%s:%s uploaded successfully!", harborURL+"/"+project, imageName, imageTag)
    fmt.Fprintf(w, "Docker image %s/%s:%s uploaded successfully!", harborURL+"/"+project, imageName, imageTag)
}
