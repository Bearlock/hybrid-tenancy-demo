"""
Kafka consumer: on tenant.signups, if 'fact-app' is in services, create tenant DB and register in TenantDB.
"""
import json
import logging
from confluent_kafka import Consumer, KafkaError
from config import Config
from tenant_pool import create_tenant_database

logger = logging.getLogger(__name__)

def run_consumer():
    c = Consumer({
        "bootstrap.servers": Config.KAFKA_BOOTSTRAP,
        "group.id": Config.APP_NAME,
        "auto.offset.reset": "earliest",
    })
    c.subscribe([Config.KAFKA_TOPIC])
    try:
        while True:
            msg = c.poll(1.0)
            if msg is None:
                continue
            if msg.error():
                if msg.error().code() == KafkaError._PARTITION_EOF:
                    continue
                logger.error("Consumer error: %s", msg.error())
                continue
            try:
                evt = json.loads(msg.value().decode())
                tenant_id = evt.get("tenant_id")
                services = evt.get("services") or []
                if tenant_id and Config.APP_NAME in services:
                    create_tenant_database(tenant_id)
                    logger.info("Created tenant DB for %s", tenant_id)
            except Exception as e:
                logger.exception("Processing event: %s", e)
    finally:
        c.close()
