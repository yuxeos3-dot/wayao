package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GenerateAIContent POST /api/v1/domains/{id}/ai-content
func (app *App) GenerateAIContent(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}

	var req struct {
		KeywordType string `json:"keyword_type"`
		Keyword     string `json:"keyword"`
		Angle       string `json:"angle"`
		Lang        string `json:"lang"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}

	// get keyword info from domain
	if req.Keyword == "" {
		var kw, kwType string
		app.DB.QueryRow("SELECT primary_keyword, keyword_type FROM domains WHERE id=?", id).Scan(&kw, &kwType)
		if req.Keyword == "" {
			req.Keyword = kw
		}
		if req.KeywordType == "" {
			req.KeywordType = kwType
		}
	}
	if req.Lang == "" {
		req.Lang = "zh-TW"
	}
	if req.Angle == "" {
		req.Angle = "review"
	}

	// get Claude API key
	apiKey := getSetting(app.DB, "claude_api_key")
	if apiKey == "" {
		jsonError(w, 400, "claude_api_key not configured in settings")
		return
	}

	prompt := buildPrompt(req.KeywordType, req.Keyword, req.Angle, req.Lang)

	content, err := callClaudeAPI(apiKey, prompt)
	if err != nil {
		jsonError(w, 500, fmt.Sprintf("AI generation failed: %v", err))
		return
	}

	// save to contents table
	app.DB.Exec("INSERT OR IGNORE INTO contents(domain_id) VALUES(?)", id)
	app.DB.Exec("UPDATE contents SET body_content=?, ai_generated=1, updated_at=CURRENT_TIMESTAMP WHERE domain_id=?", content, id)

	jsonOK(w, map[string]interface{}{
		"content":    content,
		"char_count": len(content),
	})
}

func callClaudeAPI(apiKey, prompt string) (string, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"model":      "claude-sonnet-4-20250514",
		"max_tokens": 2000,
		"messages":   []map[string]string{{"role": "user", "content": prompt}},
	})

	req, _ := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}
	if len(result.Content) == 0 {
		return "", fmt.Errorf("empty response from API")
	}
	return result.Content[0].Text, nil
}

func buildPrompt(kwType, keyword, angle, lang string) string {
	prompts := map[string]string{
		"brand":     "你是台灣資深博彩評測員。請為「%s」撰寫600-800字的深度評測文章。繁體中文，結構：平台介紹→遊戲種類→存提款體驗→安全性評估→總結。不要包含標題，直接輸出正文。",
		"game":      "你是博彩遊戲策略分析師。請為「%s」撰寫700-900字的策略攻略文章。繁體中文，結構：遊戲數學原理→核心策略→進階技巧→資金管理→常見錯誤。不要包含標題。",
		"generic":   "你是台灣娛樂城評測網站編輯。請為「%s」撰寫500-700字的推薦排行文章。繁體中文，結構：市場概況→選擇標準→2-3個推薦平台→建議。不要包含標題。",
		"sports":    "你是台灣體育投注分析師。請為「%s」撰寫600-800字的投注指南。繁體中文，結構：玩法介紹→賠率說明→投注技巧→推薦平台。不要包含標題。",
		"promo":     "你是博彩優惠活動分析師。請為「%s」撰寫500-600字的優惠攻略。繁體中文，結構：什麼是這個優惠→條件詳解→如何最大化利用→注意事項。不要包含標題。",
		"payment":   "你是博彩支付方式分析師。請為「%s」撰寫600-700字的支付教學。繁體中文，結構：方式介紹→優缺點→操作步驟→常見問題。不要包含標題。",
		"strategy":  "你是博彩策略研究員。請為「%s」撰寫700-1000字的策略分析。繁體中文，結構：策略原理→操作步驟→優點→風險評估→適合人群。不要包含標題。",
		"app":       "你是手機APP評測師。請為「%s」撰寫500-600字的APP評測和教學。繁體中文，結構：APP優勢→功能介紹→安裝說明→使用技巧。不要包含標題。",
		"register":  "你是娛樂城使用指南專家。請為「%s」撰寫500-600字的開戶教學。繁體中文，結構：為什麼要開戶→準備什麼→詳細步驟→開戶後要做什麼。不要包含標題。",
		"region":    "你是台灣各地區博彩市場分析師。請為「%s」撰寫500-600字的在地化推薦。繁體中文，結構：地區玩家特點→推薦平台→支付方式→在地建議。不要包含標題。",
		"affiliate": "你是博彩代理系統分析師。請為「%s」撰寫600-700字的代理指南。繁體中文，結構：代理介紹→佣金計算→推廣方式→注意事項。不要包含標題。",
		"credit":    "你是信用版/現金版比較專家。請為「%s」撰寫600字的版本介紹比較文章。繁體中文，包含差異、適合對象、風險提示。不要包含標題。",
		"live":      "你是真人荷官遊戲專家。請為「%s」撰寫600字的真人荷官介紹。繁體中文，包含遊戲種類、直播品質、互動功能、推薦平台。不要包含標題。",
		"community": "你是博彩社群口碑整理者。請為「%s」撰寫500-600字的口碑評測。繁體中文，模擬PTT風格整理正反評價，給出客觀分析。不要包含標題。",
		"terms":     "你是博彩術語解釋專家。請為「%s」撰寫600-800字的術語解釋。繁體中文，結構：定義→詳細解釋→計算案例→相關術語→玩家須知。不要包含標題。",
	}
	tmpl, ok := prompts[kwType]
	if !ok {
		tmpl = prompts["brand"]
	}
	return fmt.Sprintf(tmpl, keyword)
}
