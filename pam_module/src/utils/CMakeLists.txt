add_library(utils STATIC
    Utils.cpp
)

target_include_directories(utils
    PUBLIC ${CMAKE_SOURCE_DIR}/include
)