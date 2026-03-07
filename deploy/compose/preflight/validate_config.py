import argparse
import pathlib
import re
import sys


PLACEHOLDER_TOKENS = (
    "changeme",
    "change-me",
    "replace-me",
    "password",
    "example",
    "admin",
    "<set-",
    "<replace-",
)

REQUIRED_KEYS = (
    "TRAEFIK_DOMAIN",
    "TRAEFIK_ACME_EMAIL",
    "TRAEFIK_TLS_ENABLED",
    "POSTGRES_PASSWORD",
    "CLICKHOUSE_PASSWORD",
    "VALKEY_PASSWORD",
    "QDRANT_API_KEY",
    "CROWDSEC_ENROLL_KEY",
    "CROWDSEC_BOUNCER_KEY",
    "VAULT_ADDR",
    "VAULT_ROLE",
)

CREDENTIAL_KEYS = (
    "POSTGRES_PASSWORD",
    "CLICKHOUSE_PASSWORD",
    "VALKEY_PASSWORD",
    "QDRANT_API_KEY",
    "CROWDSEC_ENROLL_KEY",
    "CROWDSEC_BOUNCER_KEY",
)

SENSITIVE_NAME_PATTERN = re.compile(r"(PASSWORD|SECRET|TOKEN|API_KEY|ENROLL_KEY|BOUNCER_KEY|DSN)", re.IGNORECASE)
HARDCODED_ENV_PATTERN = re.compile(r"^\s*([A-Z0-9_]*(PASSWORD|SECRET|TOKEN|KEY)[A-Z0-9_]*)\s*:\s*(?!\$\{)(.+?)\s*$")


def parse_env_file(path: pathlib.Path) -> dict:
    values = {}
    for raw_line in path.read_text().splitlines():
        line = raw_line.strip()
        if not line or line.startswith("#"):
            continue
        if "=" not in line:
            raise ValueError(f"invalid env line: {raw_line}")
        key, value = line.split("=", 1)
        values[key.strip()] = value.strip()
    return values


def looks_placeholder(value: str) -> bool:
    normalized = value.strip().lower()
    return any(token in normalized for token in PLACEHOLDER_TOKENS)


def strong_enough(secret: str) -> bool:
    if len(secret) < 16:
        return False
    categories = 0
    categories += bool(re.search(r"[a-z]", secret))
    categories += bool(re.search(r"[A-Z]", secret))
    categories += bool(re.search(r"[0-9]", secret))
    categories += bool(re.search(r"[^A-Za-z0-9]", secret))
    return categories >= 3


def validate_env(values: dict) -> list:
    errors = []
    for key in REQUIRED_KEYS:
        if not values.get(key):
            errors.append(f"{key} is required")

    if values.get("TRAEFIK_TLS_ENABLED", "").lower() != "true":
        errors.append("TRAEFIK_TLS_ENABLED must be true for the production profile")

    if values.get("TRAEFIK_DOMAIN", "").startswith("http://") or values.get("TRAEFIK_DOMAIN", "").startswith("https://"):
        errors.append("TRAEFIK_DOMAIN must be a hostname without scheme")

    if values.get("VAULT_ADDR") and not values["VAULT_ADDR"].startswith("https://"):
        errors.append("VAULT_ADDR must use https://")

    seen_values = {}
    for key in CREDENTIAL_KEYS:
        value = values.get(key, "")
        if not value:
            continue
        if looks_placeholder(value):
            errors.append(f"{key} uses a placeholder value")
        if not strong_enough(value):
            errors.append(f"{key} must be at least 16 chars and contain 3 character classes")
        previous_key = seen_values.get(value)
        if previous_key:
            errors.append(f"{key} must not reuse the same credential as {previous_key}")
        seen_values[value] = key

    return errors


def validate_compose_file(path: pathlib.Path) -> list:
    errors = []
    contents = path.read_text()
    required_services = (
        "preflight:",
        "traefik:",
        "crowdsec:",
        "crowdsec-bouncer:",
        "trivy:",
        "postgres:",
        "clickhouse:",
        "valkey:",
        "qdrant:",
        "speclist-api:",
        "speclist-web:",
    )
    for service in required_services:
        if service not in contents:
            errors.append(f"compose file is missing required service {service[:-1]}")

    if "condition: service_completed_successfully" not in contents:
        errors.append("compose file must gate startup on preflight completion")

    for line_number, line in enumerate(contents.splitlines(), start=1):
        match = HARDCODED_ENV_PATTERN.match(line)
        if not match:
            continue
        key = match.group(1)
        value = match.group(3).strip().strip("\"'")
        if value.startswith("${"):
            continue
        if key and SENSITIVE_NAME_PATTERN.search(key) and value and not looks_placeholder(value):
            errors.append(f"compose file contains hardcoded sensitive value for {key} on line {line_number}")

    return errors


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate production compose security requirements.")
    parser.add_argument("--env-file", required=True)
    parser.add_argument("--compose-file", required=True)
    args = parser.parse_args()

    env_path = pathlib.Path(args.env_file)
    compose_path = pathlib.Path(args.compose_file)

    errors = []
    errors.extend(validate_env(parse_env_file(env_path)))
    errors.extend(validate_compose_file(compose_path))

    if errors:
        for error in errors:
            print(f"ERROR: {error}", file=sys.stderr)
        return 1

    print("preflight validation passed")
    return 0


if __name__ == "__main__":
    sys.exit(main())
