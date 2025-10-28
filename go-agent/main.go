package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

var (
	serverURL = getEnv("SERVER_URL", "https://backend:8443/metrics")
	secretKey = getEnv("JWT_SECRET", "super_secret_key")
)

func main() {
	fmt.Println("🚀 Metrics agent started...")
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}

	for {
		sendMetrics(client)
		time.Sleep(5 * time.Second)
	}
}

func sendMetrics(client *http.Client) {
	cpuPercent, _ := cpu.Percent(0, false)
	memStats, _ := mem.VirtualMemory()

	data := map[string]interface{}{
		"cpu_usage": cpuPercent[0],
		"ram_usage": memStats.UsedPercent,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	body, _ := json.Marshal(data)
	token := generateJWT()

	req, _ := http.NewRequest("POST", serverURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Error sending metrics:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("❌ Error sending metrics: server responded with", resp.Status)
		return
	}

	fmt.Println("✅ Sent metrics:", data)
}

func generateJWT() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(10 * time.Minute).Unix(),
		"iss": "go-agent",
	})
	tokenString, _ := token.SignedString([]byte(secretKey))
	return tokenString
}

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
