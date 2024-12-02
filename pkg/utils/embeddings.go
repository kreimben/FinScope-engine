package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/kreimben/FinScope-engine/pkg/logging"
)

const HF_API_URL = "https://api-inference.huggingface.co/pipeline/feature-extraction/thenlper/gte-small"

func GenerateEmbedding(apiKey string, text string) ([]float32, error) {
	// Initialize logger
	log := logging.NewLogger()

	log.WithField("text_length", len(text)).Debug("Generating embedding for text")

	// API 요청 데이터 준비
	payload := map[string]string{
		"inputs":  text,
		"options": `{"pooling": "mean", "normalize": true}`,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.WithError(err).Error("Failed to marshal payload")
		return nil, err
	}

	// HTTP 요청 생성
	req, err := http.NewRequest("POST", HF_API_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.WithError(err).Error("Failed to create HTTP request")
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-wait-for-model", "true")

	log.WithFields(map[string]interface{}{
		"url":    HF_API_URL,
		"method": "POST",
	}).Debug("Sending request to Hugging Face API")

	// 요청 전송
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error("Failed to send request to Hugging Face API")
		return nil, err
	}
	defer resp.Body.Close()

	// 응답 처리
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Failed to read response body")
		return nil, err
	}

	// 응답 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		log.WithFields(map[string]interface{}{
			"status_code": resp.StatusCode,
			"status":      resp.Status,
			"body":        string(body),
		}).Error("Received non-200 response from Hugging Face API")
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	var embedding []float32
	err = json.Unmarshal(body, &embedding)
	if err != nil {
		log.WithFields(map[string]interface{}{
			"error":         err,
			"response_body": string(body),
		}).Error("Failed to unmarshal response")
		return nil, err
	}

	log.WithFields(map[string]interface{}{
		"embedding_length": len(embedding),
		"text_length":      len(text),
	}).Debug("Successfully generated embedding")

	return embedding, nil
}

// CosineSimilarity 두 벡터 간의 코사인 유사도를 계산합니다
func CosineSimilarity(a, b []float32) (float32, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("vectors must be of same length, got %d and %d", len(a), len(b))
	}

	var dotProduct float32
	var normA float32
	var normB float32

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0, fmt.Errorf("vector normalization failed: zero vector detected")
	}

	similarity := dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))

	// 부동소수점 오차로 인해 1보다 약간 큰 값이 나올 수 있으므로 보정
	if similarity > 1 {
		similarity = 1
	} else if similarity < -1 {
		similarity = -1
	}

	return similarity, nil
}
