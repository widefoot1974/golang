package main

import (
    "fmt"
    "net"
    "net/http"
    "os"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // 명시적으로 200 OK 설정
        w.WriteHeader(http.StatusOK)

        // 호스트 이름 가져오기
        hostname, _ := os.Hostname()
        fmt.Fprintf(w, "Hostname: %s\n", hostname)

        // IP 주소 목록 가져오기 (IPv4만 필터링)
        addrs, _ := net.InterfaceAddrs()
        for _, addr := range addrs {
            if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
                if ip := ipnet.IP.To4(); ip != nil {
                    fmt.Fprintf(w, "IP: %s\n", ip)
                }
            }
        }
    })

    // 8080 포트에서 서버 시작
    http.ListenAndServe(":8080", nil)
}

