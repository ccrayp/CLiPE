add_library(pam_clipe SHARED
    pam_clipe.cpp
    ../client/Client.cpp
    ../logger/Logger.cpp
    ../utils/Utils.cpp
)

target_include_directories(pam_clipe PRIVATE
    ${CMAKE_SOURCE_DIR}/include
)

target_link_libraries(pam_clipe PRIVATE
    pam
    curl
)

set_target_properties(pam_clipe PROPERTIES
    PREFIX ""
    OUTPUT_NAME "pam_clipe"
)