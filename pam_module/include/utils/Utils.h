#include <vector>
#include <sys/types.h>
#include <nlohmann/json.hpp>

struct User {
    std::string Name;
    uid_t Uid;
    gid_t Gid;
    std::vector<std::string> Groups;
};

struct Host {
    std::string Ip;
    std::string HostName;
};

struct Time {
    std::string Timestamp;
    std::string Weekday;
};

struct Request {
    User user;
    Host host;
    std::string Service;
    std::string Action;
    Time time;
};

nlohmann::json BuildRequest(const Request& request) {
    return {
        {"user", {
            {"name", request.user.Name},
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

std::string GetUrl() {
    return "http://192.168.0.104/api/v1/decide";
}