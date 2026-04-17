add_library(pam_clipe SHARED
    pam_clipe.cpp
)

target_link_libraries(pam_clipe
    PRIVATE utils
    PRIVATE client
    PRIVATE pam
)

set_target_properties(pam_clipe PROPERTIES
    PREFIX ""
    OUTPUT_NAME "pam_clipe"
)