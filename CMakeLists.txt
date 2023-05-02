cmake_minimum_required(VERSION 3.0.0)

project(bot)

add_executable(start main.cpp)

include_directories(include)

add_subdirectory(deps/dpp)