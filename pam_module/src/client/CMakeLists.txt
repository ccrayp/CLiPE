add_library(client STATIC
    Client.cpp
)

find_package(nlohmann_json REQUIRED)

target_link_libraries(client
    PUBLIC utils
    PRIVATE nlohmann_json::nlohmann_json
)

target_include_directories(client
    PUBLIC ${CMAKE_SOURCE_DIR}/include
)