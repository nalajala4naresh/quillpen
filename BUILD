load("@rules_go//go:def.bzl", "go_binary", "go_library")


go_library(
    name = "quillpen_lib",
    srcs = ["main.go"],
    importpath = "github.com/quillpen",
    visibility = ["//visibility:private"],
    deps = [
        "//pkg/accounts",
        "//pkg/chat",
        "//pkg/posts",
        "@com_github_gorilla_handlers//:go_default_library",
        "@com_github_gorilla_mux//:go_default_library",
        "@com_github_gorilla_sessions//:go_default_library",
    ],
)

go_binary(
    name = "quillpen",
    embed = [":quillpen_lib"],
    visibility = ["//visibility:public"],
)
