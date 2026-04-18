#include "utils/Utils.h"

std::string GetValue(const std::string &key)
{
    std::ifstream file("/etc/clipe/clipe.conf");
    if (!file.is_open()) {
        throw std::runtime_error("Cannot open config");
    }

    std::string line;

    while (std::getline(file, line)) {
        if (line.rfind(key + "=", 0) == 0) {
            return line.substr(key.length() + 1);
        }
    }

    throw std::runtime_error(key + " not found in config");
}

Request BuildRequest(const std::string& username, const std::string& service, const std::string& action) {
    Request req{};

    struct passwd* pwd = getpwnam(username.c_str());
    if (!pwd) {
        throw std::runtime_error("Failed to get user info");
    }

    req.user.Name = username;
    req.user.Uid = pwd->pw_uid;
    req.user.Gid = pwd->pw_gid;

    int ngroups = 0;
    getgrouplist(username.c_str(), pwd->pw_gid, nullptr, &ngroups);

    std::vector<gid_t> groups(ngroups);
    getgrouplist(username.c_str(), pwd->pw_gid, groups.data(), &ngroups);

    for (gid_t gid : groups) {
        struct group* grp = getgrgid(gid);
        if (grp) {
            req.user.Groups.push_back(grp->gr_name);
        }
    }

    char hostname[256];
    gethostname(hostname, sizeof(hostname));
    req.host.HostName = hostname;
    req.host.Ip = GetValue("IP");

    req.Service = service;
    req.Action = action;

    std::time_t now = std::time(nullptr);
    std::tm* gmt = std::gmtime(&now);

    std::ostringstream ts;
    ts << std::put_time(gmt, "%Y-%m-%dT%H:%M:%SZ");
    req.time.Timestamp = ts.str();

    static const char* days[] = {
        "sunday", "monday", "tuesday", "wednesday",
        "thursday", "friday", "saturday"
    };
    req.time.Weekday = days[gmt->tm_wday];

    return req;
}

nlohmann::json RequestToJson(const Request& request) {
    return {
        {"user", {
            {"username", request.user.Name},
            {"uid", request.user.Uid},
            {"gid", request.user.Gid},
            {"groups", request.user.Groups}
        }},
        {"host", {
            {"ip", request.host.Ip},
            {"hostname", request.host.HostName}
        }},
        {"service", request.Service},
        {"action", request.Action},
        {"time", {
            {"timestamp", request.time.Timestamp},
            {"weekday", request.time.Weekday}
        }}
    };
}