#pragma once

#include <string>
#include <nlohmann/json.hpp>
#include "utils/Utils.h"


class Client {
    public:
        Client(const std::string &Url);
        Decision CheckAccess(const Request &request) const;
    private:
        const std::string url_;
};