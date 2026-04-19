import pwd
import socket
import subprocess
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
    "host_remove_error": {
        "ru": "Ошибка удаления хоста: {err}",
        "en": "Host uninstall error: {err}"
    },
    "user_network_error": {
        "ru": "Ошибка сети для пользователя {name}: {err}",
        "en": "Network error for user {name}: {err}"
    },
    "host_delete_fail": {
        "ru": "Не удалось удалить хост {id}: {code}",
        "en": "Failed to delete host {id}: {code}"
    },
    "user_delete_fail": {
        "ru": "Не удалось удалить пользователя {id}: {code}",
        "en": "Failed to delete user {id}: {code}"
    },
    "pam_in_use": {
    "ru": "Модуль используется в PAM конфигурации: {files}. Удалите его из этих файлов перед удалением.",
    "en": "Module is used in PAM config: {files}. Remove it from these files before uninstall."
    },
    "pam_not_found": {
        "ru": "Модуль не найден в PAM конфигурации",
        "en": "Module not found in PAM configuration"
    },
    "module_removed": {
        "ru": "Модуль удалён: {path}",
        "en": "Module removed: {path}"
    },
    "module_not_found": {
        "ru": "Модуль не найден, удалять нечего",
        "en": "Module not found, nothing to remove"
    },
    "pam_abort": {
        "ru": "Модуль используется в PAM ({files}). Удаление прервано.",
        "en": "Module is used in PAM ({files}). Uninstall aborted."
    }
}


def t(key, **kwargs):
    return MESSAGES[key][LANG].format(**kwargs)


def find_pam_usage():
    pam_dir = "/etc/pam.d"
    used_in = []

    try:
        for filename in os.listdir(pam_dir):
            file_path = os.path.join(pam_dir, filename)

            if not os.path.isfile(file_path):
                continue

            try:
                with open(file_path, "r") as f:
                    for line in f:
                        if "pam_clipe.so" in line:
                            used_in.append(filename)
                            break
            except Exception:
                continue

    except Exception as e:
        logging.error(f"Failed to scan /etc/pam.d: {e}")

    return used_in


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


def find_host(base_url, ip):
    response = requests.get(
        f"{base_url}/hosts?limit=1&offset=0",
        json={"ip": ip},
        timeout=5
    )
    response.raise_for_status()
    return response.json().get("data", {}).get("hosts", [])


def delete_host(base_url, host_id):
    response = requests.delete(
        f"{base_url}/hosts/{host_id}",
        timeout=5
    )

    if response.status_code != 200:
        logging.error(t("host_delete_fail", id=host_id, code=response.status_code))


def find_user(base_url, user):
    response = requests.get(
        f"{base_url}/users?limit=1&offset=0",
        json={
            "user_name": user.pw_name,
            "uid": user.pw_uid,
            "gid": user.pw_gid
        },
        timeout=5
    )
    response.raise_for_status()
    return response.json().get("data", {}).get("users", [])


def delete_user(base_url, user_id):
    response = requests.delete(
        f"{base_url}/users/{user_id}",
        timeout=5
    )

    if response.status_code != 200:
        logging.error(t("user_delete_fail", id=user_id, code=response.status_code))



CONFIG_DIR = "/etc/clipe"


def remove_config():
    try:
        if os.path.exists(CONFIG_DIR):
            for filename in os.listdir(CONFIG_DIR):
                file_path = os.path.join(CONFIG_DIR, filename)
                try:
                    os.remove(file_path)
                except Exception as e:
                    logging.error(f"Failed to remove file {file_path}: {e}")

            os.rmdir(CONFIG_DIR)

    except PermissionError:
        logging.error("Нет прав на удаление /etc/clipe (запусти с sudo)")
    except Exception as e:
        logging.error(f"Ошибка удаления конфига: {e}")


def get_pam_dir():
    return "/lib/aarch64-linux-gnu/security/"


def uninstall_module():
    try:
        used_files = find_pam_usage()

        if used_files:
            logging.error(
                t("pam_in_use", files=", ".join(used_files))
            )
            return
        else:
            logging.info(t("pam_not_found"))

        pam_dir = get_pam_dir()
        module_path = os.path.join(pam_dir, "pam_clipe.so")

        if not os.path.exists(module_path):
            logging.info(t("module_not_found"))
            return

        os.remove(module_path)
        logging.info(t("module_removed", path=module_path))

    except PermissionError:
        logging.error("Permission denied (run with sudo)")
    except Exception as e:
        logging.error(f"Failed to remove module: {e}")


def main():
    setup_logging()
    load_dotenv()

    used_files = find_pam_usage()
    if used_files:
        logging.error(t("pam_abort", files=", ".join(used_files)))
        return

    base_url = os.getenv("URL")

    if not base_url:
        raise ValueError(t("env_error"))

    base_url = base_url.rstrip("/")
    base_url = base_url + ":8081/api/v1/internal"

    users = pwd.getpwall()
    real_users = [u for u in users if is_real_user(u)]

    with tqdm(total=len(real_users), desc="Removing users") as user_bar:
        for user in real_users:
            try:
                found_users = find_user(base_url, user)

                for u in found_users:
                    delete_user(base_url, u["user_id"])

            except requests.exceptions.RequestException as e:
                logging.error(t("user_network_error", name=user.pw_name, err=e))

            finally:
                user_bar.update(1)

    with tqdm(total=1, desc="Removing host") as host_bar:
        try:
            ip = get_ip()
            hosts = find_host(base_url, ip)

            for host in hosts:
                delete_host(base_url, host["host_id"])

        except Exception as e:
            logging.error(t("host_remove_error", err=e))

        finally:
            host_bar.update(1)
    
    remove_config()
    uninstall_module()

if __name__ == "__main__":
    main()