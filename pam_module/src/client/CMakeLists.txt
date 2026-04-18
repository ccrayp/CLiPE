add_library(client STATIC
    Client.cpp
)

target_include_directories(client PUBLIC
    ${CMAKE_SOURCE_DIR}/include
)

find_package(CURL REQUIRED)

target_link_libraries(client PUBLIC
    CURL::libcurl
)