from flask import Flask, jsonify
from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from bs4 import BeautifulSoup
from flasgger import Swagger
import threading
import time

app = Flask(__name__)

# Swagger config
swagger_template = {
    "swagger": "2.0",
    "info": {
        "title": "Call Center Voting API",
        "description": "API that caches each ID on demand and refreshes it every 5 seconds.",
        "version": "1.0.0"
    },
    "basePath": "/",
    "schemes": ["http"]
}
swagger_config = {
    "headers": [],
    "specs": [{
        "endpoint": 'apispec_1',
        "route": '/apispec_1.json',
        "rule_filter": lambda rule: True,
        "model_filter": lambda tag: True
    }],
    "static_url_path": "/flasgger_static",
    "swagger_ui": True,
    "specs_route": "/apidocs/"
}
Swagger(app, template=swagger_template, config=swagger_config)

# Global cache and update queue
cache = {}
update_queue = set()

def fetch_votes(poll_id):
    """Fetch vote data from remote site"""
    url = f"http://vos.uztelecom.uz:3124/Televoting/counters/counters_ext.jsp?pollName={poll_id}"
    options = Options()
    options.add_argument("--headless")
    options.add_argument("--no-sandbox")
    options.add_argument("--disable-dev-shm-usage")
    options.binary_location = "/usr/bin/google-chrome"
    service = Service("/usr/bin/chromedriver")
    driver = webdriver.Chrome(service=service, options=options)

    try:
        driver.get(url)
        time.sleep(3)
        html = driver.page_source
    finally:
        driver.quit()

    soup = BeautifulSoup(html, "html.parser")
    rows = soup.select("table tr")

    results = []
    for row in rows:
        cols = row.find_all("td")
        if len(cols) >= 2 and cols[0].text.strip().isdigit():
            results.append({
                "number": cols[0].text.strip(),
                "votes": cols[1].text.strip()
            })
    return results

def background_refresher():
    """Update all cached IDs every 5 seconds"""
    while True:
        for poll_id in list(update_queue):
            try:
                result = fetch_votes(poll_id)
                cache[poll_id] = result
                print(f"[INFO] Cache refreshed: {poll_id}")
            except Exception as e:
                print(f"[ERROR] Failed to update {poll_id}: {e}")
        time.sleep(5)

@app.route("/<id>", methods=["GET"])
def get_by_id(id):
    """
    Get voting result by ID
    ---
    tags:
      - Votes
    parameters:
      - name: id
        in: path
        type: string
        required: true
        description: ID given by call center 
    responses:
      200:
        description: Vote results
        schema:
          type: array
          items:
            type: object
            properties:
              number:
                type: string
              votes:
                type: string
    """
    if id in cache:
        print(f"[INFO] Cache HIT: {id}")
        return jsonify(cache[id])
    try:
        print(f"[INFO] Cache MISS: {id}. Fetching and enabling auto-refresh.")
        result = fetch_votes(id)
        cache[id] = result
        update_queue.add(id)
        return jsonify(result)
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":
    threading.Thread(target=background_refresher, daemon=True).start()
    app.run(host="0.0.0.0", port=5000)
