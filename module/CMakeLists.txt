cmake_minimum_required(VERSION 3.16)
project(clipe)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

add_subdirectory(src/utils)
add_subdirectory(src/client)
add_subdirectory(src/pam_module)
add_subdirectory(tests)