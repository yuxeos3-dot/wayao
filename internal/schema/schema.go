package schema

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/yuxeos3-dot/wayao/internal/variation"
)

type SiteData struct {
	Domain         string
	SiteName       string
	MetaTitle      string
	MetaDesc       string
	MainContent    string
	BrandColor     string
	RedirectURL    string
	PrimaryKeyword string
	KeywordType    string
	FAQJson        string
	ExtraFields    map[string]interface{}
}

func GenerateSchema(site SiteData) string {
	schemas := []interface{}{}

	// 1. WebSite schema
	schemas = append(schemas, map[string]interface{}{
		"@type":       "WebSite",
		"name":        site.SiteName,
		"url":         "https://" + site.Domain,
		"description": site.MetaDesc,
	})

	// 2. Organization schema
	schemas = append(schemas, map[string]interface{}{
		"@type": "Organization",
		"name":  site.SiteName,
		"url":   "https://" + site.Domain,
		"logo": map[string]interface{}{
			"@type": "ImageObject",
			"url":   "https://" + site.Domain + "/images/logo.png",
		},
	})

	// 3. Article/Review schema based on keyword type
	switch site.KeywordType {
	case "brand", "generic":
		rating := getStr(site.ExtraFields, "overall_rating", "9.0")
		schemas = append(schemas, map[string]interface{}{
			"@type":       "Review",
			"itemReviewed": map[string]interface{}{
				"@type": "Organization",
				"name":  site.PrimaryKeyword,
			},
			"reviewRating": map[string]interface{}{
				"@type":       "Rating",
				"ratingValue": rating,
				"bestRating":  "10",
				"worstRating": "1",
			},
			"author": map[string]interface{}{
				"@type": "Organization",
				"name":  site.SiteName,
			},
			"datePublished": variation.RandomPublishDate(site.Domain),
			"dateModified":  time.Now().Format("2006-01-02"),
		})
	case "game", "strategy":
		schemas = append(schemas, map[string]interface{}{
			"@type":         "Article",
			"headline":      site.MetaTitle,
			"description":   site.MetaDesc,
			"datePublished": variation.RandomPublishDate(site.Domain),
			"dateModified":  time.Now().Format("2006-01-02"),
			"author": map[string]interface{}{
				"@type": "Organization",
				"name":  site.SiteName,
			},
		})
	case "sports":
		schemas = append(schemas, map[string]interface{}{
			"@type":         "SportsEvent",
			"name":          site.PrimaryKeyword,
			"description":   site.MetaDesc,
			"url":           "https://" + site.Domain,
		})
	}

	// 4. FAQ schema from faq_json
	if site.FAQJson != "" && site.FAQJson != "[]" {
		var faqs []struct {
			Q string `json:"q"`
			A string `json:"a"`
		}
		if json.Unmarshal([]byte(site.FAQJson), &faqs) == nil && len(faqs) > 0 {
			faqItems := make([]map[string]interface{}, len(faqs))
			for i, faq := range faqs {
				faqItems[i] = map[string]interface{}{
					"@type": "Question",
					"name":  faq.Q,
					"acceptedAnswer": map[string]interface{}{
						"@type": "Answer",
						"text":  faq.A,
					},
				}
			}
			schemas = append(schemas, map[string]interface{}{
				"@type":      "FAQPage",
				"mainEntity": faqItems,
			})
		}
	}

	// 5. BreadcrumbList
	schemas = append(schemas, map[string]interface{}{
		"@type": "BreadcrumbList",
		"itemListElement": []map[string]interface{}{
			{"@type": "ListItem", "position": 1, "name": "首頁", "item": "https://" + site.Domain + "/"},
			{"@type": "ListItem", "position": 2, "name": site.PrimaryKeyword},
		},
	})

	// wrap in @graph
	wrapper := map[string]interface{}{
		"@context": "https://schema.org",
		"@graph":   schemas,
	}

	// shuffle fields for fingerprint variation
	jsonBytes, _ := json.MarshalIndent(wrapper, "", "  ")
	result := variation.ShuffleJSONFields(string(jsonBytes), site.Domain)

	return result
}

func getStr(m map[string]interface{}, key, fallback string) string {
	if m == nil {
		return fallback
	}
	if v, ok := m[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return fallback
}

// GenerateSupportPages returns content for auto-generated pages
func GenerateSupportPages(siteName, domain string) map[string]string {
	year := time.Now().Year()
	pages := map[string]string{
		"about": fmt.Sprintf(`---
title: "關於我們"
url: "/about/"
---
# 關於 %s

%s 是一個獨立的資訊平台，致力於為用戶提供最全面、最客觀的線上娛樂城資訊。

我們的團隊由經驗豐富的分析師和內容作者組成，確保每篇內容都經過仔細研究和驗證。

## 我們的使命

提供透明、可信賴的評測資訊，幫助用戶做出更明智的選擇。

&copy; %d %s. All rights reserved.`, siteName, siteName, year, siteName),

		"privacy": fmt.Sprintf(`---
title: "隱私政策"
url: "/privacy/"
---
# %s 隱私政策

最後更新日期：%d年

## 資訊收集
我們可能收集的資訊包括：瀏覽數據、IP地址、瀏覽器類型。

## Cookie使用
本站使用Cookie改善用戶體驗。您可以在瀏覽器設定中管理Cookie偏好。

## 第三方連結
本站可能包含第三方網站的連結。我們不對第三方網站的隱私政策負責。

## 聯繫我們
如有任何隱私相關問題，請透過聯繫頁面與我們取得聯繫。`, siteName, year),

		"disclaimer": fmt.Sprintf(`---
title: "免責聲明"
url: "/disclaimer/"
---
# 免責聲明

## 內容聲明
%s 提供的所有內容僅供參考和教育目的。我們不保證資訊的準確性或完整性。

## 風險提示
線上博弈涉及財務風險。請務必了解您所在地區的相關法律法規，並負責任地參與。

## 年齡限制
本站內容僅適合18歲以上成年人。未成年人請勿瀏覽本站內容。`, siteName),

		"responsible-gambling": fmt.Sprintf(`---
title: "負責任博彩"
url: "/responsible-gambling/"
---
# 負責任博彩

%s 提倡負責任的博彩行為。

## 設定限額
- 設定每日/每週/每月的存款和投注限額
- 嚴格遵守自己設定的限額

## 求助資源
如果您或您認識的人有賭博問題，請聯繫：
- 台灣戒賭專線
- GamCare: www.gamcare.org.uk
- BeGambleAware: www.begambleaware.org`, siteName),

		"contact": fmt.Sprintf(`---
title: "聯繫我們"
url: "/contact/"
---
# 聯繫 %s

如果您有任何問題或建議，歡迎透過以下方式與我們聯繫。

我們將在 1-2 個工作天內回覆您的訊息。`, siteName),

		"methodology": fmt.Sprintf(`---
title: "評測方法"
url: "/methodology/"
---
# %s 評測方法

## 評測流程

每份評測報告都基於以下標準化流程：

### 第一步：帳號建立
使用真實資料完成開戶，記錄開戶流程的複雜程度，測試身份驗證要求。

### 第二步：存款測試
測試所有可用支付方式，記錄到帳速度，核查最低/最高限額。

### 第三步：遊戲體驗
測試各類遊戲（真人、電子、體育），評估遊戲流暢度和公平性。

### 第四步：出金測試
申請多次出金，記錄處理時間，測試各種提款方式。

### 第五步：客服評估
測試線上客服回應速度，評估問題解決能力。

## 評分標準

| 維度 | 權重 |
|------|------|
| 安全性與合法性 | 25%% |
| 遊戲種類與品質 | 20%% |
| 存提款體驗 | 25%% |
| 優惠方案 | 15%% |
| 客服品質 | 15%% |`, siteName),
	}
	return pages
}

func GenerateRobotsTxt(domain string) string {
	return fmt.Sprintf(`User-agent: *
Allow: /
Disallow: /api/

Sitemap: https://%s/sitemap.xml`, domain)
}

func GenerateSitemapIndex(domain string, pages []string) string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
`)
	for _, p := range pages {
		sb.WriteString(fmt.Sprintf(`  <url>
    <loc>https://%s%s</loc>
    <lastmod>%s</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
`, domain, p, time.Now().Format("2006-01-02")))
	}
	sb.WriteString("</urlset>")
	return sb.String()
}
