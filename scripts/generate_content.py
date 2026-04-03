#!/usr/bin/env python3
"""
AI內容生成腳本 - 使用Claude API為域名生成SEO內容
用法: python generate_content.py --domain-id 1 --db ../data/brandsite.db --api-key sk-ant-xxx
"""
import argparse
import json
import sqlite3
import sys

try:
    import anthropic
except ImportError:
    print("pip install anthropic")
    sys.exit(1)

PROMPTS = {
    "brand": """你是台灣線上娛樂城的專業評測作者。請為「{keyword}」撰寫一篇完整的評測內容。
要求：繁體中文、2000字以上、包含優缺點、FAQ、註冊教程、出入金方式、遊戲種類介紹。
語氣專業客觀，像一位真正體驗過的玩家在分享心得。""",

    "game": """你是線上娛樂城遊戲攻略專家。請為「{keyword}」撰寫一篇完整的遊戲攻略。
要求：繁體中文、1500字以上、包含基本規則、進階策略、常見錯誤、FAQ。""",

    "sports": """你是體育投注分析專家。請為「{keyword}」撰寫一篇體育投注指南。
要求：繁體中文、1500字以上、包含賽事介紹、投注策略、賠率分析、FAQ。""",

    "generic": """你是娛樂城排名評測專家。請為「{keyword}」撰寫一篇排名推薦文章。
要求：繁體中文、2000字以上、包含Top 5推薦、評分標準、比較表格、FAQ。""",

    "promo": """你是線上娛樂城優惠活動分析專家。請為「{keyword}」撰寫一篇優惠活動匯總。
要求：繁體中文、1500字以上、包含活動詳情、申請步驟、流水要求、注意事項、FAQ。""",

    "payment": """你是線上支付專家。請為「{keyword}」撰寫一篇完整的存提款教程。
要求：繁體中文、1500字以上、包含支付方式對比、步驟教程、手續費說明、到帳時間、FAQ。""",

    "affiliate": """你是線上娛樂城代理招募專家。請為「{keyword}」撰寫一篇代理推廣指南。
要求：繁體中文、1500字以上、包含佣金方案、申請流程、推廣技巧、FAQ。""",

    "strategy": """你是投注策略專家。請為「{keyword}」撰寫一篇投注策略攻略。
要求：繁體中文、1500字以上、包含策略原理、實戰示範、風險控制、FAQ。""",

    "app": """你是APP評測專家。請為「{keyword}」撰寫一篇APP下載安裝指南。
要求：繁體中文、1500字以上、包含下載方式、安裝步驟、功能介紹、常見問題、FAQ。""",

    "register": """你是線上註冊專家。請為「{keyword}」撰寫一篇完整的註冊開戶教程。
要求：繁體中文、1500字以上、包含註冊步驟、驗證流程、首存教程、客服聯繫方式、FAQ。""",

    "region": """你是台灣本地化博彩推薦專家。請為「{keyword}」撰寫一篇地區推薦文章。
要求：繁體中文、1500字以上、包含本地推薦、在地支付方式、地區優惠、FAQ。""",

    "credit": """你是信用版/現金版比較專家。請為「{keyword}」撰寫一篇版本介紹比較文章。
要求：繁體中文、1500字以上、包含信用版vs現金版差異、適合對象、風險提示、FAQ。""",

    "live": """你是真人荷官遊戲專家。請為「{keyword}」撰寫一篇真人荷官介紹文章。
要求：繁體中文、1500字以上、包含真人遊戲種類、直播品質、互動功能、推薦平台、FAQ。""",

    "community": """你是台灣PTT/Dcard論壇觀察者。請為「{keyword}」撰寫一篇社群口碑評測。
要求：繁體中文、1500字以上、模擬PTT風格整理網友評價、客觀總結正反面意見、FAQ。""",

    "terms": """你是博弈術語解釋專家。請為「{keyword}」撰寫一篇術語解釋文章。
要求：繁體中文、1000字以上、包含術語定義、使用場景、相關術語、FAQ。""",
}

def generate_content(client, keyword, kw_type):
    prompt_template = PROMPTS.get(kw_type, PROMPTS["brand"])
    prompt = prompt_template.format(keyword=keyword)

    response = client.messages.create(
        model="claude-sonnet-4-20250514",
        max_tokens=4096,
        messages=[{"role": "user", "content": prompt}]
    )
    return response.content[0].text

def main():
    parser = argparse.ArgumentParser(description='Generate AI content for BrandSite Pro')
    parser.add_argument('--domain-id', type=int, required=True)
    parser.add_argument('--db', default='../data/brandsite.db')
    parser.add_argument('--api-key', required=True)
    args = parser.parse_args()

    conn = sqlite3.connect(args.db)
    cur = conn.cursor()

    cur.execute("SELECT domain, keyword_type, primary_keyword FROM domains WHERE id=?", (args.domain_id,))
    row = cur.fetchone()
    if not row:
        print(f"Domain ID {args.domain_id} not found")
        sys.exit(1)

    domain, kw_type, keyword = row
    print(f"Generating content for: {domain} ({kw_type}: {keyword})")

    client = anthropic.Anthropic(api_key=args.api_key)
    content = generate_content(client, keyword, kw_type)

    # update content
    cur.execute("""INSERT INTO contents(domain_id, main_content, ai_generated) VALUES(?,?,1)
        ON CONFLICT(domain_id) DO UPDATE SET main_content=excluded.main_content, ai_generated=1, updated_at=CURRENT_TIMESTAMP""",
        (args.domain_id, content))
    conn.commit()
    print(f"Content saved ({len(content)} chars)")
    conn.close()

if __name__ == '__main__':
    main()
