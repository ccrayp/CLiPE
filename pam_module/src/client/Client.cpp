#include "client/Client.h"

#include <curl/curl.h>
#include <stdexcept>
#include <string>
#include <nlohmann/json.hpp>

using json = nlohmann::json;

static size_t WriteCallback(void* contents, size_t size, size_t nmemb, void* userp)
{
    ((std::string*)userp)->append((char*)contents, size * nmemb);
    return size * nmemb;
}

static Decision ParseDecision(const std::string& response)
{
    auto j = json::parse(response);

    Decision decision{};

    if (!j.contains("result")) {
        throw std::runtime_error("Missing field: result");
    }

    decision.allow = (j["result"] == "ALLOW");

    if (!j.contains("journal")) {
        throw std::runtime_error("Missing field: journal");
    }

    const auto& journal = j["journal"];

    decision.journal.policy_id   = journal.at("policy_id").get<int>();
    decision.journal.policy_name = journal.at("policy_name").get<std::string>();
    decision.journal.request_id  = journal.at("request_id").get<int>();
    decision.journal.decision_id = journal.at("decision_id").get<int>();

    return decision;
}

Client::Client(const std::string &Url) : url_(Url) {
    curl_global_init(CURL_GLOBAL_DEFAULT);
}

Decision Client::CheckAccess(const Request &request) const {
    CURL* curl = curl_easy_init();
    if (!curl) {
        throw std::runtime_error("curl init failed");
    }

    std::string response;
    struct curl_slist* headers = nullptr;

    try {
        json j = RequestToJson(request);
        std::string body = j.dump();

        headers = curl_slist_append(headers, "Content-Type: application/json");

        curl_easy_setopt(curl, CURLOPT_URL, url_.c_str());
        curl_easy_setopt(curl, CURLOPT_POST, 1L);
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body.c_str());
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);

        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response);

        curl_easy_setopt(curl, CURLOPT_TIMEOUT, 3L);
        curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT, 1L);

        CURLcode res = curl_easy_perform(curl);

        if (res != CURLE_OK) {
            throw std::runtime_error("HTTP request failed");
        }

        long http_code = 0;
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &http_code);

        if (http_code != 200) {
            throw std::runtime_error("Bad HTTP response: " + std::to_string(http_code));
        }

        Decision decision = ParseDecision(response);

        curl_slist_free_all(headers);
        curl_easy_cleanup(curl);

        return decision;
    }
    catch (...) {
        curl_slist_free_all(headers);
        curl_easy_cleanup(curl);
        throw;
    }
}