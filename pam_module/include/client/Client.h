#pragma once

#include <string>
#include <nlohmann/json.hpp>

#include "utils/ReturnStatus.h"

class Client {
    public:
        Client(const std::string &Url);
        Result CheckAccess(const Request &request) const;
    private:
        const std::string url_;
};