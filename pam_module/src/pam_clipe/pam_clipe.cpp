#include <string>
#include <security/pam_modules.h>

#include "logger/Logger.h"
#include "client/Client.h"

extern "C" int pam_sm_acct_mgmt(pam_handle_t *pamh, int flags, int argc, const char **argv)
{
    bool debug = false;

    for (int i = 0; i < argc; ++i) {
        if (std::string(argv[i]) == "debug") {
            debug = true;
        }
    }

    const char* user = nullptr;
    if (pam_get_user(pamh, &user, NULL) != PAM_SUCCESS || !user) {
        return PAM_AUTH_ERR;
    }

    const void* service = nullptr;
    if (pam_get_item(pamh, PAM_SERVICE, &service) != PAM_SUCCESS || !service) {
        return PAM_AUTH_ERR;
    }

    try {
        Request request = BuildRequest(
            user,
            static_cast<const char*>(service)
        );

        Logger logger(debug);

        std::string url = GetValue("URL") + "/access/api/v1/decide";
        Client client(url);

        Decision decision = client.CheckAccess(request);

        if (decision.allow) {
            logger.makeAllowLog(request, decision);
            return PAM_SUCCESS;
        } else {
            logger.makeDenyLog(request, decision);
            return PAM_PERM_DENIED;
        }
    }
    catch (const std::exception& e) {
        Logger logger(debug);

        Request fallback{};
        fallback.user.Name = user ? user : "unknown";
        fallback.Service = service ? static_cast<const char*>(service) : "unknown";

        logger.makeErrorLog(fallback, e.what());

        return PAM_PERM_DENIED;
    }
}

extern "C" int pam_sm_authenticate(pam_handle_t *pamh, int flags, int argc, const char **argv) {
    return PAM_IGNORE;
}

extern "C" int pam_sm_setcred(pam_handle_t *pamh, int flags, int argc, const char **argv) {
    return PAM_SUCCESS;
}