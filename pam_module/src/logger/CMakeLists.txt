add_library(logger STATIC
    Logger.cpp
)

target_include_directories(logger PUBLIC
    ${CMAKE_SOURCE_DIR}/include
)