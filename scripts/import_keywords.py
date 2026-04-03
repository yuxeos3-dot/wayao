#!/usr/bin/env python3
"""
關鍵詞匯入腳本 - 從CSV/Excel匯入關鍵詞到SQLite
用法: python import_keywords.py --file keywords.csv --market zh-TW --db ../data/brandsite.db
"""
import argparse
import csv
import sqlite3
import os
import sys

def main():
    parser = argparse.ArgumentParser(description='Import keywords to BrandSite Pro')
    parser.add_argument('--file', required=True, help='CSV file path')
    parser.add_argument('--market', default='zh-TW', help='Market code')
    parser.add_argument('--db', default='../data/brandsite.db', help='SQLite database path')
    parser.add_argument('--batch-size', type=int, default=5000, help='Batch insert size')
    args = parser.parse_args()

    if not os.path.exists(args.file):
        print(f"Error: file not found: {args.file}")
        sys.exit(1)

    conn = sqlite3.connect(args.db)
    cur = conn.cursor()

    # ensure table exists
    cur.execute("""CREATE TABLE IF NOT EXISTS keywords (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        keyword TEXT NOT NULL,
        category TEXT NOT NULL,
        market TEXT NOT NULL DEFAULT 'zh-TW',
        volume INTEGER DEFAULT 0,
        kd INTEGER DEFAULT 0,
        cpc REAL DEFAULT 0,
        top1_dr INTEGER DEFAULT 0,
        assigned_to INTEGER,
        is_assigned INTEGER DEFAULT 0,
        UNIQUE(keyword, market)
    )""")

    total_before = conn.total_changes
    batch = []
    total_rows = 0

    with open(args.file, 'r', encoding='utf-8-sig') as f:
        reader = csv.reader(f)
        header = next(reader, None)  # skip header
        print(f"Header: {header}")

        for row in reader:
            if len(row) < 2:
                continue
            keyword = row[0].strip()
            category = row[1].strip()
            if not keyword or not category:
                continue

            market = row[2].strip() if len(row) > 2 and row[2].strip() else args.market
            volume = int(row[3]) if len(row) > 3 and row[3].strip().isdigit() else 0
            kd = int(row[4]) if len(row) > 4 and row[4].strip().isdigit() else 0
            try:
            cpc = float(row[5]) if len(row) > 5 and row[5].strip() else 0.0
        except ValueError:
            cpc = 0.0

            batch.append((keyword, category, market, volume, kd, cpc))
            total_rows += 1

            if len(batch) >= args.batch_size:
                cur.executemany(
                    "INSERT OR IGNORE INTO keywords(keyword, category, market, volume, kd, cpc) VALUES(?,?,?,?,?,?)",
                    batch
                )
                conn.commit()
                print(f"  processed {total_rows} rows...")
                batch = []

    if batch:
        cur.executemany(
            "INSERT OR IGNORE INTO keywords(keyword, category, market, volume, kd, cpc) VALUES(?,?,?,?,?,?)",
            batch
        )
        conn.commit()

    inserted = conn.total_changes - total_before
    print(f"\nDone! Total rows: {total_rows}, Inserted: {inserted}")

    # show category stats
    cur.execute("SELECT category, COUNT(*) FROM keywords WHERE market=? GROUP BY category ORDER BY COUNT(*) DESC", (args.market,))
    print("\nCategory stats:")
    for cat, cnt in cur.fetchall():
        print(f"  {cat}: {cnt:,}")

    conn.close()

if __name__ == '__main__':
    main()
