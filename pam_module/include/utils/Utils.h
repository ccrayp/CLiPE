#pragma once

#include <pwd.h>
#include <unistd.h>
#include <grp.h>
#include <sys/types.h>
#include <vector>
#include <string>
#include <cstring>
#include <ctime>
#include <sstream>
#include <fstream>
#include <iomanip>
#include <netdb.h>
#include <arpa/inet.h>
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

struct Journal {
    int policy_id;
    std::string policy_name;
    int request_id;
    int decision_id;
};

struct Decision {
    bool allow;
    Journal journal;
};

Request BuildRequest(const std::string& username, const std::string& service);
nlohmann::json RequestToJson(const Request& request);
std::string GetValue(const std::string &key);