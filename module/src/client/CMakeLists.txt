add_library(client STATIC
    Client.cpp
)

target_link_libraries(client
    PUBLIC utils
)

target_include_directories(client
    PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/include
)