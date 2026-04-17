#include <string>
#include <security/pam_modules.h>

#include "logger/Logger.h"
#include "client/Client.h"
#include "utils/Utils.h"
#include "utils/ReturnStatus.h"

extern "C" int pam_sm_acct_mgmt(pam_handle_t *pamh, int flags, int argc, const char **argv) {

    Logger logger();

    std::string url = getUrl();
    Client client(url);
    
    Request request = buildRequest(pamh);
    Result result = client.CheckAccess(request);

    switch(result) {
        case Result.Allow: {
            logger.makeAllowLog(/*some data*/);
            return PAM_SUCCESS;
        }
        case Result.Deny: {
            logger.makeDenyLog(/*some data*/);
            return PAM_PERM_DENIED;
        }
        case Result.Error: {
            logger.makeErrorLog(/*some data*/);
            return PAM_AUTH_ERR;
        }        
    }

    return PAM_AUTH_ERR;
}

extern "C" int pam_sm_authenticate(pam_handle_t *pamh, int flags, int argc, const char **argv) {
    return PAM_IGNORE;
}

extern "C" int pam_sm_setcred(pam_handle_t *pamh, int flags, int argc, const char **argv) {
    return PAM_SUCCESS;
}