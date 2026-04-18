#pragma once

#include <string>
#include <syslog.h>
#include "utils/Utils.h"

class Logger {
public:
    Logger(bool debug);
    ~Logger();

    void logDecision(int level, const Request& request, const Decision& decision);

    void makeAllowLog(const Request& request, const Decision& decision);
    void makeDenyLog(const Request& request, const Decision& decision);
    void makeErrorLog(const Request& request, const std::string& error);
private:
    bool debug_;
};