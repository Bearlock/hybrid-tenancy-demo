from flask import Flask, request, jsonify
from config import Config
from tenant_pool import get_tenant_conn

app = Flask(__name__)

FACTS_SCHEMA = """
CREATE TABLE IF NOT EXISTS facts (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
)
"""

def tenant_conn():
    tenant_id = request.headers.get("X-Tenant-ID")
    if not tenant_id:
        return None, "missing X-Tenant-ID", 401
    conn = get_tenant_conn(tenant_id)
    if not conn:
        return None, "tenant not found", 404
    return conn, None, None

@app.route("/facts", methods=["GET"])
def list_facts():
    conn, err, code = tenant_conn()
    if err:
        return jsonify({"error": err}), code
    try:
        with conn.cursor() as cur:
            cur.execute("SELECT id, content, created_at FROM facts ORDER BY id")
            rows = cur.fetchall()
        return jsonify([{"id": r[0], "content": r[1], "created_at": str(r[2])} for r in rows])
    finally:
        conn.close()

@app.route("/facts", methods=["POST"])
def create_fact():
    conn, err, code = tenant_conn()
    if err:
        return jsonify({"error": err}), code
    data = request.get_json() or {}
    content = data.get("content") or ""
    if not content.strip():
        return jsonify({"error": "content required"}), 400
    try:
        with conn.cursor() as cur:
            cur.execute("INSERT INTO facts (content) VALUES (%s) RETURNING id, content, created_at", (content,))
            row = cur.fetchone()
        conn.commit()
        return jsonify({"id": row[0], "content": row[1], "created_at": str(row[2])}), 201
    finally:
        conn.close()

@app.route("/facts/<int:fact_id>", methods=["GET"])
def get_fact(fact_id):
    conn, err, code = tenant_conn()
    if err:
        return jsonify({"error": err}), code
    try:
        with conn.cursor() as cur:
            cur.execute("SELECT id, content, created_at FROM facts WHERE id = %s", (fact_id,))
            row = cur.fetchone()
        if not row:
            return jsonify({"error": "not found"}), 404
        return jsonify({"id": row[0], "content": row[1], "created_at": str(row[2])})
    finally:
        conn.close()

@app.route("/facts/<int:fact_id>", methods=["PUT"])
def update_fact(fact_id):
    conn, err, code = tenant_conn()
    if err:
        return jsonify({"error": err}), code
    data = request.get_json() or {}
    content = data.get("content")
    if content is None:
        return jsonify({"error": "content required"}), 400
    try:
        with conn.cursor() as cur:
            cur.execute("UPDATE facts SET content = %s WHERE id = %s RETURNING id, content, created_at", (content, fact_id))
            row = cur.fetchone()
        conn.commit()
        if not row:
            return jsonify({"error": "not found"}), 404
        return jsonify({"id": row[0], "content": row[1], "created_at": str(row[2])})
    finally:
        conn.close()

@app.route("/facts/<int:fact_id>", methods=["DELETE"])
def delete_fact(fact_id):
    conn, err, code = tenant_conn()
    if err:
        return jsonify({"error": err}), code
    try:
        with conn.cursor() as cur:
            cur.execute("DELETE FROM facts WHERE id = %s", (fact_id,))
        conn.commit()
        return "", 204
    finally:
        conn.close()

def main():
    from tenant_db import ensure_tenant_registry
    ensure_tenant_registry(Config.TENANT_DB_CONN)
    app.run(host="0.0.0.0", port=Config.HTTP_PORT, debug=True)

if __name__ == "__main__":
    main()
