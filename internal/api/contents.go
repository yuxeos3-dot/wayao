package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *App) GetContent(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}

	// auto-create if not exists
	app.DB.Exec("INSERT OR IGNORE INTO contents(domain_id) VALUES(?)", id)

	// query all content columns matching db.go schema
	rows, err := app.DB.Query(`SELECT
		id, domain_id, keyword_type, target_keyword, page_title, meta_desc, h1,
		brand_name, brand_color, cta_text, cta_sub, hero_title, hero_subtitle,
		feature_1_icon, feature_1_title, feature_1_desc,
		feature_2_icon, feature_2_title, feature_2_desc,
		feature_3_icon, feature_3_title, feature_3_desc,
		intro_text, body_content, conclusion,
		faq_title, faq_items, extra_data,
		author_name, author_title, author_bio, author_avatar,
		last_updated, last_updated_iso, review_count,
		trust_badges, disclosure, disclaimer,
		related_pages, content_angle, status, ai_generated
		FROM contents WHERE domain_id=?`, id)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	content := map[string]interface{}{}
	if rows.Next() {
		var cid, domainID int64
		var keywordType, targetKW, pageTitle, metaDesc, h1 string
		var brandName, brandColor, ctaText, ctaSub, heroTitle, heroSubtitle string
		var f1Icon, f1Title, f1Desc, f2Icon, f2Title, f2Desc, f3Icon, f3Title, f3Desc string
		var introText, bodyContent, conclusion string
		var faqTitle, faqItems, extraData string
		var authorName, authorTitle, authorBio, authorAvatar string
		var lastUpdated, lastUpdatedISO, reviewCount string
		var trustBadges, disclosure, disclaimer string
		var relatedPages, contentAngle, status string
		var aiGenerated int

		rows.Scan(&cid, &domainID, &keywordType, &targetKW, &pageTitle, &metaDesc, &h1,
			&brandName, &brandColor, &ctaText, &ctaSub, &heroTitle, &heroSubtitle,
			&f1Icon, &f1Title, &f1Desc, &f2Icon, &f2Title, &f2Desc, &f3Icon, &f3Title, &f3Desc,
			&introText, &bodyContent, &conclusion,
			&faqTitle, &faqItems, &extraData,
			&authorName, &authorTitle, &authorBio, &authorAvatar,
			&lastUpdated, &lastUpdatedISO, &reviewCount,
			&trustBadges, &disclosure, &disclaimer,
			&relatedPages, &contentAngle, &status, &aiGenerated)

		content = map[string]interface{}{
			"id": cid, "domain_id": domainID,
			"keyword_type": keywordType, "target_keyword": targetKW,
			"page_title": pageTitle, "meta_desc": metaDesc, "h1": h1,
			"brand_name": brandName, "brand_color": brandColor,
			"cta_text": ctaText, "cta_sub": ctaSub,
			"hero_title": heroTitle, "hero_subtitle": heroSubtitle,
			"feature_1_icon": f1Icon, "feature_1_title": f1Title, "feature_1_desc": f1Desc,
			"feature_2_icon": f2Icon, "feature_2_title": f2Title, "feature_2_desc": f2Desc,
			"feature_3_icon": f3Icon, "feature_3_title": f3Title, "feature_3_desc": f3Desc,
			"intro_text": introText, "body_content": bodyContent, "conclusion": conclusion,
			"faq_title": faqTitle, "faq_items": faqItems, "extra_data": extraData,
			"author_name": authorName, "author_title": authorTitle,
			"author_bio": authorBio, "author_avatar": authorAvatar,
			"last_updated": lastUpdated, "last_updated_iso": lastUpdatedISO,
			"review_count": reviewCount,
			"trust_badges": trustBadges, "disclosure": disclosure, "disclaimer": disclaimer,
			"related_pages": relatedPages, "content_angle": contentAngle,
			"status": status, "ai_generated": aiGenerated == 1,
		}
	}

	// also fetch domain info for context
	var domain, kwType, pk string
	app.DB.QueryRow("SELECT domain, keyword_type, primary_keyword FROM domains WHERE id=?", id).Scan(&domain, &kwType, &pk)

	jsonOK(w, map[string]interface{}{
		"content":         content,
		"domain":          domain,
		"keyword_type":    kwType,
		"primary_keyword": pk,
	})
}

func (app *App) SaveContent(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}

	// allowed fields must match actual db.go contents table columns
	allowed := map[string]bool{
		"keyword_type": true, "target_keyword": true, "page_title": true, "meta_desc": true,
		"h1": true, "brand_name": true, "brand_color": true, "cta_text": true, "cta_sub": true,
		"hero_title": true, "hero_subtitle": true,
		"feature_1_icon": true, "feature_1_title": true, "feature_1_desc": true,
		"feature_2_icon": true, "feature_2_title": true, "feature_2_desc": true,
		"feature_3_icon": true, "feature_3_title": true, "feature_3_desc": true,
		"intro_text": true, "body_content": true, "conclusion": true,
		"faq_title": true, "faq_items": true, "extra_data": true,
		"author_name": true, "author_title": true, "author_bio": true, "author_avatar": true,
		"last_updated": true, "last_updated_iso": true, "review_count": true,
		"trust_badges": true, "disclosure": true, "disclaimer": true,
		"related_pages": true, "content_angle": true, "status": true,
		"ai_generated": true,
	}

	// ensure row exists
	app.DB.Exec("INSERT OR IGNORE INTO contents(domain_id) VALUES(?)", id)

	for k, v := range req {
		if !allowed[k] {
			continue
		}
		app.DB.Exec(fmt.Sprintf("UPDATE contents SET %s=?, updated_at=CURRENT_TIMESTAMP WHERE domain_id=?", k), v, id)
	}

	// update content hash
	data, _ := json.Marshal(req)
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	app.DB.Exec("UPDATE contents SET content_hash=?, updated_at=CURRENT_TIMESTAMP WHERE domain_id=?", hash[:16], id)

	jsonOK(w, "saved")
}
