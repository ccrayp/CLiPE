#include "client/Client.h"
#include "utils/Utils.h"

Client::Client(const std::string &Url) : url_(Url) {}

Result Client::CheckAccess(const Request &request) const {
    
}