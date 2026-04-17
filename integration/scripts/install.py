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
        "ru": "CRUD_URL не задан",
        "en": "CRUD_URL is not set"
    },
    "host_error": {
        "ru": "Ошибка регистрации хоста: {err}",
        "en": "Host registration error: {err}"
    },
    "user_error": {
        "ru": "Ошибка создания пользователя {name}: {code}",
        "en": "User {name} failed: {code}"
    },
    "network_error": {
        "ru": "Ошибка сети для пользователя {name}: {err}",
        "en": "Network error for user {name}: {err}"
    },
    "host_id_error": {
        "ru": "Не удалось получить host_id из ответа",
        "en": "Failed to get host_id from response"
    },
    "progress_host": {
        "ru": "Регистрация хоста",
        "en": "Registering host"
    },
    "progress_users": {
        "ru": "Регистрация пользователей",
        "en": "Registering users"
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
    file_handler.setLevel(logging.INFO)

    console_handler = logging.StreamHandler()
    console_handler.setFormatter(formatter)
    console_handler.setLevel(logging.ERROR)

    logger.addHandler(file_handler)
    logger.addHandler(console_handler)


def is_real_user(user):
    if user.pw_uid < 1000:
        return False

    if "nologin" in user.pw_shell or "false" in user.pw_shell:
        return False

    return True


def get_ip():
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    try:
        s.connect(("8.8.8.8", 80))
        ip = s.getsockname()[0]
    finally:
        s.close()
    return ip


def register_host(base_url, ip):
    response = requests.post(
        f"{base_url}/hosts",
        json={"ip": ip},
        timeout=5
    )

    if response.status_code != 201:
        raise Exception(f"{response.status_code} {response.text}")

    data = response.json()

    try:
        return data["data"]["id"]
    except KeyError:
        raise Exception(t("host_id_error"))


def register_user(base_url, user, host_id):
    response = requests.post(
        f"{base_url}/users",
        json={
            "user_name": user.pw_name,
            "uid": user.pw_uid,
            "gid": user.pw_gid,
            "host_id": host_id
        },
        timeout=5
    )

    if response.status_code != 201:
        logging.error(t("user_error", name=user.pw_name, code=response.status_code))


def main():
    setup_logging()
    load_dotenv()

    base_url = os.getenv("CRUD_URL")

    if not base_url:
        raise ValueError(t("env_error"))

    base_url = base_url.rstrip("/")

    with tqdm(total=1, desc=t("progress_host")) as host_bar:
        try:
            ip = get_ip()
            host_id = register_host(base_url, ip)
        except Exception as e:
            logging.error(t("host_error", err=e))
            return
        finally:
            host_bar.update(1)

    users = pwd.getpwall()
    real_users = [u for u in users if is_real_user(u)]

    with tqdm(total=len(real_users), desc=t("progress_users")) as user_bar:
        for user in real_users:
            try:
                register_user(base_url, user, host_id)
            except requests.exceptions.RequestException as e:
                logging.error(t("network_error", name=user.pw_name, err=e))
            finally:
                user_bar.update(1)


if __name__ == "__main__":
    main()