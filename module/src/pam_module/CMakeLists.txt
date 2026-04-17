add_library(pam_module SHARED
    pam_module.cpp
)

target_link_libraries(pam_module
    PRIVATE utils
    PRIVATE client
    PRIVATE pam
)

set_target_properties(pam_module PROPERTIES
    PREFIX ""
    OUTPUT_NAME "pam_my_module"
)