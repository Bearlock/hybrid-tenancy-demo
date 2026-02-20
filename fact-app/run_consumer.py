#!/usr/bin/env python3
"""Run the Kafka consumer in a separate process to create tenant DBs on signup events."""
import logging
from config import Config
from tenant_db import ensure_tenant_registry
from consumer import run_consumer

logging.basicConfig(level=logging.INFO)
ensure_tenant_registry(Config.TENANT_DB_CONN)
run_consumer()
