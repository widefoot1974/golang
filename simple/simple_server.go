package main

import (
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
    "sync/atomic"
    "time"
)

// 여러 요청에 대한 접속 횟수(원자적 증가 사용)
var clientCount int64

func main() {
    addr := ":8080"

    // 1. 서버 시작 시점에 호스트 이름, 서버 IP 출력
    hostname, err := os.Hostname()
    if err != nil {
        log.Fatalf("호스트 이름 가져오기 실패: %v", err)
    }

    // 서버 IP 목록
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        log.Fatalf("IP 주소 가져오기 실패: %v", err)
    }
    var ipList []string
    for _, a := range addrs {
        if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
            if ip := ipNet.IP.To4(); ip != nil {
                ipList = append(ipList, ip.String())
            }
        }
    }

    log.Println("서버 시작 중...")
    log.Printf("Hostname: %s\n", hostname)
    log.Printf("서버 IP 목록: %v\n", ipList)
    log.Printf("바인딩 포트: %s\n", addr)

    // 2. 핸들러 정의
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // 요청마다 접속 횟수 원자적 증가
        newCount := atomic.AddInt64(&clientCount, 1)

        // 200 OK 상태 전송 (Go 기본값도 200이지만 명시적으로 호출)
        w.WriteHeader(http.StatusOK)

        // 클라이언트 IP 파싱 (IP:Port 형태에서 IP만 추출)
        clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            // SplitHostPort가 실패하면 그냥 r.RemoteAddr 전체를 사용
            clientIP = r.RemoteAddr
        }

        // 현재 시간 (YYYY-MM-DD HH:MM:SS)
        currentTime := time.Now().Format("2006-01-02 15:04:05")

        // 응답 메시지
        fmt.Fprintf(w, "Current Time: %s\n", currentTime)
        fmt.Fprintf(w, "Access Count: %d\n", newCount)
        fmt.Fprintf(w, "Client IP: %s\n", clientIP)
        fmt.Fprintf(w, "Hostname: %s\n", hostname)
        for _, ip := range ipList {
            fmt.Fprintf(w, "Server IP: %s\n", ip)
        }
    })

    // 3. 서버 실행
    log.Fatal(http.ListenAndServe(addr, nil))
}

