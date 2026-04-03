#!/usr/bin/env python3
"""
排名檢查腳本 - 檢查域名在Google的排名
用法: python check_rankings.py --db ../data/brandsite.db [--domain example.com] [--category brand]
"""
import argparse
import sqlite3
import json
import time
import sys
import urllib.request
import urllib.parse

def check_rank_via_api(keyword, domain, api_key=None):
    """Check keyword ranking via SerpAPI."""
    print(f"  Checking: {keyword} -> {domain}")
    if not api_key:
        return {"position": 0, "url": ""}
    try:
        import urllib.request, urllib.parse
        params = urllib.parse.urlencode({
            "q": keyword, "gl": "tw", "hl": "zh-tw",
            "num": 100, "api_key": api_key,
        })
        url = f"https://serpapi.com/search.json?{params}"
        req = urllib.request.Request(url, headers={"User-Agent": "BrandSitePro/1.0"})
        with urllib.request.urlopen(req, timeout=15) as resp:
            data = json.loads(resp.read())
        for i, r in enumerate(data.get("organic_results", []), 1):
            if domain.lower() in r.get("link", "").lower():
                return {"position": i, "url": r.get("link", "")}
    except Exception as e:
        print(f"  Warning: SerpAPI error: {e}")
    return {"position": 0, "url": ""}

def main():
    parser = argparse.ArgumentParser(description='Check keyword rankings')
    parser.add_argument('--db', default='../data/brandsite.db')
    parser.add_argument('--domain', help='Specific domain to check')
    parser.add_argument('--category', help='Keyword category filter')
    parser.add_argument('--limit', type=int, default=50, help='Max keywords to check')
    parser.add_argument('--api-key', default='', help='SERP API key')
    args = parser.parse_args()

    conn = sqlite3.connect(args.db)
    cur = conn.cursor()

    where = "d.status IN ('active','built')"
    params = []
    if args.domain:
        where += " AND d.domain = ?"
        params.append(args.domain)

    cur.execute(f"""SELECT d.id, d.domain, d.primary_keyword, d.keyword_type
        FROM domains d WHERE {where} LIMIT ?""", params + [args.limit])

    domains = cur.fetchall()
    print(f"Checking {len(domains)} domains...")

    for did, domain, keyword, kw_type in domains:
        if not keyword:
            continue
        result = check_rank_via_api(keyword, domain, args.api_key)
        if result["position"] > 0:
            cur.execute("""INSERT INTO ranking_history(domain_id, keyword, position, url)
                VALUES(?,?,?,?)""", (did, keyword, result["position"], result["url"]))
            print(f"  {domain}: #{result['position']} for '{keyword}'")
        time.sleep(2)  # rate limit

    conn.commit()
    conn.close()
    print("Done!")

if __name__ == '__main__':
    main()
