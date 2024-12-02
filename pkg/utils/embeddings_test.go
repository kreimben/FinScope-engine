package utils

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGenerateEmbedding(t *testing.T) {
	// Load environment variables from .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("HUGGINGFACE_API_KEY")
	if apiKey == "" {
		t.Fatal("HUGGINGFACE_API_KEY not found in environment")
	}

	testCases := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{
			name:    "Simple text",
			text:    "This is a test sentence.",
			wantErr: false,
		},
		{
			name:    "Empty text",
			text:    "",
			wantErr: false,
		},
		{
			name:    "Long text",
			text:    "This is a longer test sentence that contains multiple words and should still work fine with the embedding model. The model should be able to handle texts of various lengths.",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			embedding, err := GenerateEmbedding(apiKey, tc.text)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, embedding)
			assert.Greater(t, len(embedding), 0)

			// GTE-small 모델은 384 차원의 임베딩을 생성합니다
			assert.Equal(t, 384, len(embedding))

			// 임베딩 값들이 적절한 범위 내에 있는지 확인
			for _, value := range embedding {
				assert.True(t, value >= -1 && value <= 1)
			}
		})
	}
}

func TestSimilarityCheck(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("HUGGINGFACE_API_KEY")
	if apiKey == "" {
		t.Fatal("HUGGINGFACE_API_KEY not found in environment")
	}

	testCases := []struct {
		name     string
		text1    string
		text2    string
		minScore float32 // 예상되는 최소 유사도 점수
	}{
		{
			name:     "Identical texts",
			text1:    "The stock market had a significant drop today.",
			text2:    "The stock market had a significant drop today.",
			minScore: 0.99, // 동일한 텍스트는 거의 1에 가까운 유사도를 가져야 함
		},
		{
			name:     "Similar texts",
			text1:    "The stock market had a significant drop today.",
			text2:    "Today, stocks experienced a major decline.",
			minScore: 0.7, // 비슷한 의미의 텍스트는 적당한 유사도를 가져야 함
		},
		{
			name:     "Different but related texts",
			text1:    "The stock market had a significant drop today.",
			text2:    "Investors are concerned about market conditions.",
			minScore: 0.5, // 관련은 있지만 다른 텍스트는 중간 정도의 유사도를 가져야 함
		},
		{
			name:     "Completely different texts",
			text1:    "The stock market had a significant drop today.",
			text2:    "The weather is beautiful today.",
			minScore: -0.1, // 전혀 다른 주제의 텍스트는 낮은 유사도를 가져야 함
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 첫 번째 텍스트의 임베딩 생성
			embedding1, err := GenerateEmbedding(apiKey, tc.text1)
			if !assert.NoError(t, err) {
				return
			}

			// 두 번째 텍스트의 임베딩 생성
			embedding2, err := GenerateEmbedding(apiKey, tc.text2)
			if !assert.NoError(t, err) {
				return
			}

			// 유사도 계산
			similarity, err := CosineSimilarity(embedding1, embedding2)
			if !assert.NoError(t, err) {
				return
			}

			// 결과 출력
			t.Logf("Similarity between texts: %.4f", similarity)
			t.Logf("Text 1: %s", tc.text1)
			t.Logf("Text 2: %s", tc.text2)

			// 유사도가 예상 범위 내에 있는지 확인
			assert.GreaterOrEqual(t, float64(similarity), float64(tc.minScore),
				"Similarity score %.4f is lower than expected minimum %.4f", similarity, tc.minScore)
			assert.LessOrEqual(t, float64(similarity), 1.0,
				"Similarity score %.4f is higher than maximum 1.0", similarity)
		})
	}
}

func TestCosineSimilarity(t *testing.T) {
	testCases := []struct {
		name    string
		a       []float32
		b       []float32
		want    float32
		wantErr bool
	}{
		{
			name:    "Identical vectors",
			a:       []float32{1, 0, 0},
			b:       []float32{1, 0, 0},
			want:    1,
			wantErr: false,
		},
		{
			name:    "Orthogonal vectors",
			a:       []float32{1, 0, 0},
			b:       []float32{0, 1, 0},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Opposite vectors",
			a:       []float32{1, 0, 0},
			b:       []float32{-1, 0, 0},
			want:    -1,
			wantErr: false,
		},
		{
			name:    "Different lengths",
			a:       []float32{1, 0, 0},
			b:       []float32{1, 0},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Zero vector",
			a:       []float32{0, 0, 0},
			b:       []float32{1, 0, 0},
			want:    0,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CosineSimilarity(tc.a, tc.b)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.InDelta(t, tc.want, got, 1e-6)
		})
	}
}

func TestArticleStockRelevance(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("HUGGINGFACE_API_KEY")
	if apiKey == "" {
		t.Fatal("HUGGINGFACE_API_KEY not found in environment")
	}

	// 테스트할 기사들
	articles := []struct {
		title          string
		content        string
		expectedScores map[string]struct {
			min float32
			max float32
		}
	}{
		{
			title:   "Tesla's Q4 Earnings Beat Expectations",
			content: `Tesla reported strong fourth-quarter earnings...`,
			expectedScores: map[string]struct {
				min float32
				max float32
			}{
				"TSLA": {min: 0.6, max: 1.0}, // 관련 종목
				"AAPL": {min: 0.0, max: 0.4}, // 무관 종목
				"MSFT": {min: 0.0, max: 0.4}, // 무관 종목
			},
		},
		{
			title: "Recipe: How to Make Perfect Sushi",
			content: `Start with high-quality sushi rice. Wash the rice thoroughly...
			Add rice vinegar and mix well. Prepare your favorite fish and vegetables...`,
			expectedScores: map[string]struct {
				min float32
				max float32
			}{
				"TSLA": {min: 0.0, max: 0.2}, // 완전 무관
				"AAPL": {min: 0.0, max: 0.2}, // 완전 무관
				"MSFT": {min: 0.0, max: 0.2}, // 완전 무관
			},
		},
		{
			title: "Gardening Tips for Spring",
			content: `As spring approaches, it's time to prepare your garden.
			Start by testing your soil pH and adding necessary amendments.
			Choose appropriate plants for your climate zone.`,
			expectedScores: map[string]struct {
				min float32
				max float32
			}{
				"TSLA": {min: 0.0, max: 0.2},
				"AAPL": {min: 0.0, max: 0.2},
				"MSFT": {min: 0.0, max: 0.2},
			},
		},
	}

	// 종목 정보
	stocks := []struct {
		symbol      string
		description string
		minScore    float32
		maxScore    float32 // 무관한 기사의 경우 이 점수를 넘지 않아야 함
	}{
		{
			symbol: "TSLA",
			description: `Tesla, Inc. designs, develops, manufactures, leases, and sells electric vehicles, 
			and energy generation and storage systems. The company operates through automotive sales and energy generation segments.`,
			minScore: 0.3,  // 관련 기사의 최소 점수
			maxScore: 0.15, // 무관한 기사의 최대 점수
		},
		{
			symbol: "AAPL",
			description: `Apple Inc. designs, manufactures, and markets smartphones, personal computers, tablets, 
			wearables, and accessories worldwide. Its products include iPhone, Mac, iPad, and wearables like Apple Watch.`,
			minScore: 0.3,
			maxScore: 0.15,
		},
		{
			symbol: "MSFT",
			description: `Microsoft Corporation develops, licenses, and supports software, services, devices, 
			and solutions worldwide. The company operates through cloud computing, software licensing, and hardware manufacturing.`,
			minScore: 0.3,
			maxScore: 0.15,
		},
	}

	// 각 기사와 모든 종목 정보 간의 유사도 테스트
	for _, article := range articles {
		t.Run(article.title, func(t *testing.T) {
			articleEmbedding, err := GenerateEmbedding(apiKey, article.title+" "+article.content)
			if !assert.NoError(t, err) {
				return
			}

			for _, stock := range stocks {
				descEmbedding, err := GenerateEmbedding(apiKey, stock.description)
				if !assert.NoError(t, err) {
					continue
				}

				similarity, err := CosineSimilarity(articleEmbedding, descEmbedding)
				if !assert.NoError(t, err) {
					continue
				}

				expectedScore := article.expectedScores[stock.symbol]
				assert.GreaterOrEqual(t, float64(similarity), float64(expectedScore.min),
					"Similarity too low for %s: got %.4f, want >= %.4f", stock.symbol, similarity, expectedScore.min)
				assert.LessOrEqual(t, float64(similarity), float64(expectedScore.max),
					"Similarity too high for %s: got %.4f, want <= %.4f", stock.symbol, similarity, expectedScore.max)

				t.Logf("%s similarity with %s: %.4f (expected range: %.4f-%.4f)",
					article.title, stock.symbol, similarity, expectedScore.min, expectedScore.max)
			}
		})
	}
}
