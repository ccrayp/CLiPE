import pwd
import socket
import requests
import os
import logging
from dotenv import load_dotenv
from logging.handlers import RotatingFileHandler
from tqdm import tqdm


def get_lang():
    lang = os.getenv("LC_ALL") or os.getenv("LANG") or "en"
    return "ru" if lang.lower().startswith("ru") else "en"


LANG = get_lang()

MESSAGES = {
    "env_error": {
        "ru": "URL не задан",
        "en": "URL is not set"
    },
    "host_not_found": {
        "ru": "Хост не найден",
        "en": "Host not found"
    },
    "network_error": {
        "ru": "Ошибка сети: {err}",
        "en": "Network error: {err}"
    },
    "user_delete_fail": {
        "ru": "Не удалось удалить пользователя {id}: {code}",
        "en": "Failed to delete user {id}: {code}"
    },
    "user_create_fail": {
        "ru": "Не удалось создать пользователя {name}: {code}",
        "en": "Failed to create user {name}: {code}"
    },
    "progress_sync": {
        "ru": "Синхронизация пользователей",
        "en": "Syncing users"
    }
}


def t(key, **kwargs):
    return MESSAGES[key][LANG].format(**kwargs)


def setup_logging():
    logger = logging.getLogger()
    logger.setLevel(logging.INFO)

    formatter = logging.Formatter(
        "%(asctime)s [%(levelname)s] %(message)s"
    )

    file_handler = RotatingFileHandler(
        "app.log",
        maxBytes=1_000_000,
        backupCount=3
    )
    file_handler.setFormatter(formatter)

    console_handler = logging.StreamHandler()
    console_handler.setLevel(logging.ERROR)
    console_handler.setFormatter(formatter)

    logger.addHandler(file_handler)
    logger.addHandler(console_handler)


def is_real_user(user):
    return user.pw_uid >= 1000 and "nologin" not in user.pw_shell and "false" not in user.pw_shell


def get_ip():
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    try:
        s.connect(("8.8.8.8", 80))
        return s.getsockname()[0]
    finally:
        s.close()


def get_host(base_url, headers, ip):
    r = requests.post(
        f"{base_url}/hosts/search",
        params={"limit": 1, "offset": 0},
        json={
            "ip": ip
        },
        timeout=5,
        verify="/usr/local/share/ca-certificates/clipe-ca.crt",
        headers=headers
    )
    r.raise_for_status()

    hosts = r.json().get("data", {}).get("hosts", [])
    if not hosts:
        raise Exception(t("host_not_found"))

    return hosts[0]["host_id"]


def get_users(base_url, headers, host_id):
    r = requests.post(
        f"{base_url}/users/search",
        params={"limit": 1000, "offset": 0},
        json={
            "host_id": host_id
        },
        timeout=5,
        verify="/usr/local/share/ca-certificates/clipe-ca.crt",
        headers=headers
    )
    r.raise_for_status()

    return r.json().get("data", {}).get("users", [])


def create_user(base_url, headers, user, host_id):
    r = requests.post(
        f"{base_url}/users",
        json={
            "user_name": user.pw_name,
            "uid": user.pw_uid,
            "gid": user.pw_gid,
            "host_id": host_id
        },
        timeout=5,
        verify="/usr/local/share/ca-certificates/clipe-ca.crt",
        headers=headers
    )

    if r.status_code != 201:
        logging.error(t("user_create_fail", name=user.pw_name, code=r.status_code))


def delete_user(base_url, headers, user_id):
    r = requests.delete(
        f"{base_url}/users/{user_id}",
        timeout=5,
        verify="/usr/local/share/ca-certificates/clipe-ca.crt",
        headers=headers
    )

    if r.status_code != 200:
        logging.error(t("user_delete_fail", id=user_id, code=r.status_code))


def main():
    setup_logging()
    load_dotenv()

    headers = {
        "X-Internal-Token": os.getenv("INSTALLER_TOKEN"),
        "X-Caller": os.getenv("INSTALLER_ID")
    }

    base_url = os.getenv("URL")
    if not base_url:
        raise ValueError(t("env_error"))

    base_url = base_url.rstrip("/")
    base_url = base_url + "/api/v1"

    try:
        ip = get_ip()
        host_id = get_host(base_url, headers, ip)
    except Exception as e:
        logging.error(str(e))
        return

    local_users = {
        u.pw_name: u for u in pwd.getpwall() if is_real_user(u)
    }

    try:
        remote_users = get_users(base_url, headers, host_id)
    except Exception as e:
        logging.error(t("network_error", err=e))
        return

    remote_map = {u["user_name"]: u for u in remote_users}

    to_create = [u for name, u in local_users.items() if name not in remote_map]
    to_delete = [u for name, u in remote_map.items() if name not in local_users]

    total = len(to_create) + len(to_delete)

    with tqdm(total=total, desc=t("progress_sync")) as bar:

        for user in to_create:
            try:
                create_user(base_url, headers, user, host_id)
            except Exception as e:
                logging.error(str(e))
            finally:
                bar.update(1)

        for user in to_delete:
            try:
                delete_user(base_url, headers, user["user_id"])
            except Exception as e:
                logging.error(str(e))
            finally:
                bar.update(1)


if __name__ == "__main__":
    main()