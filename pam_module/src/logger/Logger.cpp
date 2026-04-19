#include "logger/Logger.h"

Logger::Logger(bool debug) : debug_(debug) {
    openlog("pam_clipe", LOG_PID, LOG_AUTHPRIV);
}

Logger::~Logger() {
    closelog();
}

void Logger::logDecision(int level, const Request& request, const Decision& decision) {
    if (debug_) {
        syslog(level,
            "user=%s uid=%d gid=%d service=%s ip=%s "
            "result=%s policy_id=%d policy_name=%s request_id=%d decision_id=%d",
            request.user.Name.c_str(),
            static_cast<int>(request.user.Uid),
            static_cast<int>(request.user.Gid),
            request.Service.c_str(),
            request.host.Ip.empty() ? "unknown" : request.host.Ip.c_str(),
            decision.allow ? "ALLOW" : "DENY",
            decision.journal.policy_id,
            decision.journal.policy_name.empty() ? "unknown" : decision.journal.policy_name.c_str(),
            decision.journal.request_id,
            decision.journal.decision_id
        );
    } else {
        syslog(level,
            "user=%s service=%s result=%s",
            request.user.Name.c_str(),
            request.Service.c_str(),
            decision.allow ? "ALLOW" : "DENY"
        );
    }
}

void Logger::makeAllowLog(const Request& request, const Decision& decision) {
    logDecision(LOG_INFO, request, decision);
}

void Logger::makeDenyLog(const Request& request, const Decision& decision) {
    logDecision(LOG_WARNING, request, decision);
}

void Logger::makeErrorLog(const Request& request, const std::string& error) {
    if (debug_) {
        syslog(LOG_ERR,
            "user=%s uid=%d service=%s ip=%s error=%s",
            request.user.Name.c_str(),
            (int)request.user.Uid,
            request.Service.c_str(),
            request.host.Ip.c_str(),
            error.c_str()
        );
    } else {
        syslog(LOG_ERR,
            "user=%s service=%s error=%s",
            request.user.Name.c_str(),
            request.Service.c_str(),
            error.c_str()
        );
    }
}